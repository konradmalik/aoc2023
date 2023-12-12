package old

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/konradmalik/aoc2023/lib"
)

type Entry struct {
	history    string
	validation []int
}

func (e Entry) IsValid() bool {
	i := 0
	vidx := 0
	for i < len(e.history) {
		if e.history[i] == '?' {
			panic("questionmark should not be validated")
		}

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

func GenerateFrom(e Entry) <-chan Entry {
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

	products := Product(len(qmarks))
	generated := make(chan Entry)
	var wg sync.WaitGroup
	for p := range products {
		wg.Add(1)
		p := p
		go func() {
			defer wg.Done()
			ne := e
			for pi, qi := range qmarks {
				ne.history = replaceAtIndex(ne.history, rune(p[pi]), qi)
			}
			if ne.IsValid() {
				generated <- ne
			}
		}()
	}
	go func() {
		wg.Wait()
		close(generated)
	}()

	return generated
}

func Product(slots int) <-chan string {
	c := make(chan string)

	go func(c chan string) {
		defer close(c)

		product(c, []rune{}, slots)
	}(c)

	return c
}

var runes = []rune{'.', '#'}

func product(c chan string, current []rune, slots int) {
	if slots <= 0 {
		c <- string(current)
		return
	}

	for _, r := range runes {
		next := append(current, r)
		product(c, next, slots-1)
	}
}

func HasConsecutiveExactly(s string, n int) bool {
	count := 0
	for _, r := range s {
		if r == '#' {
			count++
		} else {
			if count == n {
				return true
			}
			if count > 0 {
				return false
			}
		}
	}

	return count == n
}

func parseEntry(line string) Entry {
	elms := strings.Split(line, " ")

	mut := 5
	return Entry{
		unfoldHistory(elms[0], mut),
		unfoldValidations(lib.ParseNumbers(strings.ReplaceAll(elms[1], ",", " ")), mut),
	}
}

func unfoldHistory(history string, times int) string {
	nhistory := history
	for i := 0; i < times-1; i++ {
		nhistory = fmt.Sprint(nhistory, "?", history)
	}
	return nhistory
}

func unfoldValidations(validations []int, times int) []int {
	nvalidations := validations
	for i := 0; i < times-1; i++ {
		nvalidations = append(nvalidations, validations...)
	}
	return nvalidations
}

func TotalArrangementsFromFile(filepath string) int {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	ln := 0
	total := 0
	for scanner.Scan() {
		log.Println("line", ln)
		line := scanner.Text()
		entry := parseEntry(line)
		generated := GenerateFrom(entry)
		for range generated {
			total++
		}
		ln++
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
