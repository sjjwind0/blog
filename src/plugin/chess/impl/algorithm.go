package impl

import (
	"fmt"
	"time"
)

var time1 int64 = 0
var time2 int64 = 0
var time3 int64 = 0
var time4 int64 = 0
var time5 int64 = 0

type AlphaBetaAlgorithm struct {
	boardMap *BoardMap
	bestMove *ChessMove
	moveList []ChessMove
}

func NewAlphaBetaAlgorithm() *AlphaBetaAlgorithm {
	ret := &AlphaBetaAlgorithm{}
	ret.boardMap = NewBoardMap()
	ret.boardMap.InitMap(true)
	return ret
}

func (this *AlphaBetaAlgorithm) GetBoardMap() *BoardMap {
	return this.boardMap
}

func (this *AlphaBetaAlgorithm) GetBestMove() *ChessMove {
	return this.bestMove
}

func (this *AlphaBetaAlgorithm) StartNextStep(alpha int, beta int) {
	time1 = 0
	time2 = 0
	time3 = 0
	this.boardMap.ResetChessTime()
	this.AlphaBeta(SearchDeep, alpha, beta)
	fmt.Println("time1: ", time1)
	fmt.Println("time3: ", time2)
	fmt.Println("time3: ", time3)
	fmt.Printf("######################################")
	this.boardMap.ShowChessTime()
}

func (this *AlphaBetaAlgorithm) AlphaBeta(depth int, alpha int, beta int) int {
	var value int
	time1_1 := time.Now().UnixNano()
	if depth == 0 {
		value = this.boardMap.EstimateValue()
		return value
	}
	time1_2 := time.Now().UnixNano()
	time1 += time1_2 - time1_1
	var best int = -MaxValue - 1
	allMoves := this.boardMap.GetAllMoves()
	time2_1 := time.Now().UnixNano()
	time2 += time2_1 - time1_2
	// this.moveList = nil
	for _, move := range *allMoves {
		if this.boardMap.GetChess(move.targetPosition) == 0 || !this.boardMap.IsSameColor(move.sourcePosition, move.targetPosition) {
			var time3_1 int64 = time.Now().UnixNano()
			this.boardMap.MakeMove(&move)
			this.moveList = append(this.moveList, move)
			this.boardMap.SwapTurn()
			var time3_2 int64 = time.Now().UnixNano()
			value = -this.AlphaBeta(depth-1, -beta, -alpha)
			var time3_3 int64 = time.Now().UnixNano()
			this.boardMap.SwapTurn()
			this.boardMap.UnMakeMove(&move)
			if value > best {
				best = value
				if depth == SearchDeep {
					this.bestMove = NewChessMoveWithChessMove(&(this.moveList[0]))
				}
			}
			this.moveList = this.moveList[:len(this.moveList)-1]
			if value > alpha {
				alpha = best
			}
			var time3_4 int64 = time.Now().UnixNano()
			if best >= beta {
				time3 += time3_2 - time3_1 + time3_4 - time3_3
				break
			}
		}
	}
	return best
}
