package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/jpillora/puzzler/harness/aoc"
)

type Hand struct {
	Raw string
	Val int
	Bet int
}

func main() {
	aoc.Harness(run)
}

// on code change, run will be executed 4 times:
// 1. with: false (part1), and example input
// 2. with: true (part2), and example input
// 3. with: false (part1), and user input
// 4. with: true (part2), and user input
// the return value of each run is printed to stdout
func run(part2 bool, input string) any {
	// when you're ready to do part 2, remove this "not implemented" block
	sum := 0
	lines := strings.Split(strings.TrimSpace(input), "\n")
	hands := make([]Hand, len(lines))

	for i, line := range strings.Split(strings.TrimSpace(input), "\n") {
		line := strings.Split(strings.Replace(line, "\r", "", -1), " ")
		h := Convert(line[0], part2)
		v := GetStrength(h)
		b, _ := strconv.Atoi(line[1])
		hands[i] = Hand{Raw: h, Val: v, Bet: b}
		fmt.Println(h, v, b)
	}
	sort.Slice(hands, func(i, j int) bool {
		if hands[i].Val != hands[j].Val {
			return hands[i].Val < hands[j].Val // Ascending Val
		}
		return hands[i].Raw < hands[j].Raw // Ascending Raw
	})
	for i, h := range hands {
		sum += (i + 1) * h.Bet
	}
	return sum

	// 247899149 is right - part 2
}

func Convert(raw string, part2 bool) string {
	out := ""

	for _, c := range raw {
		switch c {
		case '2':
			out += "a"
		case '3':
			out += "b"
		case '4':
			out += "c"
		case '5':
			out += "d"
		case '6':
			out += "e"
		case '7':
			out += "f"
		case '8':
			out += "j"
		case '9':
			out += "k"
		case 'T':
			out += "l"
		case 'J':
			if part2 {
				out += "_"
			} else {
				out += "m"
			}
		case 'Q':
			out += "n"
		case 'K':
			out += "o"
		case 'A':
			out += "p"

		}
	}
	return out
}

func GetStrength(hand string) int {
	strength := 0
	max := 1
	cards := make(map[rune]int)
	for _, c := range hand {
		n, ok := cards[c]
		if !ok {
			cards[c] = 1
		} else {
			cards[c] = n + 1
			if c != '_' && cards[c] > max {
				max = cards[c]
			}
			if c != '_' {
				strength += cards[c]
			}
		}
	}
	j, ok := cards['_']
	if ok {
		if j == 5 {
			return 14 // all jokers
		}
		for i := 0; i < j; i++ {
			max++
			strength += max
		}
	}
	return strength
}
