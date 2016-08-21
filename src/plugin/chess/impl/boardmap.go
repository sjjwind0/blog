package impl

import (
	"fmt"
)

type ChessMove struct {
	sourcePosition *Position
	targetPosition *Position
	sourceChess    int
	targetChess    int
}

func NewChessMove(sourcePos *Position, targetPos *Position) *ChessMove {
	move := &ChessMove{}
	move.sourcePosition = sourcePos
	move.targetPosition = targetPos
	move.sourceChess = 0
	move.targetChess = 0
	return move
}

func NewChessMoveWithess(sourcePos *Position, targetPos *Position, sourceChess int, targetChess int) *ChessMove {
	move := &ChessMove{}
	move.sourcePosition = NewPositionWithPosition(sourcePos)
	move.targetPosition = NewPositionWithPosition(targetPos)
	move.sourceChess = sourceChess
	move.targetChess = targetChess
	return move
}

func NewChessMoveWithChessMove(move *ChessMove) *ChessMove {
	ret := &ChessMove{}
	ret.sourcePosition = NewPositionWithPosition(move.sourcePosition)
	ret.targetPosition = NewPositionWithPosition(move.targetPosition)
	ret.sourceChess = move.sourceChess
	ret.targetChess = move.targetChess
	return ret
}

var nextCarStepStep = [][]int{[]int{1, 0}, []int{-1, 0}, []int{0, 1}, []int{0, -1}}
var nextHorseStepStep = [][]int{[]int{2, 1}, []int{2, -1}, []int{-2, 1}, []int{-2, -1}, []int{1, 2}, []int{1, -2}, []int{-1, 2}, []int{-1, -2}}
var nextElephantStepStep = [][]int{[]int{2, 2}, []int{2, -2}, []int{-2, 2}, []int{-2, -2}}
var nextSoliderStepStep = [][]int{[]int{1, 1}, []int{1, -1}, []int{-1, 1}, []int{-1, -1}}
var nextGeneralStepStep = [][]int{[]int{1, 0}, []int{-1, 0}, []int{0, 1}, []int{0, -1}}
var nextCannonStepStep = [][]int{[]int{1, 0}, []int{-1, 0}, []int{0, 1}, []int{0, -1}}

type BoardMap struct {
	boardMap         [][]int     // 棋盘
	chessValue       [][]int     // 棋盘值
	flexibility      [][]int     // 灵活度
	safelyPos        [][]int     // 受保护的程度
	attackPos        [][]int     // 受威胁的程度
	chessTurn        ChessTurn   // 轮到哪一方
	allMoves         []ChessMove // 所有可供移动的
	currentMoveIndex int         // 当前移动到了哪一个
	valueManager     *ChessValueManager
}

func NewBoardMap() *BoardMap {
	ret := &BoardMap{}
	ret.valueManager = NewChessCalueManager()
	ret.currentMoveIndex = 0
	return ret
}

