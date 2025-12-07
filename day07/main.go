package main

import (
	"bufio"
	"fmt"
	"os"
)

type ManifoldItem byte

const (
	Empty ManifoldItem = iota
	Start
	Splitter
)

type manifold struct {
	width  int
	height int
	data   [][]ManifoldItem
}

type pos struct {
	row int
	col int
}

type cache map[pos]int

func check(err error, msg string) {
	if err != nil {
		panic(msg)
	}
}

func readManifoldDiagram(path string) *manifold {
	file, err := os.Open(path)
	check(err, "Failed to open file")
	defer file.Close()

	m := manifold{width: 0, height: 0, data: nil}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		var row []ManifoldItem
		for _, ch := range line {
			var item ManifoldItem
			switch ch {
			case '.':
				item = Empty
			case '^':
				item = Splitter
			case 'S':
				item = Start
			default:
				panic("Unknown manifold item")
			}
			row = append(row, item)
		}
		m.data = append(m.data, row)
	}

	m.width = len(m.data[0])
	m.height = len(m.data)

	return &m
}

var Cache cache

func (m *manifold) sendBeam(row int, col int) int {
	cache_count, ok := Cache[pos{row: row, col: col}]
	if ok {
		return cache_count
	}
	if row == m.height {
		return 1
	}

	item := m.data[row][col]
	var count int
	if item == Splitter {
		count = m.sendBeam(row+1, col-1) + m.sendBeam(row+1, col+1)
	} else {
		count = m.sendBeam(row+1, col)
	}

	Cache[pos{row: row, col: col}] = count
	return count
}

func main() {
	fmt.Println("Advent of Code / Day 7")

	m := readManifoldDiagram("./input/input1.txt")

	var start int = -1
	for c, i := range m.data[0] {
		if i == Start {
			start = c
		}
	}
	if start == -1 {
		panic("No start found")
	}

	beams := make(map[int]bool)
	beams[start] = true
	splits := 0
	for r, row := range m.data {
		if r == 0 {
			continue
		}

		next := beams
		for c, i := range row {
			if i == Splitter && beams[c] {
				next[c+1] = true
				next[c] = false
				next[c-1] = true

				splits += 1
			}
		}
		beams = next
	}

	fmt.Printf("1) Beam split %d times\n", splits)

	Cache = make(cache)
	fmt.Printf("2) Beam split %d times\n", m.sendBeam(1, start))
}
