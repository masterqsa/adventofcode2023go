package main

import (
	"fmt"
	"image"
	"strings"
	"sync"

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

var AllDirs = []P{
	Right: {0, 1},
	Down:  {1, 0},
	Left:  {0, -1},
	Up:    {-1, 0},
}

var Dirs = map[byte][]P{'.': AllDirs, 'v': {{1, 0}}, '^': {{-1, 0}}, '<': {{0, -1}}, '>': {{0, 1}}}

type P = image.Point

func main() {
	aoc.Harness(run)
}

func isBitSet(bitSet, bitPos int) bool {
	return (bitSet & (1 << bitPos)) != 0
}

func setBit(bitSet, bitPos int) int {
	return bitSet | (1 << bitPos)
}

var start, end P
var visited map[P]bool
var raw []string
var n, m int
var pivots []P
var crossdist map[[2]int]int
var pointtoindex map[P]int
var neighbors map[int][]int
var distances [][]int

// on code change, run will be executed 4 times:
// 1. with: false (part1), and example input
// 2. with: true (part2), and example input
// 3. with: false (part1), and user input
// 4. with: true (part2), and user input
// the return value of each run is printed to stdout
func run(part2 bool, input string) any {
	input = strings.ReplaceAll(input, "\r", "")
	raw = strings.Split(input, "\n")
	n = len(raw)
	m = len(raw[0])
	visited = make(map[P]bool)
	start = P{0, 1}
	end = P{n - 1, m - 2}
	if !part2 {
		return navigateConcurrently(start, visited, part2)
	}

	visited[start] = true
	neighbors = make(map[int][]int)
	pointtoindex = make(map[P]int)
	crossdist = make(map[[2]int]int)
	pivots = make([]image.Point, 0)
	pivots = append(pivots, start)
	pointtoindex[start] = 0
	explore(start, start.Add(AllDirs[Down]), 1)
	distances = make([][]int, len(pivots))
	for i := 0; i < len(pivots); i++ {
		distances[i] = make([]int, len(pivots))
		neighbors[i] = make([]int, 0)
	}
	for c, d := range crossdist {
		distances[c[0]][c[1]] = d
		neighbors[c[0]] = append(neighbors[c[0]], c[1])
	}
	fmt.Println(len(pivots))
	fmt.Println(pivots)
	fmt.Println(crossdist)
	build_map()
	vis := 0
	vis = setBit(vis, 0)
	result, _ := nav(0, vis)
	cur := [2]int{result, 0}
	exit := false
	fmt.Println()
	fmt.Println("Trace:")
	for !exit {
		if v, ok := trace[cur]; !ok {
			exit = true
		} else {
			fmt.Println(cur, v, pivots[v[0]], distances[v[0]][v[1]])
			if cur[0] != v[2] && cur[1] != v[1] {
				cur = [2]int{v[2], v[1]}
			} else {
				exit = true
			}
		}
	}
	return result
	// 6262
}

// preprocessing to build condensed graph
func explore(last_pivot P, s P, dist int) {
	if s == end {
		pivots = append(pivots, end)
		idx1 := pointtoindex[last_pivot]
		idx2 := len(pivots) - 1
		crossdist[[2]int{idx1, idx2}] = dist
		crossdist[[2]int{idx2, idx1}] = dist
		pointtoindex[end] = idx2
		fmt.Println("Finish:", dist)
		return
	}
	visited[s] = true
	dirs := AllDirs
	cands := make([]P, 0)
	for _, d := range dirs {
		cand := s.Add(d)
		if cand.X >= 0 && cand.X < n && cand.Y >= 0 && cand.Y < m {
			if raw[cand.X][cand.Y] != '#' {
				if _, ok := visited[cand]; !ok {
					cands = append(cands, cand)
				}
				if pivot_idx, ok := pointtoindex[cand]; ok {
					idx1 := pointtoindex[last_pivot]
					idx2 := pivot_idx
					if idx1 != idx2 {
						crossdist[[2]int{idx1, idx2}] = dist + 1
						crossdist[[2]int{idx2, idx1}] = dist + 1
					}
				}
			}
		}
	}
	if len(cands) > 1 {
		// new pivot point
		pivots = append(pivots, s)
		idx1 := pointtoindex[last_pivot]
		idx2 := len(pivots) - 1
		crossdist[[2]int{idx1, idx2}] = dist
		crossdist[[2]int{idx2, idx1}] = dist
		pointtoindex[s] = idx2
		for _, cand := range cands {
			explore(s, cand, 1)
		}
	} else if len(cands) > 0 {
		explore(last_pivot, cands[0], dist+1)
	}
}

var trace map[[2]int][3]int = make(map[[2]int][3]int)

// Part 2 navigation with improved "visited" structure - much better performance
func nav(s int, visited int) (int, int) {
	if s == pointtoindex[end] {
		return 0, setBit(visited, s)
	}
	best_dist := -10000
	best_v := visited
	best_cand := 0
	for _, d := range neighbors[s] {
		if !isBitSet(visited, d) {
			incremental, v := nav(d, setBit(visited, s))
			dist := distances[s][d] + incremental
			if dist >= best_dist {
				best_dist = dist
				best_v = v
				best_cand = d
			}
		}
	}
	// for tracing only, not needed for calculation
	trace[[2]int{best_dist, s}] = [3]int{s, best_cand, best_dist - distances[s][best_cand]}
	return best_dist, best_v
}

// Part 1 navigation
func navigate(s P, visited map[P]bool, part2 bool) int {
	if s == end {
		return 0
	}
	new_vis := make(map[P]bool)
	for key, value := range visited {
		new_vis[key] = value
	}
	new_vis[s] = true
	dirs := AllDirs
	if !part2 {
		dirs = Dirs[raw[s.X][s.Y]]
	}
	best_dist := 0
	for _, d := range dirs {
		cand := s.Add(d)
		if cand.X >= 0 && cand.X < n && cand.Y >= 0 && cand.Y < m {
			if raw[cand.X][cand.Y] != '#' {
				if _, ok := visited[cand]; !ok {
					dist := navigate(cand, new_vis, part2)
					if dist > best_dist {
						best_dist = dist
					}
				}
			}
		}
	}
	return best_dist + 1
}

// no significant improvement in performance from parallelization, but it works
func navigateConcurrently(s P, visited map[P]bool, part2 bool) int {
	if s == end {
		return 0
	}

	newVis := make(map[P]bool)
	for key, value := range visited {
		newVis[key] = value
	}
	newVis[s] = true

	dirs := AllDirs
	if !part2 {
		dirs = Dirs[raw[s.X][s.Y]]
	}

	var wg sync.WaitGroup
	results := make(chan int, len(dirs))

	for _, d := range dirs {
		wg.Add(1)
		go func(d P) {
			defer wg.Done()
			cand := s.Add(d)
			if cand.X >= 0 && cand.X < n && cand.Y >= 0 && cand.Y < m {
				if raw[cand.X][cand.Y] != '#' {
					if _, ok := visited[cand]; !ok {
						dist := navigate(cand, newVis, part2)
						results <- dist + 1
					}
				}
			}
		}(d)
	}

	wg.Wait()
	close(results)

	bestDist := 0
	for dist := range results {
		if dist > bestDist {
			bestDist = dist
		}
	}

	return bestDist
}

// made for debugging and verification - not needed for calculation
func build_map() {
	knots := make(map[P]bool)
	knots[start] = true
	for i := 1; i < n-1; i++ {
		for j := 1; j < m-1; j++ {
			if raw[i][j] != '#' {
				count := 0
				for _, d := range AllDirs {
					if raw[i+d.X][j+d.Y] != '#' {
						count++
					}
				}
				if count > 2 {
					knots[P{i, j}] = true
				}
			}
		}
	}
	knots[end] = true
	//fmt.Println("knots: ", len(knots))
	//fmt.Println(knots)
	//fmt.Println("-----------")
	knot_distance := make(map[P]map[P]int)
	bfs := func(k P) {
		steps := 0
		visited := make(map[P]bool)
		visited[k] = true
		knot_distance[k] = make(map[P]int)
		queue := []P{k}
		exit := false
		for !exit {
			steps++
			next_queue := make([]P, 0)
			for _, q := range queue {
				for _, d := range AllDirs {
					cand := q.Add(d)
					if cand.X >= 0 && cand.X < n && cand.Y >= 0 && cand.Y < m {
						if _, ok := visited[cand]; ok {
							continue
						}
						if raw[cand.X][cand.Y] != '#' {
							if _, ok := knots[cand]; ok { // found neighbor knot at a distance steps
								knot_distance[k][cand] = steps
							} else {
								next_queue = append(next_queue, cand)
								visited[cand] = true
							}
						}
					}
				}
			}
			if len(next_queue) == 0 {
				exit = true
			} else {
				queue = next_queue
			}
		}
	}
	for k, _ := range knots {
		bfs(k)
	}
	//fmt.Println("Map distances", knot_distance)
	for k := range knots {
		idx1 := pointtoindex[k]
		for p, dd := range knot_distance[k] {
			idx2 := pointtoindex[p]
			d := distances[idx1][idx2]
			if d != dd {
				fmt.Println("Found mismatch", d, dd, k, p)
			} else {
				//fmt.Println("Match for", d, dd, k, p)
			}
		}
	}
}
