package impl

const BoardHeight = 10
const BoardWidth = 9

type ChessType int

const (
	ChessTypeBlank    = iota
	ChessTypeCar      = iota
	ChessTypeHorse    = iota
	ChessTypeElephant = iota
	ChessTypeSolider  = iota
	ChessTypeGeneral  = iota
	ChessTypeCannon   = iota
	ChessTypePrivate  = iota
)

type ChessColor int

const (
	ChessColorRed   = 0x0100
	ChessColorBlack = 0x0200
)

type ChessTurn int

const (
	ChessTurnRed = iota
	ChessTurnBlack
)

type ChessBasicValue int

const (
	ChessBasicNoneValue     = 0
	ChessBasicCarValue      = 500
	ChessBasicHorseValue    = 350
	ChessBasicElephantValue = 250
	ChessBasicSoliderValue  = 250
	ChessBasicGeneralValue  = 8888
	ChessBasicCannonValue   = 350
	ChessBasicPrivateValue  = 100
)

type ChessLivelyValue int

const (
	ChessLivelyNoneValue     = 0
	ChessLivelyCarValue      = 6
	ChessLivelyHorseValue    = 12
	ChessLivelyElephantValue = 1
	ChessLivelySoliderValue  = 1
	ChessLivelyGeneralValue  = 0
	ChessLivelyCannonValue   = 6
	ChessLivelyPrivateValue  = 15
)

const MaxValue = 88888

const SearchDeep = 3
