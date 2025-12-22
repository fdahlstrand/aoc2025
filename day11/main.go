package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
)

type Graph map[string][]string

func (G Graph) Reverse() Graph {
	Grev := make(Graph)
	for v, _ := range G {
		for e := range slices.Values(G[v]) {
			Grev[e] = append(Grev[e], v)
		}
	}

	return Grev
}

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

func main() {
	fmt.Println("Advent of Code / Day 11")

	file, err := os.Open("./input/input1.txt")
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

	G = G.Reverse()
	E := make(map[string]int)
	count := G.Search("out", "you", E)
	fmt.Println("1) Paths found", count)
}
