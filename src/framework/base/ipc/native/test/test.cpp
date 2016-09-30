#include "include/ipc_mgr.h"

#include <iostream>

void func(void* args) {

}

int main(int argc, const char* argv[]) {
    std::cout << "Hello World!" << std::endl;
    pid_t fpid = fork();
    if (pid_t == 0) {

    } else {
        pthread_t tid;
        pthread_create(&tid, NULL, func, NULL);
    }
    return 0;
}