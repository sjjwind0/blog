#ifndef _FIFO_H_
#define _FIFO_H_

#include <string>

class TwoWayFifo {
public:
    TwoWayFifo(const std::string& name);
    ~TwoWayFifo();

    bool CreateServerFile();
    bool CreateClientFile();
    bool OpenServerFile();
    bool OpenClientFile();

    int Write(const std::string& data);

    int Read(std::string& data);

    void Close();

    int GetID() const { return id_; }

    int GetReadFD() const { return read_fd_; }

    int GetWriteFD() const { return write_fd_; }

    std::string GetName() const { return name_; }

private:
    int id_;
    int read_fd_;
    int write_fd_;
    std::string name_;
    static int index;
};

#endif