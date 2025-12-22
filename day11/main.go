package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
)

type Graph map[string][]string

func (G Graph) Search(v string, goal string, E map[string]int) int {
	count := 0
	for w := range slices.Values(G[v]) {
		if w == goal {
			count += 1
			continue
		}

		if c, visited := E[w]; visited {
			count += c
		} else {
			count += G.Search(w, goal, E)
		}
	}

	E[v] = count
	return count
}

func check(err error, msg string) {
	if err != nil {
		panic(msg)
	}
}

func part1(path string) {
	file, err := os.Open(path)
	check(err, "Failed to open file")
	defer file.Close()

	scanner := bufio.NewScanner(file)
	G := make(Graph)
	for scanner.Scan() {
		line := scanner.Text()
		elements := strings.Split(line, ": ")
		src := elements[0]
		dst := strings.Split(elements[1], " ")
		G[src] = append(G[src], dst...)
	}

	E := make(map[string]int)
	count := G.Search("you", "out", E)
	fmt.Println("1) Paths found", count)
}

func part2(path string) {
	file, err := os.Open(path)
	check(err, "Failed to open file")
	defer file.Close()

	scanner := bufio.NewScanner(file)
	G := make(Graph)
	for scanner.Scan() {
		line := scanner.Text()
		elements := strings.Split(line, ": ")
		src := elements[0]
		dst := strings.Split(elements[1], " ")
		G[src] = append(G[src], dst...)
	}

	E := make(map[string]int)
	a := G.Search("svr", "fft", E)
	E = make(map[string]int)
	b := G.Search("fft", "dac", E)
	E = make(map[string]int)
	c := G.Search("dac", "out", E)
	fmt.Println("2) Paths found", a*b*c)
}

func main() {
	fmt.Println("Advent of Code / Day 11")
	part1("./input/input1.txt")
	part2("./input/input1.txt")
}
