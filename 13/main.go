package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
)

type Pattern []string

func (p Pattern) String() string {
	var sb strings.Builder
	for _, r := range p {
		sb.WriteString(r)
		sb.WriteRune('\n')
	}
	return sb.String()
}

func (p Pattern) Rows() []string {
	rows := make([]string, len(p))
	copy(rows, p)
	return rows
}

func (p Pattern) Cols() []string {
	cols := make([]string, len(p[0]))
	for i := range p[0] {
		col := make([]rune, len(p))
		for j, row := range p {
			col[j] = []rune(row)[i]
		}
		cols[i] = string(col)
	}
	return cols
}

func Mirrors(s string) []int {
	ms := make([]int, 0)
	for i := 1; i < len(s); i++ {
		left, right := s[:i], s[i:]
		if prefixReversed(right, left) {
			ms = append(ms, i)
		}
	}
	return ms
}

func prefixReversed(a string, b string) bool {
	br := []rune(strings.Clone(b))
	slices.Reverse(br)

	if len(br) <= len(a) {
		return strings.HasPrefix(a, string(br))
	}
	return strings.HasPrefix(string(br), a)
}

func processPattern(pattern Pattern) int {
	fmt.Println("new pattern")

	for {
		count := 0
		colmirror := mirrorForDimension(pattern.Rows())
		fmt.Println("colmirror", colmirror)
		if colmirror >= 0 {
			count += colmirror
		}
		rowmirror := mirrorForDimension(pattern.Cols())
		fmt.Println("rowmirror", rowmirror)
		if rowmirror >= 0 {
			count += 100 * rowmirror
		}

		if count > 0 {
			return count
		}
	}
}

func mirrorForDimension(rows []string) int {
	// index -> how many times
	counts := make(map[int]int)
	for _, row := range rows {
		ms := Mirrors(row)
		for _, m := range ms {
			if _, found := counts[m]; !found {
				counts[m] = 1
			} else {
				counts[m] += 1
			}
		}
	}

	for idx, count := range counts {
		if count == len(rows) {
			return idx
		}
	}
	return -1
}

func processPatterns(filepath string) int {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	i := 0
	total := 0
	pattern := make([]string, 0)

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			if len(pattern) > 0 {
				total += processPattern(pattern)
			}
			pattern = make([]string, 0)
		} else {
			pattern = append(pattern, line)
		}
		i++
	}

	if len(pattern) > 0 {
		total += processPattern(pattern)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return total
}

func main() {
	res := processPatterns("./input.txt")
	fmt.Println(res)
}
