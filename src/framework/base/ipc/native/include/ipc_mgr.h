#ifndef _IPC_MGR_H_
#define _IPC_MGR_H_

#include <map>
#include <string>
#include <memory>
#include <functional>
#include "fifo.h"
#include "error.h"
#include "delegate.h"
#include "include/fifo.h"

class IPCManager {
public:
    typedef std::function<void(const std::string& request, std::string& response)> Method;
    typedef std::function<void(ErrorCode code, const std::string& response)> MethodCallback;
    IPCManager();
    ~IPCManager();
public:
    void StartListener();

    int CreateServer(const std::string& ipc_name, std::shared_ptr<IPCServerDelegate> delegate);
    int OpenClient(const std::string& ipc_name, std::shared_ptr<IPCClientDelegate> delegate);

    void RegisterMethod(int ipc_id, const std::string& method_name, const Method& method);
    void CallMethod(int ipc_id, const std::string& method_name, const std::string& request, const MethodCallback& callback);

    std::string GetNameByIPCID(int ipc_id);
private:
    void HandleMessage(int show_id, const std::string& data);
    void ParseData(const std::string& data);
    bool GetNextMessage(std::string& data);

    int GetShowIDByReadFD(int read_fd);

    struct IPCInfo {
        int show_fd;
        std::shared_ptr<TwoWayFifo> fifo;
        std::map<std::string, Method> method_map;
        std::map<std::string, std::map<int, MethodCallback> > method_callback_map;
        std::shared_ptr<IPCServerDelegate> server_delegate;
        std::shared_ptr<IPCClientDelegate> client_delegate;
    };
private:
    std::map<int, IPCInfo> ipc_info_map_;
    std::map<int, int> read_fd_map_;
    // 做解包用
    std::string buffer_;
    int epfd_;
    static int call_index;
};

#endif
