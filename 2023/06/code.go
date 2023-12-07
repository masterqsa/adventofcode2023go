package main

import (
	"fmt"
	"math"
	"regexp"
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
	sum := 1.0
	if part2 {
		return "not implemented"
	}
	// solve part 1 here
	split := strings.Split(strings.TrimSpace(input), "\n")
	fmt.Println()
	r := regexp.MustCompile("[a-zA-Z:\r]*")
	s := regexp.MustCompile("[ ]+")
	times := strings.Split(strings.TrimSpace(s.ReplaceAllString(r.ReplaceAllString(split[0], ""), " ")), " ")
	dist := strings.Split(strings.TrimSpace(s.ReplaceAllString(r.ReplaceAllString(split[1], ""), " ")), " ")
	for i, ts := range times {
		t, _ := strconv.Atoi(ts)
		di, _ := strconv.Atoi(dist[i])
		d := float64(di) + 0.001
		x1 := math.Ceil((float64(t) - math.Sqrt(float64(t*t)-4*d)) / 2)
		x2 := math.Floor((float64(t) + math.Sqrt(float64(t*t)-4*d)) / 2)
		fmt.Println(x1, x2)
		sum *= (x2 - x1 + 1)
	}
	return int64(sum)
}
