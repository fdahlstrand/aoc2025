package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Point struct {
	X int
	Y int
}

type Polygon []Point

func check(err error, msg string) {
	if err != nil {
		panic(msg)
	}
}

func abs(x int) int { return max(x, -x) }

func area(p, q Point) int {
	width := abs(p.X-q.X) + 1
	height := abs(p.Y-q.Y) + 1
	return width * height
}

func part1(polygon Polygon) {
	maxArea := 0
	for i := range len(polygon) {
		p := polygon[i]
		for j := i + 1; j < len(polygon); j++ {
			q := polygon[j]
			maxArea = max(maxArea, area(p, q))
		}
	}

	fmt.Printf("1) Largest area %d\n", maxArea)
}

func (P Polygon) IsPointInside(p Point) bool {
	count := 0
	n := len(P)
	for i := range n {
		p1, p2 := P[i], P[(i+1)%n]

		if p.X == p1.X && min(p1.Y, p2.Y) <= p.Y && p.Y <= max(p1.Y, p2.Y) {
			return true
		}

		if p.Y == p1.Y && min(p1.X, p2.X) <= p.X && p.X <= max(p1.X, p2.X) {
			return true
		}

		if p1.X == p2.X {
			if ((p1.Y > p.Y) != (p2.Y > p.Y)) && (p.X < p1.X) {
				count++
			}
		}
	}

	return (count % 2) == 1
}

func Intersects(h1, h2, v1, v2 Point) bool {
	hMinX, hMaxX := min(h1.X, h2.X), max(h1.X, h2.X)
	vMinY, vMaxY := min(v1.Y, v2.Y), max(v1.Y, v2.Y)

	return v1.X > hMinX && v1.X < hMaxX && h1.Y > vMinY && h1.Y < vMaxY
}

func part2(polygon Polygon) {
	var rects [][4][2]Point
	for i := range polygon {
		for j := i + 1; j < len(polygon); j++ {
			p1, p2 := polygon[i], polygon[j]
			q1 := Point{p1.X, p2.Y}
			q2 := Point{p2.X, p1.Y}

			if polygon.IsPointInside(q1) && polygon.IsPointInside(q2) {
				rects = append(rects, [4][2]Point{
					{p1, q1},
					{q1, p2},
					{p2, q2},
					{q2, p1}})
			}
		}
	}

	maxArea := 0
	for _, rc := range rects {
		n := len(polygon)
		valid := true
		for i := range n {
			p1, p2 := polygon[i], polygon[(i+1)%n]
			for j := range 4 {
				q1, q2 := rc[j][0], rc[j][1]

				if p1.X == p2.X && q1.Y == q2.Y {
					if Intersects(q1, q2, p1, p2) {
						valid = false
					}
				} else if p1.Y == p2.Y && q1.X == q2.X {
					if Intersects(p1, p2, q1, q2) {
						valid = false
					}
				}
			}
		}

		if valid {
			x1 := min(rc[0][0].X, rc[1][0].X, rc[2][0].X, rc[3][0].X)
			y1 := min(rc[0][0].Y, rc[1][0].Y, rc[2][0].Y, rc[3][0].Y)
			x2 := max(rc[0][0].X, rc[1][0].X, rc[2][0].X, rc[3][0].X)
			y2 := max(rc[0][0].Y, rc[1][0].Y, rc[2][0].Y, rc[3][0].Y)

			a := area(Point{x1, y1}, Point{x2, y2})
			maxArea = max(maxArea, a)
		}
	}

	fmt.Printf("2) Largest area %d\n", maxArea)
}

func main() {
	fmt.Println("Advent of Code / Day 9")

	file, err := os.Open("./input/input1.txt")
	check(err, "Failed to open file")
	defer file.Close()

	scanner := bufio.NewScanner(file)
	polygon := make([]Point, 0)
	for scanner.Scan() {
		line := scanner.Text()
		vals := strings.Split(line, ",")
		x, err := strconv.Atoi(vals[0])
		check(err, "Failed to convert X-coordinate")
		y, err := strconv.Atoi(vals[1])
		check(err, "Failed to convert Y-coordinate")
		polygon = append(polygon, Point{X: x, Y: y})
	}

	part1(polygon)
	part2(polygon)
}
