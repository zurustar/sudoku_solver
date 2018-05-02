package main

import (
	"fmt"
	"os"
	"strconv"
)

type BoardState int

const (
	SOLVED BoardState = iota
	NOT_SOLVED
	INVALID
)

var dic = [9 * 9][3]int{
	{0, 0, 0}, {1, 0, 0}, {2, 0, 0},
	{3, 0, 1}, {4, 0, 1}, {5, 0, 1},
	{6, 0, 2}, {7, 0, 2}, {8, 0, 2},
	{0, 1, 0}, {1, 1, 0}, {2, 1, 0},
	{3, 1, 1}, {4, 1, 1}, {5, 1, 1},
	{6, 1, 2}, {7, 1, 2}, {8, 1, 2},
	{0, 2, 0}, {1, 2, 0}, {2, 2, 0},
	{3, 2, 1}, {4, 2, 1}, {5, 2, 1},
	{6, 2, 2}, {7, 2, 2}, {8, 2, 2},
	{0, 3, 3}, {1, 3, 3}, {2, 3, 3},
	{3, 3, 4}, {4, 3, 4}, {5, 3, 4},
	{6, 3, 5}, {7, 3, 5}, {8, 3, 5},
	{0, 4, 3}, {1, 4, 3}, {2, 4, 3},
	{3, 4, 4}, {4, 4, 4}, {5, 4, 4},
	{6, 4, 5}, {7, 4, 5}, {8, 4, 5},
	{0, 5, 3}, {1, 5, 3}, {2, 5, 3},
	{3, 5, 4}, {4, 5, 4}, {5, 5, 4},
	{6, 5, 5}, {7, 5, 5}, {8, 5, 5},
	{0, 6, 6}, {1, 6, 6}, {2, 6, 6},
	{3, 6, 7}, {4, 6, 7}, {5, 6, 7},
	{6, 6, 8}, {7, 6, 8}, {8, 6, 8},
	{0, 7, 6}, {1, 7, 6}, {2, 7, 6},
	{3, 7, 7}, {4, 7, 7}, {5, 7, 7},
	{6, 7, 8}, {7, 7, 8}, {8, 7, 8},
	{0, 8, 6}, {1, 8, 6}, {2, 8, 6},
	{3, 8, 7}, {4, 8, 7}, {5, 8, 7},
	{6, 8, 8}, {7, 8, 8}, {8, 8, 8},
}

type Board struct {
	cells [9 * 9][]int
}

func NewBoard() *Board {
	b := new(Board)
	for row := 0; row < 9; row++ {
		for column := 0; column < 9; column++ {
			b.cells[row*9+column] = []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
		}
	}
	return b
}

func (b *Board) ToString() string {
	s := ""
	for row := 0; row < 9; row++ {
		if row%3 == 0 {
			s += "+---+---+---+\n"
		}
		for column := 0; column < 9; column++ {
			if column%3 == 0 {
				s += "|"
			}
			i := row*9 + column
			if len(b.cells[i]) == 1 {
				s += strconv.Itoa(b.cells[i][0])
			} else {
				s += "0"
			}
		}
		s += "|\n"
	}
	s += "+---+---+---+"
	return s
}

func (b *Board) Set(src *Board) {
	for pos := 0; pos < 9*9; pos++ {
		b.cells[pos] = src.cells[pos]
	}
}

func (b *Board) Validate() BoardState {
	for pos := 0; pos < 9*9; pos++ {
		switch len(b.cells[pos]) {
		case 0:
			return INVALID
		case 1:
		default:
			return NOT_SOLVED
		}
	}
	return SOLVED
}

func (b *Board) Show() {
	fmt.Println(b.ToString())
}

func (b *Board) Find(pos, value int) int {
	for i, v := range b.cells[pos] {
		if v == value {
			return i
		}
	}
	return -1
}

func (b *Board) Remove(pos, value int) bool {
	result := false
	cands := []int{}
	for _, v := range b.cells[pos] {
		if v == value {
			result = true
		} else {
			cands = append(cands, v)
		}
	}
	if result {
		b.cells[pos] = cands
	}
	return result
}

func (b *Board) _Update1() bool {
	updated := false
	for src_column := 0; src_column < 9; src_column++ {
		for src_row := 0; src_row < 9; src_row++ {
			src_pos := src_column*9 + src_row
			if len(b.cells[src_pos]) != 1 {
				continue
			}
			for dst_column := 0; dst_column < 9; dst_column++ {
				for dst_row := 0; dst_row < 9; dst_row++ {
					dst_pos := dst_column*9 + dst_row
					if src_pos == dst_pos {
						continue
					}
					for target := 0; target < 3; target++ {
						if dic[src_pos][target] == dic[dst_pos][target] {
							if b.Remove(dst_pos, b.cells[src_pos][0]) {
								updated = true
							}
						}
					}
				}
			}
		}
	}
	return updated
}

func (b *Board) _Update2() bool {
	updated := false
	for src_column := 0; src_column < 9; src_column++ {
		for src_row := 0; src_row < 9; src_row++ {
			src_pos := src_column*9 + src_row
			if len(b.cells[src_pos]) == 1 {
				continue
			}
			for target := 0; target < 3; target++ {
				for _, cand := range b.cells[src_pos] {
					found := false
					for dst_column := 0; dst_column < 9; dst_column++ {
						for dst_row := 0; dst_row < 9; dst_row++ {
							dst_pos := dst_column*9 + dst_row
							if src_pos == dst_pos {
								continue
							}
							if dic[src_pos][target] != dic[dst_pos][target] {
								continue
							}
							if b.Find(dst_pos, cand) >= 0 {
								found = true
							}
						}
					}
					if found == false {
						b.cells[src_pos] = []int{cand}
						updated = true
					}
				}
			}
		}
	}
	return updated
}

func (b *Board) Update() {
	updated := true
	for updated {
		updated = b._Update1()
		if !updated {
			updated = b._Update2()
		}
		if !updated {
			break
		}
	}
}

func Solve(b *Board) (BoardState, *Board) {
	b.Update()
	if b.Validate() == INVALID {
		return INVALID, nil
	}
	for pos := 0; pos < 9*9; pos++ {
		if len(b.cells[pos]) > 1 {
			for _, cand := range b.cells[pos] {
				_b := NewBoard()
				_b.Set(b)
				_b.cells[pos] = []int{cand}
				result, _b := Solve(_b)
				if result == SOLVED {
					return SOLVED, _b
				}
			}
		}
	}
	return b.Validate(), b
}

func main() {
	if len(os.Args) != 2 {
		os.Exit(1)
	}
	if len(os.Args[1]) != 81 {
		os.Exit(1)
	}
	board := NewBoard()
	for i, c := range os.Args[1] {
		v, err := strconv.Atoi(string(c))
		if err != nil {
			os.Exit(1)
		}
		if v != 0 {
			board.cells[i] = []int{v}
		}
	}
	result, board := Solve(board)
	if result == SOLVED {
		board.Show()
		os.Exit(0)
	}
	os.Exit(1)
}
