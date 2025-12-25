package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Pos struct {
	X int
	Y int
}

type Piece []Pos

type Goal struct {
	Width  int
	Length int
	Pieces []int
}

func check(err error, msg string) {
	if err != nil {
		panic(msg)
	}
}

func parsePiece(elem string) Piece {
	lines := strings.Split(elem, "\n")
	var piece Piece
	for y, line := range lines[1:] {
		for x, ch := range line {
			if ch == '#' {
				piece = append(piece, Pos{x, y})
			}
		}

	}

	return piece
}

func parseGoal(line string) Goal {
	var err error
	goal := Goal{}
	elems := strings.Fields(line)

	elems[0] = strings.Trim(elems[0], ":")
	size := strings.Split(elems[0], "x")
	goal.Width, err = strconv.Atoi(size[0])
	check(err, "Can't parse goal width")
	goal.Length, err = strconv.Atoi(size[1])
	check(err, "Can't parse goal height")

	for _, number := range elems[1:] {
		var count int
		count, err = strconv.Atoi(number)
		goal.Pieces = append(goal.Pieces, count)
	}

	return goal
}

func readFile(path string) ([]Piece, []Goal) {
	bytes, err := os.ReadFile(path)
	check(err, "Failed to open file")
	data := string(bytes)
	elem := strings.Split(data, "\n\n")

	var pieces []Piece
	for _, p := range elem[0 : len(elem)-1] {
		pieces = append(pieces, parsePiece(p))
	}

	var goals []Goal
	for line := range strings.SplitSeq(elem[len(elem)-1], "\n") {
		if len(line) == 0 {
			continue
		}
		goals = append(goals, parseGoal(line))
	}

	return pieces, goals
}

func main() {
	fmt.Println("Advent of Code / Day 12")

	pieces, goals := readFile("./input/input1.txt")

	total := 0
	for _, g := range goals {
		remainingArea := g.Length * g.Width
		remainingPieces := (g.Length / 3) * (g.Width / 3)
		for id, count := range g.Pieces {
			size := len(pieces[id])
			remainingArea -= size * count
			remainingPieces -= count
		}

		if remainingArea > 0 && remainingPieces >= 0 {
			total++
		}

	}

	fmt.Printf("1) %d regions can fit all of the listed presents\n", total)
}
