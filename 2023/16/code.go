package main

import (
	"strings"

	"github.com/jpillora/puzzler/harness/aoc"
)

type Loc struct {
	i   int
	j   int
	dir string
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
	sum := 0
	maze := make([]string, 0)
	dirs := make([][]map[string]bool, 0) // string key - NSWE format
	n := len(raw)
	m := len(raw[0])
	for _, s := range raw {
		maze = append(maze, s)
	}
	rays := make(map[int]Loc)
	ray_index := 0

	init_and_run := func() int {
		sum = 0
		for i := 0; i < n; i++ {
			dirs = append(dirs, make([]map[string]bool, m))
			for j := 0; j < m; j++ {
				dirs[i][j] = make(map[string]bool)
			}
		}
		for len(rays) > 0 {
			//fmt.Println(len(rays))
			new_rays := make(map[int]Loc)
			for k, r := range rays {
				if r.i < 0 || r.i >= n || r.j < 0 || r.j >= m {
					delete(rays, k)
					continue
				}
				_, ok := dirs[r.i][r.j][r.dir]
				if ok {
					delete(rays, k)
					//fmt.Println("deleted ray ", r.i, r.j, r.dir)
					continue
				}
				dirs[r.i][r.j][r.dir] = true
				switch maze[r.i][r.j] {
				case '.':
					{
						di, dj := d(r.dir)
						rays[k] = Loc{i: r.i + di, j: r.j + dj, dir: r.dir}
						//fmt.Println("move", r.dir)
					}
				case '-':
					{
						if r.dir == "W" || r.dir == "E" {
							di, dj := d(r.dir)
							rays[k] = Loc{i: r.i + di, j: r.j + dj, dir: r.dir}
						} else {
							di, dj := d("W")
							ray_index++
							new_ray := Loc{i: r.i + di, j: r.j + dj, dir: "W"}
							new_rays[ray_index] = new_ray
							di, dj = d("E")
							rays[k] = Loc{i: r.i + di, j: r.j + dj, dir: "E"}
						}
					}
				case '|':
					{
						if r.dir == "N" || r.dir == "S" {
							di, dj := d(r.dir)
							rays[k] = Loc{i: r.i + di, j: r.j + dj, dir: r.dir}
						} else {
							di, dj := d("N")
							ray_index++
							new_ray := Loc{i: r.i + di, j: r.j + dj, dir: "N"}
							new_rays[ray_index] = new_ray
							di, dj = d("S")
							rays[k] = Loc{i: r.i + di, j: r.j + dj, dir: "S"}
						}
					}
				case '/':
					{
						new_dir := ""
						switch r.dir {
						case "W":
							{
								new_dir = "S"
							}
						case "S":
							{
								new_dir = "W"
							}
						case "E":
							{
								new_dir = "N"
							}
						case "N":
							{
								new_dir = "E"
							}
						}
						di, dj := d(new_dir)
						rays[k] = Loc{i: r.i + di, j: r.j + dj, dir: new_dir}
					}
				case '\\':
					{
						new_dir := ""
						switch r.dir {
						case "W":
							{
								new_dir = "N"
							}
						case "N":
							{
								new_dir = "W"
							}
						case "E":
							{
								new_dir = "S"
							}
						case "S":
							{
								new_dir = "E"
							}
						}
						di, dj := d(new_dir)
						rays[k] = Loc{i: r.i + di, j: r.j + dj, dir: new_dir}
					}
				}
			}
			for k, v := range new_rays {
				rays[k] = v
			}
		}
		// solve part 1 here
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				//fmt.Print(len(dirs[i][j]))
				if len(dirs[i][j]) > 0 {
					sum++
				}
			}
			//fmt.Println()
		}
		return sum
	}
	if part2 {
		max := 0

		for i := 0; i < n; i++ {
			dir := "E"
			rays = make(map[int]Loc)
			ray_index = 0
			rays[ray_index] = Loc{i: i, j: 0, dir: dir}
			sum = init_and_run()
			if sum > max {
				max = sum
			}
			dir = "W"
			rays = make(map[int]Loc)
			ray_index = 0
			rays[ray_index] = Loc{i: i, j: m - 1, dir: dir}
			sum = init_and_run()
			if sum > max {
				max = sum
			}
		}
		for j := 0; j < m; j++ {
			dir := "S"
			rays = make(map[int]Loc)
			ray_index = 0
			rays[ray_index] = Loc{i: 0, j: j, dir: dir}
			sum = init_and_run()
			if sum > max {
				max = sum
			}
			dir = "N"
			rays = make(map[int]Loc)
			ray_index = 0
			rays[ray_index] = Loc{i: n - 1, j: j, dir: dir}
			sum = init_and_run()
			if sum > max {
				max = sum
			}
		}
		return max
	}
	rays[ray_index] = Loc{i: 0, j: 0, dir: "E"}
	sum = init_and_run()
	return sum
}

func d(s string) (int, int) {
	switch s {
	case "N":
		{
			return -1, 0
		}
	case "S":
		{
			return 1, 0
		}
	case "W":
		{
			return 0, -1
		}
	case "E":
		{
			return 0, 1
		}
	}
	return 1, 1
}
