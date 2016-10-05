CC = gcc
C++ = g++
LINK = g++

CCFLAGS = $(COMPILER_FLAGS) -c -g
C++FLAGS = $(COMPILER_FLAGS) -c -g -fPIC -std=c++11

TARGET=libipc.so

INCLUDES = -I.

C++FILES = src/fifo.cpp \
		   src/ipc_mgr.cpp \
		   src/ipc.cpp \
		   third_party/json/json.cpp

# OBJFILE = $(CFILES:.c=.o) $(C++FILES:.cpp=.o)
OBJFILE = $(C++FILES:.cpp=.o)

all:$(TARGET)

$(TARGET): $(OBJFILE)
	$(LINK) $^ $(INCLUDES) -fPIC -std=c++11 -shared -o $@

%.o:%.c
	$(CC) -o $@ $(CCFLAGS) $< $(INCLUDES)

%.o:%.cpp
	$(C++) -o $@ $(C++FLAGS) $< $(INCLUDES)

install:
	tsxs -i -o $(TARGET)

clean:
	rm -rf $(TARGET)
	rm -rf $(OBJFILE)