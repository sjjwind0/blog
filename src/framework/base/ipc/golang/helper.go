package golang

import "C"
import "unsafe"

//export GoOnAcceptNewClient
func GoOnAcceptNewClient(delegate unsafe.Pointer, manager unsafe.Pointer, ipcID C.int) {
	ptr := unsafe.Pointer(delegate)
	if d, ok := delegateMap[ptr]; ok {
		d.(IPCServerDelegate).OnAcceptNewClient(managerMap[ptr], int(ipcID))
		// delete(delegateMap, ptr)
		// delete(managerMap, ptr)
	}
}

//export GoOnClientClose
func GoOnClientClose(delegate unsafe.Pointer, manager unsafe.Pointer, ipcID C.int) {
	ptr := unsafe.Pointer(delegate)
	if d, ok := delegateMap[ptr]; ok {
		d.(IPCServerDelegate).OnClientClose(managerMap[ptr], int(ipcID))
		// delete(delegateMap, ptr)
		// delete(managerMap, ptr)
	}
}

//export GoOnConnect
func GoOnConnect(delegate unsafe.Pointer, manager unsafe.Pointer, ipcID C.int) {
	ptr := unsafe.Pointer(delegate)
	if d, ok := delegateMap[ptr]; ok {
		d.(IPCClientDelegate).OnConnect(managerMap[ptr], int(ipcID))
		// delete(delegateMap, ptr)
		// delete(managerMap, ptr)
	}
}

//export GoOnServerClose
func GoOnServerClose(delegate unsafe.Pointer, manager unsafe.Pointer) {
}

//export GoMethodFunc
func GoMethodFunc(f int, request *C.char) *C.char {
	rsp := ""
	methodMap[f](C.GoString(request), &rsp)
	return C.CString(rsp)
}

//export GoMethodCallback
func GoMethodCallback(f int, code C.int, response *C.char) {
	callbackMap[f](int(code), C.GoString(response))
}
