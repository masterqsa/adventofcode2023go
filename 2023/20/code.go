package main

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/jpillora/puzzler/harness/aoc"
)

const (
	Broadcast = iota
	FlipFlop
	Conjunction
)

type Unit struct {
	name  string
	t     int
	up    map[string]bool
	down  []string
	state bool
}

func (u *Unit) ProcessPulse(origin string, pitch bool) {
	switch u.t {
	case Broadcast:
		{
			for _, d := range u.down {
				queue.Enqueue(Pulse{origin: u.name, dest: d, pitch: pitch})
			}
		}
	case FlipFlop:
		{
			if !pitch {
				u.state = !u.state
				for _, d := range u.down {
					queue.Enqueue(Pulse{origin: u.name, dest: d, pitch: u.state})
				}
			}
		}
	case Conjunction:
		{
			u.up[origin] = pitch
			all_high := true
			for _, state := range u.up {
				all_high = all_high && state
			}
			for _, d := range u.down {
				queue.Enqueue(Pulse{origin: u.name, dest: d, pitch: !all_high})
			}
		}
	}
}

type Pulse struct {
	origin string
	dest   string
	pitch  bool // low = false, high = true
}

type Queue struct {
	elements   []Pulse
	low_count  int
	high_count int
}

// Enqueue adds an element to the end of the queue
func (q *Queue) Enqueue(element Pulse) {
	q.elements = append(q.elements, element)
	if element.pitch {
		q.high_count++
	} else {
		q.low_count++
	}
}

// Dequeue removes and returns the element from the front of the queue
func (q *Queue) Dequeue() (Pulse, bool) {
	if len(q.elements) == 0 {
		return Pulse{"", "", false}, false
	}
	frontElement := q.elements[0]
	q.elements = q.elements[1:]
	return frontElement, true
}

// IsEmpty returns true if the queue is empty
func (q *Queue) IsEmpty() bool {
	return len(q.elements) == 0
}

// Size returns the number of elements in the queue
func (q *Queue) Size() int {
	return len(q.elements)
}

var queue Queue

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
	queue.high_count = 0
	queue.low_count = 0
	start := "broadcaster"
	units := make(map[string]*Unit)

	ups := make(map[string][]string)
	input = strings.ReplaceAll(input, "\r", "")
	r := regexp.MustCompile(`([%&]*)([a-z]+) -> ([a-z, ]+)`)
	for _, s := range strings.Split(input, "\n") {
		matches := r.FindStringSubmatch(s)
		d := strings.Split(matches[3], ", ")
		var t int
		switch matches[1] {
		case "":
			{
				t = Broadcast
			}
		case "%":
			{
				t = FlipFlop
			}
		case "&":
			{
				t = Conjunction
			}
		}
		units[matches[2]] = &Unit{name: matches[2], t: t, up: make(map[string]bool), down: d, state: false}
		for _, u := range d {
			list, ok := ups[u]
			if ok {
				ups[u] = append(list, matches[2])
			} else {
				ups[u] = []string{matches[2]}
			}
		}
	}
	keys := make([]string, 0, len(ups))
	for key := range ups {
		keys = append(keys, key)
	}

	// Sort the keys
	sort.Strings(keys)
	for _, k := range keys {
		unit, ok := units[k]
		if ok {
			unit.up = make(map[string]bool)
			for _, u := range ups[k] {
				unit.up[u] = false
			}
			fmt.Println(k, ups[k])
		} else {
			fmt.Println("Destination no found", k)
		}
	}

	limit := 1000
	if part2 {
		limit = 16000000
	}
	for i := 0; i < limit; i++ {
		queue.Enqueue(Pulse{"", start, false})
		for !queue.IsEmpty() {
			pulse, _ := queue.Dequeue()
			if i%100000 == 0 {
				fmt.Println(i)
			}
			if part2 && pulse.origin == "mr" && pulse.pitch && units["qt"].up["kk"] {
				fmt.Println("mr && kk", i)
			}
			if part2 && pulse.origin == "kk" && pulse.pitch && units["qt"].up["mr"] {
				fmt.Println("kk && mr", i)
			}
			unit, ok := units[pulse.dest]
			if ok {
				unit.ProcessPulse(pulse.origin, pulse.pitch)
			}
		}
	}

	return queue.high_count * queue.low_count
}
