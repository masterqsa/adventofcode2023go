package main

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/jpillora/puzzler/harness/aoc"
)

const (
	Goto = iota
	Accept
	Reject
	Gt
	Lt
)

const (
	X = iota
	M
	A
	S
)

var rules map[string]Rule

type Action struct {
	t      int
	param1 int
	param2 int
	action *Action
	val    string
}

type Rule struct {
	actions []Action
}

type Part struct {
	measurements []int
}

type Interval struct {
	Start, End int
}

type FourDCube struct {
	dims []Interval
}

func Measure(part Part) int {
	sum := 0
	for _, v := range part.measurements {
		sum += v
	}
	return sum
}

func MeasureCube(cube FourDCube) int {
	vol := 1
	for _, v := range cube.dims {
		vol *= (v.End - v.Start + 1)
	}
	return vol
}

func (part Part) ApplyAction(action Action) int {
	switch action.t {
	case Accept:
		{
			return Measure(part)
		}
	case Reject:
		{
			return 0
		}
	case Goto:
		{
			return part.ApplyRule(rules[action.val])
		}
	case Lt:
		{
			if part.measurements[action.param1] < action.param2 {
				return part.ApplyAction(*action.action)
			} else {
				return -1
			}
		}
	case Gt:
		{
			if part.measurements[action.param1] > action.param2 {
				return part.ApplyAction(*action.action)
			} else {
				return -1
			}
		}
	}
	return -1
}

func SliceCube(cube FourDCube, dim int, val int) []FourDCube {
	var cube1 FourDCube
	var cube2 FourDCube
	for d := 0; d < 4; d++ {
		if d != dim {
			cube1.dims = append(cube1.dims, cube.dims[d])
			cube2.dims = append(cube2.dims, cube.dims[d])
		} else {
			cube1.dims = append(cube1.dims, Interval{cube.dims[d].Start, val - 1})
			cube2.dims = append(cube2.dims, Interval{val, cube.dims[d].End})
		}
	}
	return []FourDCube{cube1, cube2}
}

func (cube FourDCube) ApplyActionCube(action Action) (int, []FourDCube) {
	switch action.t {
	case Accept:
		{
			return MeasureCube(cube), nil
		}
	case Reject:
		{
			return 0, nil
		}
	case Goto:
		{
			return cube.ApplyRuleCube(rules[action.val]), nil
		}
	case Lt:
		{
			if cube.dims[action.param1].End < action.param2 {
				return cube.ApplyActionCube(*action.action)
			} else if cube.dims[action.param1].Start >= action.param2 {
				return 0, []FourDCube{cube}
			} else {
				cubes := SliceCube(cube, action.param1, action.param2)
				res, _ := cubes[0].ApplyActionCube(*action.action) // no splitting in the secondary action
				return res, []FourDCube{cubes[1]}
			}
		}
	case Gt:
		{
			if cube.dims[action.param1].Start > action.param2 {
				return cube.ApplyActionCube(*action.action)
			} else if cube.dims[action.param1].End <= action.param2 {
				return 0, []FourDCube{cube}
			} else {
				cubes := SliceCube(cube, action.param1, action.param2+1)
				res, _ := cubes[1].ApplyActionCube(*action.action) // no splitting in the secondary action
				return res, []FourDCube{cubes[0]}
			}
		}
	}
	return -1, nil
}

func (cube FourDCube) ApplyRuleCube(rule Rule) int {
	result := 0
	cubes := []FourDCube{cube}
	for _, action := range rule.actions {
		next_cubes := make([]FourDCube, 0)
		for _, c := range cubes {
			local_result, cs := c.ApplyActionCube(action)
			result += local_result
			next_cubes = append(next_cubes, cs...)
		}
		cubes = next_cubes
	}
	return result
}

func (part Part) ApplyRule(rule Rule) int {
	result := 0
	for _, action := range rule.actions {
		result = part.ApplyAction(action)
		if result >= 0 {
			return result
		}
	}
	return result
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
	xmas := map[string]int{"x": X, "m": M, "a": A, "s": S}
	gtlt := map[string]int{"<": Lt, ">": Gt}
	input = strings.ReplaceAll(input, "\r", "")
	sum := 0
	raw := strings.Split(input, "\n\n")
	rules = make(map[string]Rule)
	r := regexp.MustCompile(`([a-z]+){([ARa-z<>0-9,:]+)`)
	rr := regexp.MustCompile(`([xmas]+)([<>]+)([0-9]+):([ARa-z]+)`)
	for _, s := range strings.Split(raw[0], "\n") {
		match := r.FindStringSubmatch(s)
		name := match[1]
		rule := Rule{actions: make([]Action, 0)}
		for _, p := range strings.Split(match[2], ",") {
			cond := strings.ContainsAny(p, ":")
			if cond {
				m := rr.FindStringSubmatch(p)
				param2, _ := strconv.Atoi(m[3])
				a := getSimpleAction(m[4])
				rule.actions = append(rule.actions, Action{t: gtlt[m[2]], param1: xmas[m[1]], param2: param2, action: &a})
			} else {
				rule.actions = append(rule.actions, getSimpleAction(p))
			}
		}
		//fmt.Println(rule)
		rules[name] = rule
	}

	if part2 {
		cube := FourDCube{dims: []Interval{Interval{1, 4000}, Interval{1, 4000}, Interval{1, 4000}, Interval{1, 4000}}}
		return cube.ApplyRuleCube(rules["in"])
	}

	rp := regexp.MustCompile(`{x=([0-9]+),m=([0-9]+),a=([0-9]+),s=([0-9]+)}`)
	for _, s := range strings.Split(raw[1], "\n") {
		m := rp.FindStringSubmatch(s)
		params := make([]int, 4)
		for i := 1; i <= 4; i++ {
			v, _ := strconv.Atoi(m[i])
			params[i-1] = v
		}
		part := Part{measurements: params}
		sum += part.ApplyRule(rules["in"])
	}

	return sum
}

func getSimpleAction(p string) Action {
	if p == "A" {
		return Action{t: Accept}
	} else if p == "R" {
		return Action{t: Reject}
	}
	return Action{t: Goto, val: p}
}
