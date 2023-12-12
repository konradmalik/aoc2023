package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"

	"github.com/konradmalik/aoc2023/lib"
)

type Entry struct {
	history    string
	validation []int
}

func cacheKey(history []rune, validations []int) string {
	return string(history) + strings.Trim(strings.Join(strings.Fields(fmt.Sprint(validations)), ","), "[]")
}

var cache = make(map[string]uint64)

func Arrangements(history []rune, validations []int) uint64 {
	ck := cacheKey(history, validations)
	cached, found := cache[ck]
	if found {
		return cached
	}

	if len(history) == 0 {
		if len(validations) == 0 {
			cache[ck] = 1
			return 1
		}
		cache[ck] = 0
		return 0
	}

	if history[0] == '.' {
		res := Arrangements(history[1:], validations)
		cache[ck] = res
		return res
	}

	if history[0] == '?' {
		res := Arrangements(append([]rune{'.'}, history[1:]...), validations) +
			Arrangements(append([]rune{'#'}, history[1:]...), validations)
		cache[ck] = res
		return res
	}

	if history[0] == '#' {
		if len(validations) == 0 {
			cache[ck] = 0
			return 0
		}
		if len(history) < validations[0] {
			cache[ck] = 0
			return 0
		}
		if slices.Contains(history[:validations[0]], '.') {
			cache[ck] = 0
			return 0
		}
		if len(validations) > 1 {
			if len(history) < validations[0]+1 || history[validations[0]] == '#' {
				cache[ck] = 0
				return 0
			}
			res := Arrangements(history[validations[0]+1:], validations[1:])
			cache[ck] = res
			return res
		}
		res := Arrangements(history[validations[0]:], validations[1:])
		cache[ck] = res
		return res
	}

	panic("should not be possible")
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

func TotalArrangementsFromFile(filepath string) uint64 {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	ln := 0
	var total uint64 = 0
	for scanner.Scan() {
		log.Println("line", ln)
		line := scanner.Text()
		entry := parseEntry(line)
		total += Arrangements([]rune(entry.history), entry.validation)
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
