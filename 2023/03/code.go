package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/jpillora/puzzler/harness/aoc"
)

type Cell struct {
	Y, X int
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
	m := len(strings.Split(strings.TrimSpace(input), "\n"))
	n := len(strings.Replace(strings.Split(strings.TrimSpace(input), "\n")[0], "\r", "", -1))
	layout := make([][]string, m)
	for i, line := range strings.Split(strings.TrimSpace(input), "\n") {
		line := strings.Replace(line, "\r", "", -1)
		layout[i] = make([]string, n)
		for j := 0; j < n; j++ {
			layout[i][j] = line[j : j+1]
		}
		fmt.Println(line)
	}
	r := regexp.MustCompile("[0-9]")
	neg := regexp.MustCompile("[0-9.]")
	sum := 0
	// when you're ready to do part 2, remove this "not implemented" block
	if part2 {
		gears := make([][]int, m)
		for i := 0; i < m; i++ {
			gears[i] = make([]int, n)
		}
		gear_ratios := make([][]int, m)
		for i := 0; i < m; i++ {
			gear_ratios[i] = make([]int, n)
		}

		for i := 0; i < m; i++ {
			for j := 0; j < n; j++ {
				gears[i][j] = 0
				gear_ratios[i][j] = 1
			}
		}

		for i := 0; i < m; i++ {
			neighbors := make(map[Cell]bool)
			base := 0
			for j := 0; j < n; j++ {
				if !r.MatchString(layout[i][j]) {
					if len(neighbors) > 0 {
						for g, _ := range neighbors {
							gears[g.Y][g.X] += 1
							if gears[g.Y][g.X] > 2 {
								gear_ratios[g.Y][g.X] = 0
							} else {
								gear_ratios[g.Y][g.X] *= base
							}
						}
					}
					base = 0
					neighbors = make(map[Cell]bool)
				}
				if r.MatchString(layout[i][j]) {
					v, _ := strconv.Atoi(layout[i][j])
					base = base*10 + v
					for dy := -1; dy <= 1; dy++ {
						if i+dy >= 0 && i+dy < m {
							for dx := -1; dx <= 1; dx++ {
								if j+dx >= 0 && j+dx < n {
									if layout[i+dy][j+dx] == "*" {
										neighbors[Cell{i + dy, j + dx}] = true
									}
								}
							}
						}
					}
				}
			}
			if len(neighbors) > 0 {
				for g, _ := range neighbors {
					gears[g.Y][g.X] += 1
					if gears[g.Y][g.X] > 2 {
						gear_ratios[g.Y][g.X] = 0
					} else {
						gear_ratios[g.Y][g.X] *= base
					}
				}
			}
		}

		for i := 0; i < m; i++ {
			for j := 0; j < n; j++ {
				if gears[i][j] == 2 {
					sum += gear_ratios[i][j]
				}
			}
		}
		return sum
		// 84363105 is the right answer
	}
	// solve part 1 here

	/*for i := 0; i < m; i++ {
		fmt.Println()
		for j := 0; j < n; j++ {
			fmt.Print(layout[i][j])
		}
	}*/

	for i := 0; i < m; i++ {
		valid := false
		base := 0
		for j := 0; j < n; j++ {
			if !r.MatchString(layout[i][j]) {
				if valid {
					sum += base
					fmt.Println("Added base=" + strconv.Itoa(base))
				}
				valid = false
				base = 0
			}
			if r.MatchString(layout[i][j]) {
				v, _ := strconv.Atoi(layout[i][j])
				base = base*10 + v
				for dy := -1; dy <= 1; dy++ {
					if i+dy >= 0 && i+dy < m {
						for dx := -1; dx <= 1; dx++ {
							if j+dx >= 0 && j+dx < n {
								if !neg.MatchString(layout[i+dy][j+dx]) {
									valid = true
								}
							}
						}
					}
				}
			}
		}
		if valid {
			sum += base
			fmt.Println("Added base=" + strconv.Itoa(base))
		}
	}
	return sum
	// 13466510 is too high
	// 550853 is too low
	// 553079 is right
}
