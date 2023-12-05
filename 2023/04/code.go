package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jpillora/puzzler/harness/aoc"
)

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
	deck := make(map[int]int)
	dp := make(map[int]int)
	sum := 0
	if part2 {
		for i, line := range strings.Split(strings.TrimSpace(input), "\n") {
			winning := make(map[int]int)
			line := strings.Replace(line, "\r", "", -1)
			line = strings.Replace(line, "  ", " ", -1)
			fmt.Print(line + " -> ")
			split := strings.Split(strings.Split(line, ": ")[1], " | ")
			for j, s := range strings.Split(split[0], " ") {
				n, _ := strconv.Atoi(s)
				winning[n] = j
			}
			count := 0
			for _, s := range strings.Split(split[1], " ") {
				n, _ := strconv.Atoi(s)
				_, ok := winning[n]
				if ok {
					fmt.Print(n)
					fmt.Print(", ")
					count++
				}
			}
			if count > 0 {
				fmt.Println(count, 1<<(count-1))
				sum += (1 << (count - 1))
			} else {
				fmt.Println(count, 0)
			}
			deck[i] = count
		}
		size := len(deck)
		sum = 1
		dp[size-1] = 1
		for i := size - 2; i >= 0; i-- {
			draw := 1
			n := deck[i]
			if n > size-i-1 {
				n = size - i - 1
			}
			for j := 1; j <= n; j++ {
				draw += dp[i+j]
			}
			dp[i] = draw
			sum += draw
		}
		return sum
		//5923918
	}
	// solve part 1 here

	for _, line := range strings.Split(strings.TrimSpace(input), "\n") {
		winning := make(map[int]int)
		line := strings.Replace(line, "\r", "", -1)
		line = strings.Replace(line, "  ", " ", -1)
		fmt.Print(line + " -> ")
		split := strings.Split(strings.Split(line, ": ")[1], " | ")
		for j, s := range strings.Split(split[0], " ") {
			n, _ := strconv.Atoi(s)
			winning[n] = j
		}
		count := 0
		for _, s := range strings.Split(split[1], " ") {
			n, _ := strconv.Atoi(s)
			_, ok := winning[n]
			if ok {
				fmt.Print(n)
				fmt.Print(", ")
				count++
			}
		}
		if count > 0 {
			fmt.Println(count, 1<<(count-1))
			sum += (1 << (count - 1))
		} else {
			fmt.Println(count, 0)
		}
	}
	return sum
	// 23441
}
