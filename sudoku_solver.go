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

func (cell *Cell)set_cand(c int) {
  fmt.Println("set cand of", cell.x, ",", cell.y, "to", c)
  cell.cand = []int{c}
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

func (cell *Cell)has_cand(c int) bool {
  for _, cand := range cell.cand {
    if cand == c {
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
      board.cells = append(board.cells, NewCell(x,y,tmp[y*9+x]))
    }
  }
  return board
}

func (board *Board) solve_once_sub1(ls []int) bool {
  //fmt.Println("board.solve_once_sub1(",ls,")")
  updated := false
  for _, p := range ls {
    if len(board.cells[p].cand) == 1 {
      for _, p0 := range ls {
        if p != p0 {
          if board.cells[p0].remove_cand(board.cells[p].cand[0]) {
            updated = true
          }
        }
      }
    }
  }
  return updated
}

func (board *Board) solve_once_sub2(ls []int) bool {
  //fmt.Println("board.solve_once_sub2(", ls, ")")
  updated := false
  candcounter :=[]int{-1,0,0,0,0,0,0,0,0,0}
  for _, pos := range ls {
    for _, c := range board.cells[pos].cand {
      candcounter[c] += 1
    }
  }
  for i:=1; i<10; i++ {
    if candcounter[i] == 1 {
      for _, pos := range ls {
        if len(board.cells[pos].cand) > 1 {
          if board.cells[pos].has_cand(i) {
            board.cells[pos].set_cand(i)
            updated = true
            break
          }
        }
      }
    }
  }
  return updated
}

func (board *Board) solve_once_sub(ls []int) bool {
  //fmt.Println("solve_once_sub(", ls, ")")
  updated := false
  if board.solve_once_sub1(ls) {
    updated = true
  }
  if board.solve_once_sub2(ls) {
    updated = true
  }
  return updated
}

func (board *Board) solve_once() bool {
  updated := false
  for x:=0; x<9; x++ {
    ls := []int{}
    for pos:=0; pos<9*9; pos++ {
      if board.cells[pos].x == x {
        ls = append(ls, pos)
      }
    }
    if board.solve_once_sub(ls) {
      updated = true
    }
  }
  for y:=0; y<9; y++ {
    ls := []int{}
    for pos:=0; pos<9*9; pos++ {
      if board.cells[pos].y == y {
        ls = append(ls, pos)
      }
    }
    if board.solve_once_sub(ls) {
      updated = true
    }
  }
  for g:=0; g<9; g++ {
    ls := []int{}
    for pos:=0; pos<9*9; pos++ {
      if board.cells[pos].grp == g {
        ls = append(ls, pos)
      }
    }
    if board.solve_once_sub(ls) {
      updated = true
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
  b := NewBoard(os.Args[1])
  b.show()
  b.solve()
  b.show()
}
