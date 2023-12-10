package maze

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/konradmalik/aoc2023/lib"
)

func NewMazeFromFile(filepath string) Maze {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	maze := NewMaze()
	for scanner.Scan() {
		line := scanner.Text()
		maze.ParseInc(line)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return maze
}

// list of rows/ys
type Maze struct {
	data [][]*Tile
}

func NewMaze() Maze {
	return Maze{make([][]*Tile, 0)}
}

func (m Maze) String(tileToString func(Tile) string) string {
	var sb strings.Builder

	for _, row := range m.data {
		for _, t := range row {
			sb.WriteString(tileToString(*t))
		}
		sb.WriteRune('\n')
	}
	return sb.String()
}

func (m *Maze) ParseInc(line string) {
	row := make([]*Tile, len(line))
	for i, r := range []rune(line) {
		row[i] = &Tile{Label: r, Dist: 0}
	}
	m.data = append(m.data, row)
}

func (m Maze) GetMaxDist() int {
	res := 0
	for _, row := range m.data {
		for _, t := range row {
			if t.Dist > res {
				res = t.Dist
			}
		}
	}
	return res
}

func (m Maze) FindStart() Position {
	for y, row := range m.data {
		for x, t := range row {
			if t.Label == 'S' {
				return Position{x, y}
			}
		}
	}
	panic("no start")
}

func (m Maze) mustGet(p Position) *Tile {
	return m.data[p.Y][p.X]
}

func (m Maze) get(p Position) (*Tile, error) {
	if p.X < 0 || p.Y < 0 {
		return nil, fmt.Errorf("bad index")
	}

	if p.X > m.maxX() || p.Y > m.maxY() {
		return nil, fmt.Errorf("bad index")
	}

	return m.mustGet(p), nil
}

func (m Maze) maxX() int {
	return len(m.data[0]) - 1
}

func (m Maze) maxY() int {
	return len(m.data) - 1
}

func (m *Maze) Traverse(p Position) {
	q := lib.NewQueue[nextPosition]()
	q.Enqueue(nextPosition{p, 0})
	m.traverse(q)
}

// use queue to get BFS
func (m *Maze) traverse(q lib.Queue[nextPosition]) {
	if q.IsEmpty() {
		return
	}

	cp := q.Dequeue()
	p := cp.Pos
	ct := m.mustGet(p)

	ct.Visited = true
	ct.Dist = cp.Dist

	np := Position{p.X, max(0, p.Y-1)}
	sp := Position{p.X, min(m.maxY(), p.Y+1)}
	wp := Position{max(0, p.X-1), p.Y}
	ep := Position{min(m.maxX(), p.X+1), p.Y}

	nc := false
	sc := false
	wc := false
	ec := false

	nt := m.mustGet(np)
	if !nt.Visited && ct.ConnectedWithNorth() && nt.ConnectedWithSouth() {
		nc = true
		q.Enqueue(nextPosition{np, ct.Dist + 1})
	}

	st := m.mustGet(sp)
	if !st.Visited && ct.ConnectedWithSouth() && st.ConnectedWithNorth() {
		sc = true
		q.Enqueue(nextPosition{sp, ct.Dist + 1})
	}

	wt := m.mustGet(wp)
	if !wt.Visited && ct.ConnectedWithWest() && wt.ConnectedWithEast() {
		wc = true
		q.Enqueue(nextPosition{wp, ct.Dist + 1})
	}

	et := m.mustGet(ep)
	if !et.Visited && ct.ConnectedWithEast() && et.ConnectedWithWest() {
		ec = true
		q.Enqueue(nextPosition{ep, ct.Dist + 1})
	}

	if ct.Label == 'S' {
		if nc && sc {
			ct.Label = '|'
		} else if nc && ec {
			ct.Label = 'L'
		} else if nc && wc {
			ct.Label = 'J'
		} else if ec && wc {
			ct.Label = '-'
		} else if sc && ec {
			ct.Label = 'F'
		} else if sc && wc {
			ct.Label = '7'
		} else {
			panic("unhandled start replacing case")
		}
	}

	m.traverse(q)
}
