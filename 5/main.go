package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

func parse_numbers(line string) []int {
	words := strings.Fields(line)
	numbers := make([]int, len(words))
	for i, w := range words {
		n, err := strconv.Atoi(w)
		if err != nil {
			panic(err)
		}
		numbers[i] = n
	}
	return numbers
}

func parse_seeds(line string) []int {
	return parse_numbers(strings.Split(line, ":")[1])
}

type mapping struct {
	name    string
	mranges []mrange
}

type mrange struct {
	dest int
	src  int
	len  int
}

func (r *mrange) parse(line string) {
	nums := parse_numbers(line)
	r.dest = nums[0]
	r.src = nums[1]
	r.len = nums[2]
}

func (m *mapping) parse(lines []string) {
	m.name = lines[0]
	m.mranges = make([]mrange, len(lines)-1)
	for i, l := range lines[1:] {
		var r mrange
		r.parse(l)
		m.mranges[i] = r
	}
}

func (m mapping) lookup(v int) int {
	for _, r := range m.mranges {
		var dmin = r.dest
		var smin = r.src
		var smax = r.src + r.len
		if v >= smin && v < smax {
			return (v - smin) + dmin
		}
	}

	return v
}

func main() {
	file, err := os.Open("./input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	line_break := 0
	lines_buf := make([]string, 0, 100)
	var seeds []int
	maps := make([]mapping, 7)

	for scanner.Scan() {
		line := scanner.Text()
		if line_break == 0 {
			if len(line) == 0 {
				line_break++
				continue
			}

			seeds = parse_seeds(line)
			fmt.Println(seeds)
		} else if line_break > 0 {
			if len(line) != 0 {
				lines_buf = append(lines_buf, line)
				continue
			} else {
				var m mapping
				m.parse(lines_buf)
				fmt.Printf("map for line break %d %v\n", line_break, m)
				maps[line_break-1] = m
				line_break++
				lines_buf = nil
			}
		}
	}

	locations := make([]int, len(seeds))
	for i, seed := range seeds {
		// fmt.Println("for seed:", seed)
		value := seed
		for _, m := range maps {
			value = m.lookup(value)
			// fmt.Printf("for %s value is %d\n", m.name, value)
		}
		locations[i] = value
	}

	fmt.Println("locations", locations)
	fmt.Println("min", slices.Min(locations))

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
