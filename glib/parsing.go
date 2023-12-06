package lib

import (
	"strconv"
	"strings"
)

func ParseNumbers(line string) []int {
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