func (this *BoardMap) InitMap(needPutChess bool) {
	if this.boardMap == nil {
		this.boardMap = make([][]int, BoardHeight)
		for i := 0; i < BoardHeight; i++ {
			this.boardMap[i] = make([]int, BoardWidth)
		}
	}

	if this.flexibility == nil {
		this.flexibility = make([][]int, BoardHeight)
		for i := 0; i < BoardHeight; i++ {
			this.flexibility[i] = make([]int, BoardWidth)
		}
	}

	if this.chessValue == nil {
		this.chessValue = make([][]int, BoardHeight)
		for i := 0; i < BoardHeight; i++ {
			this.chessValue[i] = make([]int, BoardWidth)
		}
	}

	if this.safelyPos == nil {
		this.safelyPos = make([][]int, BoardHeight)
		for i := 0; i < BoardHeight; i++ {
			this.safelyPos[i] = make([]int, BoardWidth)
		}
	}

	if this.attackPos == nil {
		this.attackPos = make([][]int, BoardHeight)
		for i := 0; i < BoardHeight; i++ {
			this.attackPos[i] = make([]int, BoardWidth)
		}
	}

	for i := 0; i < BoardHeight; i++ {
		for j := 0; j < BoardWidth; j++ {
			this.boardMap[i][j] = 0
			this.flexibility[i][j] = 0
			this.chessValue[i][j] = 0
			this.safelyPos[i][j] = 0
			this.attackPos[i][j] = 0
		}
	}

	this.chessTurn = ChessTurnRed

	if needPutChess == false {
		return
	}

	var initChess = func(beginPos int, increase int, chessColor ChessColor) {
		this.boardMap[beginPos][0] = int(chessColor | ChessTypeCar)
		this.boardMap[beginPos][1] = int(chessColor | ChessTypeHorse)
		this.boardMap[beginPos][2] = int(chessColor | ChessTypeElephant)
		this.boardMap[beginPos][3] = int(chessColor | ChessTypeSolider)
		this.boardMap[beginPos][4] = int(chessColor | ChessTypeGeneral)
		this.boardMap[beginPos][5] = int(chessColor | ChessTypeSolider)
		this.boardMap[beginPos][6] = int(chessColor | ChessTypeElephant)
		this.boardMap[beginPos][7] = int(chessColor | ChessTypeHorse)
		this.boardMap[beginPos][8] = int(chessColor | ChessTypeCar)

		this.boardMap[beginPos+increase*2][1] = int(chessColor | ChessTypeCannon)
		this.boardMap[beginPos+increase*2][7] = int(chessColor | ChessTypeCannon)

		for i := 0; i < BoardWidth; i += 2 {
			this.boardMap[beginPos+increase*3][i] = int(chessColor | ChessTypePrivate)
		}
	}

	initChess(0, 1, ChessColorBlack)
	initChess(BoardHeight-1, -1, ChessColorRed)
}

func (this *BoardMap) ShowBoardMap() {
	fmt.Println("\n====================================\n")
	for i := 0; i < BoardHeight; i++ {
		for j := 0; j < BoardWidth; j++ {
			pos := NewPositionWithXY(i, j)
			chessType := this.GetChessType(pos)
			ChessColor := this.GetChessColor(pos)
			if ChessColor == ChessColorRed {
				switch chessType {
				case ChessTypeCar:
					fmt.Printf("車")
				case ChessTypeHorse:
					fmt.Printf("马")
				case ChessTypeElephant:
					fmt.Printf("象")
				case ChessTypeSolider:
					fmt.Printf("士")
				case ChessTypeGeneral:
					fmt.Printf("将")
				case ChessTypeCannon:
					fmt.Printf("炮")
				case ChessTypePrivate:
					fmt.Printf("兵")
				default:
					fmt.Printf("  ")
				}
			} else {
				switch chessType {
				case ChessTypeCar:
					fmt.Printf("车")
				case ChessTypeHorse:
					fmt.Printf("馬")
				case ChessTypeElephant:
					fmt.Printf("相")
				case ChessTypeSolider:
					fmt.Printf("仕")
				case ChessTypeGeneral:
					fmt.Printf("帅")
				case ChessTypeCannon:
					fmt.Printf("包")
				case ChessTypePrivate:
					fmt.Printf("卒")
				default:
					fmt.Printf("  ")
				}
			}
			fmt.Printf("  ")
		}
		fmt.Printf("\n")
	}
	fmt.Println("\n====================================\n")
}

func (this *BoardMap) SetChess(pos *Position, chess int) {
	if pos.IsValidPosition() {
		this.boardMap[pos.X][pos.Y] = chess
	}
}

func (this *BoardMap) MoveChess(sourcePos *Position, targetPos *Position) {
	if sourcePos.IsValidPosition() && !this.IsBlank(sourcePos) {
		if targetPos.IsValidPosition() {
			if this.IsBlank(targetPos) || this.GetChessColor(sourcePos) != this.GetChessColor(targetPos) {
				this.SetChess(targetPos, this.boardMap[sourcePos.X][sourcePos.Y])
				this.SetChess(sourcePos, 0)
			}
		}
	}
}

