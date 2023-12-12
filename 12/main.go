package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/konradmalik/aoc2023/lib"
)

type Entry struct {
	history    string
	validation []int
}

func (e Entry) IsValid() bool {
	if strings.Contains(e.history, "?") {
		panic("questionmark should not be validated")
	}

	i := 0
	vidx := 0
	for i < len(e.history) {
		if e.history[i] == '.' {
			i++
			continue
		}

		// hash after we met all validations is invalid
		if vidx >= len(e.validation) {
			return false
		}

		phist := e.history[i:]
		validation := e.validation[vidx]
		if !HasConsecutiveExactly(phist, validation) {
			return false
		}

		i += validation
		vidx++
	}

	// not all found, it's invalid
	if vidx != len(e.validation) {
		return false
	}

	return true
}

func (e Entry) FindQMarks() []int {
	q := make([]int, 0)
	for i, r := range e.history {
		if r == '?' {
			q = append(q, i)
		}
	}
	return q
}

func GenerateFrom(e Entry) []Entry {
	required := 0
	for _, v := range e.validation {
		required += v
	}
	current := strings.Count(e.history, "#")
	left := required - current
	if left < 0 {
		panic("left is negative")
	}

	qmarks := e.FindQMarks()
	if left > len(qmarks) {
		panic("more left than qmarks")
	}

	generated := make([]Entry, 0)
	products := Product(len(qmarks))
	for _, p := range products {
		ne := e
		for pi, qi := range qmarks {
			ne.history = replaceAtIndex(ne.history, rune(p[pi]), qi)
		}
		if ne.IsValid() {
			generated = append(generated, ne)
		}
	}

	return generated
}

func Product(slots int) []string {
	res := make(map[string]bool)
	current := make([]rune, 0)
	res = product(slots, current, res)

	words := make([]string, 0)
	for w := range res {
		words = append(words, w)
	}
	return words
}

func product(slots int, current []rune, res map[string]bool) map[string]bool {
	runes := []rune{'.', '#'}

	if slots == 0 {
		res[string(current)] = true
		return res
	}

	for _, r := range runes {
		next := append(current, r)
		res = product(slots-1, next, res)
	}

	return res
}

func HasConsecutiveExactly(s string, n int) bool {
	count := 0
	for _, r := range s {
		if r == '.' {
			if count == n {
				return true
			}
			if count > 0 {
				return false
			}
		} else {
			count++
		}
	}

	return count == n
}

func parseEntry(line string) Entry {
	elms := strings.Split(line, " ")
	return Entry{elms[0], lib.ParseNumbers(strings.ReplaceAll(elms[1], ",", " "))}
}

func TotalArrangementsFromFile(filepath string) int {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	total := 0
	for scanner.Scan() {
		line := scanner.Text()
		entry := parseEntry(line)
		generated := GenerateFrom(entry)
		fmt.Println(generated)
		total += len(generated)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return total
}

func replaceAtIndex(in string, r rune, i int) string {
	out := []rune(in)
	out[i] = r
	return string(out)
}

func main() {
	total := TotalArrangementsFromFile("./input.txt")
	fmt.Print(total)
}
