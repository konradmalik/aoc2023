package maze

import (
	"strings"

	"github.com/konradmalik/aoc2023/lib"
)

func (m *Maze) raycastPosition(p Position) {
	ct := m.mustGet(p)
	row := m.data[p.Y]

	var sb strings.Builder
	for _, t := range row[p.X+1:] {
		if !t.IsLoop() {
			continue
		}
		if t.Label == '-' {
			continue
		}
		sb.WriteRune(t.Label)
	}
	rowstring := sb.String()
	rowstring = strings.ReplaceAll(rowstring, "FJ", "|")
	rowstring = strings.ReplaceAll(rowstring, "L7", "|")
	ct.Intersections = strings.Count(rowstring, "|")
}

func (m Maze) findNonLoops() lib.Queue[Position] {
	q := lib.NewQueue[Position]()
	for y, row := range m.data {
		for x, t := range row {
			if !t.IsLoop() {
				q.Enqueue(Position{x, y})
			}
		}
	}
	return q
}

func (m *Maze) ProcessEnclosed() {
	q := m.findNonLoops()
	for !q.IsEmpty() {
		p := q.Dequeue()
		m.raycastPosition(p)
	}
}

func (m Maze) GetEnclosed() int {
	res := 0
	for _, row := range m.data {
		for _, t := range row {
			if t.Intersections%2 != 0 {
				res += 1
			}
		}
	}
	return res
}
