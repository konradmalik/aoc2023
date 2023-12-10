package maze

import (
	"fmt"
	"slices"
)

type Tile struct {
	Label         rune
	Visited       bool
	Dist          int
	Intersections int
}

func (t Tile) IsLoop() bool {
	return t.Visited
}

func (t Tile) String() string {
	return fmt.Sprintf("{%s, %d}", string(t.Label), t.Dist)
}

func (t Tile) ConnectedWithSouth() bool {
	return slices.Contains(
		[]rune{'S', '|', '7', 'F'},
		t.Label)
}

func (t Tile) ConnectedWithNorth() bool {
	return slices.Contains(
		[]rune{'S', '|', 'L', 'J'},
		t.Label)
}

func (t Tile) ConnectedWithWest() bool {
	return slices.Contains(
		[]rune{'S', '-', '7', 'J'},
		t.Label)
}

func (t Tile) ConnectedWithEast() bool {
	return slices.Contains(
		[]rune{'S', '-', 'L', 'F'},
		t.Label)
}
