package maze

type Position struct {
	X int
	Y int
}

type nextPosition struct {
	Pos  Position
	Dist int
}
