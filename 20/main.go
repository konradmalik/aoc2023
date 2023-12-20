package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"

	"github.com/konradmalik/aoc2023/lib"
)

type SignalType int

const (
	LOW SignalType = iota
	HIGH
)

type Signal struct {
	source string
	target string
	kind   SignalType
}

func (s Signal) String() string {
	var kind string
	switch s.kind {
	case HIGH:
		kind = "high"
	case LOW:
		kind = "low"
	}

	return fmt.Sprintf("%s -%s-> %s", s.source, kind, s.target)
}

func NewSignal(source, target string, kind SignalType) Signal {
	return Signal{source: source, target: target, kind: kind}
}

type Module interface {
	Process(Signal) []Signal
	Name() string
	Targets() []string
}

type Broadcaster struct {
	name    string
	targets []string
}

func (b Broadcaster) Process(s Signal) []Signal {
	signals := make([]Signal, len(b.targets))
	for i, t := range b.targets {
		signals[i] = NewSignal(b.name, t, s.kind)
	}
	return signals
}

func (b Broadcaster) Name() string {
	return b.name
}

func (b Broadcaster) Targets() []string {
	return b.targets
}

func NewBroadcaster(targets []string) Broadcaster {
	return Broadcaster{name: "broadcaster", targets: targets}
}

type FlipFlop struct {
	name    string
	on      bool
	targets []string
}

func NewFlipFlop(name string, targets []string) *FlipFlop {
	return &FlipFlop{name: name, on: false, targets: targets}
}

func (f FlipFlop) Name() string {
	return f.name
}

func (f FlipFlop) Targets() []string {
	return f.targets
}

func (f *FlipFlop) Process(s Signal) []Signal {
	switch s.kind {
	case HIGH:
		return []Signal{}
	case LOW:
		wasOn := f.on
		f.on = !f.on
		signals := make([]Signal, len(f.targets))
		for i, t := range f.targets {
			if !wasOn {
				signals[i] = NewSignal(f.name, t, HIGH)
			} else {
				signals[i] = NewSignal(f.name, t, LOW)
			}
		}
		return signals
	default:
		panic("FlipFlop received unknown signal")
	}
}

type Conjunction struct {
	name      string
	prevs     map[string]SignalType
	targets   []string
	numInputs int
}

func NewConjunction(name string, targets []string, numInputs int) *Conjunction {
	return &Conjunction{name: name, prevs: make(map[string]SignalType), targets: targets, numInputs: numInputs}
}

func (c Conjunction) Targets() []string {
	return c.targets
}

func (c Conjunction) Name() string {
	return c.name
}

func (c Conjunction) AllHigh() bool {
	if len(c.prevs) < c.numInputs {
		return false
	}

	for _, v := range c.prevs {
		if v != HIGH {
			return false
		}
	}

	return true
}

func (c *Conjunction) Process(s Signal) []Signal {
	signals := make([]Signal, len(c.targets))
	for i, t := range c.targets {
		c.prevs[s.source] = s.kind

		if c.AllHigh() {
			signals[i] = NewSignal(c.name, t, LOW)
		} else {
			signals[i] = NewSignal(c.name, t, HIGH)
		}
	}
	return signals
}

func NewModule(s string) Module {
	kind := rune(s[0])

	ss := strings.Split(s, " -> ")
	name := ss[0]
	targets := strings.Split(ss[1], ", ")

	switch kind {
	case '%':
		return NewFlipFlop(name[1:], targets)
	case '&':
		return NewConjunction(name[1:], targets, 0)
	default:
		return NewBroadcaster(targets)
	}
}

type Configuration struct {
	modules map[string]Module
	lows    int
	highs   int
}

func (c *Configuration) OneCycle() {
	signals := lib.NewQueue[Signal]()
	signals.Enqueue(NewSignal("button", "broadcaster", LOW))

	for !signals.IsEmpty() {
		s := signals.Dequeue()
		switch s.kind {
		case LOW:
			c.lows++
		case HIGH:
			c.highs++
		}

		for _, ns := range c.Process(s) {
			signals.Enqueue(ns)
		}
	}
}

func (c Configuration) Process(s Signal) []Signal {
	if m, found := c.modules[s.target]; found {
		return m.Process(s)
	}
	return []Signal{}
}

func (c *Configuration) configureConjunctions() {
	for _, m := range c.modules {
		switch v := m.(type) {
		case *Conjunction:
			inps := 0
			for _, m := range c.modules {
				if slices.Contains(m.Targets(), v.name) {
					inps++
				}
			}
			v.numInputs = inps
		default:
		}
	}
}

func NewConfiguration(filepath string) Configuration {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	modules := make(map[string]Module)
	for scanner.Scan() {
		line := scanner.Text()
		m := NewModule(line)
		modules[m.Name()] = m
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	c := Configuration{modules: modules}
	c.configureConjunctions()
	return c
}

func main() {
	c := NewConfiguration("./input.txt")
	pushes := 1000
	for i := 0; i < pushes; i++ {
		c.OneCycle()
	}

	fmt.Println("lows", c.lows)
	fmt.Println("highs", c.highs)
	fmt.Println("mult", c.lows*c.highs)
}
