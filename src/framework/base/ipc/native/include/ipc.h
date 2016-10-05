#ifndef _IPC_H_
#define _IPC_H_

#include "error.h"

#ifdef __cplusplus
extern "C" {
#endif
// method & method callback
typedef void (*Method)(void* func, const char* request, char** response);
typedef void (*MethodCallback)(void* func, ErrorCode code, const char* response);

// server delegate
typedef void (*OnAcceptNewClient)(void* delegate, void* ipc_ptr, int ipc_id);
typedef void (*OnClientClose)(void* delegate, void* ipc_ptr, int ipc_id);

// client delegate
typedef void (*OnConnect)(void* delegate, void* ipc_ptr, int ipc_id);
typedef void (*OnServerClose)(void* delegate, void* ipc_ptr);

struct IPCServerInterface {
	OnAcceptNewClient on_accept_new_client_ptr;
	OnClientClose on_client_close_ptr;
};

struct IPCClientInterface {
	OnConnect on_connect_ptr;
	OnServerClose on_server_close_ptr;
};

void* NewIPCManager();
void StartListener(void* ipc_ptr);
int CreateServer(void* ipc_ptr, const char* ipc_name, IPCServerInterface* delegate);
int OpenClient(void* ipc_ptr, const char* ipc_name, IPCClientInterface* delegate);
void RegisterMethod(void* ipc_ptr, int ipc_id, const char* method_name, Method method, void* param);
void CallMethod(void* ipc_ptr, int ipc_id, const char* method_name, const char* request, 
	MethodCallback callback, void* param);
char* GetNameByIPCID(void* ipc_ptr, int ipc_id);

#ifdef __cplusplus
}
#endif

#endif
