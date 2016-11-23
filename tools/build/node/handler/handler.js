'use strict';

const IPCManager = require("./ipc");
const BasicRunnerManager = require("../run/runnerMgr");

class MainRunner extends BasicRunnerManager {
	constructor(pluginName) {
		super();
		this.pluginName = pluginName;
	}

	onConnect(manager, ipcId) {
		console.log("connect server success");
	}

	onAcceptNewClient(manager, ipcId) {
		console.log("accept new client");
	}

	onClientClose(manager, ipcId) {
		console.log("client close");
	}

	onServerClose(manager) {
		console.log("server close");
	}

	handleIPCRequest(request) {
		return super.handleIPCRequest(request);
	}

	start() {
		var self = this;

		let client = new IPCManager();
		let ipcId = client.openClient(this.pluginName, this);
		client.registerMethod("HttpRequest", function(request) {
			let rsp = self.handleIPCRequest(request);
			return rsp;
		})
		client.startListener();
	}
}

module.exports = MainRunner;