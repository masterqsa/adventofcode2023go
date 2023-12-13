package main

import (
	"strconv"
	"strings"

	"github.com/jpillora/puzzler/harness/aoc"
)

type State struct {
	idx     int
	options int
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
	sum := 0
	// solve part 1 here
	raw := strings.Split(strings.TrimSpace(input), "\n")
	for _, line := range raw {
		parts := strings.Fields(strings.Replace(line, "\r", "", -1))
		p0 := parts[0]
		p1 := parts[1]
		if part2 {
			for i := 0; i < 4; i++ {
				p0 = p0 + "?" + parts[0]
				p1 = p1 + "," + parts[1]
			}
		}

		pattern := p0 + "."
		//fmt.Println(pattern)
		vals_s := strings.Split(p1, ",")
		var vals []int
		sm := 0
		for _, s := range vals_s {
			v, _ := strconv.Atoi(s)
			vals = append(vals, v)
			sm += v
		}

		q := 0
		set := 0
		brute := make(map[int](map[int](map[int]State))) // 1st dimension - position, 2nd - #s left to place, 3rd = running sum
		brute[-1] = make(map[int](map[int]State))

		for i, c := range pattern {
			brute[i] = make(map[int](map[int]State))
			for j := 0; j <= sm-set; j++ {
				brute[i][j] = make(map[int]State)
			}
			switch c {
			case '?':
				{
					q++
				}
			case '#':
				{
					set++
				}
			}
		}
		//fmt.Println(sm, q, set)
		brute[-1][sm-set] = make(map[int]State)
		brute[-1][sm-set][0] = State{idx: 0, options: 1} // on the 1st number of the list, current count of #s is 0, only one way to be here
		for i, c := range pattern {
			for key, v := range brute[i-1] {
				for val_count, val := range v {
					//fmt.Println(i, string(c), key, val_count, brute)
					switch c {
					case '?':
						{
							if (val_count == 0 || (val_count > 0 && vals[val.idx] >= val_count+1)) && key > 0 { // can or have to place #
								cu, o := ((brute[i])[key-1])[val_count+1]
								var state State
								if o {
									state = State{idx: val.idx, options: val.options + cu.options}
								} else {
									state = State{idx: val.idx, options: val.options}
								}
								brute[i][key-1][val_count+1] = state
								//fmt.Print("#")
							}
							cur, ok := (brute[i])[key][0]

							if val.idx < len(vals) && vals[val.idx] == val_count { // forced to place . advancing number set
								var state State
								if ok {
									state = State{idx: val.idx + 1, options: val.options + cur.options}
								} else {
									state = State{idx: val.idx + 1, options: val.options}
								}
								brute[i][key][0] = state
								//fmt.Print(".")
							}
							if val_count == 0 { // can place . not advancing
								var state State
								if ok {
									state = State{idx: val.idx, options: val.options + cur.options}
								} else {
									state = State{idx: val.idx, options: val.options}
								}
								brute[i][key][0] = state
								//fmt.Print(".")
							}
							//fmt.Println()
						}
					case '#':
						{
							cur, ok := (brute[i])[key][val_count+1]
							if vals[val.idx] >= val_count+1 { // valid #
								var state State
								if ok {
									state = State{idx: val.idx, options: val.options + cur.options}
								} else {
									state = State{idx: val.idx, options: val.options}
								}
								brute[i][key][val_count+1] = state
							}
						}
					case '.':
						{
							cur, ok := (brute[i])[key][0]
							if val_count > 0 && vals[val.idx] == val_count { // valid . advacing number set
								var state State
								if ok {
									state = State{idx: val.idx + 1, options: val.options + cur.options}
								} else {
									state = State{idx: val.idx + 1, options: val.options}
								}
								brute[i][key][0] = state
							} else if val_count == 0 { // valid . not advancing
								var state State
								if ok {
									state = State{idx: val.idx, options: val.options + cur.options}
								} else {
									state = State{idx: val.idx, options: val.options}
								}
								brute[i][key][0] = state
							}
						}
					}
				}
			}
		}
		//fmt.Println(brute[len(pattern)-1][0][0].options)
		sum += brute[len(pattern)-1][0][0].options
	}
	return sum
}
