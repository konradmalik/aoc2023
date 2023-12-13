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
	return p
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

func replaceAtIndex(in string, r rune, i int) string {
	out := []rune(in)
	out[i] = r
	return string(out)
}

func generateSmudges(p Pattern) <-chan Pattern {
	smudges := make(chan Pattern)
	go func() {
		defer close(smudges)

		for i := 0; i < len(p); i++ {
			for j := 0; j < len(p[0]); j++ {
				np := make(Pattern, len(p))
				copy(np, p)
				if np[i][j] == '.' {
					np[i] = replaceAtIndex(np[i], '#', j)
				} else {
					np[i] = replaceAtIndex(np[i], '.', j)
				}
				smudges <- np
			}
		}
	}()
	return smudges
}

func processPattern(pattern Pattern) (int, int) {
	colmirror := mirrorForDimension(pattern.Rows(), -1)
	rowmirror := mirrorForDimension(pattern.Cols(), -1)
	return colmirror, rowmirror
}

func processSmudges(pattern Pattern) int {
	orig_col, orig_row := processPattern(pattern)

	for smudge := range generateSmudges(pattern) {
		colmirror := mirrorForDimension(smudge.Rows(), orig_col)
		if colmirror > 0 {
			return colmirror
		}
		rowmirror := mirrorForDimension(smudge.Cols(), orig_row)
		if rowmirror > 0 {
			return 100 * rowmirror
		}
	}
	fmt.Println(pattern)
	panic("no new mirror after smudge")
}

func mirrorForDimension(rows []string, exclude int) int {
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
		if idx == exclude {
			continue
		}
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
			fmt.Println("pattern", i)
			if len(pattern) > 0 {
				total += processSmudges(pattern)
			}
			i++
			pattern = make([]string, 0)
		} else {
			pattern = append(pattern, line)
		}
	}

	if len(pattern) > 0 {
		total += processSmudges(pattern)
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
