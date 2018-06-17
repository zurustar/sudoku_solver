package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
)

var peers = [][]int{}

func Pos2Block(x, y int) int {
	return int(math.Trunc(float64(y)/3.0)*3 + math.Trunc(float64(x)/3.0))
}

func InitPeer() {
	for pos := 0; pos < 81; pos++ {
		peers = append(peers, []int{})
	}
	for sy := 0; sy < 9; sy++ {
		for sx := 0; sx < 9; sx++ {
			sp := sy*9 + sx
			for dy := 0; dy < 9; dy++ {
				for dx := 0; dx < 9; dx++ {
					dp := dy*9 + dx
					if sp == dp {
						continue
					} else if sx == dx {
						peers[sp] = append(peers[sp], dp)
					} else if sy == dy {
						peers[sp] = append(peers[sp], dp)
					} else if Pos2Block(sx, sy) == Pos2Block(dx, dy) {
						peers[sp] = append(peers[sp], dp)
					}
				}
			}
		}
	}
}

type Board [9 * 9][]int

func NewBoard() *Board {
	b := new(Board)
	for p := 0; p < 81; p++ {
		b[p] = []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	}
	return b
}

func (board *Board) ToString() string {
	s := ""
	for i := 0; i < 9*9; i++ {
		if len(board[i]) == 1 {
			s += strconv.Itoa(board[i][0])
		} else {
			s += "0"
		}
	}
	return s
}

func (board *Board) ShowDetail() {
	s := ""
	for i := 0; i < 90; i++ {
		s += "="
	}
	s += "\n"
	for y := 0; y < 9; y++ {
		s += "|"
		for x := 0; x < 9; x++ {
			p := y*9 + x
			for v := 1; v < 10; v++ {
				if board.Find(p, v) >= 0 {
					s += strconv.Itoa(v)
				} else {
					s += " "
				}
			}
			s += "|"
		}
		s += "\n"
		if y == 2 || y == 5 {
			for i := 0; i < 90; i++ {
				s += "-"
			}
			s += "\n"
		}
	}
	for i := 0; i < 90; i++ {
		s += "="
	}
	s += "\n"
	fmt.Println(s)
}

func (board *Board) Find(pos, value int) int {
	for i, v := range board[pos] {
		if v == value {
			return i
		}
	}
	return -1
}

//
// 指定された場所のセルの候補リストから指定された値を削除する
//
func (board *Board) Remove(pos, value int) int {
	found := false
	cands := []int{}
	if len(board[pos]) == 1 {
		if board[pos][0] == value {
			return -1
		}
	}
	for _, v := range board[pos] {
		if v == value {
			found = true
		} else {
			cands = append(cands, v)
		}
	}
	if found {
		board[pos] = cands
		return 1
	}
	return 0
}

func Check(board *Board) int {
	for i := 0; i < 81; i++ {
		if len(board[i]) == 1 {
			for _, j := range peers[i] {
				if len(board[j]) == 1 {
					if board[i][0] == board[j][0] {
						return -1
					}
				}
			}
		}
	}
	return 1
}

//
// 各セルの候補を絞り込む
//
func (board *Board) Update() int {
	flag, updated := true, false
	for flag {
		flag = false
		for sp := 0; sp < 81; sp++ {
			if len(board[sp]) == 1 {
				// すでに候補がひとつしかなかったら
				// お友達からその候補を削除する
				for _, dp := range peers[sp] {
					switch board.Remove(dp, board[sp][0]) {
					case 1:
						flag, updated = true, true
					case -1:
						return -1
					}
				}
			} else {
				// 候補がたくさん残っていたら、それぞれの候補について
				// 持っているお友達がいるかを調べて、いなかったら
				// その候補に確定
				for _, v := range board[sp] {
					found := false
					for _, dp := range peers[sp] {
						if board.Find(dp, v) >= 0 {
							found = true
							break
						}
					}
					if found == false {
						board[sp] = []int{v}
						flag, updated = true, true
					}
				}
			}
		}
	}
	if updated {
		return 1
	}
	return 0
}

