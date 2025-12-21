package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
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

func (m Machine) generateMatrix() Matrix {
	A := NewMatrix(len(m.JoltageRequirements), len(m.Buttons)+1)

	for j, b := range m.Buttons {
		index := 0
		var indices []int
		for b != 0 {
			if (b & 1) == 1 {
				indices = append(indices, index)
			}
			index += 1
			b >>= 1
		}

		for i := range slices.Values(indices) {
			A[i][j] = 1
		}

		for i, joltage := range m.JoltageRequirements {
			A[i][len(m.Buttons)] = float64(joltage)
		}
	}

	return A
}

func isValidSolution(D []Domain) bool {
	for k := range len(D) {
		if math.Abs(D[k].Upper-D[k].Lower) > Epsilon {
			return false
		}

		if math.Abs(D[k].Lower-math.Round(D[k].Lower)) > Epsilon {
			return false
		}

		if D[k].Lower < 0.0 {
			panic("Negative solution found")
		}
	}

	return true
}

func solve(eqns EquationSet, D []Domain, F []int) int {
	if isValidSolution(D) {
		sum := 0
		for k := range len(D) {
			val := int(math.Round(D[k].Lower))
			sum += val
		}
		return sum
	}

	if len(F) == 0 {
		return math.MaxInt
	}

	Fnew := make([]int, len(F)-1)
	f := F[0]
	copy(Fnew, F[1:])

	best := math.MaxInt
	fa, fb := int(math.Floor(D[f].Lower)), int(math.Ceil(D[f].Upper))
	for i := fa; i <= fb; i++ {
		Dnew := make([]Domain, len(D))
		copy(Dnew, D)
		Dnew[f] = NewConstDomain(float64(i))
		Dnew = UpdateDomains(Dnew, eqns)
		best = min(best, solve(eqns, Dnew, Fnew))
	}

	return best
}

func (m Machine) SolveJoltageLevels() int {
	A := m.generateMatrix()
	A = A.RowReducdEchelonForm()

	eqns := NewEquationSet(A)

	var D []Domain
	for range eqns {
		D = append(D, NewDomain(0, math.Inf(1)))
	}
	D = UpdateDomains(D, eqns)
	var F []int
	for k := range len(D) {
		if len(eqns[k]) > 1 {
			F = append(F, k)
		}
	}

	return solve(eqns, D, F)
}

type Matrix [][]float64

func NewMatrix(rows, cols int) Matrix {
	A := make([][]float64, rows)
	for i := range rows {
		A[i] = make([]float64, cols)
	}
	return A
}

func (A Matrix) Rows() int { return len(A) }
func (A Matrix) Cols() int { return len(A[0]) }

func (A Matrix) String() string {
	var b strings.Builder

	widths := make([]int, A.Cols())
	for i := range A.Rows() {
		for j := range A.Cols() {
			widths[j] = max(widths[j], len(fmt.Sprintf("%f", A[i][j])))
		}
	}

	for i := range A.Rows() {
		b.WriteString("[ ")
		for j := range A.Cols() {
			if j == A.Cols()-1 {
				b.WriteString("= ")

			}
			format := fmt.Sprintf("%%%ds ", widths[j])
			s := fmt.Sprintf("%f", A[i][j])
			b.WriteString(fmt.Sprintf(format, s))

		}
		b.WriteString("]\n")
	}

	return b.String()
}

const (
	Epsilon = 1e-9
)

