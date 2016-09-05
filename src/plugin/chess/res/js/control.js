var Control = Control || {};

Control.init = function() {
	Control.boardMap = null;
	Control.boardMapHeight = 10;
	Control.boardMapWidth = 9;
	Control.UI.init();
	Control.Message.init();
	Control.UIEvent.init();
}

Control.UIEvent = (function() {
	return {
		init: function() {
			Control.UIEvent.bindButtonEvent();
		},
		bindButtonEvent: function() {
			$(".startGame").click(function(event) {
				var hardType = $(".hardtype").val();
				var firstType = $(".first").val();
				var colorType = $(".select").val();
				var startMessage = {
					"HardType" : hardType, 
					"FirstType" : firstType,
					"colorType" : colorType
				}
				Control.Message.SendStartMessage(JSON.stringify(startMessage));
			});
		}
	}
}())

Control.Position = function(xx, yy) {
	this.x = xx;
	this.y = yy;
	this.getX = function() {
		return parseInt(this.x);
	},
	this.getY = function() {
		return parseInt(this.y);
	},
	this.setPosition = function(xx, yy) {
		this.x = x;
		this.y = yy;
	}
}

Control.Log = (function() {
	return {
		addLine: function(str) {
			$(".recodeText").append("<p>" + str + "</p>");
		},
		clean: function() {
			$(".recodeText").empty();
		}
	}
}());

Control.Map = (function() {
	return {
		init: function() {
			this.initMap();
		},
		initMapWithData: function(data) {
			Control.boardMap = {}
			for (var i = 0; i < Control.boardMapHeight; i++) {
				Control.boardMap[i] = {};
				for (var j = 0; j < Control.boardMapWidth; j++) {
					Control.boardMap[i][j] = 0;
				}
			}
			for (var i = 0; i < Control.boardMapHeight; i++) {
				for (var j = 0; j < Control.boardMapWidth; j++) {
					Control.boardMap[i][j] = data[i][j];
				}
			}
		},
		initMap: function() {
			Control.boardMap = {};
			for (var i = 0; i < Control.boardMapHeight; i++) {
				Control.boardMap[i] = {};
				for (var j = 0; j < Control.boardMapWidth; j++) {
					Control.boardMap[i][j] = 0;
				}
			}
			var initChess = function(beginPos, chessColor, increase) {
				Control.boardMap[beginPos][0] = chessColor | 0x01;
				Control.boardMap[beginPos][1] = chessColor | 0x02;
				Control.boardMap[beginPos][2] = chessColor | 0x03;
				Control.boardMap[beginPos][3] = chessColor | 0x04;
				Control.boardMap[beginPos][4] = chessColor | 0x05;
				Control.boardMap[beginPos][5] = chessColor | 0x04;
				Control.boardMap[beginPos][6] = chessColor | 0x03;
				Control.boardMap[beginPos][7] = chessColor | 0x02;
				Control.boardMap[beginPos][8] = chessColor | 0x01;

				Control.boardMap[beginPos + increase * 2][1] = chessColor | 0x06;
				Control.boardMap[beginPos + increase * 2][7] = chessColor | 0x06;

				for (var i = 0; i < Control.boardMapWidth; i += 2) {
					Control.boardMap[beginPos + increase * 3][i] = chessColor | 0x07;
				}
			}

			initChess(0, 0x0100, 1);
			initChess(Control.boardMapHeight - 1, 0x0200, -1);
		},
		isEmpty: function(pos) {
			return Control.boardMap[pos.getX()][pos.getY()] == 0;
		},
		moveChess: function(sourcePos, targetPos, callback) {
			console.log("moveChess");
			var sourceType = Control.boardMap[sourcePos.getY()][sourcePos.getX()];
			console.log(sourceType);
			if (sourceType) {
				var targetType = Control.boardMap[targetPos.getY()][targetPos.getX()];
				console.log(targetType);
				if (targetType == 0 || ((sourceType & 0xFF00) != (targetType & 0xFF00))) {
					Control.boardMap[targetPos.getY()][targetPos.getX()] = sourceType;
					Control.boardMap[sourcePos.getY()][sourcePos.getX()] = 0;
					callback();
				}
			}
		},
		getChessType: function(pos) {
			return Control.boardMap[pos.getX()][pos.getY()] & 0xFF;
		},
		getChessColor: function(pos) {
			return Control.boardMap[pos.getX()][pos.getY()] & 0xFF00;
		}, 
		getMapHeight: function() {
			return Control.boardMapHeight;
		},
		getMapWidth: function() {
			return Control.boardMapWidth;
		}
	}
}());