func (this *BoardMap) GetTurn() ChessTurn {
	return this.chessTurn
}

func (this *BoardMap) SwapTurn() {
	if this.chessTurn == ChessTurnBlack {
		this.chessTurn = ChessTurnRed
	} else {
		this.chessTurn = ChessTurnBlack
	}
}

func (this *BoardMap) GetChessMap() *[][]int {
	return &this.boardMap
}

func (this *BoardMap) GetChess(pos *Position) int {
	return this.boardMap[pos.X][pos.Y]
}

func (this *BoardMap) GetChessType(pos *Position) ChessType {
	return ChessType(this.boardMap[pos.X][pos.Y] & 0xFF)
}

func (this *BoardMap) GetChessColor(pos *Position) ChessColor {
	return ChessColor(this.boardMap[pos.X][pos.Y] & 0xFF00)
}

func (this *BoardMap) GetBasicValue(pos *Position) ChessBasicValue {
	return ChessBasicValue(this.valueManager.GetBasicValue(this.boardMap[pos.X][pos.Y]))
}

func (this *BoardMap) GetLiveValue(pos *Position) ChessLivelyValue {
	return ChessLivelyValue(this.valueManager.GetLiveValue(this.boardMap[pos.X][pos.Y]))
}

func (this *BoardMap) IsGameOver() bool {
	var hasRedGeneral bool = false
	var hasBlackGeneral bool = false
	RedGeneral := int(ChessColorRed | ChessTypeGeneral)
	BlackGeneral := int(ChessColorBlack | ChessTypeGeneral)
	for i := 0; i < 3; i++ {
		for j := 3; j < 6; j++ {
			if this.boardMap[i][j] == RedGeneral {
				hasRedGeneral = true
			}

			if this.boardMap[i][j] == BlackGeneral {
				hasBlackGeneral = true
			}
		}
	}

	for i := 7; i < BoardHeight; i++ {
		for j := 3; j < 6; j++ {
			if this.boardMap[i][j] == RedGeneral {
				hasRedGeneral = true
			}

			if this.boardMap[i][j] == BlackGeneral {
				hasBlackGeneral = true
			}
		}
	}
	return hasBlackGeneral == false || hasRedGeneral == false
}

func (this *BoardMap) IsBlank(pos *Position) bool {
	return this.boardMap[pos.X][pos.Y] == ChessTypeBlank
}

func (this *BoardMap) IsOverRiver(pos *Position) bool {
	if pos.X >= 0 && pos.X < BoardHeight/2 {
		if this.IsBlack(pos) {
			return false
		} else {
			return true
		}
	} else {
		if this.IsRed(pos) {
			return false
		} else {
			return true
		}
	}
}

func (this *BoardMap) IsSameColor(pos1 *Position, pos2 *Position) bool {
	return this.GetChessColor(pos1) == this.GetChessColor(pos2)
}

func (this *BoardMap) IsSelfChess(pos *Position) bool {
	if this.chessTurn == ChessTurnRed {
		if this.GetChessColor(pos) == ChessColorRed {
			return true
		} else {
			return false
		}
	} else {
		if this.GetChessColor(pos) == ChessColorBlack {
			return true
		} else {
			return false
		}
	}
}

func (this *BoardMap) IsBlack(pos *Position) bool {
	return this.GetChessColor(pos) == ChessColorBlack
}

func (this *BoardMap) IsRed(pos *Position) bool {
	return this.GetChessColor(pos) == ChessColorRed
}

func (this *BoardMap) IsValidMove(move *ChessMove) bool {
	if move.sourcePosition.IsValidPosition() {
		positionList := this.GetNextStep(move.sourcePosition)
		for _, movePos := range *positionList {
			if move.targetPosition.X == movePos.X && move.targetPosition.Y == movePos.Y {
				return true
			}
		}
	}

	return false
}

