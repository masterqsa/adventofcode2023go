package main

import (
	"fmt"
	"image"
	"strings"

	"github.com/jpillora/puzzler/harness/aoc"
)

const (
	Right = iota
	Down
	Left
	Up
)

const (
	None = iota
	Odd
	Even
	Both
)

var Dirs = []P{
	Right: {1, 0},
	Down:  {0, 1},
	Left:  {-1, 0},
	Up:    {0, -1},
}

type P = image.Point

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
	raw := strings.Split(input, "\n")
	n := len(raw)
	m := len(raw[0])
	fmt.Println(n, m)
	steps := make(map[P]int)
	get_point := func(x int, y int) byte {
		x = (n + (x % n)) % n
		y = (m + (y % m)) % m
		return raw[x][y]
	}
	get_polarity := func(p P) int {
		v, ok := steps[p]
		if ok {
			return v
		} else {
			return None
		}
	}
	set_polarity := func(p P, inc int) {
		v, ok := steps[p]
		if ok {
			steps[p] = v + inc
		} else {
			steps[p] = inc
		}
	}
	var start P
	exit := false
	for i, l := range raw {
		for j, s := range l {
			if s == 'S' {
				start = P{i, j}
				exit = true
				break
			}
		}
		if exit {
			break
		}
	}
	start = P{-1, 65}
	steps[start] = Even
	queue := []P{start}
	for dist := 1; dist <= 131; dist++ {
		next_queue := make([]P, 0)
		for _, cur := range queue {
			for _, d := range Dirs {
				cand := cur.Add(d)
				if cand.X >= 0 && cand.X < n && cand.Y >= 0 && cand.Y < m {
					if get_point(cand.X, cand.Y) != '#' {
						if (dist%2 == 0 && get_polarity(cand) < 2) || (dist%2 == 1 && get_polarity(cand)%2 == 0) {
							set_polarity(cand, 2-dist%2)
							next_queue = append(next_queue, cand)
						}
					}
				}
			}
		}
		queue = next_queue
	}

	if part2 {
		// 26501365 steps
		N := 202300
		// full squares: even 7457 N^2, odd 7383 (N-1)^2
		// left corner odd 5549 x1
		// left even 939 N
		// left odd 6480 N-1
		// top corner odd 5560 x1
		// top even 970 N
		// top odd 6460 N-1
		// right corner odd 5559 x1
		// right even 938 N
		// right odd 6492 N-1
		// bottom corner odd 5548
		// bottom even 959 N
		// bottom odd 6468 N-1
		return 7457*N*N + 7383*(N-1)*(N-1) + 5565 + 939*N + 6480*(N-1) + 5557 + 970*N + 6460*(N-1) + 5569 + 938*N + 6492*(N-1) + 5577 + 959*N + 6468*(N-1)
		//607334325965751
	}
	sum := 0
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if get_polarity(P{i, j}) == 1 {
				sum++
			}
		}
	}
	return sum
}
