package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
)

type position struct {
	x int
	y int
}

type tile struct {
	label rune
	dist  int
}

func (t tile) String() string {
	return fmt.Sprintf("{%s, %d}", string(t.label), t.dist)
}

func (t tile) Visited() bool {
	return t.dist > 0
}

func (t tile) ConnectedWithSouth() bool {
	return slices.Contains(
		[]rune{'|', '7', 'F'},
		t.label)
}

func (t tile) ConnectedWithNorth() bool {
	return slices.Contains(
		[]rune{'|', 'L', 'J'},
		t.label)
}

func (t tile) ConnectedWithWest() bool {
	return slices.Contains(
		[]rune{'-', '7', 'J'},
		t.label)
}

func (t tile) ConnectedWithEast() bool {
	return slices.Contains(
		[]rune{'-', 'L', 'F'},
		t.label)
}

// list of rows/ys
type maze struct {
	data [][]*tile
}

func (m maze) String() string {
	var sb strings.Builder

	for _, row := range m.data {
		for _, t := range row {
			sb.WriteRune(t.label)
		}
		sb.WriteRune('\n')
	}

	for _, row := range m.data {
		for _, t := range row {
			sb.WriteString(fmt.Sprint(t.dist))
		}
		sb.WriteRune('\n')
	}

	return sb.String()
}

func (m *maze) ParseInc(line string) {
	row := make([]*tile, len(line))
	for i, r := range []rune(line) {
		row[i] = &tile{label: r, dist: 0}
	}
	m.data = append(m.data, row)
}

func (m maze) GetMaxDist() int {
	res := 0
	for _, row := range m.data {
		for _, t := range row {
			if t.dist > res {
				res = t.dist
			}
		}
	}
	return res
}

func (m maze) findStart() position {
	for y, row := range m.data {
		for x, t := range row {
			if t.label == 'S' {
				return position{x, y}
			}
		}
	}
	panic("no start")
}

func (m maze) get(p position) *tile {
	return m.data[p.y][p.x]
}

func (m maze) MaxX() int {
	return len(m.data[0]) - 1
}

func (m maze) MaxY() int {
	return len(m.data) - 1
}

func (m *maze) StartTraversing(p position) {
	q := queue{[]nextPosition{{p, 0}}}
	m.traverse(q)
}

// use queue to get BFS
func (m *maze) traverse(q queue) {
	if q.IsEmpty() {
		return
	}

	cp := q.Dequeue()
	p := cp.p
	ct := m.get(p)

	// update distance
	ct.dist = cp.dist

	np := position{p.x, max(0, p.y-1)}
	sp := position{p.x, min(m.MaxY(), p.y+1)}
	wp := position{max(0, p.x-1), p.y}
	ep := position{min(m.MaxX(), p.x+1), p.y}

	nt := m.get(np)
	if !nt.Visited() && nt.ConnectedWithSouth() {
		q.Enqueue(nextPosition{np, ct.dist + 1})
	}

	st := m.get(sp)
	if !st.Visited() && st.ConnectedWithNorth() {
		q.Enqueue(nextPosition{sp, ct.dist + 1})
	}

	wt := m.get(wp)
	if !wt.Visited() && wt.ConnectedWithEast() {
		q.Enqueue(nextPosition{wp, ct.dist + 1})
	}

	et := m.get(ep)
	if !et.Visited() && et.ConnectedWithWest() {
		q.Enqueue(nextPosition{ep, ct.dist + 1})
	}

	m.traverse(q)
}

type nextPosition struct {
	p    position
	dist int
}

type queue struct {
	data []nextPosition
}

func (q *queue) Enqueue(p nextPosition) {
	q.data = append(q.data, p)
}

func (q *queue) Dequeue() nextPosition {
	e := q.data[0]
	q.data = q.data[1:]
	return e
}

func (q queue) IsEmpty() bool {
	return len(q.data) == 0
}

func processFile(filepath string) maze {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	maze := maze{make([][]*tile, 0)}
	for scanner.Scan() {
		line := scanner.Text()
		maze.ParseInc(line)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return maze
}

func main() {
	maze := processFile("./input.txt")

	start := maze.findStart()
	maze.StartTraversing(start)
	fmt.Println(maze)

	fmt.Println(maze.GetMaxDist())
}