func (this *BoardMap) GetAllMoves() *[]ChessMove {
	var allMoves []ChessMove = nil
	for i := 0; i < BoardHeight; i++ {
		for j := 0; j < BoardWidth; j++ {
			pos := NewPositionWithXY(i, j)
			if !this.IsBlank(pos) {
				if this.IsSelfChess(pos) {
					posList := this.GetNextStep(pos)
					for _, movePos := range *posList {
						currentMove := NewChessMoveWithess(pos, &movePos, this.GetChess(pos), this.GetChess(&movePos))
						allMoves = append(allMoves, *currentMove)
					}
				}
			}
		}
	}

	return &allMoves
}

func (this *BoardMap) MakeMove(move *ChessMove) {
	this.boardMap[move.targetPosition.X][move.targetPosition.Y] = this.boardMap[move.sourcePosition.X][move.sourcePosition.Y]
	this.boardMap[move.sourcePosition.X][move.sourcePosition.Y] = 0
}

func (this *BoardMap) UnMakeMove(move *ChessMove) {
	this.boardMap[move.sourcePosition.X][move.sourcePosition.Y] = move.sourceChess
	this.boardMap[move.targetPosition.X][move.targetPosition.Y] = move.targetChess
}

func (this *BoardMap) CopyBoard(boardMap *BoardMap) {
	for i := 0; i < BoardHeight; i++ {
		for j := 0; j < BoardWidth; j++ {
			currentPos := NewPositionWithXY(i, j)
			this.SetChess(currentPos, boardMap.GetChess(currentPos))
		}
	}
}

func (this *BoardMap) GetNextBoard() *BoardMap {
	if this.currentMoveIndex < len(this.allMoves) {
		this.currentMoveIndex++
		newBoardMap := NewBoardMap()
		newBoardMap.InitMap(false)
		newBoardMap.CopyBoard(this)
		return newBoardMap
	} else {
		return nil
	}
}

var chessTime1 int64 = 0
var chessTime2 int64 = 0
var chessTime3 int64 = 0
var chessTime4 int64 = 0
var chessTime5 int64 = 0
var chessTime6 int64 = 0
var chessTime7 int64 = 0

func (this *BoardMap) ResetChessTime() {
	chessTime1 = 0
	chessTime2 = 0
	chessTime3 = 0
	chessTime4 = 0
	chessTime5 = 0
	chessTime6 = 0
	chessTime7 = 0
}

func (this *BoardMap) GetNextStep(pos *Position) *[]Position {
	if pos.IsValidPosition() == false {
		return nil
	}
	switch this.GetChessType(pos) {
	case ChessTypeCar:
		return this.nextCarStep(pos)
	case ChessTypeHorse:
		return this.nextHorseStep(pos)
	case ChessTypeElephant:
		return this.nextElephantStep(pos)
	case ChessTypeSolider:
		return this.nextSoliderStep(pos)
	case ChessTypeGeneral:
		return this.nextGeneralStep(pos)
	case ChessTypeCannon:
		return this.nextCannonStep(pos)
	case ChessTypePrivate:
		return this.nextPrivateStep(pos)
	default:
		fmt.Println("Type Error!", pos.X, pos.Y)
	}
	return nil
}

func (this *BoardMap) nextCarStep(pos *Position) *[]Position {
	var positionList []Position = nil
	for i := 0; i < 4; i++ {
		stepX := nextCarStepStep[i][0]
		stepY := nextCarStepStep[i][1]
		currentPosition := NewPositionWithPosition(pos)
		for true {
			currentPosition.X += stepX
			currentPosition.Y += stepY
			if currentPosition.IsValidPosition() {
				if this.IsBlank(currentPosition) {
					positionList = append(positionList, *currentPosition)
				} else {
					positionList = append(positionList, *currentPosition)
					break
				}
			} else {
				break
			}
		}
	}

	return &positionList
}

