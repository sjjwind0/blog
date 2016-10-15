#include "include/ipc_mgr.h"

#include <iostream>
#include <unistd.h>
#include <sys/epoll.h>
#include "third_party/json/json.h"

#define FDSIZE      10
#define EPOLLEVENTS 10

using namespace json11;

int IPCManager::call_index = 0;

IPCManager::IPCManager() {
    epfd_ = epoll_create(FDSIZE);
}

IPCManager::~IPCManager() {
}

void* IPCManager::ThreadFunc(void* args) {
    IPCManager* self = reinterpret_cast<IPCManager*>(args);
    if (self != nullptr) {
        epoll_event events[EPOLLEVENTS];
        while (true) {
            int ev_size = epoll_wait(self->epfd_, events, EPOLLEVENTS,-1);
            for (int i = 0; i < ev_size; i++) {
                int ev_read_fd = events[i].data.fd;
                int ev_show_id = self->GetShowIDByReadFD(ev_read_fd);
                std::string recv_data = "";
                int ret = self->ipc_info_map_[ev_show_id].fifo->Read(recv_data);
                if (ret == 0) {
                    self->HandleMessage(ev_show_id, recv_data);
                }
            }
        }
    }
}

void IPCManager::StartListener() {
    // std::cout << "StartListener" << std::endl;
    pthread_t id;
    int ret = pthread_create(&id, NULL, IPCManager::ThreadFunc, this);
    if (ret) {
        std::cout << "StartListener error: " << errno << std::endl;
    }
}

int IPCManager::CreateServer(const std::string& ipc_name, std::shared_ptr<IPCServerDelegate> delegate) {
    std::shared_ptr<TwoWayFifo> fifo(new TwoWayFifo(ipc_name));
    fifo->CreateServerFile();
    ipc_info_map_[fifo->GetID()].fifo = fifo;
    ipc_info_map_[fifo->GetID()].server_delegate = delegate;
    read_fd_map_[fifo->GetReadFD()] = fifo->GetID();
    struct epoll_event ev;
    ev.events = EPOLLIN;
    ev.data.fd = fifo->GetReadFD();
    epoll_ctl(epfd_, EPOLL_CTL_ADD, fifo->GetReadFD(), &ev);
    return fifo->GetID();
}

int IPCManager::OpenClient(const std::string& ipc_name, std::shared_ptr<IPCClientDelegate> delegate) {
    std::shared_ptr<TwoWayFifo> fifo(new TwoWayFifo(ipc_name));
    fifo->CreateClientFile();
    fifo->OpenServerFile();
    ipc_info_map_[fifo->GetID()].fifo = fifo;
    ipc_info_map_[fifo->GetID()].client_delegate = delegate;
    read_fd_map_[fifo->GetReadFD()] = fifo->GetID();
    struct epoll_event ev;
    ev.events = EPOLLIN;
    ev.data.fd = fifo->GetReadFD();
    epoll_ctl(epfd_, EPOLL_CTL_ADD, fifo->GetReadFD(), &ev);
    if (fifo->Write("{\"action\": \"init\"}") == 0) {
        delegate->OnConnect(this, fifo->GetID());
    }
    return fifo->GetID();
}

void IPCManager::RegisterMethod(int ipc_id, const std::string& method_name, const Method& method) {
    ipc_info_map_[ipc_id].method_map[method_name] = method;
}

void IPCManager::CallMethod(int ipc_id, const std::string& method_name, 
        const std::string& request, const MethodCallback& callback) {
    int current_call = call_index++;
    ipc_info_map_[ipc_id].method_callback_map[method_name][current_call] = callback;
    Json obj = Json::object({
        { "action", "request" },
        { "id", current_call }, 
        { "method", method_name },
        { "request", request },
    });
    ipc_info_map_[ipc_id].fifo->Write(obj.dump());
}