//
// 解く。ルールどおりに試して答えがでなかったら仮置きして再帰呼び出し
//
func Solve(board *Board, depth int) (int, *Board) {
	fmt.Println("depth=", depth)
	//board.ShowDetail()
	if Check(board) == -1 {
		return -1, nil
	}
	switch board.Update() {
	case 0:
		return 0, board
	case -1:
		return -1, nil
	}
	// 各セルそれぞれにのこりいくつの候補が残っているかを調査
	// 残りの個数ごとにそのセルの番号を配列に格納する
	len_list := [10][]int{{}, {}, {}, {}, {}, {}, {}, {}, {}, {}}
	for pos := 0; pos < 9*9; pos++ {
		l := len(board[pos])
		len_list[l] = append(len_list[l], pos)
	}
	// 候補数がゼロのセルが存在したら異常な状態になっている
	if len(len_list[0]) > 0 {
		return -1, nil
	}
	// すべてのセルの候補数が1になっていたら解けている
	if len(len_list[1]) == 9*9 {
		return 1, board
	}
	// まだ候補数が１になっていないセルについて、
	// 適当にいっこの値に絞り込んで、ためしに解いてみる
	// 再帰呼び出しになるので、結果的にかたっぱしから試すことになる
	for i := 2; i < 10; i++ {
		for _, pos := range len_list[i] {
			if len(board[pos]) > 1 {
				for _, cand := range board[pos] {
					new_board := new(Board)
					for j := 0; j < 81; j++ {
						new_board[j] = []int{}
						for _, v := range board[j] {
							new_board[j] = append(new_board[j], v)
						}
					}
					new_board[pos] = []int{cand}
					result, new_board := Solve(new_board, depth+1)
					if result == 1 {
						return 1, new_board
					} else if result == -1 {
						board.Remove(pos, cand)
					}
				}
			}
		}
	}
	return 0, board
}

//
// 81個の数値からなる文字列を人間が見やすい形に変換
//
func ToString(src string) string {
	s := ""
	for row := 0; row < 9; row++ {
		if row%3 == 0 {
			s += "\n"
		}
		for column := 0; column < 9; column++ {
			s += string(src[row*9+column])
			if column%3 == 2 {
				s += " "
			}
		}
		s += "\n"
	}
	return s
}

//
// 問題文字列から余計な文字をカットしてから
// 問題情報を保持するBoardインスタンスに突っ込んで解く
//
func RunSolver(src string) string {
	re := regexp.MustCompile(`[^0-9]`)
	q := re.ReplaceAllString(src, "")
	if len(q) != 9*9 {
		log.Println("invalid question")
		return ""
	}
	board := NewBoard()
	for i, c := range q {
		v, err := strconv.Atoi(string(c))
		if err != nil {
			return "" // bug?
		}
		if v != 0 {
			board[i] = []int{v}
		}
	}
	result, board := Solve(board, 0)
	if result == 1 {
		return board.ToString()
	}
	return ""
}

// 問題が書いてあるファイルを読み込む
func Load(filename string) string {
	log.Println("load", filename)
	fp, err := os.Open(filename)
	if err != nil {
		return ""
	}
	defer fp.Close()
	result := ""
	buf := make([]byte, 1024)
	for {
		n, err := fp.Read(buf)
		if n == 0 {
			break
		}
		if err != nil {
			return ""
		}
		result += string(buf[:n])
	}
	return result
}

func main() {
	if len(os.Args) == 2 {
		InitPeer()
		s := ""
		_, err := os.Stat(os.Args[1])
		if err == nil {
			s = Load(os.Args[1])
		} else {
			s = os.Args[1]
		}
		fmt.Println(s)
		fmt.Println(ToString(RunSolver(s)))
	}
}
