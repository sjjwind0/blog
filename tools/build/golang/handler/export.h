#ifndef _GOLANG_IPC_EXPORT_H_
#define _GOLANG_IPC_EXPORT_H_

#include <stdio.h>
#include <stdlib.h>
#include <dlfcn.h>

#ifndef NULL
#define NULL 0
#endif

// method & method callback
typedef void (*Method)(void* param, const char* request, char** response);
typedef void (*MethodCallback)(void* param, int code, const char* response);

// server delegate
typedef void (*OnAcceptNewClient)(void* delegate, void* ipc_ptr, int ipc_id);
typedef void (*OnClientClose)(void* delegate, void* ipc_ptr, int ipc_id);

// client delegate
typedef void (*OnConnect)(void* delegate, void* ipc_ptr, int ipc_id);
typedef void (*OnServerClose)(void* delegate, void* ipc_ptr);

extern void GoOnAcceptNewClient(void* delegate, void* ipc_ptr, int ipc_id);
extern void GoOnClientClose(void* delegate, void* ipc_ptr, int ipc_id);
extern void GoOnConnect(void* delegate, void* ipc_ptr, int ipc_id);
extern void GoOnServerClose(void* delegate, void* ipc_ptr);
extern char* GoMethodFunc(int param, const char* request);
extern void GoMethodCallback(int param, int code, const char* response);

struct IPCServerInterface {
	OnAcceptNewClient on_accept_new_client_ptr;
	OnClientClose on_client_close_ptr;
};

struct IPCClientInterface {
	OnConnect on_connect_ptr;
	OnServerClose on_server_close_ptr;
};

void COnAcceptNewClient(void* delegate, void* ipc_ptr, int ipc_id) {
	GoOnAcceptNewClient(delegate, ipc_ptr, ipc_id);
}

void COnClientClose(void* delegate, void* ipc_ptr, int ipc_id) {
	GoOnClientClose(delegate, ipc_ptr, ipc_id);
}

void COnConnect(void* delegate, void* ipc_ptr, int ipc_id) {
	GoOnConnect(delegate, ipc_ptr, ipc_id);
}

void COnServerClose(void* delegate, void* ipc_ptr) {
	GoOnServerClose(delegate, ipc_ptr);
}

void CMethodFunc(void* param_ptr, const char* request, char** response) {
	int param = *((int*)param_ptr);
	// free(param_ptr);
	*response = GoMethodFunc(param, request);
}

void CMethodCallback(void* param_ptr, int code, const char* response) {
	int param = *((int*)param_ptr);
	// free(param_ptr);
	GoMethodCallback(param, code, response);
}

typedef void* (*Func_NewIPCManager)();
typedef void (*Func_StartListener)(void* ipc_ptr);
typedef int (*Func_CreateServer)(void* ipc_ptr, const char* ipc_name, struct IPCServerInterface* delegate);
typedef int (*Func_OpenClient)(void* ipc_ptr, const char* ipc_name, struct IPCClientInterface* delegate);
typedef void (*Func_RegisterMethod)(void* ipc_ptr, int ipc_id, const char* method_name, Method method, void* param);
typedef void (*Func_CallMethod)(void* ipc_ptr, int ipc_id, const char* method_name,
	const char* request, MethodCallback callback, void* param);
typedef char* (*Func_GetNameByIPCID)(void* ipc_ptr, int ipc_id);

Func_NewIPCManager _c_newIPCManager_ptr = NULL;
Func_StartListener _c_startListener_ptr = NULL;
Func_CreateServer _c_createServer_ptr = NULL;
Func_OpenClient _c_openClient_ptr = NULL;
Func_RegisterMethod _c_registerMethod_ptr = NULL;
Func_CallMethod _c_callMethod_ptr = NULL;
Func_GetNameByIPCID _c_getNameByIPCID_ptr = NULL;

const char* lib_path = "/home/wind/Project/blog/src/framework/base/ipc/native/libipc.so";
void* handle = NULL;
void LoadLibrary() {
	handle = dlopen(lib_path, 1);
	if (handle != NULL) {
		_c_newIPCManager_ptr = (Func_NewIPCManager)dlsym(handle, "NewIPCManager");
		_c_startListener_ptr = (Func_StartListener)dlsym(handle, "StartListener");
		_c_createServer_ptr = (Func_CreateServer)dlsym(handle, "CreateServer");
		_c_openClient_ptr = (Func_OpenClient)dlsym(handle, "OpenClient");
		_c_registerMethod_ptr = (Func_RegisterMethod)dlsym(handle, "RegisterMethod");
		_c_callMethod_ptr = (Func_CallMethod)dlsym(handle, "CallMethod");
		_c_getNameByIPCID_ptr = (Func_GetNameByIPCID)dlsym(handle, "GetNameByIPCID");
		printf("load so success!\n");
	} else {
		printf("load so failed!\n");
	}
}

void Close() {
	if (handle != NULL) {
		dlclose(handle);
	}
}

void* NewServerInterface() {
	struct IPCServerInterface* instance = (struct IPCServerInterface*)malloc(sizeof(struct IPCServerInterface));
	instance->on_accept_new_client_ptr = COnAcceptNewClient;
	instance->on_client_close_ptr = COnClientClose;
	return (void*)instance;
}

void* NewClientInterface() {
	struct IPCClientInterface* instance = (struct IPCClientInterface*)malloc(sizeof(struct IPCClientInterface));
	instance->on_connect_ptr = COnConnect;
	instance->on_server_close_ptr = COnServerClose;
	return (void*)instance;
}

void* NewIPCManager() {
	if (_c_newIPCManager_ptr != NULL) {
		return _c_newIPCManager_ptr();
	}
	return NULL;
}

void StartListener(void* ipc_ptr) {
	if (_c_startListener_ptr != NULL) {
		return _c_startListener_ptr(ipc_ptr);
	}
}

int CreateServer(void* ipc_ptr, const char* ipc_name, void* delegate) {
	if (_c_createServer_ptr != NULL) {
		return _c_createServer_ptr(ipc_ptr, ipc_name, delegate);
	}
	return -1;
}

int OpenClient(void* ipc_ptr, const char* ipc_name, void* delegate) {
	if (_c_openClient_ptr != NULL) {
		return _c_openClient_ptr(ipc_ptr, ipc_name, delegate);
	}
	return -1;
}

void RegisterMethod(void* ipc_ptr, int ipc_id, const char* method_name, void* method, int param) {
	if (_c_registerMethod_ptr != NULL) {
		int* param_ptr = (int*)malloc(sizeof(int));
		*param_ptr = param;
		_c_registerMethod_ptr(ipc_ptr, ipc_id, method_name, (Method)method, (void*)param_ptr);
	}
}

void CallMethod(void* ipc_ptr, int ipc_id, const char* method_name, const char* request,
		void* callback, int param) {
	if (_c_callMethod_ptr != NULL) {
		int* param_ptr = (int*)malloc(sizeof(int));
		*param_ptr = param;
		_c_callMethod_ptr(ipc_ptr, ipc_id, method_name, request, (MethodCallback)callback, 
			(void*)param_ptr);
	}
}

char* GetNameByIPCID(void* ipc_ptr, int ipc_id) {
	if (_c_getNameByIPCID_ptr != NULL) {
		return _c_getNameByIPCID_ptr(ipc_ptr, ipc_id);
	}
	return NULL;
}

#endif