std::string IPCManager::GetNameByIPCID(int ipc_id) {
    if (ipc_info_map_.find(ipc_id) != ipc_info_map_.end()) {
        return ipc_info_map_[ipc_id].fifo->GetName();
    }
    return "";
}

void IPCManager::HandleMessage(int show_id, const std::string& data) {
    ParseData(data);
    while (true) {
        std::string next_message = "";
        if (GetNextMessage(next_message)) {
            std::string err;
            auto json = Json::parse(next_message, err);
            if (!err.empty()) {
                std::cout << "error: " << err << std::endl;
                continue;
            }
            // std::cout << "message: " << next_message << std::endl;
            std::string action = json["action"].string_value();
            // std::cout << "action: " << action << std::endl;
            if (action == "request") {
                // 请求
                int req_id = json["id"].int_value();
                std::string method = json["method"].string_value();
                std::string req = json["request"].string_value();
                if (ipc_info_map_[show_id].method_map.find(method) != 
                    ipc_info_map_[show_id].method_map.end()) {
                    std::string response = "";
                    (ipc_info_map_[show_id].method_map[method])(req, response);
                    Json obj = Json::object({
                        { "action", "response" },
                        { "id", req_id }, 
                        { "method", method },
                        { "response", response },
                        { "code", ErrorOK },
                    });
                    ipc_info_map_[show_id].fifo->Write(obj.dump());
                } else {
                    Json obj = Json::object({
                        { "action", "response" },
                        { "id", req_id }, 
                        { "method", method },
                        { "error", "no such api." },
                        { "code", ErrorNoSuchAPI },
                    });
                    ipc_info_map_[show_id].fifo->Write(obj.dump());
                }
            } else if (action == "init") {
                // 收到客户端发过来的初始化完成的消息
                if (ipc_info_map_[show_id].fifo->OpenClientFile()) {
                    if (ipc_info_map_[show_id].server_delegate != nullptr) {
                        ipc_info_map_[show_id].server_delegate->OnAcceptNewClient(this, show_id);
                    }
                } else {
                    std::cout << "open client failed: " << errno;
                }

            } else  if (action == "response") {
                // 请求的回应
                std::string method = json["method"].string_value();
                if (ipc_info_map_[show_id].method_callback_map.find(method) != ipc_info_map_[show_id].method_callback_map.end()) {
                    int req_id = json["id"].int_value();
                    if (ipc_info_map_[show_id].method_callback_map[method].find(req_id) != ipc_info_map_[show_id].method_callback_map[method].end()) {
                        ErrorCode code = ErrorCode((json["code"].int_value()));
                        if (code == ErrorOK) {
                            std::string req = json["response"].string_value();
                            (ipc_info_map_[show_id].method_callback_map[method][req_id])(code, req);
                        } else {
                            std::string error = json["error"].string_value();
                            (ipc_info_map_[show_id].method_callback_map[method][req_id])(code, error);
                        }
                    }
                }
            } else if (action == "close") {
                // 关闭
            }
        } else {
            break;
        }
    }
}

void IPCManager::ParseData(const std::string& data) {
    // 后续可以根据情况改为循环队列+内存池取数据
    buffer_ += data;
}

bool IPCManager::GetNextMessage(std::string& data) {
    long unsigned int data_size = buffer_.size();
    if (data_size == 0) {
        return false;
    }
    const char* buf_ptr = reinterpret_cast<const char*>(buffer_.c_str());
    long unsigned int next_message_size = *(reinterpret_cast<int*>(const_cast<char*>(buf_ptr)));
    if (next_message_size > data_size) {
        std::cout << "not enough size" << std::endl;
        return false;
    } else {
        data.assign(buf_ptr + 4, next_message_size);
        if (sizeof(int) + next_message_size == buffer_.size()) {
            buffer_.clear();
        } else {
            buffer_ = buffer_.substr(4 + next_message_size);
        }
        return true;
    }
}

int IPCManager::GetShowIDByReadFD(int read_fd) {
    return read_fd_map_[read_fd];
}