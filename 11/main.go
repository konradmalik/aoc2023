package main

import (
	"fmt"
	"github.com/konradmalik/aoc2023/universe"
)

func main() {
	u := universe.NewUniverseFromFile("./input.txt")
	fmt.Println("Original Galaxies:")
	galaxies := u.Galaxies()
	fmt.Println(galaxies)
	u.Expand(2)

	fmt.Println("Expanded Galaxies:")
	galaxies = u.Galaxies()
	fmt.Println(galaxies)

	fmt.Println("Combinations:")
	combinations := universe.Combinations(galaxies)
	fmt.Println(len(combinations))

	fmt.Println("Distances:")
	distances := universe.Distances(combinations)

	sum := 0
	for _, d := range distances {
		sum += d
	}
	fmt.Println(sum)
}
