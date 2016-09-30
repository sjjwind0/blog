#include "include/ipc_mgr.h"

#include <sys/epoll.h>
#include "third_party/json/json.h"

#define FDSIZE      10
#define EPOLLEVENTS 10

IPCManager::IPCManager() {
}

IPCManager::~IPCManager() {
}

void IPCManager::StartListener() {
    epfd_ = epoll_create(FDSIZE);
    epoll_event events[EPOLLEVENTS];
    while (true) {
        int ev_size = epoll_wait(epollfd, events, EPOLLEVENTS,-1);
        for (int i = 0; i < ev_size; i++) {
            int ev_read_fd = events[i].data.fd;
            int ev_show_id = GetShowIDByReadFD(ev_read_fd);
            std::string recv_data = "";
            int ret = ipc_info_map_[ev_show_id].pipe->Read(recv_data);
            if (ret == 0) {
                HandleMessage(recv_data);
            }
        }
    }
}

int IPCManager::CreateServer(const std::string& ipc_name) {
    shared_ptr<TwoWayFifo> fifo = new TwoWayFifo(ipc_name);
    fifo->CreateServerFile();
    ipc_info_map_[fifo->GetShowID()].fifo.reset(fifo);
    return fifo_->GetShowID();
}

int IPCManager::OpenClient(int show_id) {
    shared_ptr<TwoWayFifo> fifo = new TwoWayFifo(ipc_name);
    fifo->CreateClientFile();
    fifo->OpenServerFile();
    fifo->Write("{\"param\": \"init\"}");
    ipc_info_map_[fifo->GetShowID()].fifo.reset(fifo);
}

void IPCManager::RegisterMethod(int ipc_id, const std::string& method_name, const Method& method) {
    ipc_info_map_[ipc_id].method_map_[method_name] = method;
}

void IPCManager::CallMethod(int ipc_id, const std::string& method_name, const std::string& request, const MethodCallback& callback) {
}

int IPCManager::CreateIPCChannel(const std::string& channel_name, IPCDelegate* delegate) {
}

int IPCManager::OpenIPCChannel(const std::string& channel_name, IPCDelegate* delegate) {
}

void IPCManager::HandleMessage(const std::string& data) {
    ParseData(data);
    while (true) {
        std::string next_message = "";
        if (GetNextMessage(next_message)) {
            string err;
            auto json = Json::parse(next_message, err);
            std::string action = json["action"].string_value();
            if (action == "request") {
                // 请求
            } else if (action == "init") {
                // 收到客户端发过来的初始化完成的消息
                int show_id = json["param"]["id"];
                std::cout << "init: " << show_id;
                ipc_info_map_[show_id].fifo->OpenClientFile();
            } else  if (action == "response") {
                // 请求的回应
                int show_id = json["param"]["id"];
            } else if (action == "close") {
                // 关闭
            }
        }
    }
}

void IPCManager::ParseData(const std::string& data) {
    // 后续可以根据情况改为循环队列+内存池取数据
    buffer_ += data;
}

bool IPCManager::GetNextMessage(std::string& data) {
    unsigned int data_size = buffer_.size();
    unsigned char* buf_ptr = reinterpret_cast<unsigned char*>(buffer_.c_str());
    int next_message_size = int(reinterpret_cast<int*>(buf_ptr));
    std::cout << "next message size: " << next_message_size;
    if (next_message_size > data_size) {
        std::cout << "not enough size";
    } else {
        data.assign(buf_ptr + 4, next_message_size);
        buffer_ = buffer_.substr(4);
    }
}

int IPCManager::GetShowIDByReadFD(int read_fd) {
    return read_fd_map_[read_fd];
}