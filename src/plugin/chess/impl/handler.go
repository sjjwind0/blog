package impl

import (
	"encoding/json"
)

type DataPackage struct {
	Type  string    `json:"MessageType"`
	Param DataParam `json:"Param"`
}

type Option struct {
	HardType  string `json:HardType`
	FirstType string `json:FirstType`
	ColorType string `json:ColorType`
}

type DataParam struct {
	ChessType    string  `json:"ChessType"`
	Data         [][]int `json:BoardMap`
	SourcePos    []int   `json:SourcePos`
	TargetPos    []int   `json:TargetPos`
	ShowMessage  string  `json:ShowMessage`
	OptionalType string  `json:OptionalType`
}

type MessageManager struct {
	ws        *MessageController
	algorithm *AlphaBetaAlgorithm
}

func NewMessageManager(ws *MessageController) *MessageManager {
	manager := &MessageManager{}
	manager.ws = ws
	manager.algorithm = NewAlphaBetaAlgorithm()
	return manager
}

func (this *MessageManager) MessageProc(msgType string, message string) {
	dataPackage := &DataPackage{}
	json.Unmarshal([]byte(message), dataPackage)
	switch msgType {
	case "RecvData":
		{
			switch dataPackage.Type {
			case "MoveChess":
				{
					// fmt.Println("收到消息")
					// fmt.Printf("Web端移动: (%d, %d) --> (%d,%d)\n", dataPackage.Param.SourcePos[0], dataPackage.Param.SourcePos[1],
					// 	dataPackage.Param.TargetPos[0], dataPackage.Param.TargetPos[1])
					sourcePos := NewPositionWithXY(dataPackage.Param.SourcePos[0], dataPackage.Param.SourcePos[1])
					targetPos := NewPositionWithXY(dataPackage.Param.TargetPos[0], dataPackage.Param.TargetPos[1])
					turn := this.algorithm.GetBoardMap().GetTurn()
					sourceColor := this.algorithm.GetBoardMap().GetChessColor(sourcePos)
					if (turn == ChessTurnRed && sourceColor == ChessColorRed) || (turn == ChessTurnBlack && sourceColor == ChessColorBlack) {
						if this.algorithm.GetBoardMap().IsValidMove(NewChessMoveWithess(sourcePos, targetPos, 0, 0)) {
							if this.algorithm.GetBoardMap().IsGameOver() {
								this.SendAlertMessage("你赢了")
							}
							this.SendValidMove(NewChessMove(sourcePos, targetPos))
							this.algorithm.GetBoardMap().MoveChess(sourcePos, targetPos)
							this.algorithm.GetBoardMap().SwapTurn()
							move := this.CalcNextStep()
							this.algorithm.GetBoardMap().MakeMove(move)
							this.algorithm.GetBoardMap().SwapTurn()
							this.SendCurrentMove(move)
							this.SendBoardMap()
							if this.algorithm.GetBoardMap().IsGameOver() {
								this.SendAlertMessage("你输了")
							}
						} else {
							this.SendInvalidMove(NewChessMove(sourcePos, targetPos))
						}
					}

				}
			case "StartChess":
				{
					option := &Option{}
					json.Unmarshal([]byte(dataPackage.Param.OptionalType), option)
					switch option.HardType {
					case "1":
						{
							if option.FirstType == "1" {
								// 红方先手
								if option.ColorType == "1" {
									// 选择红方
									this.SendBoardMap()
								} else {
									// this.algorithm.GetBoardMap().SwapTurn()
									move := this.CalcNextStep()
									this.algorithm.GetBoardMap().MakeMove(move)
									this.algorithm.GetBoardMap().SwapTurn()
									this.SendBoardMap()
								}
							} else {
								// 黑方先手
								if option.ColorType == "1" {
									this.algorithm.GetBoardMap().SwapTurn()
									move := this.CalcNextStep()
									this.algorithm.GetBoardMap().MakeMove(move)
									this.algorithm.GetBoardMap().SwapTurn()
									this.SendBoardMap()
								} else {
									this.algorithm.GetBoardMap().SwapTurn()
									this.SendBoardMap()
								}
							}
						}
					}
				}
			}
		}
	}
}

func (this *MessageManager) SendInvalidMove(move *ChessMove) {
	boardMap := this.algorithm.GetBoardMap().GetChessMap()
	dataPackage := &DataPackage{"InvalidMove", DataParam{"BoardMap", *boardMap, []int{move.sourcePosition.X,
		move.sourcePosition.Y}, []int{move.targetPosition.X, move.targetPosition.Y}, "无效的移动", ""}}
	this.sendPackage(dataPackage)
}

func (this *MessageManager) SendCurrentMove(move *ChessMove) {
	boardMap := this.algorithm.GetBoardMap().GetChessMap()
	dataPackage := &DataPackage{"ServerMove", DataParam{"BoardMap", *boardMap, []int{move.sourcePosition.X,
		move.sourcePosition.Y}, []int{move.targetPosition.X, move.targetPosition.Y}, "服务端移动", ""}}
	this.sendPackage(dataPackage)
}

func (this *MessageManager) SendValidMove(move *ChessMove) {
	boardMap := this.algorithm.GetBoardMap().GetChessMap()
	dataPackage := &DataPackage{"ValidMove", DataParam{"BoardMap", *boardMap, []int{move.sourcePosition.X,
		move.sourcePosition.Y}, []int{move.targetPosition.X, move.targetPosition.Y}, "有效的移动", ""}}
	this.sendPackage(dataPackage)
}

func (this *MessageManager) SendAlertMessage(showMessage string) {
	dataPackage := &DataPackage{"ShowMessage", DataParam{"AlertMessage", nil, nil, nil, showMessage, ""}}
	this.sendPackage(dataPackage)
}

func (this *MessageManager) SendBoardMap() {
	boardMap := this.algorithm.GetBoardMap().GetChessMap()
	dataPackage := &DataPackage{"SetChess", DataParam{"BoardMap", *boardMap, nil, nil, "", ""}}
	this.sendPackage(dataPackage)
}

func (this *MessageManager) CalcNextStep() *ChessMove {
	var alpha int = -MaxValue
	var beta int = MaxValue
	this.algorithm.StartNextStep(alpha, beta)
	return this.algorithm.GetBestMove()
}

func (this *MessageManager) sendPackage(dataPackage *DataPackage) {
	jsonStr, _ := json.Marshal(dataPackage)
	// if err == nil {
	// 	fmt.Println("发送消息: ", string(jsonStr))
	// } else {
	// 	fmt.Println("SendError!", err)
	// }
	this.ws.SendMessage(string(jsonStr))
}

func (this *MessageManager) StartNewChess() {
	go func(this *MessageManager) {
		// if move != nil {
		// 	var message string
		// 	message = fmt.Sprintf("最佳的move: (%d, %d) --> (%d, %d)", move.sourcePosition.X, move.sourcePosition.Y,
		// 		move.targetPosition.X, move.targetPosition.Y)
		// 	this.SendAlertMessage(message)
		// }

	}(this)
}
