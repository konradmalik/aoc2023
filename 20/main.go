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
	// increases every button push
	counter int
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

func NewSignal(source, target string, kind SignalType, counter int) Signal {
	return Signal{source: source, target: target, kind: kind, counter: counter}
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
		signals[i] = NewSignal(b.name, t, s.kind, s.counter)
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
				signals[i] = NewSignal(f.name, t, HIGH, s.counter)
			} else {
				signals[i] = NewSignal(f.name, t, LOW, s.counter)
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
	firstHigh int
}

func NewConjunction(name string, targets []string, numInputs int) *Conjunction {
	return &Conjunction{name: name, prevs: make(map[string]SignalType), targets: targets, numInputs: numInputs, firstHigh: -1}
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
			signals[i] = NewSignal(c.name, t, LOW, s.counter)
		} else {
			signals[i] = NewSignal(c.name, t, HIGH, s.counter)
			if c.firstHigh < 0 {
				c.firstHigh = s.counter
			}
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

func (c *Configuration) OneCycle(i int) bool {
	signals := lib.NewQueue[Signal]()
	signals.Enqueue(NewSignal("button", "broadcaster", LOW, i))

	for !signals.IsEmpty() {
		s := signals.Dequeue()
		switch s.kind {
		case LOW:
			// finish early if rx
			if s.target == "rx" {
				return true
			}
			c.lows++
		case HIGH:
			c.highs++
		}

		for _, ns := range c.Process(s) {
			signals.Enqueue(ns)
		}
	}

	return false
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

// greatest common divisor (GCD) via Euclidean algorithm
func GCD(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

// find Least Common Multiple (LCM) via GCD
func LCM(a, b int, integers ...int) int {
	result := a * b / GCD(a, b)

	for i := 0; i < len(integers); i++ {
		result = LCM(result, integers[i])
	}

	return result
}

func main() {
	c := NewConfiguration("./input.txt")

	// NOTE: We are being asked to figure out when a single low pulse is
	// sent to "rx". Certain things are true about the input data that
	// make this easier to figure out without brute force:
	// 1. The only input to "rx" is a single conjunction module (in my
	// case, "ft").
	// 2. The only input to each input to "ft" is also a single
	// conjunction module. (vz, bq, qh, lt)
	// 3. Each input to "ft" only sends a high pulse periodically (and
	// sends a low pulse at all other times).
	// This means that the correct number of button presses is the LCM
	// (lowest common multiple) of the period lengths for each input to
	// "ft".

	modules := make([]*Conjunction, 0)
	for _, m := range []string{"vz", "bq", "qh", "lt"} {
		modules = append(modules, c.modules[m].(*Conjunction))
	}

	// press until all fired
	i := 0
main:
	for {
		i++
		c.OneCycle(i)
		for _, m := range modules {
			if m.firstHigh < 0 {
				continue main
			}
		}
		break
	}

	nums := make([]int, 0)
	for _, m := range modules {
		nums = append(nums, m.firstHigh)
		fmt.Println(m.name, m.firstHigh)
	}

	fmt.Println("LCM", LCM(nums[0], nums[1], nums[2:]...))
}
