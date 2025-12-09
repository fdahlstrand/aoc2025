package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Pos struct {
	X int
	Y int
}

func check(err error, msg string) {
	if err != nil {
		panic(msg)
	}
}

func abs(x int) int { return max(x, -x) }

func area(p, q Pos) int {
	width := abs(p.X-q.X) + 1
	height := abs(p.Y-q.Y) + 1
	return width * height
}

func main() {
	fmt.Println("Advent of Code / Day 9")

	file, err := os.Open("./input/input1.txt")
	check(err, "Failed to open file")
	defer file.Close()

	scanner := bufio.NewScanner(file)
	positions := make([]Pos, 0)
	for scanner.Scan() {
		line := scanner.Text()
		vals := strings.Split(line, ",")
		x, err := strconv.Atoi(vals[0])
		check(err, "Failed to convert X-coordinate")
		y, err := strconv.Atoi(vals[1])
		check(err, "Failed to convert Y-coordinate")
		positions = append(positions, Pos{X: x, Y: y})
	}

	maxArea := 0
	for i := range len(positions) {
		p := positions[i]
		for j := i + 1; j < len(positions); j++ {
			q := positions[j]
			maxArea = max(maxArea, area(p, q))
		}
	}

	fmt.Printf("1) Largest area %d\n", maxArea)
}
