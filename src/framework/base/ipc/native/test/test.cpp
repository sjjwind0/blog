#include "include/ipc_mgr.h"

#include <unistd.h>
#include <iostream>
#include <time.h>

class TestDelegate : public IPCServerDelegate, public IPCClientDelegate {
public:
	void OnAcceptNewClient(IPCManager* manager, int ipc_id) {
		std::cout << "accept ipc name: " << manager->GetNameByIPCID(ipc_id) << std::endl;
		manager->CallMethod(ipc_id, "Method3", "aaa", [manager, ipc_id](ErrorCode code, const std::string& rsp) {
			if (code == ErrorOK) {
	    		std::cout << "server pid: " << getppid() << ", name: " << manager->GetNameByIPCID(ipc_id) << ", CallMethod3 rsp: " << rsp  << std::endl;
	    	} else {
	    		std::cout << "server pid: " << getppid() << ", name: " << manager->GetNameByIPCID(ipc_id) << ", CallMethod3 error: " << rsp  << std::endl;
	    	}
    	});
    	manager->CallMethod(ipc_id, "Method4", "aaa", [manager, ipc_id](ErrorCode code, const std::string& rsp) {
    		if (code == ErrorOK) {
	    		std::cout << "server pid: " << getppid() << ", name: " << manager->GetNameByIPCID(ipc_id) << ", CallMethod4 rsp: " << rsp  << std::endl;
	    	} else {
	    		std::cout << "server pid: " << getppid() << ", name: " << manager->GetNameByIPCID(ipc_id) << ", CallMethod4 error: " << rsp  << std::endl;
	    	}
    	});
	}
	void OnClientClose(IPCManager* manager, int ipc_id) {
	}
	void OnConnect(IPCManager* manager, int ipc_id) {
		std::cout << "connect ipc name: " << manager->GetNameByIPCID(ipc_id) << std::endl;
		manager->CallMethod(ipc_id, "Method1", "HelloWorld", [manager, ipc_id](ErrorCode code, const std::string& rsp) {
			if (code == ErrorOK) {
				std::cout << "client pid: " << getppid() << ", name: " << manager->GetNameByIPCID(ipc_id) << ", CallMethod1 rsp: " << rsp  << std::endl;
	    	} else {
	    		std::cout << "client pid: " << getppid() << ", name: " << manager->GetNameByIPCID(ipc_id) << ", CallMethod1 error: " << rsp  << std::endl;
	    	}
        });
        manager->CallMethod(ipc_id, "Method2", "HelloWorld", [manager, ipc_id](ErrorCode code, const std::string& rsp) {
        	if (code == ErrorOK) {
	    		std::cout << "client pid: " << getppid() << ", name: " << manager->GetNameByIPCID(ipc_id) << ", CallMethod2 rsp: " << rsp  << std::endl;
	    	} else {
	    		std::cout << "client pid: " << getppid() << ", name: " << manager->GetNameByIPCID(ipc_id) << ", CallMethod2 error: " << rsp  << std::endl;
	    	}
        });
	}
	void OnServerClose(IPCManager* manager) {
	}
};

int main(int argc, const char* argv[]) {
    pid_t fpid = fork();
    if (fpid == 0) {
    	fpid = fork();
    	if (fpid != 0) {
	    	// waiting for server start.
	    	usleep(40 * 1000);
	    	std::cout << "client connect ..." << std::endl;
	        IPCManager* ipc = new IPCManager();
	        std::shared_ptr<IPCClientDelegate> delegate(static_cast<IPCClientDelegate*>(new TestDelegate()));
	        int ipc_id = ipc->OpenClient("test", delegate);
	        ipc->RegisterMethod(ipc_id, "Method3", [](const std::string& req, std::string& rsp) {
	        	rsp = "天之道，损有余而补不足，是故需胜实。";
	        });
	        ipc->StartListener();
	    } else {
	    	// waiting for server start.
	    	usleep(40 * 1000);
	    	std::cout << "client connect ..." << std::endl;
	        IPCManager* ipc = new IPCManager();
	        std::shared_ptr<IPCClientDelegate> delegate(static_cast<IPCClientDelegate*>(new TestDelegate()));
	        int ipc_id = ipc->OpenClient("test1", delegate);
	        ipc->RegisterMethod(ipc_id, "Method4", [](const std::string& req, std::string& rsp) {
	        	rsp = "冰，水为之而寒于水，青，出于蓝而胜于蓝。";
	        });
	        ipc->StartListener();
	    }
    } else {
    	std::cout << "server start ..." << std::endl;
    	IPCManager* ipc = new IPCManager();
    	std::shared_ptr<IPCServerDelegate> delegate(static_cast<IPCServerDelegate*>(new TestDelegate()));
    	int ipc_id_1 = ipc->CreateServer("test", delegate);
    	ipc->RegisterMethod(ipc_id_1, "Method1", [](const std::string& req, std::string& rsp) {
    		rsp = req;
    	});
    	int ipc_id_2 = ipc->CreateServer("test1", delegate);
    	ipc->RegisterMethod(ipc_id_2, "Method2", [](const std::string& req, std::string& rsp) {
    		rsp = "我是一个粉刷匠，粉刷本领强";
    	});
    	ipc->StartListener();
    }
    return 0;
}
