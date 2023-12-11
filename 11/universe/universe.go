package universe

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"golang.org/x/exp/maps"
)

type Galaxy struct {
	Label    string
	Position Position
}

type Universe struct {
	galaxies map[string]Galaxy
}

func (g Galaxy) String() string {
	return fmt.Sprintf("%s(%d,%d)", g.Label, g.Position.X, g.Position.Y)
}

func (g Galaxy) moveX(x int) Galaxy {
	g.Position = Position{X: x, Y: g.Position.Y}
	return g
}

func (g Galaxy) moveY(y int) Galaxy {
	g.Position = Position{X: g.Position.X, Y: y}
	return g
}

func NewUniverseFromFile(filepath string) Universe {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	linenr := 0
	u := NewUniverse()
	for scanner.Scan() {
		line := scanner.Text()
		u.parseInc(line, linenr)
		linenr++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return u
}

func (u *Universe) parseInc(line string, linenr int) {
	for i, r := range []rune(line) {
		if r == '#' {
			num := fmt.Sprintf("%d", len(u.galaxies)+1)
			u.galaxies[num] = Galaxy{Label: num, Position: Position{X: i, Y: linenr}}
		}
	}
}

func NewUniverse() Universe {
	return Universe{make(map[string]Galaxy)}
}

func (u Universe) Galaxies() []Galaxy {
	return maps.Values(u.galaxies)
}

func (u Universe) galaxiesAfter(cond func(Galaxy) bool) []Galaxy {
	gs := make([]Galaxy, 0)
	for _, g := range u.galaxies {
		if cond(g) {
			gs = append(gs, g)
		}
	}
	return gs
}

func (u Universe) Expand(times int) {
	u.expand(times-1, func(g Galaxy) int { return g.Position.X }, func(g Galaxy, i int) Galaxy { return g.moveX(i) })
	u.expand(times-1, func(g Galaxy) int { return g.Position.Y }, func(g Galaxy, i int) Galaxy { return g.moveY(i) })
}

func (u Universe) expand(extra int, getIndex func(Galaxy) int, move func(Galaxy, int) Galaxy) {
	x := 0
	for {
		maxx := 0
		xsWithGalaxies := make(map[int]bool)
		for _, g := range u.galaxies {
			i := getIndex(g)
			xsWithGalaxies[i] = true
			if i > maxx {
				maxx = i
			}
		}

		if x >= maxx {
			return
		}

		for x < maxx {
			_, found := xsWithGalaxies[x]
			if found {
				x++
				continue
			}

			for _, g := range u.galaxiesAfter(func(g Galaxy) bool { return getIndex(g) > x }) {
				u.galaxies[g.Label] = move(g, getIndex(g)+extra)
			}

			x = x + 1 + extra
			break
		}
	}
}