func (this *BoardMap) nextHorseStep(pos *Position) *[]Position {
	// 标记为是否走过
	var positionList []Position = nil
	for i := 0; i < 8; i++ {
		stepX := nextHorseStepStep[i][0]
		stepY := nextHorseStepStep[i][1]
		currentPosition := NewPositionWithPosition(pos)
		currentPosition.X += stepX
		currentPosition.Y += stepY
		if currentPosition.IsValidPosition() {
			footPosition := NewPosition()
			if stepX == 2 || stepX == -2 {
				footPosition.X = pos.X + stepX/2
				footPosition.Y = pos.Y
			} else {
				footPosition.X = pos.X
				footPosition.Y = pos.Y + stepY/2
			}
			// 判断蹩脚
			if this.IsBlank(footPosition) {
				positionList = append(positionList, *currentPosition)
			}
		}

	}
	return &positionList
}

func (this *BoardMap) nextElephantStep(pos *Position) *[]Position {
	var positionList []Position = nil
	for i := 0; i < 4; i++ {
		stepX := nextElephantStepStep[i][0]
		stepY := nextElephantStepStep[i][1]
		currentPosition := NewPositionWithPosition(pos)
		currentPosition.X += stepX
		currentPosition.Y += stepY
		if currentPosition.IsValidPosition() {
			midPosition := NewPosition()
			midPosition.X = pos.X + stepX/2
			midPosition.Y = pos.Y + stepY/2
			if (pos.X >= 5 && currentPosition.X >= 5) || (pos.X < 5 && currentPosition.X < 5) {
				if this.IsBlank(midPosition) {
					positionList = append(positionList, *currentPosition)
				}
			}
		}
	}
	return &positionList
}

func (this *BoardMap) nextSoliderStep(pos *Position) *[]Position {
	var positionList []Position = nil
	for i := 0; i < 4; i++ {
		stepX := nextSoliderStepStep[i][0]
		stepY := nextSoliderStepStep[i][1]
		currentPosition := NewPositionWithPosition(pos)
		currentPosition.X += stepX
		currentPosition.Y += stepY
		if currentPosition.IsValidPosition() {
			if this.GetChessColor(pos) == ChessColorRed {
				if currentPosition.X < BoardHeight && currentPosition.X >= 7 && currentPosition.Y <= 5 && currentPosition.Y >= 3 {
					positionList = append(positionList, *currentPosition)
				}
			} else {
				if currentPosition.X < 3 && currentPosition.X >= 0 && currentPosition.Y <= 5 && currentPosition.Y >= 3 {
					positionList = append(positionList, *currentPosition)
				}
			}
		}
	}
	return &positionList
}

func (this *BoardMap) nextGeneralStep(pos *Position) *[]Position {
	var positionList []Position = nil
	for i := 0; i < 4; i++ {
		stepX := nextGeneralStepStep[i][0]
		stepY := nextGeneralStepStep[i][1]
		currentPosition := NewPositionWithPosition(pos)
		currentPosition.X += stepX
		currentPosition.Y += stepY
		if ((currentPosition.X <= 2 && currentPosition.X >= 0) || (currentPosition.X <= 9 && currentPosition.X >= 7)) && currentPosition.Y <= 5 && currentPosition.Y >= 3 {
			if currentPosition.IsValidPosition() {
				positionList = append(positionList, *currentPosition)
			}
		}
	}

	currentPosition := NewPositionWithPosition(pos)
	if this.GetChessColor(currentPosition) == ChessColorRed {
		currentPosition.X--
		for currentPosition.X >= 0 && currentPosition.IsValidPosition() {
			if this.IsBlank(currentPosition) {
				// do nothing
			} else {
				if this.GetChess(currentPosition) != (ChessColorBlack | ChessTypeGeneral) {
					break
				} else {
					positionList = append(positionList, *currentPosition)
					break
				}
			}
			currentPosition.X--
		}
	} else {
		currentPosition.X++
		for currentPosition.X < BoardHeight && currentPosition.IsValidPosition() {
			if this.IsBlank(currentPosition) {
				// do nothing
			} else {
				if this.GetChess(currentPosition) != (ChessColorRed | ChessTypeGeneral) {
					break
				} else {
					positionList = append(positionList, *currentPosition)
					break
				}
			}
			currentPosition.X++
		}
	}

	return &positionList
}

