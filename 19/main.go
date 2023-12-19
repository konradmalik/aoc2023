package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func DecideOnPart(ws map[string]Workflow, wn string, p Part) bool {
	for {
		if wn == "A" {
			return true
		} else if wn == "R" {
			return false
		}

		w := ws[wn]
		for _, r := range w.rules {
			pv := p.params[r.field]
			switch r.op {
			case "":
				return DecideOnPart(ws, r.then, p)
			case ">":
				if pv > r.value {
					return DecideOnPart(ws, r.then, p)
				}
			case "<":
				if pv < r.value {
					return DecideOnPart(ws, r.then, p)
				}
			}
		}
	}
}

type Workflow struct {
	tag   string
	rules []Rule
}

type Rule struct {
	field string
	op    string
	value int
	then  string
}

func NewWorkflow(s string) Workflow {
	s = strings.ReplaceAll(s, "{", " ")
	s = strings.ReplaceAll(s, "}", "")
	elems := strings.Split(s, " ")
	w := Workflow{tag: elems[0], rules: make([]Rule, 0)}

	for _, e := range strings.Split(elems[1], ",") {
		rs := strings.Split(e, ":")
		if len(rs) == 1 {
			w.rules = append(w.rules, Rule{field: "", value: 0, op: "", then: rs[0]})
		} else {
			then := rs[1]
			makeRule := func(op string) Rule {
				r := strings.Split(rs[0], op)
				f := r[0]
				v, err := strconv.Atoi(r[1])
				if err != nil {
					panic(err)
				}
				return Rule{field: f, value: v, op: op, then: then}
			}

			if strings.Contains(rs[0], ">") {
				w.rules = append(w.rules, makeRule(">"))
			} else if strings.Contains(rs[0], "<") {
				w.rules = append(w.rules, makeRule("<"))
			} else {
				panic(rs)
			}
		}
	}

	return w
}

type Part struct {
	params map[string]int
}

func NewPart(s string) Part {
	s = strings.ReplaceAll(s, "{", "")
	s = strings.ReplaceAll(s, "}", "")
	elems := strings.Split(s, ",")

	p := Part{params: make(map[string]int)}

	for _, e := range elems {
		kv := strings.Split(e, "=")
		v, err := strconv.Atoi(kv[1])
		if err != nil {
			panic(err)
		}
		k := kv[0]
		p.params[k] = v
	}

	return p
}

type PartRange struct {
	min int
	max int
}

func partRangesProduct(ranges []PartRange) int {
	result := 1
	for _, r := range ranges {
		result *= r.max - r.min + 1
	}
	return result
}

func ProcessPartRanges(ws map[string]Workflow, wn string, ranges []PartRange) int {
	fieldToIdx := map[string]int{"x": 0, "m": 1, "a": 2, "s": 3}

	if wn == "A" {
		return partRangesProduct(ranges)
	} else if wn == "R" {
		return 0
	}

	combinations := 0
	w := ws[wn]
	for _, r := range w.rules {
		idx := fieldToIdx[r.field]
		newRanges := make([]PartRange, len(ranges))
		copy(newRanges, ranges[:])

		switch r.op {
		case "":
			combinations += ProcessPartRanges(ws, r.then, ranges)
		case ">":
			newRanges[idx].min = r.value + 1
			ranges[idx].max = r.value
			combinations += ProcessPartRanges(ws, r.then, newRanges)
		case "<":
			newRanges[idx].max = r.value - 1
			ranges[idx].min = r.value
			combinations += ProcessPartRanges(ws, r.then, newRanges)
		}
	}

	return combinations
}

func parseFile(filepath string) (map[string]Workflow, []Part) {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	part := false
	workflows := make(map[string]Workflow)
	parts := make([]Part, 0)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			part = true
			continue
		}

		if !part {
			workflow := NewWorkflow(line)
			workflows[workflow.tag] = workflow
		} else {
			part := NewPart(line)
			parts = append(parts, part)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return workflows, parts
}

func main() {
	workflows, parts := parseFile("./input.txt")

	accepted := make([]Part, 0)
	for _, p := range parts {
		if DecideOnPart(workflows, "in", p) {
			accepted = append(accepted, p)
		}
	}

	sum := 0
	for _, p := range accepted {
		for _, v := range p.params {
			sum += v
		}
	}
	fmt.Println(sum)

	ranges := []PartRange{{1, 4000}, {1, 4000}, {1, 4000}, {1, 4000}}
	fmt.Println(ProcessPartRanges(workflows, "in", ranges))
}
