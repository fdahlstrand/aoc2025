package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Operand byte

const (
	Add Operand = iota
	Mul
	None
)

type problem struct {
	values []int
	op     Operand
}

func check(err error, msg string) {
	if err != nil {
		panic(msg)
	}
}

func main() {
	fmt.Println("Advent of Code / Day 6")

	file, err := os.Open("input/input1.txt")
	check(err, "Failed to open file")
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var rows [][]string = nil
	for scanner.Scan() {
		line := scanner.Text()
		cols := strings.Fields(line)
		rows = append(rows, cols)
	}

	var problems []problem = nil
	for col := range len(rows[0]) {
		cur := problem{values: nil, op: None}
		for row := range len(rows) {
			val := rows[row][col]
			switch val {
			case "+":
				cur.op = Add
			case "*":
				cur.op = Mul
			default:
				v, err := strconv.Atoi(val)
				check(err, "Failed to convert number")
				cur.values = append(cur.values, v)
			}

		}
		problems = append(problems, cur)
	}

	sum := 0
	for _, p := range problems {
		var ans int
		switch p.op {
		case Add:
			ans = 0
			for _, v := range p.values {
				ans += v
			}
			fmt.Println(ans)
		case Mul:
			ans = 1
			for _, v := range p.values {
				ans *= v
			}
			fmt.Println(ans)
		default:
			panic("Unexpected operand")
		}
		sum += ans
	}
	fmt.Printf("1) Grand total %d\n", sum)
}
