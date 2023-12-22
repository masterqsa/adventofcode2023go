package main

import (
	"container/heap"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/jpillora/puzzler/harness/aoc"
)

type Brick struct {
	index int
	x     [2]int
	y     [2]int
	z     [2]int
	up    map[int]bool
	down  map[int]bool
}

type PriorityQueue []*Brick

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the lowest, not highest, z so we use less than here.
	return pq[i].z[1] < pq[j].z[1]
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(*Brick))
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

func atoi(s string) int {
	value, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return value
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
	n := len(raw)
	bricks := make([]Brick, n+1)
	min_x, max_x, min_y, max_y, min_z, max_z := 0, 0, 0, 0, 99999, 0
	r := regexp.MustCompile(`([0-9]+),([0-9]+),([0-9]+)~([0-9]+),([0-9]+),([0-9]+)`)
	for i, l := range raw {
		m := r.FindStringSubmatch(l)
		b := Brick{index: i + 1, x: [2]int{atoi(m[1]), atoi(m[4])}, y: [2]int{atoi(m[2]), atoi(m[5])}, z: [2]int{atoi(m[3]), atoi(m[6])}, up: make(map[int]bool), down: make(map[int]bool)}
		bricks[i+1] = b
		if max_x < b.x[1] {
			max_x = b.x[1]
		}
		if max_y < b.y[1] {
			max_y = b.y[1]
		}
		if min_z > b.z[0] {
			min_z = b.z[0]
		}
		if max_z < b.z[1] {
			max_z = b.z[1]
		}
	}
	bricks[0] = Brick{index: 0, x: [2]int{min_x, max_x}, y: [2]int{min_y, max_y}, z: [2]int{0, 0}, up: make(map[int]bool), down: make(map[int]bool)} // artificial foundation
	fmt.Println(max_x, max_y, "(", min_z, max_z, ")")
	glass := make([][][]int, max_x-min_x+1)
	above := make([][][]int, max_x-min_x+1)
	for i := 0; i < len(glass); i++ {
		glass[i] = make([][]int, max_y-min_y+1)
		above[i] = make([][]int, max_y-min_y+1)
		for j := 0; j < len(glass[i]); j++ {
			glass[i][j] = make([]int, max_z+2)
			above[i][j] = make([]int, max_z+2)
			for k := 0; k < max_z+2; k++ {
				glass[i][j][k] = -1 // no brick present by default
			}
		}
	}
	for _, b := range bricks {
		//fmt.Println("brick", b)
		for x := b.x[0]; x <= b.x[1]; x++ {
			for y := b.y[0]; y <= b.y[1]; y++ {
				for z := b.z[0]; z <= b.z[1]; z++ {
					glass[x][y][z] = b.index
				}
			}
		}
	}
	set_neighbors := func(b Brick) {
		for x := b.x[0]; x <= b.x[1]; x++ {
			for y := b.y[0]; y <= b.y[1]; y++ {
				if glass[x][y][b.z[0]-1] > 0 {
					b.down[glass[x][y][b.z[0]-1]] = true
				}
				if glass[x][y][b.z[1]+1] > 0 {
					b.up[glass[x][y][b.z[1]+1]] = true
				}
			}
		}
	}
	for i := 1; i < len(bricks); i++ {
		set_neighbors(bricks[i])
	}
	for x := min_x; x <= max_x; x++ {
		for y := min_y; y <= max_y; y++ {
			cur := 0
			for z := 1; z < max_z; z++ {
				if glass[x][y][z] > 0 {
					cur = glass[x][y][z]
				}
				above[x][y][z] = cur
			}
		}
	}
	drop_brick := func(b Brick) int {
		max_floor := 0
		brick_height := b.z[1] - b.z[0]
		for x := b.x[0]; x <= b.x[1]; x++ {
			for y := b.y[0]; y <= b.y[1]; y++ {
				if max_floor < bricks[above[x][y][b.z[0]-1]].z[1]+1 {
					max_floor = bricks[above[x][y][b.z[0]-1]].z[1] + 1 // can't fall any further
				}
			}
		}
		for x := b.x[0]; x <= b.x[1]; x++ {
			for y := b.y[0]; y <= b.y[1]; y++ {
				for z := b.z[0]; z <= b.z[1]; z++ {
					glass[x][y][z] = -1 // cleanup
				}
			}
		}
		for x := b.x[0]; x <= b.x[1]; x++ {
			for y := b.y[0]; y <= b.y[1]; y++ {
				for z := max_floor; z <= max_floor+brick_height; z++ {
					glass[x][y][z] = b.index // paint new location
				}
				for z := max_floor; z <= b.z[0]-1; z++ {
					above[x][y][z] = b.index // paint new empty space above the brick
				}
			}
		}
		for u, _ := range b.up {
			delete(bricks[u].down, b.index) // bricks above can no longer rest on this brick
			//fmt.Println("Deleted", b.index, "from", u, "down map")
		}
		b.z[0], b.z[1], b.up = max_floor, max_floor+brick_height, make(map[int]bool)
		bricks[b.index] = b

		return max_floor
	}
	sum := 0

	/*fmt.Println("above")
	for z := 1; z <= max_z; z++ {
		for x := min_x; x <= max_x; x++ {
			for y := min_y; y <= max_y; y++ {
				fmt.Print(above[x][y][z])
			}
			fmt.Println()
		}
		fmt.Println()
	}*/
	for z := 2; z <= max_z; z++ {
		for x := min_x; x <= max_x; x++ {
			for y := min_y; y <= max_y; y++ {
				i := glass[x][y][z]
				if i > 0 {
					if bricks[i].z[0] < z {
						continue // this brick is already processed on the way up
					}
					if len(bricks[i].down) > 0 {
						//fmt.Println("skipping brick", i)
						continue // this brick has support
					}
					drop_brick(bricks[i])
					//fmt.Println("Brick", i, "dropped", z-new_z, "floors")

				}
			}
		}
	}
	for i := 1; i < len(bricks); i++ {
		set_neighbors(bricks[i])
	}
	if part2 {
		var demo map[int]bool
		check_demo := func(i int) bool {
			_, ok := demo[i]
			if ok {
				return true
			} else {
				return false
			}
		}
		falls_map := make(map[int]int)
		for z := max_z; z > 0; z-- {
			for x := min_x; x <= max_x; x++ {
				for y := min_y; y <= max_y; y++ {
					i := glass[x][y][z]
					if i > 0 {
						if _, ok := falls_map[i]; ok {
							continue // already processed
						}
						b := bricks[i]
						demo = map[int]bool{i: true}

						if b.z[0] > z {
							continue // already processed
						}
						local_sum := 0
						pq := make(PriorityQueue, 1)
						pq[0] = &b

						heap.Init(&pq)

						for pq.Len() > 0 {
							cur := heap.Pop(&pq).(*Brick)
							for u, _ := range cur.up {
								fall := true
								for d, _ := range bricks[u].down {
									if !check_demo(d) {
										fall = false
										break
									}
								}
								if fall {
									heap.Push(&pq, &(bricks[u]))
									if !check_demo(u) {
										local_sum++
										demo[u] = true
									}
								}
							}
						}

						falls_map[i] = local_sum
					}
				}
			}
		}
		for _, v := range falls_map {
			sum += v
		}
		return sum
	}

	fmt.Println(bricks)
	for i := 1; i < len(bricks); i++ {
		b := bricks[i]
		safe := true
		for u, _ := range b.up {
			if len(bricks[u].down) == 1 {
				safe = false
			}
		}
		if safe {
			sum++
			fmt.Println(i, " ", b)
		}
	}

	// solve part 1 here
	return sum
}
