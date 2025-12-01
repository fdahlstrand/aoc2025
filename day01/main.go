package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func mod(a, b int) int {
	return (a%b + b) % b
}

func main() {
	fmt.Println("Advent of Code / Day 1")

	file, err := os.Open("input/input1.txt")
	if err != nil {
		panic("Failed to open file!")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	pos := 50
	count := 0
	pass_count := 0
	for scanner.Scan() {
		line := scanner.Text()

		i, err := strconv.Atoi(line[1:])
		if err != nil {
			panic("Failed to convert integer!")
		}

		if line[0] == 'L' {
			i = -i
		}

		next := pos + i

		if next >= 100 {
			pass_count += next / 100
		} else if next <= 0 {
			pass_count += (next / -100) + 1
			if pos == 0 {
				pass_count -= 1
			}
		}

		pos = mod(next, 100)
		if pos == 0 {
			count += 1
		}
	}

	fmt.Printf("Password: %d\n", count)
	fmt.Printf("Password method 0x434C49434B: %d\n", pass_count)
}
