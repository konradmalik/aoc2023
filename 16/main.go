package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
)

type direction int

const (
	N direction = iota
	W
	S
	E
)

type position struct {
	x int
	y int
	d direction
}

type grid struct {
	tiles [][]tile
	moves map[position]bool
}

func (g grid) CountEnergized() int {
	count := 0
	for _, row := range g.tiles {
		for _, t := range row {
			if t.energized {
				count++
			}
		}
	}
	return count
}

func (g grid) IsOutOfBounds(p position) bool {
	if p.x < 0 || p.y < 0 {
		return true
	}
	if p.y >= len(g.tiles) {
		return true
	}
	if p.x >= len(g.tiles[0]) {
		return true
	}
	return false
}

type tile struct {
	value     rune
	energized bool
}

func (g *grid) Move(p position) {
	if g.IsOutOfBounds(p) {
		return
	}
	tile := g.tiles[p.y][p.x]
	if _, found := g.moves[p]; found {
		return
	}
	g.moves[p] = true
	g.tiles[p.y][p.x].energized = true

	switch p.d {
	case N:
		switch tile.value {
		case '.', '|':
			g.Move(position{x: p.x, y: p.y - 1, d: N})
		case '\\':
			g.Move(position{x: p.x - 1, y: p.y, d: W})
		case '/':
			g.Move(position{x: p.x + 1, y: p.y, d: E})
		case '-':
			g.Move(position{x: p.x - 1, y: p.y, d: W})
			g.Move(position{x: p.x + 1, y: p.y, d: E})
		}
	case W:
		switch tile.value {
		case '.', '-':
			g.Move(position{x: p.x - 1, y: p.y, d: W})
		case '\\':
			g.Move(position{x: p.x, y: p.y - 1, d: N})
		case '/':
			g.Move(position{x: p.x, y: p.y + 1, d: S})
		case '|':
			g.Move(position{x: p.x, y: p.y - 1, d: N})
			g.Move(position{x: p.x, y: p.y + 1, d: S})
		}
	case S:
		switch tile.value {
		case '.', '|':
			g.Move(position{x: p.x, y: p.y + 1, d: S})
		case '\\':
			g.Move(position{x: p.x + 1, y: p.y, d: E})
		case '/':
			g.Move(position{x: p.x - 1, y: p.y, d: W})
		case '-':
			g.Move(position{x: p.x - 1, y: p.y, d: W})
			g.Move(position{x: p.x + 1, y: p.y, d: E})
		}
	case E:
		switch tile.value {
		case '.', '-':
			g.Move(position{x: p.x + 1, y: p.y, d: E})
		case '\\':
			g.Move(position{x: p.x, y: p.y + 1, d: S})
		case '/':
			g.Move(position{x: p.x, y: p.y - 1, d: N})
		case '|':
			g.Move(position{x: p.x, y: p.y - 1, d: N})
			g.Move(position{x: p.x, y: p.y + 1, d: S})
		}
	}
}

func (g grid) String(tileToString func(tile) string) string {
	var sb strings.Builder

	for _, row := range g.tiles {
		for _, t := range row {
			sb.WriteString(tileToString(t))
		}
		sb.WriteRune('\n')
	}
	return sb.String()
}

func parseLine(l string) []tile {
	tiles := make([]tile, len(l))
	for i, r := range l {
		tiles[i] = tile{value: r, energized: false}
	}
	return tiles
}

func parseGrid(filepath string) grid {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	tiles := make([][]tile, 0)
	for scanner.Scan() {
		line := scanner.Text()
		row := parseLine(line)
		tiles = append(tiles, row)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return grid{tiles: tiles, moves: make(map[position]bool)}
}

func getPossibleStarts(g grid) []position {
	rs := make([]position, 0)
	xmax := len(g.tiles[0]) - 1
	ymax := len(g.tiles) - 1
	for i := 0; i <= ymax; i++ {
		rs = append(rs, position{x: 0, y: i, d: E})
		rs = append(rs, position{x: xmax, y: i, d: W})
	}
	for i := 0; i <= xmax; i++ {
		rs = append(rs, position{x: i, y: 0, d: S})
		rs = append(rs, position{x: i, y: ymax, d: N})
	}
	return rs
}

func main() {
	grid := parseGrid("./input.txt")
	fmt.Println(grid.String(func(t tile) string { return string(t.value) }))
	possibleStarts := getPossibleStarts(grid)
	counts := make([]int, len(possibleStarts))
	for i, p := range possibleStarts {
		grid := parseGrid("./input.txt")
		grid.Move(p)
		counts[i] = grid.CountEnergized()
	}

	// fmt.Println(grid.String(func(t tile) string {
	// 	if t.energized {
	// 		return "#"
	// 	}
	// 	return "."
	// }))
	fmt.Println(slices.Max(counts))
}
