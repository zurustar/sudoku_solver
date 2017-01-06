package main

import "os"
import "fmt"
import "strings"
import "strconv"

// =======================================================================
type Cell struct {
	y     int
	x     int
	cands []int
}

// -----------------------------------------------------------------------
func NewCell(x, y, v int) *Cell {
	cell := new(Cell)
	cell.x = x
	cell.y = y
	if v == 0 {
		cell.cands = []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	} else {
		cell.cands = []int{v}
	}
	return cell
}

// -----------------------------------------------------------------------
func (cell *Cell) remove_cand(cand int) int {
	tmp := []int{}
	for _, c := range cell.cands {
		if c != cand {
			tmp = append(tmp, c)
		}
	}
	if len(tmp) == 0 {
		return -1
	}
	if len(tmp) == len(cell.cands) {
		return 0
	}
	cell.cands = tmp
	return 1
}

// -----------------------------------------------------------------------
func (cell *Cell) has_cand(cand int) bool {
	for _, c := range cell.cands {
		if c == cand {
			return true
		}
	}
	return false
}

// =======================================================================
type Board struct {
	cells []*Cell
}

// -----------------------------------------------------------------------
func NewBoard(str string) *Board {
	if len(str) != 9*9 {
		return nil
	}
	board := new(Board)
	tmp := strings.Split(str, "")
	for y := 0; y < 9; y++ {
		for x := 0; x < 9; x++ {
			i, _ := strconv.Atoi(tmp[y*9+x])
			board.cells = append(board.cells, NewCell(x, y, i))
		}
	}
	return board
}

// -----------------------------------------------------------------------
func (b *Board) clone() *Board {
	s := ""
	for i:=0; i<9*9; i++ {
		s += "0"
	}
	newboard := NewBoard(s)
	for i:=0; i<9*9; i++ {
		newboard.cells[i].cands = b.cells[i].cands
	}
	return newboard
}

// -----------------------------------------------------------------------
func (b *Board) is_solved() bool {
	for i := 0; i < 9*9; i++ {
		if len(b.cells[i].cands) != 1 {
			return false
		}
	}
	return true
}

// -----------------------------------------------------------------------
func (board *Board) show() {
	for y := 0; y < 9; y++ {
		if y%3 == 0 {
			fmt.Println("+---+---+---+")
		}
		s := ""
		for x := 0; x < 9; x++ {
			if x%3 == 0 {
				s += "|"
			}
			if len(board.cells[y*9+x].cands) == 1 {
				s += strconv.Itoa(board.cells[y*9+x].cands[0])
			} else {
				s += " "
			}
		}
		fmt.Println(s + "|")
	}
	fmt.Println("+---+---+---+")
}

// ---------------------------------------------------------------------------
func solve_sub1(b *Board, x, y int) int {
	result := 0
	cand := b.cells[y*9+x].cands[0]
	for i := 0; i < 9; i++ {
		if i != x {
			ret := b.cells[y*9+i].remove_cand(cand)
			if ret == -1 {
				return -1
			}
			if ret == 1 {
				result = 1
			}
		}
		if i != y {
			ret := b.cells[i*9+x].remove_cand(cand)
			if ret == -1 {
				return -1
			}
			if ret == 1 {
				result = 1
			}
		}
	}
	mx := x - x%3
	my := y - y%3
	for ty := my; ty < my+3; ty++ {
		for tx := mx; tx < mx+3; tx++ {
			if !(x == tx && y == ty) {
				ret := b.cells[ty*9+tx].remove_cand(cand)
				if ret == -1 {
					return -1
				}
				if ret == 1 {
					result = 1
				}
			}
		}
	}
	return result
}

// ---------------------------------------------------------------------------
func solve_sub2(b *Board, x, y int) int {
	for _, cand := range b.cells[y*9+x].cands {
		found := false
		for tx := 0; tx < 9; tx++ {
			if tx != x {
				if b.cells[y*9+tx].has_cand(cand) {
					found = true
					break
				}
			}
		}
		if !found {
			b.cells[y*9+x].cands = []int{cand}
			return 1
		}

		found = false
		for ty := 0; ty < 9; ty++ {
			if ty != y {
				if b.cells[ty*9+x].has_cand(cand) {
					found = true
					break
				}
			}
		}
		if !found {
			b.cells[y*9+x].cands = []int{cand}
			return 1
		}

		found = false
		mx := x - x%3
		my := y - y%3
		for ty := my; ty < my+3; ty++ {
			for tx := mx; tx < mx+3; tx++ {
				if !(x == tx && y == ty) {
					if b.cells[ty*9+tx].has_cand(cand) {
						found = true
						break
					}
				}
			}
		}
		if !found {
			b.cells[y*9+x].cands = []int{cand}
			return 1
		}
	}
	return 0
}

// ---------------------------------------------------------------------------
func solve_once(b *Board) int {
	result := 0
	for y := 0; y < 9; y++ {
		for x := 0; x < 9; x++ {
			res := 0
			if len(b.cells[y*9+x].cands) == 1 {
				res = solve_sub1(b, x, y)
			} else {
				res = solve_sub2(b, x, y)
			}
			if res == -1 {
				return -1
			}
			if res == 1 {
				result = 1
			}
		}
	}
	return result
}

// ---------------------------------------------------------------------------
func guess(b *Board) int {
	for y := 0; y < 9; y++ {
		for x := 0; x < 9; x++ {
			if len(b.cells[y*9+x].cands) > 1 {
				for _, cand := range b.cells[y*9+x].cands {
					tmp := b.clone()
					tmp.cells[y*9+x].cands = []int{cand}
					if solve_once(tmp) == -1 {
						b.cells[y*9+x].remove_cand(cand)
					}
				}
			}
		}
	}
	return 0
	
}

// ---------------------------------------------------------------------------
func solve(b *Board) bool {
	for {
		res := solve_once(b)
		if res == -1 {
			return false
		}
		if b.is_solved() {
			return true
		}
	}
	return false
}

func main() {
	b := NewBoard(os.Args[1])
	if b != nil {
		b.show()
		solve(b)
		b.show()
	}
}
