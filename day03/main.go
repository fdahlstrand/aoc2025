package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strconv"
)

func check(err error, msg string) {
	if err != nil {
		panic(msg)
	}
}

func find_max(s string, start int, end int) (int, byte) {
	var max byte = 0
	max_ix := -1
	for i := start; i < end; i++ {
		b := s[i] - '0'
		if b > max {
			max = b
			max_ix = i
		}
	}

	return max_ix, s[max_ix]
}

func get_joltage(s string, batteries int) int {
	ix := -1
	buf := bytes.NewBufferString("")
	for n := range batteries {
		var ch byte = 0
		ix, ch = find_max(s, ix+1, len(s)-(batteries-1)+n)
		buf.WriteByte(ch)
	}
	joltage, err := strconv.Atoi(buf.String())
	check(err, "Failed to convert joltage string")

	return joltage
}

func main() {
	fmt.Println("Advent of Code / Day 3")

	file, err := os.Open("input/input1.txt")
	check(err, "Failed to open file")
	defer file.Close()

	scanner := bufio.NewScanner(file)

	joltage_sum1 := 0
	joltage_sum2 := 0
	for scanner.Scan() {
		line := scanner.Text()
		joltage_sum1 += get_joltage(line, 2)
		joltage_sum2 += get_joltage(line, 12)
	}

	fmt.Printf("1) Total joltage = %d\n", joltage_sum1)
	fmt.Printf("2) Total joltage = %d\n", joltage_sum2)
}
