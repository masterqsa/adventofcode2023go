package main

import (
	"math"
	"strings"

	"github.com/jpillora/puzzler/harness/aoc"
)

type Galaxy struct {
	i, j int
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
	galaxies := make([]Galaxy, 0)
	raw := strings.Split(strings.TrimSpace(input), "\n")
	counts_i := make([]int, len(raw))
	counts_j := make([]int, len(strings.Replace(raw[0], "\r", "", -1)))
	for i, _ := range counts_i {
		counts_i[i] = 0
	}
	for j, _ := range counts_j {
		counts_j[j] = 0
	}

	for i, s := range raw {
		s = strings.Replace(s, "\r", "", -1)
		for j, c := range s {
			if c == '#' {
				galaxies = append(galaxies, Galaxy{i: i, j: j})
				counts_i[i]++
				counts_j[j]++
			}
		}
	}
	f_i := make([]int, len(counts_i))
	f_j := make([]int, len(counts_j))
	if counts_i[0] == 0 {
		f_i[0] = 1
	} else {
		f_i[0] = 0
	}
	if counts_j[0] == 0 {
		f_j[0] = 1
	} else {
		f_j[0] = 0
	}
	for i := 1; i < len(f_i); i++ {
		if counts_i[i] == 0 {
			f_i[i] = f_i[i-1] + 1
		} else {
			f_i[i] = f_i[i-1]
		}
	}
	for j := 1; j < len(f_j); j++ {
		if counts_j[j] == 0 {
			f_j[j] = f_j[j-1] + 1
		} else {
			f_j[j] = f_j[j-1]
		}
	}
	sum := 0
	for n := 0; n < len(galaxies)-1; n++ {
		for m := n + 1; m < len(galaxies); m++ {
			sum += int(math.Abs(float64(galaxies[n].i) - float64(galaxies[m].i)))
			sum += int(math.Abs(float64(galaxies[n].j) - float64(galaxies[m].j)))
			sum += int(math.Abs(float64(f_i[galaxies[n].i]) - float64(f_i[galaxies[m].i])))
			sum += int(math.Abs(float64(f_j[galaxies[n].j]) - float64(f_j[galaxies[m].j])))
		}
	}
	if part2 {
		var sum int64
		sum = 0
		for n := 0; n < len(galaxies)-1; n++ {
			for m := n + 1; m < len(galaxies); m++ {
				sum += int64(math.Abs(float64(galaxies[n].i) - float64(galaxies[m].i)))
				sum += int64(math.Abs(float64(galaxies[n].j) - float64(galaxies[m].j)))
				sum += (999999 * int64(math.Abs(float64(f_i[galaxies[n].i])-float64(f_i[galaxies[m].i]))))
				sum += (999999 * int64(math.Abs(float64(f_j[galaxies[n].j])-float64(f_j[galaxies[m].j]))))
			}
		}
		return sum
		// 382980107092 too high
	}
	// solve part 1 here
	return sum
}
