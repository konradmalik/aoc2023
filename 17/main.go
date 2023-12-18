package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
)

type dimension int

const (
	NONE dimension = iota
	HORIZONTAL
	VERTICAL
)

type grid [][]int

func parseLine(l string) []int {
	row := make([]int, len(l))
	for i, r := range l {
		cost, err := strconv.Atoi(string(r))
		if err != nil {
			panic(err)
		}
		row[i] = cost
	}
	return row
}

func parseGrid(filepath string) grid {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	grid := make([][]int, 0)
	for scanner.Scan() {
		line := scanner.Text()
		row := parseLine(line)
		grid = append(grid, row)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return grid
}

type Node struct {
	x int
	y int
	d dimension
}

func (n Node) String() string {
	return fmt.Sprintf("(%d,%d)", n.x, n.y)
}

type Edge struct {
	source Node
	target Node
	cost   int
}

type Graph struct {
	nodes       []Node
	nodeToEdges map[Node][]Edge
}

func (g grid) GetEdgeTo(n Node, x, y int, d dimension) (Edge, error) {
	if x < 0 || y < 0 {
		return Edge{}, fmt.Errorf("less than zero")
	}

	if x >= len(g[0]) || y >= len(g) {
		return Edge{}, fmt.Errorf("more than grid")
	}

	cost := 0
	if n.x < x {
		for i := n.x + 1; i <= x; i++ {
			cost += g[y][i]
		}
	} else if n.x > x {
		for i := n.x - 1; i >= x; i-- {
			cost += g[y][i]
		}
	} else if n.y < y {
		for i := n.y + 1; i <= y; i++ {
			cost += g[i][x]
		}
	} else if n.y > y {
		for i := n.y - 1; i >= y; i-- {
			cost += g[i][x]
		}
	} else {
		panic("cost problem")
	}

	return Edge{n, Node{x, y, d}, cost}, nil
}

func (g grid) GetEdges(n Node, minStep, maxStep int) []Edge {
	edges := make([]Edge, 0)
	if n.d == VERTICAL || n.d == NONE {
		for i := minStep; i <= maxStep; i++ {
			e, err := g.GetEdgeTo(n, n.x+i, n.y, HORIZONTAL)
			if err == nil {
				edges = append(edges, e)
			}
		}
		for i := minStep; i <= maxStep; i++ {
			e, err := g.GetEdgeTo(n, n.x-i, n.y, HORIZONTAL)
			if err == nil {
				edges = append(edges, e)
			}
		}
	}

	if n.d == HORIZONTAL || n.d == NONE {
		for i := minStep; i <= maxStep; i++ {
			e, err := g.GetEdgeTo(n, n.x, n.y+i, VERTICAL)
			if err == nil {
				edges = append(edges, e)
			}
		}
		for i := minStep; i <= maxStep; i++ {
			e, err := g.GetEdgeTo(n, n.x, n.y-i, VERTICAL)
			if err == nil {
				edges = append(edges, e)
			}
		}
	}

	return edges
}

func (g grid) toGraph() Graph {
	nodes := make([]Node, 0)
	nodeToEdges := make(map[Node][]Edge)
	for y, row := range g {
		for x := range row {
			dims := []dimension{HORIZONTAL, VERTICAL}
			if x == 0 && y == 0 {
				dims = append(dims, NONE)
			}

			for _, d := range dims {
				cnode := Node{x: x, y: y, d: d}
				nodes = append(nodes, cnode)

				cedges := g.GetEdges(cnode, 4, 10)
				for _, edge := range cedges {
					n := edge.source
					saved, found := nodeToEdges[n]
					if found {
						nodeToEdges[n] = append(saved, edge)
					} else {
						nodeToEdges[n] = []Edge{edge}
					}
				}
			}
		}
	}

	return Graph{nodes: nodes, nodeToEdges: nodeToEdges}
}

func (g Graph) Dijkstra(start Node, end Node) int {
	nodes := g.nodes

	dist := make(map[Node]int)
	visited := make(map[Node]bool)

	for _, n := range nodes {
		dist[n] = math.MaxInt32
		visited[n] = false
	}

	dist[start] = 0

	// TODO use a priority queue
	next := func() Node {
		lowest := Node{-1, -1, NONE}
		for _, n := range nodes {
			if !visited[n] && (lowest.x == -1 || dist[n] < dist[lowest]) {
				lowest = n
			}

		}
		return lowest
	}

	for {
		u := next()

		if u.x == end.x && u.y == end.y {
			return dist[u]
		}

		if u.x == -1 {
			return math.MaxInt32
		}

		visited[u] = true

		for _, edge := range g.nodeToEdges[u] {
			v := edge.target
			if !visited[v] {
				if alt := dist[u] + edge.cost; alt < dist[v] {
					dist[v] = alt
				}
			}
		}
	}
}

func main() {
	grid := parseGrid("./input.txt")

	graph := grid.toGraph()

	start := Node{x: 0, y: 0, d: NONE}
	end := Node{x: len(grid[0]) - 1, y: len(grid) - 1, d: VERTICAL}

	dist := graph.Dijkstra(start, end)

	fmt.Printf("Shortest path from node %d to %d\n", start, end)
	fmt.Println("dist", dist)
}
