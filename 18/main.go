package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

type Entry struct {
	d     string
	count int
}

func NewEntry(s string) Entry {
	ss := strings.Split(s, " ")
	hex := strings.ReplaceAll(ss[2], "(", "")
	hex = strings.ReplaceAll(hex, ")", "")
	hex = hex[1:]
	count := hex[:len(hex)-1]
	d := rune(hex[len(hex)-1])
	var dir string
	if d == '0' {
		dir = "R"
	} else if d == '1' {
		dir = "D"
	} else if d == '2' {
		dir = "L"
	} else {
		dir = "U"
	}

	n, err := strconv.ParseInt(count, 16, 64)
	fmt.Println(n)
	if err != nil {
		panic(err)
	}
	return Entry{dir, int(n)}
}

type Point struct {
	x int
	y int
}

func (p Point) Dist(q Point) int {
	return int(math.Abs(float64(p.x-q.x)) + math.Abs(float64(p.y-q.y)))
}

func EntriesToPoints(entries []Entry) []Point {
	points := make([]Point, len(entries)+1)
	points[0] = Point{x: 0, y: 0}
	for i, e := range entries {
		l := points[i]
		switch e.d {
		case "U":
			points[i+1] = Point{x: l.x, y: l.y - e.count}
		case "D":
			points[i+1] = Point{x: l.x, y: l.y + e.count}
		case "L":
			points[i+1] = Point{x: l.x - e.count, y: l.y}
		case "R":
			points[i+1] = Point{x: l.x + e.count, y: l.y}
		}
	}

	if points[0] != points[len(points)-1] {
		panic("did not close")
	}
	return points[:len(points)-1]
}

func Perimeter(pts []Point) int {
	sum := 0
	for i := 1; i < len(pts); i++ {
		sum += pts[i-1].Dist(pts[i])
	}
	sum += pts[len(pts)-1].Dist(pts[0])
	return sum
}

// shoelace
func Area(pts []Point) int {
	area := 0
	for i, j := 0, len(pts)-1; i < len(pts); i++ {
		pi := pts[i]
		pj := pts[j]
		area += (pj.x + pi.x) * (pj.y - pi.y)
		j = i
	}

	return int(math.Abs(float64(area / 2)))
}

func parseEntries(filepath string) []Entry {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	entries := make([]Entry, 0)
	for scanner.Scan() {
		line := scanner.Text()
		entry := NewEntry(line)
		entries = append(entries, entry)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return entries
}

func main() {
	entries := parseEntries("./input.txt")
	fmt.Println(entries)
	points := EntriesToPoints(entries)
	fmt.Println(points)
	perimeter := Perimeter(points)
	fmt.Println(perimeter)
	area := Area(points)
	fmt.Println(area)
	result := area + perimeter/2 + 1
	fmt.Println(result)
}
