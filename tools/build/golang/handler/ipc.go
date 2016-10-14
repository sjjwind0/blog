package handler

// #cgo LDFLAGS: -ldl
// #include "export.h"
import "C"
import "unsafe"
import "sync"

var delegateMap map[unsafe.Pointer]interface{} = nil
var managerMap map[unsafe.Pointer]*IPCManager = nil
var callbackMap map[int]MethodCallback = nil
var methodMap map[int]Method = nil
var currentIndex int = 0

func goServerDelegateToCServerDelegate(i *IPCManager, delegate IPCServerDelegate) unsafe.Pointer {
	if delegateMap == nil {
		delegateMap = make(map[unsafe.Pointer]interface{})
	}
	if managerMap == nil {
		managerMap = make(map[unsafe.Pointer]*IPCManager)
	}
	ptr := unsafe.Pointer(C.NewServerInterface())
	delegateMap[ptr] = delegate
	managerMap[ptr] = i
	return (unsafe.Pointer)(ptr)
}

func goClientDelegateToCClientDelegate(i *IPCManager, delegate IPCClientDelegate) unsafe.Pointer {
	if delegateMap == nil {
		delegateMap = make(map[unsafe.Pointer]interface{})
	}
	if managerMap == nil {
		managerMap = make(map[unsafe.Pointer]*IPCManager)
	}
	ptr := unsafe.Pointer(C.NewClientInterface())
	delegateMap[ptr] = delegate
	managerMap[ptr] = i
	return ptr
}

func goMethodCallbackToCMethodCallback(callback MethodCallback) int {
	if callbackMap == nil {
		callbackMap = make(map[int]MethodCallback)
	}
	index := currentIndex
	currentIndex = currentIndex + 1
	callbackMap[index] = callback
	return index
}

func goMethodToCMethod(method Method) int {
	if methodMap == nil {
		methodMap = make(map[int]Method)
	}
	index := currentIndex
	currentIndex = currentIndex + 1
	methodMap[index] = method
	return index
}

type IPCServerDelegate interface {
	OnAcceptNewClient(manager *IPCManager, ipcID int)
	OnClientClose(manager *IPCManager, ipcID int)
}

type IPCClientDelegate interface {
	OnConnect(manager *IPCManager, ipcID int)
	OnServerClose(manager *IPCManager)
}

type Method func(request string, response *string)
type MethodCallback func(code int, response string)

var ipcManagerOnce sync.Once

type IPCManager struct {
	ipcPtr unsafe.Pointer
}

func NewIPCManager() *IPCManager {
	ipcManagerOnce.Do(func() {
		C.LoadLibrary()
	})
	i := &IPCManager{}
	i.ipcPtr = unsafe.Pointer(C.NewIPCManager())
	return i
}

func (i *IPCManager) StartListener() {
	C.StartListener(i.ipcPtr)
}

func (i *IPCManager) CreateServer(ipcName string, delegate IPCServerDelegate) int {
	return int(C.CreateServer(i.ipcPtr, C.CString(ipcName),
		goServerDelegateToCServerDelegate(i, delegate)))
}

func (i *IPCManager) OpenClient(ipcName string, delegate IPCClientDelegate) int {
	return int(C.OpenClient(i.ipcPtr, C.CString(ipcName),
		goClientDelegateToCClientDelegate(i, delegate)))
}

func (i *IPCManager) RegisterMethod(ipcID int, methodName string, method Method) {
	C.RegisterMethod(i.ipcPtr, C.int(ipcID), C.CString(methodName),
		C.CMethodFunc, C.int(goMethodToCMethod(method)))
}

func (i *IPCManager) CallMethod(ipcID int, methodName string, request string, callback MethodCallback) {
	C.CallMethod(i.ipcPtr, C.int(ipcID), C.CString(methodName), C.CString(request),
		C.CMethodCallback, C.int(goMethodCallbackToCMethodCallback(callback)))
}

func (i *IPCManager) GetNameByIPCID(ipcID int) string {
	return C.GoString(C.GetNameByIPCID(i.ipcPtr, C.int(ipcID)))
}
