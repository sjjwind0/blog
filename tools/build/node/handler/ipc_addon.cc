#include <node.h>
#include <v8.h>
#include <uv.h>

#include <map>
#include <string>
#include <cstring>
#include <thread>
#include <functional>
#include <iostream>
#include <unistd.h>
#include <semaphore.h>
#include <pthread.h>
#include "export.h"

using v8::FunctionCallbackInfo;
using v8::Isolate;
using v8::Local;
using v8::Object;
using v8::String;
using v8::Number;
using v8::Integer;
using v8::Value;
using v8::Function;
using v8::Handle;
using v8::Context;
using v8::HandleScope;
using v8::Persistent;
using v8::MaybeLocal;
using v8::Exception;

namespace ipc {

struct NodeUserInfo {
  uv_work_t* req;
  std::function<void()> cb;
  int status;
  std::string result;
  pthread_mutex_t mutex;
};

const static char* kOnAcceptNewClient = "onAcceptNewClient";
const static char* kOnClientClose = "onClientClose";
const static char* kOnConnect = "onConnect";
const static char* kOnServerClose = "onServerClose";

std::map<void*, Persistent<Object> > delegateMap;
std::map<int32_t, Persistent<Function> > methodMap;
std::map<void*, int32_t> ipcMap;
std::map<int32_t, void*> indexMap;
int currentIndex = 0;

void doing_work (uv_work_t *req) {
	// do nothing
}

void after_doing_work(uv_work_t *req, int status) {
    NodeUserInfo* info = (NodeUserInfo*)req->data;
    info->cb();
    delete info->req;
    info->req = NULL;
    delete info;
    info = NULL;
}

void after_doing_work_unlock(uv_work_t *req, int status) {
    NodeUserInfo* info = (NodeUserInfo*)req->data;
    info->cb();
    pthread_mutex_unlock(&info->mutex);
    pthread_mutex_destroy(&info->mutex);
    delete info->req;
    info->req = NULL;
    delete info;
    info = NULL;
}

void timer_cb(uv_timer_s *handle) {
	// do nothing
}

void keepNodeRun() {
	// node没有事件挂在uv_loop上会退出进程，找点东西强行挂上去，免得进程退出了。
	// uv_pipe居然是用socket实现的，坑，先用timer占坑。
    uv_timer_t timer;  
    uv_timer_init(uv_default_loop(), &timer);
    uv_timer_start(&timer, timer_cb, 1, 1);
}

void postTaskToMainThread(const std::function<void()>& cb) {
	NodeUserInfo* info = new NodeUserInfo;
	info->req = new uv_work_t;
	info->cb = cb;
    info->req->data = info;
    uv_queue_work(uv_default_loop(), info->req, doing_work, after_doing_work);
}

void postSyncTaskToMainThread(const std::function<void()>& cb) {
	NodeUserInfo* info = new NodeUserInfo;
	pthread_mutex_init(&info->mutex, NULL);
	info->req = new uv_work_t;
    info->req->data = info;
    uv_queue_work(uv_default_loop(), info->req, doing_work, after_doing_work_unlock);
    pthread_mutex_lock(&info->mutex);
}

extern "C" {

void JSOnAcceptNewClient(void* delegate, void* ipc_ptr, int ipc_id) {
	if (delegateMap.find(delegate) != delegateMap.end()) {
		postTaskToMainThread([delegate, ipc_ptr, ipc_id]() {
			Isolate* isolate = Isolate::GetCurrent();
			HandleScope handle_scope(isolate);
			Local<Context> context = isolate->GetCurrentContext();
			Handle<Object> obj = Handle<Object>::New(isolate, delegateMap[delegate]);
			Local<String> key = String::NewFromUtf8(isolate, kOnAcceptNewClient);
			if (obj->Has(key)) {
				MaybeLocal<Value> maybeOnAcceptNewClientCallback = obj->Get(context, key);
				Local<Value> onAcceptNewClientCallback;
				if (maybeOnAcceptNewClientCallback.ToLocal(&onAcceptNewClientCallback) && 
					ipcMap.find(ipc_ptr) != ipcMap.end()) {
					int index = ipcMap[ipc_ptr];
					Local<Integer> ipcManager = Integer::New(isolate, reinterpret_cast<int32_t>(index));
					Local<Integer> ipcID = Integer::New(isolate, reinterpret_cast<int32_t>(ipc_id));
					if (onAcceptNewClientCallback->IsObject()) {
						Local<Object> onAcceptNewClientCallbackObj = onAcceptNewClientCallback.As<Object>();
						if (onAcceptNewClientCallbackObj->IsCallable()) {
							Local<Value> args[2] = { ipcManager, ipcID };
							Handle<Function> cb = Handle<Function>::Cast(onAcceptNewClientCallbackObj);
							cb->Call(context->Global(), 2, args);
						}
					}
				}
			}
		});
	}
}

void JSOnClientClose(void* delegate, void* ipc_ptr, int ipc_id) {
	if (delegateMap.find(delegate) != delegateMap.end()) {
		Isolate* isolate = Isolate::GetCurrent();
		HandleScope handle_scope(isolate);
		Local<Context> context = isolate->GetCurrentContext();
		Handle<Object> obj = Handle<Object>::New(isolate, delegateMap[delegate]);
		Local<String> key = String::NewFromUtf8(isolate, kOnClientClose);
		if (obj->Has(key)) {
			MaybeLocal<Value> maybeOnAcceptNewClientCallback = obj->Get(context, key);
			Local<Value> onAcceptNewClientCallback;
			if (maybeOnAcceptNewClientCallback.ToLocal(&onAcceptNewClientCallback) && 
				ipcMap.find(ipc_ptr) != ipcMap.end()) {
				int index = ipcMap[ipc_ptr];
				Local<Integer> ipcManager = Integer::New(isolate, reinterpret_cast<int32_t>(index));
				Local<Integer> ipcID = Integer::New(isolate, reinterpret_cast<int32_t>(ipc_id));
				if (onAcceptNewClientCallback->IsObject()) {
					Local<Object> onAcceptNewClientCallbackObj = onAcceptNewClientCallback.As<Object>();
					if (onAcceptNewClientCallbackObj->IsCallable()) {
						Local<Value> args[2] = { ipcManager, ipcID };
						Handle<Function> cb = Handle<Function>::Cast(onAcceptNewClientCallbackObj);
						cb->Call(context->Global(), 2, args);
					}
				}
			}
		}
	}
}

void JSOnConnect(void* delegate, void* ipc_ptr, int ipc_id) {
	if (delegateMap.find(delegate) != delegateMap.end()) {
		Isolate* isolate = Isolate::GetCurrent();
		HandleScope handle_scope(isolate);
		Local<Context> context = isolate->GetCurrentContext();
		Handle<Object> obj = Handle<Object>::New(isolate, delegateMap[delegate]);
		Local<String> key = String::NewFromUtf8(isolate, kOnConnect);
		if (obj->Has(key)) {
			MaybeLocal<Value> maybeOnAcceptNewClientCallback = obj->Get(context, key);
			Local<Value> onAcceptNewClientCallback;
			if (maybeOnAcceptNewClientCallback.ToLocal(&onAcceptNewClientCallback) && 
				ipcMap.find(ipc_ptr) != ipcMap.end()) {
				int index = ipcMap[ipc_ptr];
				Local<Integer> ipcManager = Integer::New(isolate, reinterpret_cast<int32_t>(index));
				Local<Integer> ipcID = Integer::New(isolate, reinterpret_cast<int32_t>(ipc_id));
				if (onAcceptNewClientCallback->IsObject()) {
					Local<Object> onAcceptNewClientCallbackObj = onAcceptNewClientCallback.As<Object>();
					if (onAcceptNewClientCallbackObj->IsCallable()) {
						Local<Value> args[2] = { ipcManager, ipcID };
						Handle<Function> cb = Handle<Function>::Cast(onAcceptNewClientCallbackObj);
						cb->Call(context->Global(), 2, args);
					}
				}
			}
		}
	}
}

void JSOnServerClose(void* delegate, void* ipc_ptr) {
	if (delegateMap.find(delegate) != delegateMap.end()) {
		Isolate* isolate = Isolate::GetCurrent();
		HandleScope handle_scope(isolate);
		Local<Context> context = isolate->GetCurrentContext();
		Handle<Object> obj = Handle<Object>::New(isolate, delegateMap[delegate]);
		Local<String> key = String::NewFromUtf8(isolate, kOnServerClose);
		if (obj->Has(key)) {
			MaybeLocal<Value> maybeOnAcceptNewClientCallback = obj->Get(context, key);
			Local<Value> onAcceptNewClientCallback;
			if (maybeOnAcceptNewClientCallback.ToLocal(&onAcceptNewClientCallback) && 
				ipcMap.find(ipc_ptr) != ipcMap.end()) {
				int index = ipcMap[ipc_ptr];
				Local<Integer> ipcManager = Integer::New(isolate, reinterpret_cast<int32_t>(index));
				if (onAcceptNewClientCallback->IsObject()) {
					Local<Object> onAcceptNewClientCallbackObj = onAcceptNewClientCallback.As<Object>();
					if (onAcceptNewClientCallbackObj->IsCallable()) {
						Local<Value> args[2] = { ipcManager };
						Handle<Function> cb = Handle<Function>::Cast(onAcceptNewClientCallbackObj);
						cb->Call(context->Global(), 1, args);
					}
				}
			}
		}
	}
}

char* JSMethodFunc(void* param, const char* request) {
	int* intParam = reinterpret_cast<int*>(param);
	int index = *intParam;
	// delete intParam;
	std::string response = "";
	std::string req = request;
	sem_t sem;
	sem_init(&sem, 0, 0);
	postTaskToMainThread([index, req, &response, &sem]() {
		if (methodMap.find(index) != methodMap.end()) {
			Isolate* isolate = Isolate::GetCurrent();
			HandleScope handle_scope(isolate);
			Local<Context> context = isolate->GetCurrentContext();
			Handle<Function> cb = Handle<Function>::New(isolate, methodMap[index]);

			Local<String> localRequest = String::NewFromUtf8(isolate, req.c_str());
			Local<Value> args[1] = { localRequest };
			Local<Value> retValue = cb->Call(context->Global(), 1, args);
			Local<String> retString = retValue->ToString();
			response = *(String::Utf8Value(retString));
		}
		sem_post(&sem);
	});
	sem_wait(&sem);
	sem_destroy(&sem);
	char* ret = new char[sizeof(char) * response.size() + 1];
	memset(ret, 0, sizeof(char) * response.size() + 1);
	strcpy(ret, response.c_str());
	return ret;
}

void JSMethodCallback(void* param, int code, const char* response) {
	int* intParam = reinterpret_cast<int*>(param);
	int index = *intParam;
	std::cout << "index: " << index << std::endl;
	// delete intParam;
	std::string rsp = response;
	postTaskToMainThread([code, index, rsp]() {
		if (methodMap.find(index) != methodMap.end()) {
			Isolate* isolate = Isolate::GetCurrent();
			HandleScope handle_scope(isolate);
			Local<Context> context = isolate->GetCurrentContext();
			Handle<Function> cb = Handle<Function>::New(isolate, methodMap[index]);

			Local<Integer> localCode = Integer::New(isolate, code);
			Local<String> localResponse = String::NewFromUtf8(isolate, rsp.c_str());
			Local<Value> args[2] = { localCode, localResponse };
			cb->Call(context->Global(), 2, args);
		}
	});
}

}

void newIPCManager(const FunctionCallbackInfo<Value>& args) {
	Isolate* isolate = args.GetIsolate();
	if (args.Length() != 0) {
		isolate->ThrowException(Exception::TypeError(
	        String::NewFromUtf8(isolate, "Wrong number of arguments")));
		return;
	}
	void* ipcManager = NewIPCManager();
	int32_t ipcIndex = currentIndex++;
	ipcMap[ipcManager] = ipcIndex;
	indexMap[ipcIndex] = ipcManager;
	int32_t pointer = reinterpret_cast<int32_t>(ipcIndex);
	Local<Integer> ret_pointer = Integer::New(isolate, pointer);
	args.GetReturnValue().Set(ret_pointer);
}



void startListener(const FunctionCallbackInfo<Value>& args) {
	Isolate* isolate = args.GetIsolate();
	if (args.Length() != 1) {
		isolate->ThrowException(Exception::TypeError(
	        String::NewFromUtf8(isolate, "Wrong number of arguments")));
		return;
	}
	if (!args[0]->IsNumber()) {
		isolate->ThrowException(Exception::TypeError(
		    String::NewFromUtf8(isolate, "Wrong arguments")));
		return;
	}
	int32_t index = args[0]->IntegerValue();
	if (indexMap.find(index) != indexMap.end()) {
		void* ipcManager = indexMap[index];
		// node do not support multi-thread, create a listen thread in C++ level.
		std::thread([ipcManager]() {
			StartListener(ipcManager);
		}).detach();
	}
}

void createServer(const FunctionCallbackInfo<Value>& args) {
	Isolate* isolate = args.GetIsolate();
	if (args.Length() != 3) {
		isolate->ThrowException(Exception::TypeError(
	        String::NewFromUtf8(isolate, "Wrong number of arguments")));
		return;
	}
	if (!args[0]->IsNumber()) {
		isolate->ThrowException(Exception::TypeError(
		    String::NewFromUtf8(isolate, "Wrong arguments")));
		return;
	}
	int32_t index = args[0]->IntegerValue();
	if (indexMap.find(index) != indexMap.end()) {
		void* ipcManager = indexMap[index];
		if (!args[1]->IsString()) {
			isolate->ThrowException(Exception::TypeError(
			    String::NewFromUtf8(isolate, "Wrong arguments")));
			return;
		}
		String::Utf8Value str(args[1]);
	    const char* ipcName = *str;

	    void* cDelegate = NewServerInterface();

	    Local<Object> delegate = Local<Object>::Cast(args[2]);
	    delegateMap[cDelegate].Reset(isolate, Persistent<Object>(isolate, delegate));

	    int server_id = CreateServer(ipcManager, ipcName, cDelegate);
		args.GetReturnValue().Set(Integer::New(isolate, server_id));
	}
}

void openClient(const FunctionCallbackInfo<Value>& args) {
	Isolate* isolate = args.GetIsolate();
	if (args.Length() != 3) {
		isolate->ThrowException(Exception::TypeError(
	        String::NewFromUtf8(isolate, "Wrong number of arguments")));
		return;
	}
	if (!args[0]->IsNumber()) {
		isolate->ThrowException(Exception::TypeError(
		    String::NewFromUtf8(isolate, "Wrong arguments")));
		return;
	}
	int32_t index = args[0]->IntegerValue();
	if (indexMap.find(index) != indexMap.end()) {
		void* ipcManager = indexMap[index];
		if (!args[1]->IsString()) {
			isolate->ThrowException(Exception::TypeError(
			    String::NewFromUtf8(isolate, "Wrong arguments")));
			return;
		}
		String::Utf8Value str(args[1]);
	    const char* ipcName = *str;

	    void* cDelegate = NewClientInterface();

	    Local<Object> delegate = Local<Object>::Cast(args[2]);
	    delegateMap[cDelegate].Reset(isolate, Persistent<Object>(isolate, delegate));

	    int client_id = OpenClient(ipcManager, ipcName, cDelegate);
	    args.GetReturnValue().Set(Integer::New(isolate, client_id));
	}
}

void registerMethod(const FunctionCallbackInfo<Value>& args) {
	Isolate* isolate = args.GetIsolate();
	if (args.Length() != 4) {
		isolate->ThrowException(Exception::TypeError(
	        String::NewFromUtf8(isolate, "Wrong number of arguments")));
		return;
	}
	if (!args[0]->IsNumber()) {
		isolate->ThrowException(Exception::TypeError(
		    String::NewFromUtf8(isolate, "Wrong arguments")));
		return;
	}
	int32_t index = args[0]->IntegerValue();
	if (indexMap.find(index) != indexMap.end()) {
		void* ipcManager = indexMap[index];
		if (!args[1]->IsNumber()) {
			isolate->ThrowException(Exception::TypeError(
		    String::NewFromUtf8(isolate, "Wrong arguments 1")));
			return;
		}
		int32_t ipcID = args[1]->IntegerValue();
		if (!args[2]->IsString()) {
			isolate->ThrowException(Exception::TypeError(
			    String::NewFromUtf8(isolate, "Wrong arguments 2")));
			return;
		}
		String::Utf8Value str(args[2]);
	    const char* methodName = *str;

	    if (!args[3]->IsFunction()) {
	    	isolate->ThrowException(Exception::TypeError(
			    String::NewFromUtf8(isolate, "Wrong arguments 3")));
			return;
	    }

	    Handle<Function> cb = Handle<Function>::Cast(args[3]);

	    int32_t delegateIndex = currentIndex++;
	    methodMap[delegateIndex].Reset(isolate, Persistent<Function>(isolate, cb));

	    int* intParam = new int;
	    *intParam = delegateIndex;

	    //void RegisterMethod(void* ipc_ptr, int ipc_id, const char* method_name, void* method, void* param) {
	    RegisterMethod(ipcManager, ipcID, methodName, reinterpret_cast<void*>(CMethodFunc), 
	    	reinterpret_cast<void*>(intParam));
	}
}

// CallMethod)(void* ipc_ptr, int ipc_id, const char* method_name, const char* request, MethodCallback callback, void* param);
void callMethod(const FunctionCallbackInfo<Value>& args) {
	Isolate* isolate = args.GetIsolate();
	if (args.Length() != 5) {
		isolate->ThrowException(Exception::TypeError(
	        String::NewFromUtf8(isolate, "Wrong number of arguments")));
		return;
	}
	if (!args[0]->IsNumber()) {
		isolate->ThrowException(Exception::TypeError(
		    String::NewFromUtf8(isolate, "Wrong arguments")));
		return;
	}
	int32_t index = args[0]->IntegerValue();
	if (indexMap.find(index) != indexMap.end()) {
		void* ipcManager = indexMap[index];
		if (!args[1]->IsNumber()) {
			isolate->ThrowException(Exception::TypeError(
		    String::NewFromUtf8(isolate, "Wrong arguments")));
			return;
		}
		int32_t ipcID = args[1]->IntegerValue();
		if (!args[2]->IsString()) {
			isolate->ThrowException(Exception::TypeError(
			    String::NewFromUtf8(isolate, "Wrong arguments")));
			return;
		}
		String::Utf8Value methodStr(args[2]);
	    const char* methodName = *methodStr;

	    if (!args[3]->IsString()) {
	    	isolate->ThrowException(Exception::TypeError(
			    String::NewFromUtf8(isolate, "Wrong arguments")));
			return;
	    }

	    String::Utf8Value requestStr(args[3]);
	    const char* cRequest = *requestStr;

	    if (!args[4]->IsFunction()) {
	    	isolate->ThrowException(Exception::TypeError(
			    String::NewFromUtf8(isolate, "Wrong arguments")));
			return;
	    }

	    Handle<Function> cb = Handle<Function>::Cast(args[4]);

	    int32_t delegateIndex = currentIndex++;
	    methodMap[delegateIndex].Reset(isolate, Persistent<Function>(isolate, cb));

	    int* intParam = new int;
	    *intParam = delegateIndex;

	    CallMethod(ipcManager, ipcID, methodName, cRequest, 
	    	reinterpret_cast<void*>(CMethodCallback), reinterpret_cast<void*>(intParam));
	}
}

void init(Local<Object> exports) {
	LoadLibrary();
	NODE_SET_METHOD(exports, "newIPCManager", newIPCManager);
	NODE_SET_METHOD(exports, "startListener", startListener);
	NODE_SET_METHOD(exports, "createServer", createServer);
	NODE_SET_METHOD(exports, "openClient", openClient);
	NODE_SET_METHOD(exports, "registerMethod", registerMethod);
	NODE_SET_METHOD(exports, "callMethod", callMethod);
	// keepNodeRun();
}

NODE_MODULE(addon, init)

} // namespace ipc;
