package main

import (
	"bufio"
	"fmt"
	"maps"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"
)

type Node struct {
	X      int
	Y      int
	Z      int
	Parent *Node
}

type NodePair struct {
	A *Node
	B *Node
}

func check(err error, msg string) {
	if err != nil {
		panic(msg)
	}
}

func distance(a, b *Node) int {
	return (a.X-b.X)*(a.X-b.X) + (a.Y-b.Y)*(a.Y-b.Y) + (a.Z-b.Z)*(a.Z-b.Z)
}

func find(x *Node) *Node {
	if x.Parent == x {
		return x
	}

	x.Parent = find(x.Parent)
	return x.Parent
}

func union(a, b *Node) {
	a = find(a)
	b = find(b)

	if a == b {
		return
	}

	b.Parent = a
}

func readNodes(path string) []*Node {
	file, err := os.Open(path)
	check(err, "Failed to open file")
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var nodes []*Node
	for scanner.Scan() {
		line := scanner.Text()
		vals := strings.Split(line, ",")
		var v [3]int
		for i := range 3 {
			v[i], err = strconv.Atoi(vals[i])
		}
		n := Node{X: v[0], Y: v[1], Z: v[2]}
		n.Parent = &n
		nodes = append(nodes, &n)
	}

	return nodes
}

func collectPairs(nodes []*Node) []NodePair {
	var pairs []NodePair
	for i := range len(nodes) {
		for j := i + 1; j < len(nodes); j++ {
			pairs = append(pairs, NodePair{A: nodes[i], B: nodes[j]})

		}
	}

	sort.Slice(pairs, func(i, j int) bool {
		return distance(pairs[i].A, pairs[i].B) < distance(pairs[j].A, pairs[j].B)
	})

	return pairs
}

func connect(path string, boxCount int) []int {
	nodes := readNodes(path)
	pairs := collectPairs(nodes)

	for i := range boxCount {
		union(pairs[i].A, pairs[i].B)
	}

	circuits := make(map[*Node]int)
	for _, x := range nodes {
		circuits[find(x)]++
	}

	s := slices.Collect(maps.Values(circuits))
	sort.Slice(s, func(i, j int) bool { return s[i] > s[j] })

	return s
}

func maxCircuit(path string) (*Node, *Node) {
	nodes := readNodes(path)
	pairs := collectPairs(nodes)

	circuits := len(nodes)
	var a *Node
	var b *Node
	for _, p := range pairs {
		a = p.A
		b = p.B
		if find(a) == find(b) {
			continue
		}

		circuits -= 1
		union(a, b)

		if circuits == 1 {
			break
		}
	}

	return a, b
}

func main() {
	fmt.Println("Advent of Code / Day 8")

	s := connect("./input/input1.txt", 1000)
	fmt.Printf("1) %d * %d * %d = %d\n", s[0], s[1], s[2], s[0]*s[1]*s[2])

	a, b := maxCircuit("./input/input1.txt")
	fmt.Printf("2) %d * %d = %d\n", a.X, b.X, a.X*b.X)
}
