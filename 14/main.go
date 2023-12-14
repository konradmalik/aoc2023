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

// assume rocks alwys roll down
func (p Dish) Load() int {
	load := 0
	for i, row := range p.Rows() {
		for _, cell := range row {
			if cell == 'O' {
				load += i + 1
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
	})
}

func (p Dish) rollToSide(getCols func(Dish) []string) Dish {
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

	return NewDishFromCols(cols)
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

func main() {
	dish := parseDish("./input.txt")
	fmt.Println(dish)
	dish = dish.RollNorth()
	fmt.Println(dish)
	fmt.Println(dish.Load())
}
