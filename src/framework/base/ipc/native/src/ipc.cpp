#include "include/ipc.h"

#include <stdio.h>
#include <memory>
#include <cstring>
#include <iostream>
#include "include/ipc_mgr.h"

class IPCServerProxyDelegate : public IPCServerDelegate {
public:
	IPCServerProxyDelegate(IPCServerInterface* delegate) 
		: delegate_(delegate) {
	}
	virtual void OnAcceptNewClient(IPCManager* manager, int ipc_id) override {
		if (delegate_ != nullptr) {
			void* ptr = reinterpret_cast<void*>(manager);
			delegate_->on_accept_new_client_ptr(delegate_, ptr, ipc_id);
		}
	}
	virtual void OnClientClose(IPCManager* manager, int ipc_id) override {
		if (delegate_ != nullptr) {
			void* ptr = reinterpret_cast<void*>(manager);
			delegate_->on_client_close_ptr(delegate_, ptr, ipc_id);
		}
	}
private:
	IPCServerInterface* delegate_;
};

class IPCClientProxyDelegate : public IPCClientDelegate {
public:
	IPCClientProxyDelegate(IPCClientInterface* delegate) 
		: delegate_(delegate) {
	}
	virtual void OnConnect(IPCManager* manager, int ipc_id) override {
		if (delegate_ != nullptr) {
			void* ptr = reinterpret_cast<void*>(manager);
			delegate_->on_connect_ptr(delegate_, ptr, ipc_id);
		}
	}
	virtual void OnServerClose(IPCManager* manager) override {
		if (delegate_ != nullptr) {
			void* ptr = reinterpret_cast<void*>(manager);
			delegate_->on_server_close_ptr(delegate_, ptr);
		}
	}
private:
	IPCClientInterface* delegate_;
};

void* NewIPCManager() {
	return reinterpret_cast<void*>(new IPCManager());
}

void StartListener(void* ipc_ptr) {
	IPCManager* manager = reinterpret_cast<IPCManager*>(ipc_ptr);
	manager->StartListener();
}

int CreateServer(void* ipc_ptr, const char* ipc_name, IPCServerInterface* delegate) {
	IPCManager* manager = reinterpret_cast<IPCManager*>(ipc_ptr);
	IPCServerDelegate* proxy_delegate = new IPCServerProxyDelegate(delegate);
	std::shared_ptr<IPCServerDelegate> ptr(proxy_delegate);
	return manager->CreateServer(ipc_name, ptr);
}

int OpenClient(void* ipc_ptr, const char* ipc_name, IPCClientInterface* delegate) {
	IPCManager* manager = reinterpret_cast<IPCManager*>(ipc_ptr);
	IPCClientDelegate* proxy_delegate = new IPCClientProxyDelegate(delegate);
	std::shared_ptr<IPCClientDelegate> ptr(proxy_delegate);
	return manager->OpenClient(ipc_name, ptr);
}

void RegisterMethod(void* ipc_ptr, int ipc_id, const char* method_name, Method method, void* param) {
	IPCManager* manager = reinterpret_cast<IPCManager*>(ipc_ptr);
	auto cb = [method, param](const std::string& request, std::string& response) {
		if (method != nullptr) {
			char* rsp = nullptr;
			method(param, request.c_str(), &rsp);
			response = rsp;
			delete[] rsp;
		}
	};
	manager->RegisterMethod(ipc_id, method_name, cb);
}

void CallMethod(void* ipc_ptr, int ipc_id, const char* method_name, 
		const char* request, MethodCallback callback, void* param) {
	IPCManager* manager = reinterpret_cast<IPCManager*>(ipc_ptr);
	auto cb = [callback, param](ErrorCode code, const std::string& response) {
		if (callback != nullptr) {
			callback(param, code, response.c_str());
		}
	};
	manager->CallMethod(ipc_id, method_name, request, cb);
}

char* GetNameByIPCID(void* ipc_ptr, int ipc_id) {
	IPCManager* manager = reinterpret_cast<IPCManager*>(ipc_ptr);
	std::string name = manager->GetNameByIPCID(ipc_id);
	char* ret_name = new char[name.size()];
	strcmp(ret_name, name.c_str());
	return ret_name;
}

bool StopClient(void* ipc_ptr, int ipc_id) {
	IPCManager* manager = reinterpret_cast<IPCManager*>(ipc_ptr);
	return manager->StopClient(ipc_id);
}
