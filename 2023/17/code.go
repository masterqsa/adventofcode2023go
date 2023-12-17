package main

import (
	"image"
	"strconv"
	"strings"

	"github.com/jpillora/puzzler/harness/aoc"
)

const (
	Right = iota
	Down
	Left
	Up
)

var Dirs = []P{
	Right: {1, 0},
	Down:  {0, 1},
	Left:  {-1, 0},
	Up:    {0, -1},
}

type P = image.Point

type Pos struct {
	coords P
	dir    int
	streak int
}

const FAR int = 9999999

func main() {
	aoc.Harness(run)
}

func opposite(d int) int {
	switch d {
	case Right:
		{
			return Left
		}
	case Left:
		{
			return Right
		}
	case Up:
		{
			return Down
		}
	case Down:
		{
			return Up
		}
	}
	return -1
}

// on code change, run will be executed 4 times:
// 1. with: false (part1), and example input
// 2. with: true (part2), and example input
// 3. with: false (part1), and user input
// 4. with: true (part2), and user input
// the return value of each run is printed to stdout
func run(part2 bool, input string) any {
	min_streak := 0
	max_streak := 2
	if part2 {
		min_streak = 3
		max_streak = 9
	}
	input = strings.ReplaceAll(input, "\r", "")
	raw := strings.Split(input, "\n")
	maze := make([]string, 0)
	dirs := make(map[Pos]int)
	n := len(raw)
	m := len(raw[0])
	for _, s := range raw {
		maze = append(maze, s)
	}
	front := make([]Pos, 0)
	p := Pos{P{0, 0}, Right, 0}
	front = append(front, p)
	dirs[p] = 0
	p = Pos{P{0, 0}, Down, 0}
	front = append(front, p)
	dirs[p] = 0

	get_dist := func(cur Pos) int {
		min_dist := FAR
		if cur.streak >= min_streak {
			for streak := min_streak; streak <= cur.streak; streak++ {
				v, ok := dirs[Pos{cur.coords, cur.dir, streak}]
				if ok {
					if min_dist > v {
						min_dist = v
					}
				}
			}
		} else {
			v, ok := dirs[cur]
			if ok {
				if min_dist > v {
					min_dist = v
				}
			}
		}
		return min_dist
	}
	set_dist := func(cur Pos, dist int) bool { // return true if improvement was made
		improvement := false
		if cur.streak >= min_streak {
			for streak := max_streak; streak >= cur.streak; streak-- {
				v, ok := dirs[Pos{cur.coords, cur.dir, streak}]
				if ok {
					if v > dist {
						dirs[Pos{cur.coords, cur.dir, streak}] = dist
						improvement = true
					}
				} else {
					dirs[Pos{cur.coords, cur.dir, streak}] = dist
					improvement = true
				}
			}
		} else {
			v, ok := dirs[cur]
			if ok {
				if v > dist {
					dirs[cur] = dist
					improvement = true
				}
			} else {
				dirs[cur] = dist
				improvement = true
			}
		}
		return improvement
	}
	var new_front []Pos

	visit := func(cur Pos, dist int) {
		if cur.coords.X < 0 || cur.coords.X > m-1 || cur.coords.Y < 0 || cur.coords.Y > n-1 {
			return
		}
		loss, _ := strconv.Atoi(string(maze[cur.coords.Y][cur.coords.X]))
		if set_dist(cur, dist+loss) {
			new_front = append(new_front, cur)
		}
	}
	for true {
		new_front = make([]Pos, 0)
		for _, cur := range front {
			//fmt.Println(cur)
			for d, dxy := range Dirs {
				if d == opposite(cur.dir) {
					continue
				}
				if d == cur.dir && cur.streak == max_streak {
					continue
				}
				if d != cur.dir && cur.streak < min_streak {
					continue
				}
				streak := 0
				if d == cur.dir {
					streak = cur.streak + 1
				}
				new_point := Pos{cur.coords.Add(dxy), d, streak}
				visit(new_point, get_dist(cur))
			}
		}
		if len(new_front) == 0 {
			break
		} else {
			front = new_front
		}
	}

	min := FAR
	for d, _ := range Dirs {
		dist := get_dist(Pos{P{m - 1, n - 1}, d, max_streak})
		if min > dist {
			min = dist
		}
	}
	return min
}
