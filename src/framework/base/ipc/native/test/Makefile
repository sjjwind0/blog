CC = gcc-4.9
C++ = g++-4.9
LINK = g++-4.9

LIBS=-lpthread
CCFLAGS=$(COMPILER_FLAGS) -c -g
C++FLAGS=$(COMPILER_FLAGS) -c -g -std=c++11 -Wc++11-extensions

TARGET=ipc_test

INCLUDES=-I. -I..

C++FILES=../src/fifo.cpp \
		 ../src/ipc_mgr.cpp \
		 test.cpp

OBJFILE=$(C++FILES:.cc=.o)

all:$(TARGET)


$(TARGET): $(OBJFILE)
	$(LINK) $^ $(LIBS) $(INCLUDES) -std=c++11 -o $@

%.o:%.cpp
	$(C++) -o $@ $(C++FLAGS) $< $(INCLUDES)

install:
	tsxs -i -o $(TARGET)

clean:
	rm -rf $(TARGET)