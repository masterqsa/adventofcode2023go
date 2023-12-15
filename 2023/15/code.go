package main

import (
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/jpillora/puzzler/harness/aoc"
)

type Lens struct {
	focal int
	pos   int
}

type LensSet struct {
	lens    map[string]Lens
	max_pos int
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
	input = strings.ReplaceAll(input, "\r", "")
	sum := 0
	r := regexp.MustCompile(`([a-z]+)([=\-]+)([0-9]*)`)
	if part2 {
		boxes := make([]LensSet, 256)
		for i := 0; i < 256; i++ {
			boxes[i] = LensSet{lens: make(map[string]Lens), max_pos: -1}
		}
		for _, s := range strings.Split(input, ",") {
			match := r.FindStringSubmatch(s)
			label := match[1]
			cur := hash(label)
			switch match[2] {
			case "-":
				{
					_, ok := boxes[cur].lens[label]
					if ok {
						delete(boxes[cur].lens, label)
					}
				}
			case "=":
				{
					focal, _ := strconv.Atoi(match[3])
					c, ok := boxes[cur].lens[label]
					if ok {
						boxes[cur].lens[label] = Lens{focal: focal, pos: c.pos}
					} else {
						boxes[cur].max_pos += 1
						boxes[cur].lens[label] = Lens{focal: focal, pos: boxes[cur].max_pos}
					}

				}
			}
		}
		for i := 0; i < 256; i++ {
			j := 0
			lens := boxes[i].lens
			// Extract the keys into a slice
			keys := make([]string, 0, len(lens))
			for k := range lens {
				keys = append(keys, k)
			}

			// Sort the keys slice based on the pos field
			sort.Slice(keys, func(i, j int) bool {
				return lens[keys[i]].pos < lens[keys[j]].pos
			})

			// Now keys are sorted based on the pos field
			for _, k := range keys {
				val := (i + 1) * (j + 1) * lens[k].focal
				sum += val
				//fmt.Println(val, i+1, j+1, lens[k].focal)
				j++
			}

		}
		return sum
		// 257355 too low
	}

	for _, s := range strings.Split(input, ",") {
		cur := hash(s)
		sum += cur
	}
	return sum
}

func hash(input string) int {
	cur := 0
	for _, c := range input {
		cur += int(c)
		cur *= 17
		cur %= 256
	}
	return cur
}
