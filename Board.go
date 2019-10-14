package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

var pos2grp []int = []int{
	0, 0, 0, 1, 1, 1, 2, 2, 2,
	0, 0, 0, 1, 1, 1, 2, 2, 2,
	0, 0, 0, 1, 1, 1, 2, 2, 2,
	3, 3, 3, 4, 4, 4, 5, 5, 5,
	3, 3, 3, 4, 4, 4, 5, 5, 5,
	3, 3, 3, 4, 4, 4, 5, 5, 5,
	6, 6, 6, 7, 7, 7, 8, 8, 8,
	6, 6, 6, 7, 7, 7, 8, 8, 8,
	6, 6, 6, 7, 7, 7, 8, 8, 8}

var grp2pos [][]int = [][]int{
	{0, 1, 2, 9, 10, 11, 18, 19, 20},
	{3, 4, 5, 12, 13, 14, 21, 22, 23},
	{6, 7, 8, 15, 16, 17, 24, 25, 26},
	{27, 28, 29, 36, 37, 38, 45, 46, 47},
	{30, 31, 32, 39, 40, 41, 48, 49, 50},
	{33, 34, 35, 42, 43, 44, 51, 52, 53},
	{54, 55, 56, 63, 64, 65, 72, 73, 74},
	{57, 58, 59, 66, 67, 68, 75, 76, 77},
	{60, 61, 68, 69, 70, 71, 78, 79, 80}}

func IsPeer(a, b int) bool {
	if a == b {
		return false
	}
	ax, ay, ag := GetPosInfo(a)
	bx, by, bg := GetPosInfo(b)
	if ax == bx || ay == by || ag == bg {
		return true
	}
	return false
}

func GetPosInfo(pos int) (int, int, int) {
	sx := pos % 9
	sy := (pos - sx) / 9
	return sx, sy, pos2grp[pos]
}

type Board struct {
	Cells [][]int
}

func NewBoard(filename string) *Board {
	p := new(Board)
	p.Cells = [][]int{}
	fp, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer fp.Close()
	b, err := ioutil.ReadAll(fp)
	if err != nil {
		log.Fatal(err)
	}
	for _, c := range b {
		v, err := strconv.Atoi(string(c))
		if err == nil {
			if v == 0 {
				p.Cells = append(p.Cells, []int{1, 2, 3, 4, 5, 6, 7, 8, 9})
			} else {
				p.Cells = append(p.Cells, []int{v})
			}
		}
	}
	if len(p.Cells) != 9*9 {
		log.Fatal("failed to load", filename, "(invalid format)")
	}
	return p
}

func (p *Board) Solved() uint64 {
	var result uint64
	result = 1
	for _, cands := range p.Cells {
		result *= uint64(len(cands))
	}
	return result
}

func (p *Board) Duplicate() *Board {
	b := new(Board)
	for _, cands := range p.Cells {
		b.Cells = append(b.Cells, cands)
	}
	return b
}

func (p *Board) HasCandOf(pos, cand int) bool {
	for _, c := range p.Cells[pos] {
		if c == cand {
			return true
		}
	}
	return false
}

func (p *Board) Remove(pos, cand int) bool {
	newcands := []int{}
	found := false
	for _, c := range p.Cells[pos] {
		if c == cand {
			found = true
		} else {
			newcands = append(newcands, c)
		}
	}
	p.Cells[pos] = newcands
	return found
}

func (p *Board) Update() bool {
	updated := true
	for updated {
		updated = false
		if p.Update1() {
			updated = true
		}
		if p.Update2() {
			updated = true
		}
	}
	return updated
}

func (p *Board) Update1() bool {
	updated := false
	for pos := 0; pos < 9*9; pos++ {
		if len(p.Cells[pos]) == 1 {
			for peer := 0; peer < 9*9; peer++ {
				if IsPeer(pos, peer) {
					if p.Remove(peer, p.Cells[pos][0]) {
						updated = true
					}
				}
			}
		}
	}
	return updated
}

func (p *Board) Update2() bool {
	updated := false
	for pos := 0; pos < 9*9; pos++ {
		if len(p.Cells[pos]) > 1 {
			for _, v := range p.Cells[pos] {
				found := false
				for peer := 0; peer < 9*9; peer++ {
					if IsPeer(pos, peer) {
						if p.HasCandOf(peer, v) {
							found = true
							break
						}
					}
				}
				if !found {
					p.Cells[pos] = []int{v}
					updated = true
					break
				}
			}
		}
	}
	return updated
}

func (p *Board) ToS() string {
	s := ""
	border := "+---+---+---+---+---+---+---+---+---+\n"
	s += border
	for y := 0; y < 9; y++ {
		for j := 0; j < 3; j++ {
			s += "|"
			for x := 0; x < 9; x++ {
				for i := 0; i < 3; i++ {
					if p.HasCandOf(y*9+x, j*3+i+1) {
						s += strconv.Itoa(j*3 + i + 1)
					} else {
						if len(p.Cells[y*9+x]) == 1 {
							s += " "
						} else {
							s += "."
						}
					}
				}
				s += "|"
			}
			s += "\n"
		}
		s += border
	}
	return fmt.Sprintf("%s\n%d\n", s, p.Solved())
}
