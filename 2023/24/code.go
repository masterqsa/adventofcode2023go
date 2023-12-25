package main

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/jpillora/puzzler/harness/aoc"
)

type Vector struct {
	X, Y, Z int
}

type Line struct {
	Point     Vector // A point on the line
	Direction Vector // Direction vector of the line
}

func main() {
	aoc.Harness(run)
}
func atoi(s string) int {
	value, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return value
}

var minv int = -1000
var maxv int = 1000

// on code change, run will be executed 4 times:
// 1. with: false (part1), and example input
// 2. with: true (part2), and example input
// 3. with: false (part1), and user input
// 4. with: true (part2), and user input
// the return value of each run is printed to stdout
func run(part2 bool, input string) any {
	low := 200000000000000.0
	high := 400000000000000.0
	//low = 7
	//high = 27
	input = strings.ReplaceAll(input, "\r", "")
	raw := strings.Split(input, "\n")
	n := len(raw)
	rays := make([][]int, 0)
	r := regexp.MustCompile(`-?\d+`)
	for _, s := range raw {
		m := r.FindAllString(s, -1)
		rays = append(rays, []int{atoi(m[0]), atoi(m[1]), atoi(m[2]), atoi(m[3]), atoi(m[4]), atoi(m[5])})
		//fmt.Println(rays[len(rays)-1])
	}
	sum := 0
	for i := 0; i < n-1; i++ {
		for j := i + 1; j < n; j++ {
			a1 := rays[i][0]
			b1 := rays[i][3]
			c1 := rays[i][1]
			d1 := rays[i][4]
			a2 := rays[j][0]
			b2 := rays[j][3]
			c2 := rays[j][1]
			d2 := rays[j][4]
			if b2*d1-d2*b1 == 0 {
				// parallel
				//fmt.Println(i, j, "parallel")
				continue
			}
			t1 := float64((b2*(c1-c2) - d2*(a1-a2))) / float64(b1*d2-d1*b2)
			t2 := float64((b1*(c2-c1) - d1*(a2-a1))) / float64(b2*d1-d2*b1)
			if t2 <= 0 || t1 <= 0 {
				// past
				//fmt.Println(i, j, "past")
				continue
			}
			x := float64(a2) + float64(b2)*t2
			if x < low || x > high {
				// out of test area
				//fmt.Println(i, j, "x out", x)
				continue
			}
			y := float64(c2) + float64(d2)*t2
			if y < low || y > high {
				// out of test area
				//fmt.Println(i, j, "y out", y)
				continue
			}
			//fmt.Println(i, j, "intersect")
			sum++
		}
	}
	if part2 {
		xVelocities, yVelocities, zVelocities := make([]int, 0), make([]int, 0), make([]int, 0)
		for i := 0; i < n-1; i++ {
			for j := i + 1; j < n; j++ {
				a1 := rays[i][0]
				b1 := rays[i][3]
				c1 := rays[i][1]
				d1 := rays[i][4]
				e1 := rays[i][2]
				f1 := rays[i][5]
				a2 := rays[j][0]
				b2 := rays[j][3]
				c2 := rays[j][1]
				d2 := rays[j][4]
				e2 := rays[j][2]
				f2 := rays[j][5]
				v1 := Vector{b1, d1, f1}
				v2 := Vector{b2, d2, f2}
				p1 := Vector{a1, c1, e1}
				p2 := Vector{a2, c2, e2}
				// checking if any lines are parallel or coplanar, as it might simplify the task
				if areLinesParallel(b1, d1, f1, b2, d2, f2) {
					fmt.Println(i, j, "parallel", a1, c1, e1, a2, c2, e2)
				}
				if areLinesCoPlanar(v1, v2, p1, p2) {
					fmt.Println(i, j, "co-planar")
				}
				if v1.X == v2.X {
					valid := validVelocities(p2.X-p1.X, v1.X)
					if len(xVelocities) == 0 {
						xVelocities = valid
					} else {
						xVelocities = intersect(xVelocities, valid)
					}
				}
				if v1.Y == v2.Y {
					valid := validVelocities(p2.Y-p1.Y, v1.Y)
					if len(yVelocities) == 0 {
						yVelocities = valid
					} else {
						yVelocities = intersect(yVelocities, valid)
					}
				}
				if v1.Z == v2.Z {
					valid := validVelocities(p2.Z-p1.Z, v1.Z)
					if len(xVelocities) == 0 {
						zVelocities = valid
					} else {
						zVelocities = intersect(zVelocities, valid)
					}
				}
			}
		}
		fmt.Println(xVelocities, yVelocities, zVelocities)
		v := Vector{xVelocities[0], yVelocities[0], zVelocities[0]}
		// convenience code just to populate v1, v2, p1, p2
		a1 := rays[0][0]
		b1 := rays[0][3]
		c1 := rays[0][1]
		d1 := rays[0][4]
		e1 := rays[0][2]
		f1 := rays[0][5]
		a2 := rays[1][0]
		b2 := rays[1][3]
		c2 := rays[1][1]
		d2 := rays[1][4]
		e2 := rays[1][2]
		f2 := rays[1][5]
		v1 := Vector{b1, d1, f1}
		v2 := Vector{b2, d2, f2}
		p1 := Vector{a1, c1, e1}
		p2 := Vector{a2, c2, e2}
		a := v1.X - v.X
		b := v1.Y - v.Y
		c := v2.X - v.X
		d := v2.Y - v.Y
		t2 := (a*(p2.Y-p1.Y) - b*(p2.X-p1.X)) / (b*c - a*d)
		x := p2.X + v2.X*t2 - v.X*t2
		y := p2.Y + v2.Y*t2 - v.Y*t2
		z := p2.Z + v2.Z*t2 - v.Z*t2
		return x + y + z
	}
	// solve part 1 here
	return sum
}

