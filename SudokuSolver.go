package main

import (
	"fmt"
	"os"
)

var initial_complexibility uint64
var try_counter uint64

func main() {
	for i, filename := range os.Args {
		if i != 0 {
			b := NewBoard(filename)
			fmt.Println("\n", b.ToS())
			b.Update()
			initial_complexibility = b.Solved()
			try_counter = 0
			Solve(b, 0)
			fmt.Println("\n", b.ToS())
		}
	}
}

func Solve(b *Board, depth uint64) uint64 {
	try_counter += 1
	fmt.Println("depth", depth, "counter", try_counter)
	//	fmt.Println("\n", b.ToS())
	b.Update()
	result := b.Solved()
	if result == 1 {
		fmt.Println("\n", b.ToS())
		os.Exit(1)
	}
	if result == 0 {
		return 0
	}
	fmt.Println("initial:", initial_complexibility, " current:", result)
	tmppos := []int{}
	for l := 2; l <= 9; l++ {
		for pos := 0; pos < 9*9; pos++ {
			if len(b.Cells[pos]) == l {
				tmppos = append(tmppos, pos)
			}
		}
	}
	for _, pos := range tmppos {
		cands := b.Cells[pos]
		for _, c := range cands {
			tmpb := b.Duplicate()
			fmt.Println("guess", pos, "is", c)
			tmpb.Cells[pos] = []int{c}
			Solve(tmpb, depth+1)
		}
	}
	return 0
}
