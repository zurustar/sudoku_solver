package main

import "os"
import "fmt"
import "strings"
import "strconv"



type Cell struct {
  y int
  x int
  cands []int
}

//
// Cellのコンストラクタ 
//
func NewCell(x, y, v int) *Cell {
  cell := new(Cell)
  cell.x = x
  cell.y = y
  if(v == 0) {
    cell.cands = []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
  } else {
    cell.cands = []int{v}
  }
  return cell
}

func (cell *Cell)remove_cand(cand int) int {
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

func (cell *Cell)has_cand(cand int) bool {
  for _, c := range cell.cands {
    if c == cand {
      return true
    }
  }
  return false
}

type Board struct {
  cells []*Cell
}

func NewBoard(str string) *Board {
  board := new(Board)
  tmp := strings.Split(str, "")
  for y := 0; y<9; y++ {
    for x := 0; x<9; x++ {
      i, _ := strconv.Atoi(tmp[y*9+x])
      board.cells = append(board.cells, NewCell(x,y,i))
    }
  }
  return board
}


func (board *Board)solve_once() int {
  result := 0
  for y := 0; y<9; y++ {
    for x := 0; x<9; x++ {
      if(len(board.cells[y*9+x].cands) == 1) {
        cand := board.cells[y*9+x].cands[0]
        for i:=0; i<9; i++ {
          if i != x {
            ret := board.cells[y*9+i].remove_cand(cand)
            if ret == -1 {
              return -1
            }
            if ret == 1 {
              result = 1
            }
          }
          if i != y {
            ret := board.cells[i*9+x].remove_cand(cand)
            if ret == -1 {
              return -1
            }
            if ret == 1 {
              result = 1
            }
          }
        }
        mx := x - x % 3
        my := y - y % 3
        for ty := my; ty < my+3; ty++ {
          for tx := mx; tx < mx+3; tx++ {
            if !(x == tx && y == ty) {
              ret := board.cells[ty*9+tx].remove_cand(cand)
              if ret == -1 {
                return -1
              }
              if ret == 1 {
                result = 1
              }
            }
          }
        }
      } else {
        for _, cand := range board.cells[y*9+x].cands {
          //
          found := false
          for tx := 0; tx < 9; tx++ {
            if tx != x {
              if board.cells[y*9+tx].has_cand(cand) {
                found = true
                break
              }
            }
          }
          if !found {
            board.cells[y*9+x].cands = []int{cand}
            result = 1
            break
          }
          //
          found = false
          for ty := 0; ty < 9; ty++ {
            if ty != y {
              if board.cells[ty*9+x].has_cand(cand) {
                found = true
                break
              }
            }
          }
          if !found {
            board.cells[y*9+x].cands = []int{cand}
            result = 1
            break
          }
          //
          found = false
          mx := x - x % 3
          my := y - y % 3
          for ty := my; ty < my+3; ty++ {
            for tx := mx; tx < mx+3; tx++ {
              if !(x == tx && y == ty) {
                if board.cells[ty*9+tx].has_cand(cand) {
                  found = true
                  break
                }
              }
            }
          }
          if !found {
            board.cells[y*9+x].cands = []int{cand}
            result = 1
            break
          }
        }
      }
    }
  }
  return result
}

func (board *Board) solve() {
  result := 0
  for {
    result = board.solve_once()
    if result != 1 {
      break
    }
  }
}

func (board *Board) show() {
  for y := 0; y<9; y++ {
    if y % 3 == 0 {
      fmt.Println("+---+---+---+")
    }
    s := ""
    for x := 0; x<9; x++ {
      if x % 3 == 0 {
        s += "|"
      }
      if len(board.cells[y*9+x].cands) == 1 {
        s += strconv.Itoa(board.cells[y*9+x].cands[0])
      } else {
        s += " "
      }
    }
    fmt.Println(s+"|")
  }
  fmt.Println("+---+---+---+")
}


func main() {
  b := NewBoard(os.Args[1])
  b.show()
  b.solve()
  b.show()
}

