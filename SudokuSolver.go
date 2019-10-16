package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

var p2x []int
var p2y []int
var p2g []int

var p2xf [][]int
var p2yf [][]int
var p2gf [][]int

func IsPeer(a, b int) bool {
	if a == b {
		return false
	}
	if p2x[a] == p2x[b] || p2y[a] == p2y[b] || p2g[a] == p2g[b] {
		return true
	}
	return false
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
				for _, peer := range p2xf[pos] {
					if p.HasCandOf(peer, v) {
						found = true
						break
					}
				}
				if !found {
					p.Cells[pos] = []int{v}
					break
				}
				found = false
				for _, peer := range p2yf[pos] {
					if p.HasCandOf(peer, v) {
						found = true
						break
					}
				}
				if !found {
					p.Cells[pos] = []int{v}
					break
				}
				found = false
				for _, peer := range p2gf[pos] {
					if p.HasCandOf(peer, v) {
						found = true
						break
					}
				}
				if !found {
					p.Cells[pos] = []int{v}
					break
				}
			}
		}
	}
	return updated
}

func (p *Board) ToS() string {
	border := "+---+---+---+---+---+---+---+---+---+\n"
	s := border
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

func main() {
	// initialize global variables
	n2g := []int{0, 0, 0, 1, 1, 1, 2, 2, 2}
	for y := 0; y < 9; y++ {
		for x := 0; x < 9; x++ {
			p2x = append(p2x, x)
			p2y = append(p2y, y)
			p2g = append(p2g, n2g[y]*3+n2g[x])
		}
	}
	for pos := 0; pos < 9*9; pos++ {
		p2xf = append(p2xf, []int{})
		p2yf = append(p2yf, []int{})
		p2gf = append(p2gf, []int{})
		for peer := 0; peer < 9*9; peer++ {
			if pos != peer {
				if p2x[pos] == p2x[peer] {
					p2xf[pos] = append(p2xf[pos], peer)
				}
				if p2y[pos] == p2y[peer] {
					p2yf[pos] = append(p2yf[pos], peer)
				}
				if p2g[pos] == p2g[peer] {
					p2gf[pos] = append(p2gf[pos], peer)
				}
			}
		}
	}
	// load and solve
	for i, filename := range os.Args {
		if i != 0 {
			b := NewBoard(filename)
			b.Update()
			fmt.Println(b.ToS())
			Solve(b, 0)
		}
	}
}

func Solve(b *Board, depth uint64) uint64 {
	b.Update()
	result := b.Solved()
	if result == 1 {
		fmt.Println(b.ToS())
		os.Exit(0)
	}
	if result == 0 {
		return 0
	}
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
			tmpb.Cells[pos] = []int{c}
			Solve(tmpb, depth+1)
		}
	}
	return 0
}
