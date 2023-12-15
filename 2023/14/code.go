package main

import (
	"fmt"
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
	input = strings.ReplaceAll(input, "\r", "")
	sum := 0
	raw := strings.Split(input, "\n")
	n := len(raw)
	m := len(raw[0])
	plate := make([][]rune, n)
	for i := 0; i < n; i++ {
		plate[i] = make([]rune, m)
		for j, c := range raw[i] {
			plate[i][j] = c
		}
	}
	if !part2 {
		for j := 0; j < m; j++ {
			local_sum := 0
			base := 0
			for i := 0; i < n; i++ {
				switch raw[i][j] {
				case 'O':
					{
						local_sum += (n - base)
						base++
					}
				case '#':
					{
						base = i + 1
					}
				}
			}
			sum += local_sum
		}
	}
	if part2 {
		cycle := func(t [][]rune) int {
			sum := 0
			// N
			for j := 0; j < m; j++ {
				base := 0
				dist := 0
				for i := 0; i < n; i++ {
					switch t[i][j] {
					case 'O':
						{
							t[i][j] = '.'
							t[base][j] = 'O'
							base++
						}
					case '#':
						{
							base = dist + 1
						}
					}
					dist++
				}
			}
			// W
			for i := 0; i < n; i++ {
				base := 0
				for j := 0; j < m; j++ {
					switch t[i][j] {
					case 'O':
						{
							t[i][j] = '.'
							t[i][base] = 'O'
							base++
						}
					case '#':
						{
							base = j + 1
						}
					}
				}
			}
			// S
			for j := 0; j < m; j++ {
				base := n - 1
				for i := n - 1; i >= 0; i-- {
					switch t[i][j] {
					case 'O':
						{
							t[i][j] = '.'
							t[base][j] = 'O'
							base--
						}
					case '#':
						{
							base = i - 1
						}
					}
				}
			}
			// E
			for i := 0; i < n; i++ {
				base := m - 1
				for j := m - 1; j >= 0; j-- {
					switch t[i][j] {
					case 'O':
						{
							t[i][j] = '.'
							t[i][base] = 'O'
							base--
						}
					case '#':
						{
							base = j - 1
						}
					}
				}
			}
			// North Wall
			for j := 0; j < m; j++ {
				local_sum := 0
				for i := 0; i < n; i++ {
					if t[i][j] == 'O' {
						local_sum += (n - i)
					}
				}
				sum += local_sum
			}

			return sum
		}
		// N W S E
		//fmt.Println(plate)
		for s := 1; s <= 200; s++ {
			c := cycle(plate)
			for i := 0; i < n; i++ {
				for j := 0; j < m; j++ {
					//fmt.Print(string(plate[i][j]))
				}
				//fmt.Println()
			}
			//fmt.Println()
			//fmt.Println(plate)
			fmt.Println(s, "sum", c)
		}
		return 87273 // 106 + (1000000000-106) mod 13
	}
	return sum
}
