package main

import (
	"fmt"
	"regexp"
	"strconv"
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
	// when you're ready to do part 2, remove this "not implemented" block
	if part2 {
		r_1 := regexp.MustCompile("one")
		r_2 := regexp.MustCompile("two")
		r_3 := regexp.MustCompile("three")
		r_4 := regexp.MustCompile("four")
		r_5 := regexp.MustCompile("five")
		r_6 := regexp.MustCompile("six")
		r_7 := regexp.MustCompile("seven")
		r_8 := regexp.MustCompile("eight")
		r_9 := regexp.MustCompile("nine")
		sum := 0
		for _, line := range strings.Split(strings.TrimSpace(input), "\n") {
			fmt.Print(line + " ")
			line = r_1.ReplaceAllString(line, "o1ne")
			line = r_2.ReplaceAllString(line, "t2wo")
			line = r_3.ReplaceAllString(line, "th3ree")
			line = r_4.ReplaceAllString(line, "fo4ur")
			line = r_5.ReplaceAllString(line, "fi5ve")
			line = r_6.ReplaceAllString(line, "s6ix")
			line = r_7.ReplaceAllString(line, "sev7en")
			line = r_8.ReplaceAllString(line, "ei8ght")
			line = r_9.ReplaceAllString(line, "ni9ne")

			r := regexp.MustCompile("[a-z\r]*")
			line = r.ReplaceAllString(line, "")
			fmt.Print(line + " " + strconv.Itoa(len(line)) + " " + line[len(line)-1:] + " ")
			v1, _ := strconv.Atoi(line[:1])
			v2, _ := strconv.Atoi(line[len(line)-1:])
			v := (v1*10 + v2)
			fmt.Println(v)
			sum += v
		}
		return sum
	}
	// solve part 1 here
	sum := 0
	for _, line := range strings.Split(strings.TrimSpace(input), "\n") {
		fmt.Print(line + " ")
		r := regexp.MustCompile("[a-z\r]*")
		line = r.ReplaceAllString(line, "")
		fmt.Print(line + " " + strconv.Itoa(len(line)) + " " + line[len(line)-1:] + " ")
		v1, _ := strconv.Atoi(line[:1])
		v2, _ := strconv.Atoi(line[len(line)-1:])
		v := (v1*10 + v2)
		fmt.Println(v)
		sum += v
	}
	return sum
}
