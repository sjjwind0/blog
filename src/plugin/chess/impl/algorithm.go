package impl

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
	this.moveList = nil
	this.AlphaBeta(SearchDeep, alpha, beta)
}

func (this *AlphaBetaAlgorithm) AlphaBeta(depth int, alpha int, beta int) int {
	var value int
	if depth == 0 {
		value = this.boardMap.EstimateValue()
		return value
	}
	var best int = -MaxValue - 1
	allMoves := this.boardMap.GetAllMoves()
	for _, move := range *allMoves {
		if this.boardMap.GetChess(move.targetPosition) == 0 || !this.boardMap.IsSameColor(move.sourcePosition, move.targetPosition) {
			this.boardMap.MakeMove(&move)
			this.moveList = append(this.moveList, move)
			this.boardMap.SwapTurn()
			value = -this.AlphaBeta(depth-1, -beta, -alpha)
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
			if best >= beta {
				break
			}
		}
	}
	return best
}
