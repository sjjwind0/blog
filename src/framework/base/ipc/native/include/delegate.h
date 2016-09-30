#ifndef _DELEGATE_H_
#define _DELEGATE_H_

class IPCDelegate {
public:
    virtual void OnReceiveData(int fd, int code, const std::string& data);
};

#endif