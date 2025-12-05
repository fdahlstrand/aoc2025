package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type ReadState byte

const (
	ReadRanges ReadState = iota
	ReadIngredients
)

type Range struct {
	lower int
	upper int
}

type RangeSlice []Range

func (r RangeSlice) Len() int      { return len(r) }
func (r RangeSlice) Swap(i, j int) { r[i], r[j] = r[j], r[i] }
func (r RangeSlice) Less(i, j int) bool {
	if r[i].lower == r[j].lower {
		return r[i].upper < r[j].upper
	}
	return r[i].lower < r[j].lower
}

func (r Range) String() string {
	return fmt.Sprintf("(%d, %d)", r.lower, r.upper)
}

func check(err error, msg string) {
	if err != nil {
		panic(msg)
	}
}

func rangeOverlaps(a Range, b Range) bool {
	if a.upper < b.lower || a.lower > b.upper {
		return false
	}
	return true
}

func merge(a Range, b Range) Range {
	if !rangeOverlaps(a, b) {
		panic("Ranges do not overlap")
	}

	return Range{min(a.lower, b.lower), max(a.upper, b.upper)}
}

func main() {
	fmt.Println("Advent of Code / Day 5")

	file, err := os.Open("./input/input1.txt")
	check(err, "Failed to open file")

	scanner := bufio.NewScanner(file)

	state := ReadRanges
	var ranges []Range = nil
	var ingredients []int = nil
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			state = ReadIngredients
		} else if state == ReadRanges {
			r := strings.Split(line, "-")
			lower, err := strconv.Atoi(r[0])
			check(err, "Failed to parse range")
			upper, err := strconv.Atoi(r[1])
			check(err, "Failed to parse range")
			ranges = append(ranges, Range{lower, upper})
		} else {
			i, err := strconv.Atoi(line)
			check(err, "Failed to parse ingredient")
			ingredients = append(ingredients, i)
		}
	}

	sum := 0
	for _, i := range ingredients {
		for _, r := range ranges {
			if r.lower <= i && i <= r.upper {
				sum += 1
				break
			}

		}
	}
	fmt.Printf("1) Fresh ingredients: %d\n", sum)

	sort.Sort(RangeSlice(ranges))

	var combined []Range = nil
	i := 1
	cur := ranges[0]
	for i < len(ranges) {
		if rangeOverlaps(cur, ranges[i]) {
			cur = merge(cur, ranges[i])
		} else {
			combined = append(combined, cur)
			cur = ranges[i]
		}
		i++
	}
	combined = append(combined, cur)

	count := 0
	for _, r := range combined {
		count += r.upper - r.lower + 1
	}
	fmt.Printf("2) Total fresh ingredients: %d\n", count)
}
