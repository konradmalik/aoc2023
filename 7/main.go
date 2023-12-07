package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"
)

var cards = []rune{'A', 'K', 'Q', 'T', '9', '8', '7', '6', '5', '4', '3', '2', 'J'}

type hand struct {
	text  string
	cards [5]rune
	bid   int
}

type handtype int

const (
	Five handtype = iota
	Four
	Full
	Three
	TwoPairs
	OnePair
	High
)

func sortHands(hands []hand) {
	sort.Slice(hands, func(i, j int) bool {
		h1 := hands[i]
		h2 := hands[j]
		cmp := compareHands(h1, h2)
		// h2 is stronger, so h1 is "less"
		return cmp == 1
	})
}

func (h hand) classify() handtype {
	if slices.Contains(h.cards[:], 'J') {
		swaps := cards[:len(cards)-1]
		newhands := make([]hand, len(swaps))
		for i, swap := range swaps {
			// laaazy but... ¯\_(ツ)_/¯
			line := fmt.Sprintf("%s %d", h.text, h.bid)
			line = strings.ReplaceAll(line, "J", string(swap))
			newhand := parseLine(line)
			newhands[i] = newhand
		}

		sortHands(newhands)
		bestHand := newhands[len(newhands)-1]
		return classifyClassic(bestHand)
	}

	return classifyClassic(h)
}

func classifyClassic(h hand) handtype {
	cards := make(map[rune]int)
	for _, c := range h.cards {
		count := cards[c]
		cards[c] = count + 1
	}

	if len(cards) == 1 {
		return Five
	}

	if len(cards) == 2 {
		for _, count := range cards {
			if count == 4 {
				return Four
			}
		}
		return Full
	}

	if len(cards) == 3 {
		for _, count := range cards {
			if count == 3 {
				return Three
			}
		}
		return TwoPairs
	}

	if len(cards) == 4 {
		return OnePair
	}

	return High
}

func compareHandCards(h1 hand, h2 hand) int {
	for i := 0; i < len(h1.cards); i++ {
		c1 := h1.cards[i]
		c2 := h2.cards[i]
		c1i := slices.Index(cards, c1)
		c2i := slices.Index(cards, c2)
		if c1i < c2i {
			return -1
		} else if c1i > c2i {
			return 1
		}
	}
	return 0
}

func compareHands(h1 hand, h2 hand) int {
	h1t := h1.classify()
	h2t := h2.classify()

	if h1t < h2t {
		return -1
	} else if h1t > h2t {
		return 1
	}

	return compareHandCards(h1, h2)
}

func parseLine(line string) hand {
	elems := strings.Split(line, " ")
	cards := []rune(elems[0])
	bid, _ := strconv.Atoi(elems[1])
	return hand{elems[0], ([5]rune)(cards), bid}
}

func main() {
	file, err := os.Open("./input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	hands := make([]hand, 0, 100)
	for scanner.Scan() {
		line := scanner.Text()

		hnd := parseLine(line)
		fmt.Println(hnd.text, hnd.bid)
		fmt.Println()

		hands = append(hands, hnd)
	}

	sortHands(hands)
	fmt.Println(hands)

	winnings := 0
	for i, hand := range hands {
		rank := i + 1
		winnings += rank * hand.bid
	}

	fmt.Println(winnings)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
