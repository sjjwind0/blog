package impl

type Position struct {
	X int
	Y int
}

func NewPosition() *Position {
	ret := &Position{}
	return ret
}

func NewPositionWithXY(x int, y int) *Position {
	ret := &Position{x, y}
	return ret
}

func NewPositionWithPosition(pos *Position) *Position {
	ret := &Position{pos.X, pos.Y}
	return ret
}

func (this *Position) IsValidPosition() bool {
	if this.X >= 0 && this.X < BoardHeight && this.Y >= 0 && this.Y < BoardWidth {
		return true
	}

	return false
}

func (this *Position) DescriptionPosition() {
	// fmt.Printf("Position: (%d, %d)\n", this.X, this.Y)
}
