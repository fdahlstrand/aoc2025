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

func readProblems(file *os.File) []problem {
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

	return problems
}

func readCephalopodProblems(file *os.File) []problem {
	scanner := bufio.NewScanner(file)
	var rows [][]string = nil
	for scanner.Scan() {
		line := scanner.Text()
		cols := strings.Fields(line)
		rows = append(rows, cols)
	}

	widths := make([]int, len(rows[0]))
	for _, row := range rows {
		for col, v := range row {
			widths[col] = max(widths[col], len(v))
		}
	}

	file.Seek(0, 0)
	scanner = bufio.NewScanner(file)
	values := make([][]string, len(rows[0]))
	for scanner.Scan() {
		line := scanner.Text()
		pos := 0
		for col, w := range widths {
			v := line[pos : pos+w]
			values[col] = append(values[col], v)
			pos += w + 1
		}
	}

	var problems []problem
	for col, v := range values {
		width := widths[col]
		var problem problem
		for n := range width {

			var b strings.Builder
			for _, str := range v {
				switch str[n] {
				case '+':
					problem.op = Add
				case '*':
					problem.op = Mul
				default:
					b.WriteByte(str[n])
				}
			}
			v, err := strconv.Atoi(strings.TrimSpace(b.String()))
			check(err, "Failed to convert number")
			problem.values = append(problem.values, v)
		}
		problems = append(problems, problem)
	}
	return problems
}

func solve(problems []problem) int {
	sum := 0
	for _, p := range problems {
		var ans int
		switch p.op {
		case Add:
			ans = 0
			for _, v := range p.values {
				ans += v
			}
		case Mul:
			ans = 1
			for _, v := range p.values {
				ans *= v
			}
		default:
			panic("Unexpected operand")
		}
		sum += ans
	}

	return sum
}

func main() {
	fmt.Println("Advent of Code / Day 6")

	file, err := os.Open("input/input1.txt")
	check(err, "Failed to open file")
	defer file.Close()

	sum := solve(readProblems(file))
	fmt.Printf("1) Grand total %d\n", sum)

	file.Seek(0, 0)
	sum = solve(readCephalopodProblems(file))
	fmt.Printf("2) Grand total %d\n", sum)
}
