#ifndef _FIFO_H_
#define _FIFO_H_

#include <string>

class TwoWayFifo {
public:
    TwoWayFifo(const std::string& name);
    ~TwoWayFifo();

    bool CreatServerFile();
    bool CreatClientFile();
    bool OpenServerFile();
    bool OpenClientFile();

    int Write(const std::string& data);

    int Read(std::string& data);

    void Close();

    int GetFD() { return show_fd_; }

    int GetReadFD() { return read_fd_; }

    int GetWriteFD() { return write_fd_; }

private:
    int show_fd_;
    int read_fd_;
    int write_fd_;
    std::string name_;
    static int index;
};

#endif