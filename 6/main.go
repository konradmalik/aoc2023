package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strings"

	"github.com/konradmalik/aoc2023/lib"
)

func CeilNe(n float64) int {
	nc := math.Ceil(n)
	if nc == n {
		return int(nc + 1)
	}
	return int(nc)
}

func FloorNe(n float64) int {
	nc := math.Floor(n)
	if nc == n {
		return int(nc - 1)
	}
	return int(nc)
}

func main() {
	file, err := os.Open("./input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	nline := 0
	times := make([]int, 0, 4)
	distances := make([]int, 0, 4)
	for scanner.Scan() {
		line := scanner.Text()
		if nline == 0 {
			times = lib.ParseNumbers(strings.Split(line, ":")[1])
		}
		if nline == 1 {
			distances = lib.ParseNumbers(strings.Split(line, ":")[1])
		}
		nline++
	}

	res := 1
	for i, t := range times {
		dmin := distances[i]
		fmt.Println("iteration:", i+1)
		fmt.Printf("t:%d\n", t)
		fmt.Printf("dmin:%d\n", dmin)
		// -tv^2+t*tv -dmin > 0
		// disc = t^2-4dmin
		// tv1 = 1/2*(t-sqrt(disc))
		// tv2 = 1/2*(t+sqrt(disc))
		tf := float64(t)
		dminf := float64(dmin)
		disc_sqrt := math.Sqrt(tf*tf - 4*dminf)
		tv1 := 0.5 * (tf - disc_sqrt)
		tv2 := 0.5 * (tf + disc_sqrt)
		fmt.Printf("tv1: %f, tv2: %f\n", tv1, tv2)
		tv1i := CeilNe(tv1)
		tv2i := FloorNe(tv2)
		fmt.Printf("tv1i: %d, tv2i: %d\n", tv1i, tv2i)
		ways := tv2i - tv1i + 1
		fmt.Println(ways)
		fmt.Println()
		res *= ways
	}

	fmt.Println("result:", res)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
