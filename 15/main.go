package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

type lens struct {
	label string
	op    rune
	focal int
}

func Parse(st string) lens {
	s := []rune(st)
	slen := len(s)
	if s[slen-1] == '-' {
		return lens{string(s[:slen-1]), s[slen-1], 0}
	}

	foc, err := strconv.Atoi(string(s[slen-1]))
	if err != nil {
		panic(err)
	}
	return lens{string(s[:slen-2]), s[slen-2], foc}
}

func hash(s string) int {
	curr := 0
	for _, r := range s {
		curr += int(r)
		curr *= 17
		curr %= 256
	}

	return curr
}

func sameAs(l1 lens) func(lens) bool {
	return func(l lens) bool {
		return l1.label == l.label
	}
}

func insert(boxes [][]lens, l lens) {
	h := hash(l.label)
	ls := boxes[h]
	if len(ls) > 0 {
		// update focal
		if i := slices.IndexFunc(ls, sameAs(l)); i != -1 {
			boxes[h][i].focal = l.focal
		} else {
			ls = append(ls, l)
			boxes[h] = ls
		}
	} else {
		boxes[h] = []lens{l}
	}
}

func remove(boxes [][]lens, l lens) {
	h := hash(l.label)
	ls := boxes[h]
	if len(ls) == 0 {
		return
	}

	if i := slices.IndexFunc(ls, sameAs(l)); i != -1 {
		ls = append(ls[:i], ls[i+1:]...)
		boxes[h] = ls
	}

	if len(ls) == 0 {
		boxes[h] = nil
	}
}

func processFile(filepath string) [][]lens {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	boxes := make([][]lens, 256)
	f := SplitAt(",")
	scanner.Split(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		l := Parse(line)

		if l.op == '-' {
			remove(boxes, l)
		} else {
			insert(boxes, l)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return boxes
}

func SplitAt(substring string) func(data []byte, atEOF bool) (advance int, token []byte, err error) {
	searchBytes := []byte(substring)
	searchLen := len(searchBytes)
	return func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		dataLen := len(data)

		// Return nothing if at end of file and no data passed
		if atEOF && dataLen == 0 {
			return 0, nil, nil
		}

		// Find next separator and return token
		if i := bytes.Index(data, searchBytes); i >= 0 {
			return i + searchLen, data[0:i], nil
		}

		// If we're at EOF, we have a final, non-terminated line. Return it.
		if atEOF {
			return dataLen, data, nil
		}

		// Request more data.
		return 0, nil, nil
	}
}

func main() {
	boxes := processFile("./input.txt")
	fmt.Println(boxes)

	total := 0
	for b, ls := range boxes {
		for i, l := range ls {
			total += (1 + b) * (1 + i) * l.focal
		}
	}

	fmt.Println("the answer", total)
}
