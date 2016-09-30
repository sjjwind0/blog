#ifndef _IPC_MGR_H_
#define _IPC_MGR_H_

#include <map>
#include <string>
#include <memory>
#include <functional>
#include "fifo.h"
#include "delegate.h"
#include "include/fifo.h"

class IPCManager {
public:
    typedef std::function<void(const std::string& request, std::string& response)> Method;
    typedef std::function<void(const std::string& response)> MethodCallback;
    IPCManager();
    ~IPCManager();
public:
    void StartListener();

    int CreateServer(const std::string& ipc_name);
    int OpenClient(const std::string& ipc_name);

    void RegisterMethod(int ipc_id, const std::string& method_name, const Method& method);
    void CallMethod(int ipc_id, const std::string& method_name, const std::string& request, const MethodCallback& callback);
private:
    int CreateIPCChannel(const std::string& channel_name, IPCDelegate* delegate);
    int OpenIPCChannel(const std::string& channel_name, IPCDelegate* delegate);

    void HandleMessage(const std::string& data);
    void ParseData(const std::string& data);
    bool GetNextData(std::string& data);

    int GetShowIDByReadFD(int read_fd);

    struct IPCInfo {
        int show_fd;
        std::shared_ptr<TwoWayFifo> fifo;
        std::shared_ptr<IPCDelegate> delegate;
        std::map<std::string, const Method&> method_map_;
        std::map<std::string, const MethodCallback&> method_callback_map_;
    };
private:
    std::map<int, IPCInfo> ipc_info_map_;
    std::map<int, int> read_fd_map_;
    // 做解包用
    std::string buffer_;
    int epfd_;
};

#endif
