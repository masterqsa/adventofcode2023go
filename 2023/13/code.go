package main

import (
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
	raw := strings.Split(input, "\n\n")
	for _, s := range raw {
		split := strings.Split(s, "\n")
		rows := make([]string, 0)
		col_builders := make([]strings.Builder, len(split[0]))
		for _, l := range split {
			line := l
			rows = append(rows, line)
			for j, c := range line {
				col_builders[j].WriteRune(c)
			}
		}
		cols := make([]string, len(col_builders))
		for i, b := range col_builders {
			cols[i] = b.String()
		}
		//fmt.Println(rows)
		//fmt.Println(cols)
		horiz := findMirror(rows, part2)
		vert := findMirror(cols, part2)
		sum += (100*horiz + vert)
	}
	return sum
}

func findMirror(lines []string, part2 bool) int {
	cur := 0
	depth := 1
	diff := 0
	if !part2 {
		for i := 1; i < len(lines); i++ {
			if i-2*depth+1 < 0 {
				return cur
			}
			if lines[i] == lines[i-2*depth+1] {
				if depth == 1 {
					cur = i
				}
				depth++
			} else {
				cur = 0
				depth = 1
			}
		}
	} else {
		for cur = 1; cur < len(lines); cur++ {
			diff = 0
			for depth = 1; depth <= len(lines)/2; depth++ {
				if (cur+depth-1) > len(lines)-1 || cur-depth < 0 {
					if diff == 1 {
						return cur
					} else {
						break
					}
				}
				diff += compareStrings(lines[cur+depth-1], lines[cur-depth])
				if diff > 1 {
					break
				}
			}
			if diff == 1 {
				return cur
			}
		}
		return 0
	}
	//fmt.Println(cur)
	if part2 {
		if diff == 0 {
			return 0
		}
	}
	return cur
}

func compareStrings(s1 string, s2 string) int {
	diff := 0
	for i := 0; i < len(s1); i++ {
		if s1[i] != s2[i] {
			diff++
		}
	}
	return diff
}
