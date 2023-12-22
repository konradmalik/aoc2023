package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type Cell struct {
	value rune
}

func NewCellsFromLine(s string) []Cell {
	cells := make([]Cell, len(s))
	for i, r := range s {
		cells[i] = Cell{value: r}
	}
	return cells
}

func (c Cell) IsGarden() bool {
	return c.value == '.' || c.value == 'S'
}

func (c Cell) String() string {
	return string(c.value)
}

type Position struct {
	X     int
	Y     int
	steps int
}

type Grid struct {
	data [][]Cell
}

func (g Grid) MustGet(p Position) Cell {
	return g.data[p.Y][p.X]
}

func (g Grid) Get(p Position) (Cell, bool) {
	if p.X < 0 || p.Y < 0 {
		return Cell{}, false
	}

	if p.X >= len(g.data[0]) || p.Y >= len(g.data) {
		return Cell{}, false
	}

	return g.MustGet(p), true
}

func (g Grid) FindStart() Position {
	var sb strings.Builder
	for y, row := range g.data {
		for x, c := range row {
			if c.value == 'S' {
				return Position{X: x, Y: y}
			}
		}
		sb.WriteRune('\n')
	}
	panic("no start")
}

func (g Grid) WalkAndCountGardens(steps int) int {
	s := g.FindStart()
	s.steps = steps
	cache := make(map[Position][]Position)
	return len(g.traverse(s, cache))
}

func (g Grid) traverse(p Position, cache map[Position][]Position) []Position {
	if cached, found := cache[p]; found {
		return cached
	}

	if p.steps <= 0 {
		ret := []Position{p}
		cache[p] = ret
		return ret
	}

	gardens := make([]Position, 0)
	positions := []Position{{p.X - 1, p.Y, p.steps - 1}, {p.X, p.Y - 1, p.steps - 1}, {p.X + 1, p.Y, p.steps - 1}, {p.X, p.Y + 1, p.steps - 1}}

	for _, p := range positions {
		if c, ok := g.Get(p); ok && c.IsGarden() {
			gardens = append(gardens, g.traverse(p, cache)...)
		}
	}

	// unique
	unique := make([]Position, 0)
	tracker := make(map[Position]struct{})
	for _, g := range gardens {
		if _, found := tracker[g]; !found {
			unique = append(unique, g)
			tracker[g] = struct{}{}
		}
	}

	cache[p] = unique
	return unique
}

func (g Grid) String() string {
	var sb strings.Builder
	for _, row := range g.data {
		for _, c := range row {
			sb.WriteRune(c.value)
		}
		sb.WriteRune('\n')
	}
	return sb.String()
}

func NewGrid(filepath string) Grid {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	data := make([][]Cell, 0)
	for scanner.Scan() {
		line := scanner.Text()
		data = append(data, NewCellsFromLine(line))
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return Grid{data: data}
}

func main() {
	g := NewGrid("./input.txt")
	fmt.Println(g)
	// part 1
	fmt.Println(g.WalkAndCountGardens(64))
}
