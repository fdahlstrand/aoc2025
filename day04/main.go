package main

import (
	"fmt"
)

func main() {
	fmt.Println("Advent of Code / Day 4")

	room := loadRoom("input/input1.txt")
	a := room.getAccessible()
	fmt.Printf("1) %d rolls of paper can be accessed\n", len(a))

	sum := len(a)
	for len(a) > 0 {
		for _, p := range a {
			room.removePaperRoll(p.row, p.col)
		}
		a = room.getAccessible()

		sum += len(a)
	}

	fmt.Printf("2) Total of %d rolls of paper removed\n", sum)
}