Control.Event = (function() {
	var currentX = -1;
	var currentY = -1;
	var logX = -1;
	var logY = -1;
	var isYourTurn = false;
	return {
		initPositin: function() {
			// currentX = -1;
			// currentY = -1;
			isYourTurn = true;
		},
		getX: function() {
			return currentX;
		},
		getY: function() {
			return currentY;
		},
		getLogSquareX: function() {
			return logX;
		},
		getLogSquareY: function() {
			return logY;
		},
		changeTurn: function() {
			if (isYourTurn == false) {
				currentX = currentY = -1;
			}
			isYourTurn = isYourTurn ? false : true;
		},
		mouseUp: function(evt) {
			if (!isYourTurn) {
				return;
			}
			console.log("mouse up");
			var nowX = parseInt(evt.offsetX / Control.chessHeight);
			var nowY = parseInt(evt.offsetY / Control.chessWidth);
			console.log(nowX + " " + nowY + " " + currentX + " " + currentY);
			if ((nowX != currentX || nowY != currentY) && (currentX != -1 && currentY != -1)) {
				var sourcePos = new Control.Position(currentX, currentY);
				var targetPos = new Control.Position(nowX, nowY);
				Control.Message.SendMoveEvent(sourcePos, targetPos);
				// Control.Map.moveChess(sourcePos, targetPos, function() {
				// 	console.log("callback");
				// 	Control.Event.changeTurn();
				// 	Control.Message.SendMoveEvent(sourcePos, targetPos);
				// });
			}
			currentX = nowX;
			currentY = nowY;
			logX = nowX;
			logY = nowY;
			Control.UI.render();
		},
		mouseDown: function(evt) {
			// currentX = parseInt(evt.offsetX / Control.chessHeight);
			// currentY = parseInt(evt.offsetY / Control.chessWidth);
		}
	}
}());

Control.UI = (function() {
	return {
		init: function() {
			Control.Map.init();
			this.initUI();
			this.bingEvent();
			// this.render();
		},
		bingEvent: function() {
			Control.canvas.addEventListener("mousedown", Control.Event.mouseDown, false);
			Control.canvas.addEventListener("mouseup", Control.Event.mouseUp, false);
		},
		preloadImage: function(){  
			this.loadImage();  
		},
		loadImage: function() {
			Control.boardImageUrl = "chess/res/img/bg.png";
			
			Control.chessImageUrlList = new Array("chess/res/img/redCar.png", "chess/res/img/redHorse.png", "chess/res/img/redElephant.png", 
				"chess/res/img/redSolider.png", "chess/res/img/redGeneral.png", "chess/res/img/redCannon.png", "chess/res/img/redPrivate.png", 
				"chess/res/img/blackCar.png", "chess/res/img/blackHorse.png", "chess/res/img/blackElephant.png", "chess/res/img/blackSolider.png", 
				"chess/res/img/blackGeneral.png", "chess/res/img/blackCannon.png", "chess/res/img/blackPrivate.png");

			Control.chessWidth = 57;
			Control.chessHeight = 58;
		},
		initUI: function() {
			this.preloadImage()
			Control.canvas = document.getElementById("board");
			Control.context = Control.canvas.getContext('2d');
		},
		render: function() {
			console.log("render UI");
			Control.context.fillStyle = "#FFF";
			Control.context.fillRect(0, 0, $(Control.canvas).attr("height"),$(Control.canvas).attr("width"));
			this.renderBoardMap();
			// this.renderSquare();
		},
		renderBoardMap: function() {
			this.preImage(Control.boardImageUrl, 0, 0, function() {
				Control.context.drawImage(this, 6, 6);
				Control.UI.renderChess();
				Control.UI.renderSquare();
			});
		},
		preImage: function(url, x, y, callback){

			var img = new Image();
			img.src = url;
			img.imageX = x;
			img.imageY = y;
			if (img.complete) {
				callback.call(img);
				return;
			}

			img.onload = function () {
				callback.call(img);
			};
		},
		renderSquare: function() {
			var x = Control.Event.getLogSquareX();
			var y = Control.Event.getLogSquareY();
			if (x >= 0 && y >= 0 && x < Control.boardMapWidth && y < Control.boardMapHeight) {
				this.preImage("chess/res/img/redBox.png", x, y, function() {
					Control.context.drawImage(this, this.imageX * Control.chessHeight, this.imageY * Control.chessWidth);
				});
			}
		},
		renderChess: function() {
			var pos = 0;
			for (var i = 0; i < Control.Map.getMapHeight(); i++) {
				for (var j = 0; j < Control.Map.getMapWidth(); j++) {
					var currentPostion = new Control.Position(i, j);
					var currentChessType = Control.Map.getChessType(currentPostion);
					if (currentChessType != 0) {
						imageIndex = Control.Map.getChessColor(currentPostion) == 0x0200 ? currentChessType + 7 : currentChessType;
						imageIndex -= 1;
						var imageX = j * Control.chessWidth;
						var imageY = i * Control.chessHeight;
						this.preImage(Control.chessImageUrlList[imageIndex], imageX, imageY, function() {
							Control.context.drawImage(this, this.imageX, this.imageY);
						});
					}
				}
			}
		}
	}
}());