func (this *BoardMap) nextCannonStep(pos *Position) *[]Position {
	var positionList []Position = nil
	for i := 0; i < 4; i++ {
		stepX := nextCannonStepStep[i][0]
		stepY := nextCannonStepStep[i][1]
		currentPosition := NewPositionWithPosition(pos)
		for true {
			currentPosition.X += stepX
			currentPosition.Y += stepY
			if currentPosition.IsValidPosition() {
				if this.IsBlank(currentPosition) {
					positionList = append(positionList, *currentPosition)
				} else {
					// 继续往前判断，隔山打牛
					for currentPosition.IsValidPosition() {
						currentPosition.X += stepX
						currentPosition.Y += stepY
						if currentPosition.IsValidPosition() && !this.IsBlank(currentPosition) {
							positionList = append(positionList, *currentPosition)
							break
						}
					}
					break
				}
			} else {
				break
			}
		}
	}
	return &positionList
}

func (this *BoardMap) nextPrivateStep(pos *Position) *[]Position {
	var positionList []Position = nil
	if this.GetChessColor(pos) == ChessColorRed {
		// red private，only one point to move
		if pos.X >= 5 {
			currentPosition := NewPositionWithPosition(pos)
			currentPosition.X--
			positionList = append(positionList, *currentPosition)
		} else {
			// over river, it can move forworad, left, right
			step := [][]int{[]int{-1, 0}, []int{0, 1}, []int{0, -1}}
			for i := 0; i < 3; i++ {
				currentPosition := NewPositionWithPosition(pos)
				currentPosition.X += step[i][0]
				currentPosition.Y += step[i][1]
				if currentPosition.IsValidPosition() {
					positionList = append(positionList, *currentPosition)
				}
			}

		}
	} else {
		if pos.X < 5 {
			currentPosition := NewPositionWithPosition(pos)
			currentPosition.X++
			positionList = append(positionList, *currentPosition)
		} else {
			// over river, it can move forworad, left, right
			step := [][]int{[]int{1, 0}, []int{0, 1}, []int{0, -1}}
			for i := 0; i < 3; i++ {
				currentPosition := NewPositionWithPosition(pos)
				currentPosition.X += step[i][0]
				currentPosition.Y += step[i][1]
				if currentPosition.IsValidPosition() {
					positionList = append(positionList, *currentPosition)
				}
			}

		}
	}
	return &positionList
}

func (this *BoardMap) GetPrivateValue(pos *Position) int {
	if this.GetChessColor(pos) == ChessColorRed {
		return this.valueManager.GetRedPrivateValue(pos)
	} else {
		return this.valueManager.GetBlackPrivateValue(pos)
	}
	return 0
}

func (this *BoardMap) CleanAllValues() {
	for i := 0; i < BoardHeight; i++ {
		for j := 0; j < BoardWidth; j++ {
			this.chessValue[i][j] = 0
			this.flexibility[i][j] = 0
			this.attackPos[i][j] = 0
			this.safelyPos[i][j] = 0
		}
	}
}

