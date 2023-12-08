package main

import (
	"bufio"
	"log"
	"os"
	"strings"
)

type Input struct {
	moves string
	nodes map[string]Node
}

func gatherNodes(mnodes map[string]Node) {
	for _, n := range mnodes {
		for i, nb := range n.neighbors {
			nbi := mnodes[nb.name]
			n.neighbors[i] = nbi
		}
	}
}

type Node struct {
	name      string
	neighbors []Node
}

func (n *Node) isStart() bool {
	return strings.HasSuffix(n.name, "A")
}

func (n *Node) isEnd() bool {
	return strings.HasSuffix(n.name, "Z")
}

func (n *Node) parse(line string) {
	elems := strings.Split(line, "=")
	n.name = strings.TrimSpace(elems[0])
	nbs := strings.TrimSpace(elems[1])
	nbs = strings.ReplaceAll(nbs, "(", "")
	nbs = strings.ReplaceAll(nbs, ")", "")

	for _, nb := range strings.Split(nbs, ",") {
		n.neighbors = append(n.neighbors, Node{strings.TrimSpace(nb), make([]Node, 0)})
	}
}

func parseFile(filepath string) Input {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	lnum := 0
	var moves string
	nodes := make(map[string]Node)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}
		if lnum == 0 {
			moves = line
		} else {
			var node Node
			node.parse(line)
			nodes[node.name] = node
		}

		lnum++
	}
	gatherNodes(nodes)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return Input{
		moves: moves,
		nodes: nodes,
	}
}

func getPositions(nodes map[string]Node) []Node {
	positions := make([]Node, 0)
	for _, n := range nodes {
		if n.isStart() {
			positions = append(positions, n)
		}
	}
	return positions
}

func areAllEnds(positions []Node) bool {
	for _, p := range positions {
		if !p.isEnd() {
			return false
		}
	}
	return true
}

// greatest common divisor (GCD) via Euclidean algorithm
func GCD(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

// find Least Common Multiple (LCM) via GCD
func LCM(a, b int, integers ...int) int {
	result := a * b / GCD(a, b)

	for i := 0; i < len(integers); i++ {
		result = LCM(result, integers[i])
	}

	return result
}

func main() {
	input := parseFile("./input.txt")
	moves := []rune(input.moves)
	nodes := input.nodes

	positions := getPositions(nodes)
	// remembers full path starting at key and ending at value
	// to not pass it again
	memo := make(map[string]string)

	iterations := 0

	// position idx to interations until end
	endings := make([]int, len(positions))
	for i := range endings {
		endings[i] = 0
	}

	for !areAllEnds(positions) {
		iterations++

		for i, p := range positions {
			node := p
			i := i
			if endings[i] > 0 {
				continue
			}

			// ..and i've seen it all befoooo..
			if end, found := memo[node.name]; found {
				positions[i] = nodes[end]
				continue
			}

			start := p
			for _, move := range moves {
				// log.Printf("node: %s, moving: %s", node.name, string(move))
				if move == 'L' {
					node = nodes[node.neighbors[0].name]
				} else {
					node = nodes[node.neighbors[1].name]
				}
				positions[i] = node
			}
			end := node
			if end.isEnd() {
				log.Println("end for", end.name)
				endings[i] = iterations * len(moves)
			}
			memo[start.name] = end.name
		}
	}

	log.Println(endings)
	total_moves := LCM(endings[0], endings[1], endings[2:]...)
	log.Println(total_moves)
}
