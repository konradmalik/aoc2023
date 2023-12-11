package universe

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"

	"golang.org/x/exp/maps"
)

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

type Universe struct {
	galaxies map[string]Galaxy
}

func NewUniverse() Universe {
	return Universe{make(map[string]Galaxy)}
}

func (u Universe) Galaxies() []Galaxy {
	return maps.Values(u.galaxies)
}

func (u Universe) galaxiesAfterX(x int) []Galaxy {
	gs := make([]Galaxy, 0)
	for _, g := range u.galaxies {
		if g.Position.X > x {
			gs = append(gs, g)
		}
	}
	return gs
}

func (u Universe) galaxiesAfterY(y int) []Galaxy {
	gs := make([]Galaxy, 0)
	for _, g := range u.galaxies {
		if g.Position.Y > y {
			gs = append(gs, g)
		}
	}
	return gs
}

func (u Universe) Expand(times int) {
	u.expandX(times - 1)
	u.expandY(times - 1)
}

func (u Universe) expandX(extra int) {
	x := 0
	for {
		xsWithGalaxies := make(map[int]bool)
		for _, g := range u.galaxies {
			xsWithGalaxies[g.Position.X] = true
		}
		maxx := slices.Max(maps.Keys(xsWithGalaxies))

		if x >= maxx {
			return
		}

		for x < maxx {
			_, found := xsWithGalaxies[x]
			if found {
				x++
				continue
			}

			for _, g := range u.galaxiesAfterX(x) {
				g.moveX(g.Position.X + extra)
				u.galaxies[g.Label] = g
			}

			x = x + 1 + extra
			break
		}
	}
}

func (u Universe) expandY(extra int) {
	y := 0
	for {
		ysWithGalaxies := make(map[int]bool)
		for _, g := range u.galaxies {
			ysWithGalaxies[g.Position.Y] = true
		}
		maxy := slices.Max(maps.Keys(ysWithGalaxies))
		if y >= maxy {
			return
		}

		for y < maxy {
			_, found := ysWithGalaxies[y]
			if found {
				y++
				continue
			}

			for _, g := range u.galaxiesAfterY(y) {
				g.moveY(g.Position.Y + extra)
				u.galaxies[g.Label] = g
			}

			y = y + 1 + extra
			break
		}
	}
}

type Galaxy struct {
	Label    string
	Position Position
}

func (g Galaxy) String() string {
	return fmt.Sprintf("%s(%d,%d)", g.Label, g.Position.X, g.Position.Y)
}

func (g *Galaxy) moveX(x int) {
	g.Position = Position{X: x, Y: g.Position.Y}
}

func (g *Galaxy) moveY(y int) {
	g.Position = Position{X: g.Position.X, Y: y}
}
