package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func is_valid(id int) bool {
	s := fmt.Sprintf("%d", id)

	if len(s)%2 != 0 {
		return false
	}

	m := len(s) / 2
	return s[0:m] == s[m:]
}

func is_valid2(id int) bool {
	s := fmt.Sprintf("%d", id)

	for length := 1; length <= len(s)/2; length++ {
		if len(s)%length != 0 {
			continue
		}

		valid := true
		for ix := length; ix < len(s)-length+1; ix += length {
			if s[0:length] != s[ix:ix+length] {
				valid = false
				break
			}
		}
		if valid {
			return true
		}

	}

	return false
}

func to_int(s string) int {
	i, e := strconv.Atoi(s)
	if e != nil {
		panic("Can't convert string to integer")
	}

	return i
}

func main() {
	fmt.Println("Advent of Code / Day 2")

	file, err := os.Open("input/input1.txt")
	if err != nil {
		panic("Failed to open file")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	sum_1 := 0
	sum_2 := 0
	for scanner.Scan() {
		line := scanner.Text()
		for rng := range strings.SplitSeq(line, ",") {
			r := strings.Split(rng, "-")
			lower := to_int(r[0])
			upper := to_int(r[1])
			for id := lower; id <= upper; id++ {
				if is_valid(id) {
					sum_1 += id
				}
				if is_valid2(id) {
					sum_2 += id
				}
			}

		}
	}

	fmt.Printf("1) Sum of invalid IDs: %d\n", sum_1)
	fmt.Printf("2) Sum of invalid IDs: %d\n", sum_2)
}
