package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jpillora/puzzler/harness/aoc"
)

type Cell struct {
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
	var maze [][]rune
	var dist [][]int
	distance := 0
	start_j, start_i := 0, 0
	neighbors := make(map[Cell][]rune)
	neighbors[Cell{i: -1, j: 0}] = []rune{'|', '7', 'F'}
	neighbors[Cell{i: 1, j: 0}] = []rune{'|', 'J', 'L'}
	neighbors[Cell{i: 0, j: -1}] = []rune{'-', 'L', 'F'}
	neighbors[Cell{i: 0, j: 1}] = []rune{'-', '7', 'J'}

	get_neighbors := func(start Cell) []Cell {
		res := make([]Cell, 0)
		for d, options := range neighbors {
			i := start.i + d.i
			j := start.j + d.j
			//fmt.Println("looking at", i, " ", j, " options are ", options)
			if i >= 0 && i < len(maze) && j >= 0 && j < len(maze[0]) {
				//fmt.Println("Value is ", string(maze[i][j]))
				for _, o := range options {
					if o == maze[i][j] {
						res = append(res, d)
					}
				}
			}
		}
		return res
	}

	get_next := func(start Cell) (Cell, Cell, error) {
		for d, _ := range neighbors {
			i := start.i + d.i
			j := start.j + d.j
			if i >= 0 && i < len(maze) && j >= 0 && j < len(maze[0]) {
				if dist[i][j] == dist[start.i][start.j]+1 || (dist[start.i][start.j] == distance && dist[i][j] == 0) {
					return Cell{i: i, j: j}, d, nil
				}
			}
		}
		return Cell{-1, -1}, Cell{-1, -1}, errors.New("Unable to find next cell")
	}

	check := func(c Cell) bool {
		if c.i >= 0 && c.i < len(maze) && c.j >= 0 && c.j < len(maze[0]) {
			if dist[c.i][c.j] == -1 {
				return true
			}
		}
		return false
	}

	var cur1 Cell
	raw := strings.Split(strings.TrimSpace(input), "\n")
	for i, s := range raw {
		s = strings.Replace(s, "\r", "", -1)
		row := make([]rune, 0)
		dist = append(dist, make([]int, len(s)))
		for j, c := range s {
			row = append(row, c)
			dist[i][j] = -1
			if c == 'S' {
				start_j = j
				start_i = i
				dist[start_i][start_j] = 0

			}
		}
		maze = append(maze, row)
	}
	cs := get_neighbors(Cell{i: start_i, j: start_j})
	fmt.Println("Found ", len(cs), " valid neighbors")
	cur1 = Cell{i: start_i + cs[0].i, j: start_j + cs[0].j}
	d1 := Cell{i: cs[0].i, j: cs[0].j}
	start := Cell{i: start_i, j: start_j}
	anchor := Cell{-1, -1}
	var inward string
	for cur1 != start {
		//fmt.Println(string(maze[cur1.i][cur1.j]), string(maze[cur2.i][cur2.j]))
		if dist[cur1.i][cur1.j] > 0 {
			fmt.Println("Break invoked", dist[cur1.i][cur1.j], " ", cur1)
			fmt.Println(dist)
			break
		}
		distance++
		dist[cur1.i][cur1.j] = distance
		d1, _ = getDir(d1, maze[cur1.i][cur1.j])
		if anchor == (Cell{-1, -1}) {
			if d1.i == 0 && cur1.i == 0 {
				fmt.Print("sanity")

				if maze[cur1.i][cur1.j] == '-' || maze[cur1.i][cur1.j] == '7' || maze[cur1.i][cur1.j] == 'F' {
					anchor = cur1
					if d1.j == 1 {
						inward = "R"
					}
					if d1.j == -1 {
						inward = "L"
					}
					fmt.Print(distance)
				}
			}
			if d1.j == 0 && cur1.j == 0 {
				fmt.Print("sanity")
				if maze[cur1.i][cur1.j] == '|' || maze[cur1.i][cur1.j] == 'L' || maze[cur1.i][cur1.j] == 'F' {
					anchor = cur1
					if d1.i == 1 {
						inward = "L"
					}
					if d1.i == -1 {
						inward = "R"
					}
					fmt.Print(distance)
				}
			}
		}
		cur1 = Cell{i: cur1.i + d1.i, j: cur1.j + d1.j}

	}

	if part2 {
		fmt.Println(anchor, inward)
		/*for i := 0; i < len(maze); i++ {
			for j := 0; j < len(maze[0]); j++ {
				if dist[i][j] >= 0 {
					fmt.Print("*")
				} else {
					fmt.Print(".")
				}
			}
			fmt.Println()
		}*/
		seeds := make([]Cell, 0)
		cur := anchor
		next, d, _ := get_next(cur)
		step := 0
		for step <= distance {
			cands := make([]Cell, 0)
			switch maze[cur.i][cur.j] {
			case 'I':
				if inward == "L" {
					cands = append(cands, Cell{i: cur.i, j: cur.j + d.i})
				} else {
					cands = append(cands, Cell{i: cur.i, j: cur.j - d.i})
				}
			case '-':
				if inward == "L" {
					cands = append(cands, Cell{i: cur.i - d.j, j: cur.j})
				} else {
					cands = append(cands, Cell{i: cur.i + d.j, j: cur.j})
				}
			case '7':
				if (d.i == 1 && inward == "L") || (d.j == -1 && inward == "R") {
					cands = append(cands, Cell{i: cur.i - 1, j: cur.j})
					cands = append(cands, Cell{i: cur.i - 1, j: cur.j + 1})
					cands = append(cands, Cell{i: cur.i, j: cur.j + 1})
				}
				if (d.j == -1 && inward == "L") || (d.i == 1 && inward == "R") {
					cands = append(cands, Cell{i: cur.i + 1, j: cur.j - 1})
				}
			case 'F':
				if (d.i == 1 && inward == "R") || (d.j == 1 && inward == "L") {
					cands = append(cands, Cell{i: cur.i - 1, j: cur.j})
					cands = append(cands, Cell{i: cur.i - 1, j: cur.j - 1})
					cands = append(cands, Cell{i: cur.i, j: cur.j - 1})
				}
				if (d.j == 1 && inward == "R") || (d.i == 1 && inward == "L") {
					cands = append(cands, Cell{i: cur.i + 1, j: cur.j + 1})
				}
			case 'L':
				if (d.i == -1 && inward == "L") || (d.j == 1 && inward == "R") {
					cands = append(cands, Cell{i: cur.i + 1, j: cur.j})
					cands = append(cands, Cell{i: cur.i + 1, j: cur.j - 1})
					cands = append(cands, Cell{i: cur.i, j: cur.j - 1})
				}
				if (d.j == 1 && inward == "L") || (d.i == -1 && inward == "R") {
					cands = append(cands, Cell{i: cur.i - 1, j: cur.j + 1})
				}
			case 'J':
				if (d.i == -1 && inward == "R") || (d.j == -1 && inward == "L") {
					cands = append(cands, Cell{i: cur.i + 1, j: cur.j})
					cands = append(cands, Cell{i: cur.i + 1, j: cur.j + 1})
					cands = append(cands, Cell{i: cur.i, j: cur.j + 1})
				}
				if (d.j == -1 && inward == "R") || (d.i == -1 && inward == "L") {
					cands = append(cands, Cell{i: cur.i - 1, j: cur.j - 1})
				}
			}

			for _, c := range cands {
				if check(c) {
					seeds = append(seeds, c)
				}
			}
			cur = next
			next, d, _ = get_next(cur)
			step++
		}
		count := 0
		lookup := []Cell{{i: -1, j: -1}, {i: -1, j: 0}, {i: -1, j: +1},
			{i: 0, j: -1}, {i: 0, j: +1},
			{i: 1, j: -1}, {i: 1, j: 0}, {i: 1, j: +1}}
		for len(seeds) > 0 {
			x := seeds[0]
			if check(x) {
				dist[x.i][x.j] = -2
				count++
				for _, d := range lookup {
					z := Cell{i: x.i + d.i, j: x.j + d.j}
					if check(z) {
						seeds = append(seeds, z)
					}
				}
			}
			seeds = seeds[1:]
		}
		return count
	}

	return (distance + 1) / 2
}

