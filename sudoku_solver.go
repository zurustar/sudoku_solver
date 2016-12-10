package main

import "os"
import "fmt"
import "strings"
import "strconv"



type Cell struct {
  y int      // このセルのx座標
  x int      // このセルのy座標
  pos int    // このセルの番号（ y * 9 + x の答え )
  grp int    // グループ番号
  cand []int // このセルにおけるかもしれない数字
}

//
// Cellのコンストラクタ 
//
func NewCell(x, y int, value string) *Cell {
  cell := new(Cell)
  cell.x = x
  cell.y = y
  cell.pos = y * 9 + x
  cell.grp = (cell.y/3)*3+cell.x/3
  c,_ := strconv.Atoi(value)
  if c == 0 {
    cell.cand = []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
  } else {
    cell.cand = []int{c}
  }
  return cell
}

func (cell *Cell) dump() {
  fmt.Println("DUMP", cell.x, cell.y, cell.pos, cell.grp, "--", cell.cand)
}

func (cell *Cell) to_s() string {
  result := ""
  for i := 0; i<len(cell.cand); i++ {
    result += strconv.Itoa(cell.cand[i])
  }
  return result;
}

func (cell *Cell) remove_cand(c int) bool {
  updated := false
  newcand := []int{}
  for _, n := range cell.cand {
    if n == c {
      updated = true
    } else {
     newcand = append(newcand, n)
    }
  }
  if len(newcand) == 0 {
    fmt.Println("ERROR")
    cell.dump()
    fmt.Println("c", c)
    os.Exit(1)
  }
  cell.cand = newcand
  return updated
}


type Board struct {
  cells []*Cell
}

func NewBoard(str string) *Board {
  board := new(Board)
  tmp := strings.Split(str, "")
  for y := 0; y<9; y++ {
    for x := 0; x<9; x++ {
      board.cells = append(board.cells, NewCell(x,y,tmp[y*9+x]))
    }
  }
  return board
}

func (board *Board) to_s() string {
  result := ""
  for y := 0; y<9; y++ {
    for x := 0; x<9; x++ {
      if result != "" {
        result += ","
      }
      result += board.cells[x + y * 9].to_s()
    }
  }
  return result
}

func (board *Board) solve_once() bool {
  updated := false
  for p:=0; p<9*9; p++ {
    if len(board.cells[p].cand) == 1 {
      c := board.cells[p].cand[0]
      y := board.cells[p].y
      for x:=0;x<9;x++ {
        if board.cells[p].x != x {
          if board.cells[y*9+x].remove_cand(c) {
            updated = true
          }
        }
      }
      x := board.cells[p].x
      for y:=0;y<9;y++ {
        if board.cells[p].y != y {
          if board.cells[y*9+x].remove_cand(c) {
            updated = true
          }
        }
      }
      for p0:=0; p0<9*9; p0++ {
        if p != p0 {
          if board.cells[p].grp == board.cells[p0].grp {
            if board.cells[p0].remove_cand(c) {
              updated = true
            }
          }
        }
      }
    }
  }
  return updated
}

func (board *Board) solve() {
  for {
    result := board.solve_once()
    if result == false {
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
      if len(board.cells[y*9+x].cand) == 1 {
        s += strconv.Itoa(board.cells[y*9+x].cand[0])
      } else {
        s += " "
      }
    }
    fmt.Println(s+"|")
  }
  fmt.Println("+---+---+---+")
}

func (board *Board) dump() {
  fmt.Println("DUMP x y p g --  cands")
  for y:=0;y<9;y++ {
    for x:=0;x<9;x++ {
      board.cells[y*9+x].dump()
    }
  }
}

func main() {
  fmt.Println("--------------------")
  b := NewBoard(os.Args[1])
  b.dump()
  b.show()
  b.solve()
  b.show()
  b.dump()
}
