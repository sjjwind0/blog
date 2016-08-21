package impl

type ChessValueManager struct {
	basicValue        []int
	liveValue         []int
	redPrivateValue   [][]int
	blackPrivateValue [][]int
}

func NewChessCalueManager() *ChessValueManager {
	ret := &ChessValueManager{}
	ret.basicValue = []int{ChessBasicNoneValue, ChessBasicCarValue, ChessBasicHorseValue, ChessBasicElephantValue,
		ChessBasicSoliderValue, ChessBasicGeneralValue, ChessBasicCannonValue, ChessBasicPrivateValue}
	ret.liveValue = []int{ChessLivelyNoneValue, ChessLivelyCarValue, ChessLivelyHorseValue, ChessLivelyElephantValue,
		ChessLivelySoliderValue, ChessLivelyGeneralValue, ChessLivelyCannonValue,
		ChessLivelyPrivateValue}

	ret.redPrivateValue = [][]int{
		[]int{0, 0, 0, 0, 0, 0, 0, 0, 0},
		[]int{90, 90, 110, 120, 120, 120, 110, 90, 90},
		[]int{90, 90, 110, 120, 120, 120, 110, 90, 90},
		[]int{70, 90, 110, 110, 110, 110, 110, 90, 70},
		[]int{70, 70, 70, 70, 70, 70, 70, 70, 70},
		[]int{0, 0, 0, 0, 0, 0, 0, 0, 0},
		[]int{0, 0, 0, 0, 0, 0, 0, 0, 0},
		[]int{0, 0, 0, 0, 0, 0, 0, 0, 0},
		[]int{0, 0, 0, 0, 0, 0, 0, 0, 0},
		[]int{0, 0, 0, 0, 0, 0, 0, 0, 0},
	}

	ret.blackPrivateValue = [][]int{
		[]int{0, 0, 0, 0, 0, 0, 0, 0, 0},
		[]int{0, 0, 0, 0, 0, 0, 0, 0, 0},
		[]int{0, 0, 0, 0, 0, 0, 0, 0, 0},
		[]int{0, 0, 0, 0, 0, 0, 0, 0, 0},
		[]int{0, 0, 0, 0, 0, 0, 0, 0, 0},
		[]int{70, 70, 70, 70, 70, 70, 70, 70, 70},
		[]int{70, 90, 110, 110, 110, 110, 110, 90, 70},
		[]int{90, 90, 110, 120, 120, 120, 110, 90, 90},
		[]int{90, 90, 110, 120, 120, 120, 110, 90, 90},
		[]int{0, 0, 0, 0, 0, 0, 0, 0, 0},
	}

	return ret
}

func (this *ChessValueManager) GetRedPrivateValue(pos *Position) int {
	return this.redPrivateValue[pos.X][pos.Y]
}

func (this *ChessValueManager) GetBlackPrivateValue(pos *Position) int {
	return this.blackPrivateValue[pos.X][pos.Y]
}

func (this *ChessValueManager) GetBasicValue(chess int) ChessBasicValue {
	return ChessBasicValue(this.basicValue[chess&0xFF])
}

func (this *ChessValueManager) GetLiveValue(chess int) ChessLivelyValue {
	return ChessLivelyValue(this.liveValue[chess&0xFF])
}

func (this *ChessValueManager) GetPrivateValue(pos Position, chessColor ChessColor) ChessBasicValue {
	if chessColor == ChessColorRed {
		return ChessBasicValue(this.redPrivateValue[pos.X][pos.Y])
	} else {
		return ChessBasicValue(this.blackPrivateValue[pos.X][pos.Y])
	}
}