Control.Message = (function() {
	var wsUri = "wss://windyx.com/message";
	var websocket;
	return {
		init: function() {
			websocket = new WebSocket(wsUri);
			websocket.onopen = function(evt) {
				Control.Message.OnOpen(evt);
			};
			websocket.onclose = function(evt) {
				Control.Message.OnClose(evt);
			};
			websocket.onmessage = function(evt) {
				Control.Message.OnMessage(evt);
			};
			websocket.onerror = function(evt) {
				Control.Message.OnError(evt);
			}
		},
		OnOpen: function(evt) {
			this.SendMessage("Hello World!");
			// do nothing
		},
		OnClose: function(evt) {
			// do nothing
		},
		OnMessage: function(evt) {
			this.parseMessage(evt);
		},
		OnError: function(evt) {
			// TODO: close
		},
		SendMessage: function(message) {
			console.log(message);
			websocket.send(message);
		},
		SendStartMessage: function(startMessage) {
			var message = {
				"MessageType" : "StartChess",
				"Param" : {
 					"OptionalType" : startMessage
				}
			}
			this.SendMessage(JSON.stringify(message));
		},
		SendMoveEvent: function(sourcePos, targetPos) {
			var message = {
				"MessageType": "MoveChess",
				"Param" : {
					"SourcePos" : [sourcePos.getY(), sourcePos.getX()],
					"TargetPos" : [targetPos.getY(), targetPos.getX()]
				}
			}
			this.SendMessage(JSON.stringify(message));
		},
		parseMessage: function(evt) {
			var data = JSON.parse(evt.data);
			console.log(data);
			switch (data["MessageType"]) {
			case "SetChess":
				Control.Event.changeTurn();
				console.log("SetChess");
				Control.Map.initMapWithData(data["Param"]["Data"]);
				Control.UI.render();
				break;
			case "ShowMessage":
				console.log("ShowMessage");
				Control.Log.addLine(data["Param"]["ShowMessage"]);
				var showMessage = data["Param"]["ShowMessage"];
				// alert(showMessage);
				console.log("ShowMessage: " + showMessage);
				alert(showMessage);
				break;
			case "InvalidMove":
				console.log("InvalidMove");
				Control.Event.initPositin();
				break;
			case "ValidMove":
				{
					console.log("ValidMove");
					var sourcePos = new Control.Position(data["Param"]["SourcePos"][1], data["Param"]["SourcePos"][0]);
					var targetPos = new Control.Position(data["Param"]["TargetPos"][1], data["Param"]["TargetPos"][0]);
					console.log(sourcePos);
					console.log(targetPos);
					Control.Log.addLine("自己移动: (" + sourcePos.getX() + ", " + sourcePos.getY() + ") -- >" + "(" + targetPos.getX() + ", " + targetPos.getY() + ")")
					Control.Map.moveChess(sourcePos, targetPos, function() {
						console.log("callback");
						Control.Event.changeTurn();
						// Control.Message.SendMoveEvent(sourcePos, targetPos);
						Control.UI.render();
					});
				}
				break;
			case "ServerMove":
				{
					console.log("ServerMove");
					var sourcePos = new Control.Position(data["Param"]["SourcePos"][1], data["Param"]["SourcePos"][0]);
					var targetPos = new Control.Position(data["Param"]["TargetPos"][1], data["Param"]["TargetPos"][0]);
					console.log(sourcePos);
					console.log(targetPos);
					Control.Log.addLine("敌方移动: (" + sourcePos.getY() + ", " + sourcePos.getX() + ") -- >" + "(" + targetPos.getY() + ", " + targetPos.getX() + ")")
				}	
				break;
			default:
				break;
			}
		}
	}
}());

window.onload = function() {
	Control.init();
}