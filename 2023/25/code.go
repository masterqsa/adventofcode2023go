package main

import (
	"fmt"
	"strings"

	"github.com/jpillora/puzzler/harness/aoc"
)

func main() {
	aoc.Harness(run)
}

var counts map[[2]string]int = make(map[[2]string]int)
var edges [][2]string = make([][2]string, 0)

func getKey(s1, s2 string) [2]string {
	var key [2]string
	if s1 < s2 {
		key = [2]string{s1, s2}
	} else {
		key = [2]string{s2, s1}
	}
	return key
}

func addEdge(s1, s2 string) {
	edges = append(edges, getKey(s1, s2))
}

func incrementCount(s1, s2 string) {
	key := getKey(s1, s2)
	if v, ok := counts[key]; ok {
		counts[key] = v + 1
	} else {
		counts[key] = 1
	}
}

type UnionFind struct {
	parent map[string]string
	rank   map[string]int
	size   map[string]int
}

func NewUnionFind() *UnionFind {
	return &UnionFind{
		parent: make(map[string]string),
		rank:   make(map[string]int),
		size:   make(map[string]int),
	}
}

func (uf *UnionFind) Find(s string) string {
	if uf.parent[s] != s {
		uf.parent[s] = uf.Find(uf.parent[s])
	}
	return uf.parent[s]
}

func (uf *UnionFind) Union(s1, s2 string) {
	root1 := uf.Find(s1)
	root2 := uf.Find(s2)

	if root1 != root2 {
		if uf.rank[root1] < uf.rank[root2] {
			uf.parent[root1] = root2
			uf.size[root2] += uf.size[root1]
		} else if uf.rank[root1] > uf.rank[root2] {
			uf.parent[root2] = root1
			uf.size[root1] += uf.size[root2]
		} else {
			uf.parent[root2] = root1
			uf.rank[root1]++
			uf.size[root1] += uf.size[root2]
		}
	}
}

func (uf *UnionFind) MakeSet(s string) {
	if _, exists := uf.parent[s]; !exists {
		uf.parent[s] = s
		uf.rank[s] = 0
		uf.size[s] = 1
	}
}

func (uf *UnionFind) GetSetSize(s string) int {
	root := uf.Find(s)
	return uf.size[root]
}

func (uf *UnionFind) GetUniqueSets() map[string]int {
	sets := make(map[string]int)
	for element := range uf.parent {
		root := uf.Find(element)
		if _, ok := sets[root]; !ok {
			sets[root] = uf.size[root]
		}
	}
	return sets
}

// on code change, run will be executed 4 times:
// 1. with: false (part1), and example input
// 2. with: true (part2), and example input
// 3. with: false (part1), and user input
// 4. with: true (part2), and user input
// the return value of each run is printed to stdout
func run(part2 bool, input string) any {
	input = strings.ReplaceAll(input, "\r", "")
	input = strings.ReplaceAll(input, ":", "")
	conns := make(map[string][]string)
	unique := make(map[string]bool)
	raw := strings.Split(input, "\n")
	edges = make([][2]string, 0)
	for _, s := range raw {
		parts := strings.Fields(s)
		o := parts[0]
		conns[o] = parts[1:]
		unique[o] = true
		for _, n := range parts[1:] {
			addEdge(o, n)
			unique[n] = true
			if c, ok := conns[n]; ok {
				c = append(c, o)
				conns[n] = c
			} else {
				conns[n] = []string{o}
			}
		}
	}
	fmt.Println(len(unique))
	//fmt.Println(conns)
	nodes := make([]string, 0)
	for n, _ := range unique {
		nodes = append(nodes, n)
	}
	fmt.Println("Edges:", len(edges))
	fmt.Println("Nodes:", len(nodes))

	for i := 0; i < len(nodes)-1; i++ {
		for j := i + 1; j < len(nodes); j++ {
			start := nodes[i]
			end := nodes[j]
			removed := make(map[[2]string]bool)

			addToRemoved := func(s1, s2 string) {
				key := getKey(s1, s2)
				removed[key] = true
			}
			isRemoved := func(s1, s2 string) bool {
				_, ok := removed[getKey(s1, s2)]
				return ok
			}

			for r := 0; r < 4; r++ { // remove 3 paths between points and check if they are still connected 4th time after that
				visited := map[string]int{start: 0}
				dist := 0
				queue := []string{start}
				exit := false
				for !exit && len(queue) > 0 {
					next_queue := make([]string, 0)
					for _, cur := range queue {
						visited[cur] = dist
						if cur == end {
							exit = true
						} else {
							for _, cand := range conns[cur] {
								if isRemoved(cur, cand) {
									continue
								}
								if _, ok := visited[cand]; !ok {
									next_queue = append(next_queue, cand)
								}
							}
						}
					}
					dist++
					queue = next_queue
				}
				if exit && r < 3 { // found a path
					// back trace
					cur := end
					cur_dist := visited[end]
					for cur_dist > 0 {
						cur_dist--
						for _, n := range conns[cur] {
							if visited[n] == cur_dist {
								addToRemoved(cur, n)
								cur = n
								break
							}
						}
					}
				}
				if !exit { // found partition
					uf := NewUnionFind()

					// Initialize the sets
					for _, s := range nodes {
						uf.MakeSet(s)
					}
					for _, e := range edges {
						if !isRemoved(e[0], e[1]) {
							uf.Union(e[0], e[1])
						}
					}
					u := uf.GetUniqueSets()
					//fmt.Println("Sets:", u, len(u), edges[i], edges[j], edges[k])
					if len(u) == 2 {

						product := 1
						for _, v := range u {
							product *= v
						}
						return product
					}
				}
			}
		}
	}
	fmt.Println(counts)

	fmt.Println("~~~")

	if part2 {
		return "not implemented"
	}
	// solve part 1 here
	return 42
}
