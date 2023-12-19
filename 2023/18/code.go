package main

import (
	"fmt"
	"image"
	"sort"
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

type Interval struct {
	Start, End int
	Horiz      bool
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
	input = strings.ReplaceAll(input, "\r", "")
	raw := strings.Split(input, "\n")
	min_y, max_y := 0, 0
	cond := make(map[int]map[int]bool) // condensed representation, first key is y coord, second key is x coord, last value is false for vertical, true for horizontal
	d := map[string]int{"U": Up, "D": Down, "R": Right, "L": Left}
	cur := P{0, 0}
	sum := 0
	cond[cur.Y] = make(map[int]bool)
	cond[cur.Y][cur.X] = false
	for _, l := range raw {
		parts := strings.Fields(l)
		dir_idx := d[parts[0]]
		dir := Dirs[dir_idx]
		dist, _ := strconv.Atoi(parts[1])
		if part2 {
			hexNumber := parts[2][2:7]
			hex, _ := strconv.ParseInt(hexNumber, 16, 32)
			dist = int(hex)
			d_idx, _ := strconv.Atoi(parts[2][7:8])
			dir_idx = d_idx
			dir = Dirs[dir_idx]
		}
		if dir_idx == Right {
			cond[cur.Y][cur.X+dir.X*dist] = true // marking rightmost point of horizontals with "true"
			cur = P{X: cur.X + dir.X*dist, Y: cur.Y}
		}
		if dir_idx == Left {
			cond[cur.Y][cur.X+dir.X*dist] = false
			cond[cur.Y][cur.X] = true
			cur = P{X: cur.X + dir.X*dist, Y: cur.Y}
		}
		if dir_idx == Up || dir_idx == Down {
			for i := 1; i <= dist; i++ {
				_, ok := cond[cur.Y+i*dir.Y]
				if !ok {
					cond[cur.Y+i*dir.Y] = make(map[int]bool)
				}
				cond[cur.Y+i*dir.Y][cur.X] = false
			}
			cur = P{X: cur.X, Y: cur.Y + dir.Y*dist}
		}
		if cur.Y < min_y {
			min_y = cur.Y
		}
		if cur.Y > max_y {
			max_y = cur.Y
		}
	}
	get_intervals := func(points map[int]bool, initial bool) []Interval {
		step := 1
		if initial {
			step = 2
		}
		result := make([]Interval, 0)
		sorted_points := make([]int, 0, len(points))
		for key := range points {
			sorted_points = append(sorted_points, key)
		}
		// Sort the keys
		sort.Ints(sorted_points)
		for i := 1; i < len(sorted_points); i += step {
			if sorted_points[i]-sorted_points[i-1] > 1 {
				result = append(result, Interval{sorted_points[i-1] + 1, sorted_points[i] - 1, points[sorted_points[i]]})
			}
		}
		return result
	}

	fmt.Println(min_y, max_y)
	sum += len(cond[min_y])
	prev := get_intervals(cond[min_y], true)
	fmt.Println(prev)
	for y := min_y + 1; y <= max_y; y++ {
		sum += calcLength(prev)
		sum += len(cond[y])
		//fmt.Print("points ", len(cond[y]), "+ ")
		ints := get_intervals(cond[y], false)
		intersecting, not_intersecting := findIntersectingIntervals(prev, ints)
		prev = make([]Interval, 0)
		//fmt.Println("Intersecting", intersecting)
		for _, i := range intersecting {
			if !i.Horiz {
				prev = append(prev, i)
			} else {
				l := calcIntervalLength(i)
				sum += l
				//fmt.Print(l, "+")
			}
		}
		//fmt.Println("Not Intersecting", not_intersecting)
		for _, i := range not_intersecting {
			if i.Horiz {
				prev = append(prev, i)
			}
		}
		//fmt.Println(prev)
		if len(prev) == 0 {
			fmt.Println("Break at ", y, cond[y])
			break
		}
	}
	sum += calcLength(prev)
	fmt.Println("last row", sum)

	return sum
}

func doIntervalsIntersect(a, b Interval) bool {
	return a.Start <= b.End && b.Start <= a.End
}

func findIntersectingIntervals(set1, set2 []Interval) ([]Interval, []Interval) {
	used := make(map[int]bool)
	var intersectingIntervals []Interval
	var notIntersectingIntervals []Interval
	for _, interval1 := range set1 {
		for i, interval2 := range set2 {
			if doIntervalsIntersect(interval1, interval2) {
				_, ok := used[i]
				if !ok {
					intersectingIntervals = append(intersectingIntervals, interval2)
					used[i] = true
				}
			}

		}
	}
	for i, interval2 := range set2 {
		_, ok := used[i]
		if !ok {
			notIntersectingIntervals = append(notIntersectingIntervals, interval2)
		}
	}
	return intersectingIntervals, notIntersectingIntervals
}

func calcLength(ints []Interval) int {
	length := 0
	for _, v := range ints {
		length += calcIntervalLength(v)
	}
	return length
}

func calcIntervalLength(v Interval) int {
	return v.End - v.Start + 1
}