func getDir(from_origin Cell, tile rune) (Cell, error) {
	switch from_origin {
	case Cell{i: -1, j: 0}:
		{
			switch tile {
			case '|':
				return Cell{i: -1, j: 0}, nil
			case '7':
				return Cell{i: 0, j: -1}, nil
			case 'F':
				return Cell{i: 0, j: 1}, nil
			}
		}
	case Cell{i: 1, j: 0}:
		{
			switch tile {
			case '|':
				return Cell{i: 1, j: 0}, nil
			case 'J':
				return Cell{i: 0, j: -1}, nil
			case 'L':
				return Cell{i: 0, j: 1}, nil
			}
		}
	case Cell{i: 0, j: -1}:
		{
			switch tile {
			case '-':
				return Cell{i: 0, j: -1}, nil
			case 'F':
				return Cell{i: 1, j: 0}, nil
			case 'L':
				return Cell{i: -1, j: 0}, nil
			}
		}
	case Cell{i: 0, j: 1}:
		{
			switch tile {
			case '-':
				return Cell{i: 0, j: 1}, nil
			case '7':
				return Cell{i: 1, j: 0}, nil
			case 'J':
				return Cell{i: -1, j: 0}, nil
			}
		}
	}
	return Cell{0, 0}, errors.New("Unexpected return")
}
