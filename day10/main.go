package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
)

type Machine struct {
	GoalLights          int
	Buttons             []int
	JoltageRequirements []int
}

func check(err error, msg string) {
	if err != nil {
		panic(msg)
	}
}

func advance(s string) string {
	if len(s) > 1 {
		return s[1:]
	}

	if len(s) == 1 {
		return ""
	}

	panic("Syntax Error: Unexpected end of string")
}

func match(s string, ch byte) string {
	if len(s) > 0 && s[0] == ch {
		return s[1:]
	}

	msg := ""
	if len(s) > 0 {
		msg = fmt.Sprintf("Syntax Error: expected '%c' got '%c'", ch, s[0])
	} else {
		msg = "Got empty string"
	}
	panic(msg)
}

func parseLightDiagram(s string) (string, int) {
	s = match(s, '[')

	// TODO: Wrong order
	lights := 0
	pos := 0
	for len(s) > 0 && s[0] != ']' {
		switch s[0] {
		case '.':
			lights |= 0 << pos
			pos++
			s = advance(s)
		case '#':
			lights |= 1 << pos
			pos++
			s = advance(s)
		default:
			panic("Syntax Error: Light Diagram")
		}
	}

	s = match(s, ']')

	return s, lights
}

func parseWiringSchematic(s string) (string, int) {
	s = match(s, '(')

	wiring := 0
	done := false
	for len(s) > 0 && !done {
		switch s[0] {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			v, err := strconv.Atoi(s[0:1])
			check(err, "Syntax Error: Failed to convert digit")
			wiring |= 1 << v
		default:
			panic("Syntax Error: Wiring Schematics")
		}

		s = advance(s)

		switch s[0] {
		case ',':
			s = advance(s)
		case ')':
			s = advance(s)
			done = true
		default:
			panic("Syntax Error: Wiring Schematics")
		}
	}

	return s, wiring
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func parseNumber(s string) (string, int) {
	if len(s) == 0 || !isDigit(s[0]) {
		panic("Syntax Error: Expected Number")
	}
	pos := 0

	for len(s) > 0 && isDigit(s[pos+1]) {
		pos++
	}

	val, err := strconv.Atoi(s[0 : pos+1])
	check(err, "Syntax Error: Failed to parse number")

	return s[pos+1:], val
}

func parseJoltageRequirements(s string) (string, []int) {
	s = match(s, '{')

	var reqs []int = nil
	done := false
	for len(s) > 0 && !done {
		var val int
		s, val = parseNumber(s)
		reqs = append(reqs, val)

		switch s[0] {
		case ',':
			s = advance(s)
		case '}':
			s = advance(s)
			done = true
		default:
			panic("Syntax Error: Joltage requirements")
		}
	}

	return s, reqs
}

func parseLine(line string) Machine {
	s := line[:]
	var m Machine
	for len(s) > 0 {
		switch s[0] {
		case '[':
			s, m.GoalLights = parseLightDiagram(s)
		case '(':
			wiring := 0
			s, wiring = parseWiringSchematic(s)
			m.Buttons = append(m.Buttons, wiring)
		case '{':
			s, m.JoltageRequirements = parseJoltageRequirements(s)
		case ' ':
			s = s[1:]
		default:
			msg := fmt.Sprintf("Unepected character '%d'", s[0])
			panic(msg)
		}

	}

	return m
}

func (m Machine) search(state int, depth int, seen map[int]int) int {
	if state == m.GoalLights {
		return depth
	}

	seenDepth, isSeen := seen[state]
	if isSeen {
		if seenDepth <= depth {
			return math.MaxInt
		}
	}

	seen[state] = depth
	p := math.MaxInt
	for _, b := range m.Buttons {
		next_state := state ^ b
		p = min(p, m.search(next_state, depth+1, seen))
	}

	return p
}

func main() {
	fmt.Println("Advent of Code / Day 10")

	file, err := os.Open("./input/input1.txt")
	check(err, "Failed to open file")
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var machines []Machine
	for scanner.Scan() {
		line := scanner.Text()
		machines = append(machines, parseLine(line))
	}

	sum := 0
	for _, m := range machines {
		seen := make(map[int]int)
		sum += m.search(0, 0, seen)
	}

	fmt.Printf("1) Need at least %d button presses\n", sum)
}
