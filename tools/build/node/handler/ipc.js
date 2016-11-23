'use strict';

const ipc = require("ipc.node");

class IPCManager {
	constructor() {
		this.mgr = null;
		this.ipcId = null;
	}

	createServer(name, delegate) {
		this.mgr = ipc.newIPCManager();
		this.ipcId = ipc.createServer(this.mgr, name, delegate);
	}

	openClient(name, delegate) {
		this.mgr = ipc.newIPCManager();
		this.ipcId = ipc.openClient(this.mgr, name, delegate);	
	}

	startListener() {
		ipc.startListener(this.mgr);
	}

	registerMethod(methodName, handler) {
		ipc.registerMethod(this.mgr, this.ipcId, methodName, handler);
	}

	call(methodName, request, callback) {
		ipc.callMethod(this.mgr, this.ipcId, methodName, request, callback);
	}
}

module.exports = IPCManager;
