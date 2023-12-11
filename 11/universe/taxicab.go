package universe

import "math"

func Distances(cs [][]Galaxy) []int {
	dists := make([]int, len(cs))
	for i, c := range cs {
		dists[i] = taxicabDistance(c[0].Position, c[1].Position)
	}
	return dists
}

func Combinations(gs []Galaxy) [][]Galaxy {
	cs := make([][]Galaxy, 0)
	for i := 0; i < len(gs); i++ {
		for j := 0; j < i; j++ {
			cs = append(cs, []Galaxy{gs[i], gs[j]})
		}
	}
	return cs
}

func taxicabDistance(p1 Position, p2 Position) int {
	return int(math.Abs(float64(p1.X)-float64(p2.X)) + math.Abs(float64(p1.Y)-float64(p2.Y)))
}
