package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/konradmalik/aoc2023/lib"
)

func processLine(line string) int {
	numbers := lib.ParseNumbers(line)
	diffs := recDiff(numbers)

	data := append([][]int{numbers}, diffs...)
	fmt.Println("data", data)

	recExtrapolate(data)
	fmt.Println("extrapolated", data)
	return data[0][len(numbers)]
}

// in-place
func recExtrapolate(data [][]int) {
	for i := len(data) - 1; i >= 0; i-- {
		if i == len(data)-1 {
			data[i] = append(data[i], 0)
			continue
		}

		prow := data[i+1]
		row := data[i]
		// prev + diff
		extrapolated := row[len(row)-1] + prow[len(prow)-1]
		data[i] = append(data[i], extrapolated)
	}

}

func recDiff(nums []int) [][]int {
	lvls := [][]int{newtonDiff(nums)}
	for !allZeros(lvls[len(lvls)-1]) {
		lvls = append(lvls, newtonDiff(lvls[len(lvls)-1]))
	}
	return lvls
}

// implicit step = 1
func newtonDiff(nums []int) []int {
	der := make([]int, len(nums)-1)
	for i := range der {
		der[i] = nums[i+1] - nums[i]
	}
	return der
}

func allZeros(nums []int) bool {
	for _, v := range nums {
		if v != 0 {
			return false
		}
	}
	return true
}

func processFile(filepath string) int {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	sum := 0
	for scanner.Scan() {
		line := scanner.Text()
		value := processLine(line)
		sum += value
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return sum
}

func main() {
	sum := processFile("./input.txt")
	fmt.Println(sum)
}
