package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/jpillora/puzzler/harness/aoc"
)

type Range struct {
	Start, End int
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
	seeds_str := strings.Split(strings.Split(strings.Split(strings.ReplaceAll(input, "\r\n", "\n"), "\n\n")[0], ": ")[1], " ")
	num := len(seeds_str)
	if part2 {
		num = num / 2
		current := make(map[int]int)
		next := make(map[int]int)
		for i := 0; i < num; i++ {
			n1, _ := strconv.Atoi(seeds_str[i*2])
			n2, _ := strconv.Atoi(seeds_str[i*2+1])
			current[n1] = n1 + n2 - 1
		}
		for i, lines := range strings.Split(strings.ReplaceAll(input, "\r\n", "\n"), "\n\n") {
			fmt.Println(i, lines)
			if i == 0 {
				continue
			}
			for j, line := range strings.Split(lines, "\n") {
				if j == 0 {
					continue
				}
				split := strings.Split(line, " ")
				dest, _ := strconv.Atoi(split[0])
				orig, _ := strconv.Atoi(split[1])
				leng, _ := strconv.Atoi(split[2])
				end := orig + leng - 1
				for s, e := range current {
					if s > end || e < orig {
						continue
					}
					if s >= orig && e <= end {
						next[s-orig+dest] = e - orig + dest
						delete(current, s)
					}
					if s >= orig && s <= end && e > end {
						current[end+1] = e
						next[s-orig+dest] = end - orig + dest
						delete(current, s)
					}
					if s < orig && e > end {
						current[s] = orig - 1
						current[end+1] = e
						next[dest] = dest + leng - 1
					}
					if s < orig && e >= orig && e <= end {
						current[s] = orig - 1
						next[dest] = e - orig + dest
					}
				}
			}

			fmt.Println()
			for s, e := range current {
				next[s] = e
			}
			current = next
			next = make(map[int]int)
		}
		min := math.MaxInt32
		for i, _ := range current {
			if i < min {
				min = i
			}
		}
		return min
	}
	/* 0 seeds: 79 14 55 13
	1 seed-to-soil map:
	2 soil-to-fertilizer map:
	3 fertilizer-to-water map:
	4 water-to-light map:
	5 light-to-temperature map:
	6 temperature-to-humidity map:
	7 humidity-to-location map: */

	current := make([]int, num)
	done := make([]bool, num)
	for i := 0; i < num; i++ {
		current[i], _ = strconv.Atoi(seeds_str[i])
		done[i] = false
	}
	for i, lines := range strings.Split(strings.ReplaceAll(input, "\r\n", "\n"), "\n\n") {
		fmt.Println(i, lines)
		if i == 0 {
			continue
		}
		for j, line := range strings.Split(lines, "\n") {
			if j == 0 {
				continue
			}
			split := strings.Split(line, " ")
			dest, _ := strconv.Atoi(split[0])
			orig, _ := strconv.Atoi(split[1])
			leng, _ := strconv.Atoi(split[2])
			for i := 0; i < num; i++ {
				if done[i] {
					continue
				}
				if current[i] >= orig && current[i] < orig+leng {
					current[i] = current[i] - orig + dest
					done[i] = true
				}
			}
		}
		for i := 0; i < num; i++ {
			done[i] = false
			fmt.Print(current[i], " ")
		}
		fmt.Println()
	}
	min := math.MaxInt32
	for i := 0; i < num; i++ {
		if current[i] < min {
			min = current[i]
		}
	}
	return min
	// 312406816 too low
}
