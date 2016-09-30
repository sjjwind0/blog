#include "include/fifo.h"

#include <errno.h>
#include <fcntl.h>
#include <unistd.h>
#include <sys/types.h>
#include <sys/stat.h>

#define BUFFER 4096

int TwoWayFifo::index = 0;

TwoWayFifo::TwoWayFifo(const std::string& name)
    : name_(name), show_fd_(-1), read_fd_(-1), write_fd_(-1) {
    show_fd_ = index++;
}

TwoWayFifo::~TwoWayFifo() {
    Close();
}

bool TwoWayFifo::CreatServerFile() {
    std::string file_path = "/tmp/com.ipc." + name_ + ".server";
    if (access(file_path.c_str(), F_OK) != 0) {
        if (0 != mkfifo(file_path.c_str(), 0664)) {
            perror("mkfifo error");
        }
    }
    read_fd_ = open(file_path.c_str(), O_RDONLY | O_NONBLOCK);
    return read_fd_ != 0;
}

bool TwoWayFifo::CreatClientFile() {
    std::string file_path = "/tmp/com.ipc." + name_ + ".client";
    if (access(file_path.c_str(), F_OK) != 0) {
        if (0 != mkfifo(file_path.c_str(), 0664)) {
            perror("mkfifo error");
        }
    }
    read_fd_ = open(file_path.c_str(), O_RDONLY | O_NONBLOCK);
    return read_fd_ != 0;
}

bool TwoWayFifo::OpenServerFile() {
    std::string file_path = "/tmp/com.ipc." + name_ + ".server";
    if (access(file_path.c_str(), F_OK) != 0) {
        return false;
    }
    write_fd_ = open(file_path.c_str(), O_WRONLY | O_NONBLOCK);
    return write_fd_ != 0;
}

bool TwoWayFifo::OpenClientFile() {
    std::string file_path = "/tmp/com.ipc." + name_ + ".client";
    if (access(file_path.c_str(), F_OK) != 0) {
        return false;
    }
    write_fd_ = open(file_path.c_str(), O_WRONLY | O_NONBLOCK);
    return write_fd_ != 0;
}

int TwoWayFifo::Write(const std::string& data) {
    const unsigned char* write_buf = reinterpret_cast<const unsigned char*>(data.c_str());
    int writed_size = 0;
    while (true) {
        int receive_size = write(write_fd_, write_buf + writed_size, BUFFER);
        if ((receive_size == -1 && errno == EOF) || receive_size == 0) {
            return 0;
        } else if (receive_size == -1) {
            perror("write failed");
        }
        writed_size += receive_size;
    }
}

int TwoWayFifo::Read(std::string& data) {
    unsigned char buf[BUFFER];
    while (true) {
        int read_size = read(read_fd_, buf, BUFFER);
        if (read_size == -1 && errno == EOF) {
            return 0;
        } else if (read_size == -1) {
            perror("write failed");
        }
        data.append((const char*)buf, read_size);
    }
}

void TwoWayFifo::Close() {
    if (read_fd_ != -1) {
        close(read_fd_);
    }
    if (write_fd_ != -1) {
        close(write_fd_);
    }
}