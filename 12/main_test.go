package main

import (
	"strconv"
	"testing"
)

func TestProduct(t *testing.T) {
	var tests = []struct {
		n        int
		expected []string
	}{
		{1, []string{".", "#"}},
		{2, []string{".#", "#.", "..", "##"}},
		{3, []string{"..#", ".#.", "#..", "#.#", "##.", ".##", "###", "..."}},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			actual := Product(tt.n)
			if !sameStringSlice(actual, tt.expected) {
				t.Errorf("product: got %v, want %v", actual, tt.expected)
			}

		})
	}
}

func TestIsValid(t *testing.T) {
	var tests = []struct {
		s        string
		n        []int
		expected bool
	}{
		{"#", []int{1}, true},
		{"#", []int{2}, false},
		{"#.#.###", []int{1, 1, 3}, true},
		{".#...#....###.", []int{1, 1, 3}, true},
		{".#.###.#.######", []int{1, 3, 1, 6}, true},
		{"####.#...#...", []int{4, 1, 1}, true},
		{"#....######..#####.", []int{1, 6, 5}, true},
		{".###.##....#", []int{3, 2, 1}, true},
		{".###.##....#", []int{1, 2, 3}, false},
		{".###........", []int{3, 2, 1}, false},
	}

	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			e := Entry{tt.s, tt.n}
			actual := e.IsValid()
			if actual != tt.expected {
				t.Errorf("validation: got %t, want %t", actual, tt.expected)
			}

		})
	}
}

func TestHasConsecutiveExactly(t *testing.T) {
	var tests = []struct {
		txt      string
		num      int
		expected bool
	}{
		{"#", 1, true},
		{"##", 2, true},
		{"##", 3, false},
		{"####", 3, false},
		{".##", 3, false},
		{"##.", 3, false},
		{".#.#.#.", 3, false},
		{".#.#.#.###", 3, false},
		{".###.#.#.###", 3, true},
	}

	for _, tt := range tests {
		t.Run(tt.txt, func(t *testing.T) {
			actual := HasConsecutiveExactly(tt.txt, tt.num)
			if actual != tt.expected {
				t.Errorf("consecutive: got %t, want %t", actual, tt.expected)
			}

		})
	}
}

func TestArrangements(t *testing.T) {
	var tests = []struct {
		fname            string
		wantArrangements int
	}{
		{"./cases/case1.txt", 21},
	}

	for _, tt := range tests {
		testname := tt.fname
		t.Run(testname, func(t *testing.T) {
			runFor(t, tt.fname, tt.wantArrangements)
		})
	}
}

func runFor(t *testing.T, fname string, wantArrangements int) {
	total := TotalArrangementsFromFile(fname)
	if total != wantArrangements {
		t.Errorf("arrangements: got %d, want %d", total, wantArrangements)
	}
}

func sameStringSlice(x, y []string) bool {
	if len(x) != len(y) {
		return false
	}
	// create a map of string -> int
	diff := make(map[string]int, len(x))
	for _, _x := range x {
		// 0 value for int is 0, so just increment a counter for the string
		diff[_x]++
	}
	for _, _y := range y {
		// If the string _y is not in diff bail out early
		if _, ok := diff[_y]; !ok {
			return false
		}
		diff[_y]--
		if diff[_y] == 0 {
			delete(diff, _y)
		}
	}
	return len(diff) == 0
}
