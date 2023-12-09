package main

import (
	"fmt"
	"math/big"
	"regexp"
	"strings"

	"github.com/jpillora/puzzler/harness/aoc"
)

func main() {
	aoc.Harness(run)
}

type Node struct {
	L, R string
}

// on code change, run will be executed 4 times:
// 1. with: false (part1), and example input
// 2. with: true (part2), and example input
// 3. with: false (part1), and user input
// 4. with: true (part2), and user input
// the return value of each run is printed to stdout
func run(part2 bool, input string) any {
	steps := 0
	var starts []string
	lines := strings.Split(strings.TrimSpace(input), "\n")
	dirs := strings.ReplaceAll(lines[0], "\r", "")
	fmt.Println(dirs)
	maze := make(map[string]Node)
	r := regexp.MustCompile(`([A-Z]+) = \(([A-Z]+), ([A-Z]+)\)`)
	for i := 2; i < len(lines); i++ {
		l := strings.ReplaceAll(lines[i], "\r", "")
		match := r.FindStringSubmatch(l)
		maze[match[1]] = Node{L: match[2], R: match[3]}
		if match[1][2] == 'A' {
			starts = append(starts, match[1])
			fmt.Println(match[1], match[2], match[3])
		}
	}
	if part2 {
		currents := make([]string, len(starts))
		cycles := make([]int64, len(starts))
		for i := 0; i < len(starts); i++ {
			currents[i] = starts[i]
			cycles[i] = 0
		}
		pos := 0
		done := 0
		for done < len(starts) {
			steps++
			for i := 0; i < len(currents); i++ {
				switch dirs[pos] {
				case 'L':
					currents[i] = maze[currents[i]].L
				case 'R':
					currents[i] = maze[currents[i]].R
				}
				if currents[i][2] == 'Z' {
					if cycles[i] == 0 {
						cycles[i] = int64(steps)
						fmt.Println(starts[i], steps)
						done++
					}
				}
			}
			pos++
			if pos == len(dirs) {
				pos = 0
			}
		}

		return lcmOfList(cycles)
		// 274987485330670259 is too high
	}
	// solve part 1 here

	pos := 0
	current := "AAA"
	finish := "ZZZ"
	for current != finish {
		switch dirs[pos] {
		case 'L':
			current = maze[current].L
		case 'R':
			current = maze[current].R
		}
		pos++
		if pos == len(dirs) {
			pos = 0
		}
		steps++
	}
	return steps
}

// Function to find the Greatest Common Divisor (GCD)
func gcd(a, b int64) int64 {
	return new(big.Int).GCD(nil, nil, big.NewInt(a), big.NewInt(b)).Int64()
}

// Function to find the Lowest Common Multiple (LCM) of two numbers
func lcm(a, b int64) int64 {
	return a * b / gcd(a, b)
}

// Function to find the LCM of a list of numbers
func lcmOfList(nums []int64) int64 {
	result := nums[0]
	for _, num := range nums[1:] {
		result = lcm(result, num)
	}
	return result
}