func validVelocities(dist int, velocity int) []int {
	result := make([]int, 0)
	for v := minv; v <= maxv; v++ {
		if v != velocity && dist%(v-velocity) == 0 {
			result = append(result, v)
		}
	}
	return result
}

func intersect(ar1 []int, ar2 []int) []int {
	result := make([]int, 0)
	for _, v := range ar1 {
		if slices.Contains(ar2, v) {
			result = append(result, v)
		}
	}
	return result
}

func subtractVectors(v1, v2 Vector) Vector {
	return Vector{v1.X - v2.X, v1.Y - v2.Y, v1.Z - v2.Z}
}

// Calculate the cross product of two vectors
func crossProduct(v1, v2 Vector) Vector {
	return Vector{X: v1.Y*v2.Z - v1.Z*v2.Y, Y: v1.Z*v2.X - v1.X*v2.Z, Z: v1.X*v2.Y - v1.Y*v2.X}
}

func dotProduct(v1, v2 Vector) int {
	return v1.X*v2.X + v1.Y*v2.Y + v1.Z*v2.Z
}

// Handles division in the context of integers
func divideIntegers(numerator, denominator int) int {
	// This is a simple way to handle division; you might need a more sophisticated approach
	// depending on specific requirements.
	return numerator / denominator
}

func shortestSegmentBetweenLines(line1, line2 Line) (Vector, Vector) {
	p1 := line1.Point
	d1 := line1.Direction
	p2 := line2.Point
	d2 := line2.Direction

	n := crossProduct(d1, d2)
	n1 := crossProduct(d2, n)
	n2 := crossProduct(d1, n)

	c1 := divideIntegers(dotProduct(subtractVectors(p2, p1), n1), dotProduct(d1, n1))
	c2 := divideIntegers(dotProduct(subtractVectors(p1, p2), n2), dotProduct(d2, n2))

	closestPointOnLine1 := Vector{p1.X + c1*d1.X, p1.Y + c1*d1.Y, p1.Z + c1*d1.Z}
	closestPointOnLine2 := Vector{p2.X + c2*d2.X, p2.Y + c2*d2.Y, p2.Z + c2*d2.Z}

	return closestPointOnLine1, closestPointOnLine2
}

// Function to calculate the midpoint of a line segment
func midpoint(p1, p2 Vector) Vector {
	return Vector{
		X: (p1.X + p2.X) / 2,
		Y: (p1.Y + p2.Y) / 2,
		Z: (p1.Z + p2.Z) / 2,
	}
}

// Calculate the determinant of three vectors
func determinant(v1, v2, v3 Vector) int {
	return v1.X*(v2.Y*v3.Z-v2.Z*v3.Y) - v1.Y*(v2.X*v3.Z-v2.Z*v3.X) + v1.Z*(v2.X*v3.Y-v2.Y*v3.X)
}

func areLinesParallel(b1, d1, f1, b2, d2, f2 int) bool {
	// Check for cases where any of the direction components for the second line are 0
	if b2 == 0 || d2 == 0 || f2 == 0 {
		return (b1 == 0 && b2 == 0) && (d1 == 0 && d2 == 0) && (f1 == 0 && f2 == 0)
	}

	k1 := float64(b1) / float64(b2)
	k2 := float64(d1) / float64(d2)
	k3 := float64(f1) / float64(f2)

	// If k1, k2, and k3 are equal, the lines are parallel
	return k1 == k2 && k2 == k3
}

// Check if two vectors are scalar multiples of each other (parallel)
func areParallel(v1, v2 Vector) bool {
	if v1.X != 0 && v2.X != 0 {
		ratio := float64(v1.X) / float64(v2.X)
		return float64(v1.Y)/ratio == float64(v2.Y) && float64(v1.Z)/ratio == float64(v2.Z)
	}
	if v1.Y != 0 && v2.Y != 0 {
		ratio := float64(v1.Y) / float64(v2.Y)
		return float64(v1.X)/ratio == float64(v2.X) && float64(v1.Z)/ratio == float64(v2.Z)
	}
	if v1.Z != 0 && v2.Z != 0 {
		ratio := float64(v1.Z) / float64(v2.Z)
		return float64(v1.X)/ratio == float64(v2.X) && float64(v1.Y)/ratio == float64(v2.Y)
	}
	return v1 == (Vector{0, 0, 0}) && v2 == (Vector{0, 0, 0})
}

// Check if two lines are co-planar
func areLinesCoPlanar(v1, v2, p1, p2 Vector) bool {
	if areParallel(v1, v2) {
		// Parallel lines are always co-planar
		return true
	}
	connectingVector := Vector{X: p2.X - p1.X, Y: p2.Y - p1.Y, Z: p2.Z - p1.Z}
	return determinant(v1, v2, connectingVector) == 0
}
