package main

import (
	"log"
)

type Cell struct {
	Cands []int
	Pos   int
	Peers []int
}

func toG(n int) int {
	if 0 <= n && n <= 2 {
		return 0
	}
	if 3 <= n && n <= 5 {
		return 1
	}
	if 6 <= n && n <= 8 {
		return 2
	}
	log.Fatal("invalid value")
	return -1
}

func NewCell() *Cell {
	p := new(Cell)
	p.Cands = []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	p.Peers = []int{}
	return p
}

func (p *Cell) Init(pos int) {
	p.Cands = []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	p.Pos = pos
	sx := pos % 9
	sy := (pos - sx) / 9
	sg := toG(sy)*3 + toG(sx)
	for opos := 0; opos < 9*9; opos++ {
		if pos != opos {
			ox := opos % 9
			oy := (opos - ox) / 9
			og := toG(oy)*3 + toG(ox)
			if ox == sx || oy == sy || og == sg {
				p.Peers = append(p.Peers, opos)
			}
		}
	}
}

func (p *Cell) Set(v int) {
	if (0 < v) && (v <= 9) {
		p.Cands = []int{v}
	}
}

func (p *Cell) HasCandOf(v int) bool {
	for _, c := range p.Cands {
		if v == c {
			return true
		}
	}
	return false
}

func (p *Cell) Remove(v int) bool {
	newcands := []int{}
	found := false
	for _, c := range p.Cands {
		if c == v {
			found = true
		} else {
			newcands = append(newcands, c)
		}
	}
	p.Cands = newcands
	return found
}
