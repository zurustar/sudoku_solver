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

func (board *Board) ToString() string {
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
			if len(board.cells[i]) == 1 {
				s += strconv.Itoa(board.cells[i][0])
			} else {
				s += "0"
			}
		}
		s += "|\n"
	}
	s += "+---+---+---+"
	return s
}

func (board *Board) CopyFrom(src *Board) {
	for pos := 0; pos < 9*9; pos++ {
		board.cells[pos] = src.cells[pos]
	}
}

func (board *Board) Validate() BoardState {
	for pos := 0; pos < 9*9; pos++ {
		switch len(board.cells[pos]) {
		case 0:
			return INVALID
		case 1:
		default:
			return NOT_SOLVED
		}
	}
	return SOLVED
}

func (board *Board) Show() {
	fmt.Println(board.ToString())
}

func (board *Board) Find(pos, value int) int {
	for i, v := range board.cells[pos] {
		if v == value {
			return i
		}
	}
	return -1
}

func (board *Board) Remove(pos, value int) bool {
	result := false
	cands := []int{}
	for _, v := range board.cells[pos] {
		if v == value {
			result = true
		} else {
			cands = append(cands, v)
		}
	}
	if result {
		board.cells[pos] = cands
	}
	return result
}

func (board *Board) _Update1() bool {
	updated := false
	for src_pos := 0; src_pos < 9*9; src_pos++ {
		if len(board.cells[src_pos]) != 1 {
			continue
		}
		for dst_pos := 0; dst_pos < 9*9; dst_pos++ {
			if src_pos == dst_pos {
				continue
			}
			for target := 0; target < 3; target++ {
				if dic[src_pos][target] == dic[dst_pos][target] {
					if board.Remove(dst_pos, board.cells[src_pos][0]) {
						updated = true
					}
				}
			}
		}
	}
	return updated
}

func (board *Board) _Update2() bool {
	updated := false
	for src_pos := 0; src_pos < 9*9; src_pos++ {
		if len(board.cells[src_pos]) == 1 {
			continue
		}
		for target := 0; target < 3; target++ {
			for _, cand := range board.cells[src_pos] {
				found := false
				for dst_pos := 0; dst_pos < 9*9; dst_pos++ {
					if src_pos == dst_pos {
						continue
					}
					if dic[src_pos][target] != dic[dst_pos][target] {
						continue
					}
					if board.Find(dst_pos, cand) >= 0 {
						found = true
					}
				}
				if found == false {
					board.cells[src_pos] = []int{cand}
					updated = true
				}
			}
		}
	}
	return updated
}

func (board *Board) Update() {
	updated := true
	for updated {
		updated = board._Update1()
		if !updated {
			updated = board._Update2()
		}
		if !updated {
			break
		}
	}
}

func Solve(board *Board) (BoardState, *Board) {
	board.Update()
	if board.Validate() == INVALID {
		return INVALID, nil
	}
	for pos := 0; pos < 9*9; pos++ {
		if len(board.cells[pos]) > 1 {
			for _, cand := range board.cells[pos] {
				_b := NewBoard()
				_b.CopyFrom(board)
				_b.cells[pos] = []int{cand}
				result, _b := Solve(_b)
				if result == SOLVED {
					return SOLVED, _b
				}
			}
		}
	}
	return board.Validate(), board
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
