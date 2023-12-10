package main

import (
	"fmt"

	"github.com/konradmalik/aoc2023/maze"
)

func main() {
	m := maze.NewMazeFromFile("./input.txt")

	start := m.FindStart()
	m.Traverse(start)
	m.ProcessEnclosed()
	fmt.Println(m.String(func(t maze.Tile) string {
		if t.Intersections%2 != 0 {
			return "I"
		}
		if t.Visited {
			return "*"
		}
		return "."
	}))
	fmt.Println("enclosed:", m.GetEnclosed())
}
