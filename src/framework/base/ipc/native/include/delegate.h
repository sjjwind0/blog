#ifndef _DELEGATE_H_
#define _DELEGATE_H_

class IPCManager;
class IPCServerDelegate {
public:
	virtual void OnAcceptNewClient(IPCManager* manager, int ipc_id) = 0;
	virtual void OnClientClose() = 0;
};

class IPCClientDelegate {
public:
	virtual void OnConnect(IPCManager* manager, int ipc_id) = 0;
	virtual void OnServerClose() = 0;
};

#endif