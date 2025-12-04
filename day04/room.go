package main

import (
	"bufio"
	"fmt"
	"os"
)

type RoomItem byte

const (
	Empty RoomItem = iota
	PaperRoll
)

type room struct {
	width  int
	height int
	data   [][]RoomItem
}

type pos struct {
	row int
	col int
}

func check(err error, msg string) {
	if err != nil {
		panic(msg)
	}
}

func loadRoom(path string) *room {
	file, err := os.Open(path)
	check(err, "Failed to open file")

	scanner := bufio.NewScanner(file)

	r := room{width: 0, height: 0, data: nil}
	row := 0
	for scanner.Scan() {
		line := scanner.Text()
		r.data = append(r.data, nil)
		for _, ch := range line {
			switch ch {
			case '.':
				r.data[row] = append(r.data[row], Empty)
			case '@':
				r.data[row] = append(r.data[row], PaperRoll)
			default:
				fmt.Printf("> '%c'\n", ch)
				panic("Unknown room item")
			}
		}

		row += 1
	}

	r.height = len(r.data)
	r.width = len(r.data[0])

	return &r
}

func (r *room) countAtPos(row int, col int) int {
	if row < 0 || row >= r.height || col < 0 || col >= r.width {
		return 0
	}

	if r.data[row][col] == PaperRoll {
		return 1
	}

	return 0
}

func (r *room) countSurroundingAt(row int, col int) int {
	return r.countAtPos(row-1, col-1) +
		r.countAtPos(row-1, col) +
		r.countAtPos(row-1, col+1) +
		r.countAtPos(row, col-1) +
		r.countAtPos(row, col+1) +
		r.countAtPos(row+1, col-1) +
		r.countAtPos(row+1, col) +
		r.countAtPos(row+1, col+1)
}

func (r *room) isPosAccessible(row int, col int) bool {
	return r.countSurroundingAt(row, col) < 4
}

func (r *room) getAccessible() []pos {
	var a []pos = nil
	for row := range r.height {
		for col := range r.width {
			if r.data[row][col] == PaperRoll && r.isPosAccessible(row, col) {
				a = append(a, pos{row, col})
			}
		}
	}
	return a
}

func (r *room) removePaperRoll(row int, col int) {
	if r.data[row][col] == PaperRoll {
		r.data[row][col] = Empty
	}
}
