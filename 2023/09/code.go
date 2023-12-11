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
	sum := 0
	sum_p := 0
	for _, line := range strings.Split(strings.TrimSpace(input), "\n") {
		val_s := strings.Fields(strings.Replace(line, "\r", "", -1))
		fmt.Println(val_s)
		n := len(val_s)
		vals := make([]int, n)
		degree := 0
		current := make([]int, n-degree)
		for i, s := range val_s {
			vals[i], _ = strconv.Atoi(s)
			current[i] = vals[i]
		}
		ds := make([]int, n-1)
		pds := make([]int, n-1)
		exit := false
		for !exit {
			exit = true

			for i := 0; i < n-degree-1; i++ {
				current[i] = current[i+1] - current[i]
				//fmt.Print(current[i], " ")
				if current[i] != 0 {
					exit = false
				}
				if i == n-degree-2 {
					ds[degree] = current[i]
					//fmt.Print(" ds= ", ds[degree])
				}
			}
			pds[degree] = current[0]
			//fmt.Println()
			if !exit {
				degree++
			}
		}
		next := vals[n-1]
		prev := vals[0]
		sign := -1
		for i := 0; i <= degree; i++ {
			next += ds[i]
			prev += sign * pds[i]
			sign *= -1
		}
		//fmt.Println("next= ", next)
		sum += next
		sum_p += prev
		fmt.Println(prev)
	}
	if part2 {
		return sum_p
	}
	// solve part 1 here
	return sum
	// 1814179399 is too low
	// 1819125966
}
