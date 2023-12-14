package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
)

type Dish []string

func NewDishFromRows(rows []string) Dish {
	return rows
}

func NewDishFromCols(cols []string) Dish {
	rows := make([]string, len(cols[0]))
	for i := range cols[0] {
		row := make([]rune, len(cols))
		for j, col := range cols {
			row[j] = []rune(col)[i]
		}
		rows[i] = string(row)
	}
	return rows

}

func (p Dish) String() string {
	var sb strings.Builder
	for _, r := range p {
		sb.WriteString(r)
		sb.WriteRune('\n')
	}
	return sb.String()
}

func (p Dish) Rows() []string {
	return p
}

func (p Dish) Cols() []string {
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

func (p Dish) NorthLoad() int {
	load := 0
	for i, row := range p.Rows() {
		for _, cell := range row {
			if cell == 'O' {
				load += len(p) - i
			}
		}
	}
	return load
}

func (p Dish) RollNorth() Dish {
	return p.rollToSide(func(d Dish) []string {
		cs := d.Cols()
		for i := range cs {
			row := []rune(cs[i])
			slices.Reverse(row)
			cs[i] = string(row)
		}
		return cs
	},
		func(s []string) Dish {
			for i := range s {
				row := []rune(s[i])
				slices.Reverse(row)
				s[i] = string(row)
			}
			return NewDishFromCols(s)
		})
}

func (p Dish) RollWest() Dish {
	return p.rollToSide(func(d Dish) []string {
		cs := d.Rows()
		for i := range cs {
			row := []rune(cs[i])
			slices.Reverse(row)
			cs[i] = string(row)
		}
		return cs
	},
		func(s []string) Dish {
			for i := range s {
				row := []rune(s[i])
				slices.Reverse(row)
				s[i] = string(row)
			}
			return NewDishFromRows(s)
		})
}

func (p Dish) RollSouth() Dish {
	return p.rollToSide(func(d Dish) []string {
		return p.Cols()
	},
		func(s []string) Dish {
			return NewDishFromCols(s)
		})
}

func (p Dish) RollEast() Dish {
	return p.rollToSide(func(d Dish) []string {
		return p.Rows()
	},
		func(s []string) Dish {
			return NewDishFromRows(s)
		})
}

func (p Dish) rollToSide(getCols func(Dish) []string, getDish func([]string) Dish) Dish {
	cols := getCols(p)

	for i, col := range cols {
		rolled := true
		for rolled {
			rolled = false
			for c, cell := range col {
				if cell == 'O' {
					rns, didroll := roll([]rune(col), c)
					if !rolled {
						rolled = didroll
					}
					cols[i] = string(rns)
					col = cols[i]
				}
			}
		}
	}

	return getDish(cols)
}

func roll(row []rune, srci int) ([]rune, bool) {
	src := row[srci]
	if src != 'O' {
		panic("cannot roll " + string(src))
	}

	rolled := false
	for (srci + 1) < len(row) {
		nxt := row[srci+1]
		if nxt == '#' || nxt == 'O' {
			break
		}

		row[srci+1] = 'O'
		row[srci] = '.'
		srci++
		rolled = true
	}

	return row, rolled
}

func parseDish(filepath string) Dish {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	dish := make([]string, 0)

	for scanner.Scan() {
		line := scanner.Text()
		dish = append(dish, line)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return dish
}

// assumes that seq contains only that pattern
func guessRepeatingPatternLength(seq []int) int {
	maxLen := len(seq) / 2
	for i := 2; i < maxLen; i++ {
		if slices.Compare(seq[:i], seq[i:i*2]) == 0 {
			return i
		}
	}
	return maxLen
}

func guessRepeatingPatternStartIdx(seq []int) int {
	// mark true where repeats
	tracker := make([]bool, len(seq))
	for i, v := range seq {
		if past := slices.Index(seq[:i], v); past != -1 {
			tracker[i] = true
		} else {
			tracker[i] = false
		}
	}

	// reverse to find last idx where false
	slices.Reverse(tracker)
	revIdx := slices.Index(tracker, false)
	return len(seq) - revIdx
}

func cycles(dish Dish, count int) []int {
	loads := make([]int, count)
	i := 1
	for i <= count {
		dish = dish.RollNorth()
		dish = dish.RollWest()
		dish = dish.RollSouth()
		dish = dish.RollEast()
		load := dish.NorthLoad()
		log.Println("cycle", i, "load", load)
		loads[i-1] = load
		i++
	}

	return loads
}

func main() {
	dish := parseDish("./input.txt")
	fmt.Println(dish)
	loads := cycles(dish, 200)
	cycleStart := guessRepeatingPatternStartIdx(loads)
	fmt.Println("cycle starts at", cycleStart)
	repeatingLoads := loads[cycleStart:]
	fmt.Println("repeatingLoads", repeatingLoads)
	cycleLen := guessRepeatingPatternLength(repeatingLoads)
	fmt.Println("cycle length is", cycleLen)
	resultIdx := ((1000000000 - cycleStart) % cycleLen) - 1
	fmt.Println("answer index is", resultIdx)
	fmt.Println("the answer", repeatingLoads[resultIdx])
}
