package main

import (
	"fmt"
	"testing"

	"github.com/konradmalik/aoc2023/maze"
)

func TestMaze(t *testing.T) {
	var tests = []struct {
		fname        string
		wantMaxDist  int
		wantEnclosed int
	}{
		{"./cases/case1.txt", 4, 1},
		{"./cases/case2.txt", 8, 1},
		{"./cases/case3.txt", 23, 4},
		{"./cases/case4.txt", 22, 4},
		{"./cases/case5.txt", 70, 8},
		{"./cases/case6.txt", 80, 10},
	}

	for _, tt := range tests {
		testname := tt.fname
		t.Run(testname, func(t *testing.T) {
			runFor(t, tt.fname, tt.wantMaxDist, tt.wantEnclosed)
		})
	}
}

func runFor(t *testing.T, fname string, expectedMaxDist int, expectedEnclosed int) {
	m := maze.NewMazeFromFile(fname)
	fmt.Println("Maze")
	fmt.Println(m.String(func(t maze.Tile) string {
		return string(t.Label)
	}))

	start := m.FindStart()
	m.Traverse(start)
	maxDist := m.GetMaxDist()
	if maxDist != expectedMaxDist {
		t.Errorf("distance: got %d, want %d", maxDist, expectedMaxDist)
	}

	fmt.Println("Loop")
	fmt.Println(m.String(func(t maze.Tile) string {
		if t.Visited {
			return "x"
		}
		return "o"
	}))

	m.ProcessEnclosed()
	fmt.Println("Intersections")
	fmt.Println(m.String(func(t maze.Tile) string {
		return fmt.Sprint(t.Intersections)
	}))

	fmt.Println("Enclosed")
	fmt.Println(m.String(func(t maze.Tile) string {
		if t.Intersections%2 != 0 {
			return "I"
		}
		if t.Visited {
			return "X"
		}
		return "O"
	}))

	enclosed := m.GetEnclosed()
	if enclosed != expectedEnclosed {
		t.Errorf("enclosed: got %d, want %d", enclosed, expectedEnclosed)
	}
}