func (A Matrix) RowReducdEchelonForm() Matrix {
	result := make(Matrix, A.Rows())
	for i := range A.Rows() {
		result[i] = make([]float64, A.Cols())
		copy(result[i], A[i])
	}

	pivotRow := 0
	for j := 0; j < A.Cols() && pivotRow < A.Rows(); j++ {
		i := pivotRow
		for i < A.Rows() && math.Abs(result[i][j]) < Epsilon {
			i++
		}

		if i < A.Rows() {
			result[i], result[pivotRow] = result[pivotRow], result[i]

			pivotElement := result[pivotRow][j]
			if math.Abs(pivotElement) > Epsilon {
				scalar := 1.0 / pivotElement
				for k := range A.Cols() {
					result[pivotRow][k] *= scalar
				}
			}

			for i := range A.Rows() {
				if i != pivotRow {
					factor := result[i][j]
					for k := j; k < A.Cols(); k++ {
						result[i][k] -= factor * result[pivotRow][k]
					}
				}
			}

			pivotRow++
		}
	}

	for i := range A.Rows() {
		for j := range A.Cols() {
			if math.Abs(result[i][j]) < Epsilon {
				result[i][j] = 0.0
			}
		}
	}

	return result
}

type EquationSet map[int][][]float64

func NewEquationSet(A Matrix) EquationSet {
	eqns := make(EquationSet)
	for k := range A.Cols() - 1 {
		for i := range A.Rows() {
			if math.Abs(A[i][k]) > Epsilon {
				factor := 1.0 / A[i][k]
				eqn := make([]float64, A.Cols())
				for j := range A.Cols() {
					if j == k {
						continue
					}

					if j == A.Cols()-1 {
						eqn[j] = factor * A[i][j]
					} else {
						eqn[j] = -factor * A[i][j]
					}

					if math.Abs(eqn[j]) < Epsilon {
						eqn[j] = 0.0
					}
				}
				eqns[k] = append(eqns[k], eqn)
			}
		}
	}

	return eqns
}

func (eqns EquationSet) String() string {
	var b strings.Builder
	for i, eq := range eqns {
		for _, e := range eq {
			b.WriteString(fmt.Sprintf("x[%d] =", i))
			b.WriteString(fmt.Sprintf(" %f", e[len(e)-1]))
			for j := range len(e) - 1 {
				f := e[j]
				if math.Abs(f) < Epsilon {
					continue
				}
				cj := math.Abs(f)
				sign := "+"
				if f < 0 {
					sign = "-"
				}
				b.WriteString(fmt.Sprintf(" %s %fx[%d]", sign, cj, j))
			}
			b.WriteString("\n")
		}
	}

	return b.String()
}

type Domain struct {
	Lower float64
	Upper float64
}

func NewDomain(lower, upper float64) Domain {
	return Domain{Lower: lower, Upper: upper}
}

func NewConstDomain(c float64) Domain {
	return NewDomain(c, c)
}

func (d Domain) String() string {
	return fmt.Sprintf("{%f..%f}", d.Lower, d.Upper)
}

func UpdateDomains(D []Domain, eqns EquationSet) []Domain {
	Dnew := make([]Domain, len(D))
	copy(Dnew, D)
	changed := true
	for changed {
		changed = false
		for i, d := range Dnew {
			for _, eq := range eqns[i] {
				C := NewConstDomain(eq[len(eq)-1])
				for j := range len(eq) - 1 {
					a := eq[j]
					d := Dnew[j]
					if math.Abs(a) < Epsilon {
						continue
					}

					if a < 0 {
						d.Lower, d.Upper = a*d.Upper, a*d.Lower
					} else {
						d.Lower, d.Upper = a*d.Lower, a*d.Upper
					}

					C.Lower += d.Lower
					C.Upper += d.Upper
				}
				Dnew[i].Lower = max(Dnew[i].Lower, C.Lower)
				Dnew[i].Upper = min(Dnew[i].Upper, C.Upper)
				if math.Abs(Dnew[i].Lower-d.Lower) > Epsilon || math.Abs(Dnew[i].Upper-d.Upper) > Epsilon {
					changed = true
				}
			}
		}
	}

	return Dnew
}

func (d Domain) AsIntRange() (int, int) {
	return int(math.Floor(d.Lower)), int(math.Ceil(d.Upper))
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

	sum = 0
	for k := range len(machines) {
		v := machines[k].SolveJoltageLevels()
		sum += v
	}

	fmt.Printf("2) Need at least %d button presses\n", sum)
}