func (this *BoardMap) EstimateValue() int {
	this.CleanAllValues()
	for i := 0; i < BoardHeight; i++ {
		for j := 0; j < BoardWidth; j++ {
			currentPos := NewPositionWithXY(i, j)
			if currentPos.IsValidPosition() && !this.IsBlank(currentPos) {
				positionList := this.GetNextStep(currentPos)
				for _, pos := range *positionList {
					if pos.IsValidPosition() {
						if this.IsBlank(&pos) {
							// 能够到达的位置，该点的灵活度+1
							this.flexibility[i][j]++
						} else {
							if this.GetChessColor(currentPos) == this.GetChessColor(&pos) {
								this.safelyPos[pos.X][pos.Y] += int((30 + (this.GetBasicValue(&pos)-this.GetBasicValue(currentPos))>>3) >> 3)
							} else {
								// this.attackPos[pos.X][pos.Y]++
								this.flexibility[i][j]++
								// 判断威胁的棋子
								switch this.boardMap[pos.X][pos.Y] {
								case (ChessTypeGeneral | ChessColorRed):
									if this.chessTurn == ChessTurnBlack {
										return MaxValue
									}
								case (ChessTypeGeneral | ChessColorBlack):
									if this.chessTurn == ChessTurnRed {
										return MaxValue
									}
								default:
									this.attackPos[pos.X][pos.Y] += int((30 + (this.GetBasicValue(currentPos)-this.GetBasicValue(&pos))>>3) >> 3)
								}
							}
						}
					}
				}
			}
		}
	}

	for i := 0; i < BoardHeight; i++ {
		for j := 0; j < BoardWidth; j++ {
			currentPos := NewPositionWithXY(i, j)
			if !this.IsBlank(currentPos) {
				this.chessValue[i][j]++
				this.chessValue[i][j] += int(this.GetLiveValue(currentPos)) * this.flexibility[i][j]
				if this.GetChessType(currentPos) == ChessTypePrivate {
					this.chessValue[i][j] += int(this.GetPrivateValue(currentPos))
				}
			}
		}
	}

	halfValue := 0
	for i := 0; i < BoardHeight; i++ {
		for j := 0; j < BoardWidth; j++ {
			currentPos := NewPositionWithXY(i, j)
			if !this.IsBlank(currentPos) {
				halfValue = int(this.GetBasicValue(currentPos)) >> 4
				this.chessValue[i][j] += int(this.GetBasicValue(currentPos))
				currentChessType := this.GetChessType(currentPos)
				var judgeValue = func(chessTurn ChessTurn) int {
					if this.attackPos[i][j] != 0 {
						if this.chessTurn == chessTurn {
							if currentChessType == ChessTypeGeneral {
								this.chessValue[i][j] -= 20
							} else {
								this.chessValue[i][j] -= halfValue << 1
								if this.safelyPos[i][j] != 0 {
									this.chessValue[i][j] += halfValue
								}
							}
						} else {
							if currentChessType == ChessTypeGeneral {
								return MaxValue
							}
							this.chessValue[i][j] -= halfValue * 10

							if this.safelyPos[i][j] != 0 {
								this.chessValue[i][j] += halfValue * 9
							}
						}
						this.chessValue[i][j] -= this.attackPos[i][j]
					} else {
						if this.safelyPos[i][j] != 0 {
							this.chessValue[i][j] += 5
						}
					}

					return 0
				}

				if this.IsRed(currentPos) {
					if MaxValue == judgeValue(ChessTurnRed) {
						return MaxValue
					}
				} else {
					if MaxValue == judgeValue(ChessTurnBlack) {
						return MaxValue
					}
				}
			}
		}
	}

	redValue := 0
	blackValue := 0
	for i := 0; i < BoardHeight; i++ {
		for j := 0; j < BoardWidth; j++ {
			currentPos := NewPositionWithXY(i, j)
			if !this.IsBlank(currentPos) {
				if this.IsRed(currentPos) {
					redValue += this.chessValue[i][j]
				} else {
					blackValue += this.chessValue[i][j]
				}
			}
		}
	}

	if this.chessTurn == ChessTurnRed {
		return redValue - blackValue
	} else {
		return blackValue - redValue
	}

	return 0
}
