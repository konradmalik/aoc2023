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

func main() {
	input := parseFile("./input.txt")
	moves := []rune(input.moves)
	nodes := input.nodes

	node := nodes["AAA"]
	move_idx := 0
	nmoves := len(input.moves)
	total_moves := 0
	for node.name != "ZZZ" {
		move := moves[move_idx]
		log.Printf("node: %s, moving: %s", node.name, string(move))
		if move == 'L' {
			node = nodes[node.neighbors[0].name]
		} else {
			node = nodes[node.neighbors[1].name]
		}
		move_idx = (move_idx + 1) % nmoves
		total_moves++
	}
	log.Println(total_moves)
}
