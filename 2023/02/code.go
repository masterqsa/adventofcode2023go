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
	if part2 {
		sum := 0
		for _, line := range strings.Split(strings.TrimSpace(input), "\n") {
			line := strings.Replace(line, "\r", "", -1)
			r_max := 0
			g_max := 0
			b_max := 0
			fmt.Println(line + " ")
			parts := strings.Split(line, ": ")

			for _, attempt := range strings.Split(parts[1], "; ") {
				for _, rock := range strings.Split(attempt, ", ") {
					num, _ := strconv.Atoi(strings.Split(rock, " ")[0])
					color := strings.Replace(strings.Split(rock, " ")[1], "\r", "", -1)
					switch color {
					case "red":
						if num > r_max {
							r_max = num
						}
					case "green":
						if num > g_max {
							g_max = num
						}
					case "blue":
						if num > b_max {
							b_max = num
						}
					default:
						fmt.Println("___ALARM___ " + color + " ____!!!")
					}
				}
			}
			fmt.Println(r_max * g_max * b_max)
			sum += (r_max * g_max * b_max)

		}
		return sum
		// 65547
	}
	// solve part 1 here
	r_max := 12
	g_max := 13
	b_max := 14
	sum := 0
	for _, line := range strings.Split(strings.TrimSpace(input), "\n") {
		fmt.Print(line + " ")
		parts := strings.Split(line, ": ")
		game_id, _ := strconv.Atoi(strings.Split(parts[0], " ")[1])
		fmt.Print("Game#=" + strconv.Itoa(game_id) + " ")
		possible := true
		for _, attempt := range strings.Split(parts[1], "; ") {
			for _, rock := range strings.Split(attempt, ", ") {
				num, _ := strconv.Atoi(strings.Split(rock, " ")[0])
				color := strings.Split(rock, " ")[1]
				switch color {
				case "red":
					if num > r_max {
						possible = false
					}
				case "green":
					if num > g_max {
						possible = false
					}
				case "blue":
					if num > b_max {
						possible = false
					}
				}
			}
		}
		fmt.Println()
		if possible {
			sum += game_id
		}
	}
	return sum
}
