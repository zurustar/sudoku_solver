package main

import (
	"io/ioutil"
	"log"
	"os"
	"strconv"
)


func GetPosInfo(pos int) (int, int, int) {
	sx := pos % 9
	sy := (pos - sx) / 9
	gx := 2
	if 0 <= sx && sx < 3 {
		gx = 0
	} else if 3 <= sx && sx < 6 {
		gx = 1
	}
	gy := 2
	if 0 <= sy && sy < 3 {
		gy = 0
	} else if 3 <= sy && sy < 6 {
		gy = 1
	}
	sg := gy * 3 + gx
	return sx, sy, sg
}

type Board struct {
	Cells []*Cell
}

func NewBoard() *Board {
	p := new(Board)
	p.Cells = []*Cell{}
	for pos := 0; pos < 9*9; pos++ {
		p.Cells = append(p.Cells, NewCell())
	}
	for pos := 0; pos < 9*9; pos++ {
		p.Cells[pos].Init(pos)
	}
	return p
}

func (p *Board) Load(filename string) {
	fp, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer fp.Close()
	b, err := ioutil.ReadAll(fp)
	if err != nil {
		log.Fatal(err)
	}
	n := 0
	for i := 0; i < len(b); i++ {
		v, err := strconv.Atoi(string(b[i]))
		if err == nil {
			p.Cells[n].Set(v)
			n += 1
			if n >= 81 {
				break
			}
		}
	}
	if n != 81 {
		log.Fatal("invalid format")
	}
	p.Update()
	log.Println("\n" + p.ToS())
}

func (p *Board) Update() {
	updated := true
	for updated {
		updated = false
		if p.UpdateSub1() {
			updated = true
		}
		if p.UpdateSub2() {
			updated = true
		}
	}
}

func (p *Board) UpdateSub1() bool {
	updated := false
	for pos := 0; pos < 9*9; pos++ {
		if len(p.Cells[pos].Cands) == 1 {
			sx, sy, sg := GetPosInfo(pos)
			for peer := 0; peer < 9*9; peer ++ {
				if pos == peer {
					continue
				}
				px, py, pg := GetPosInfo(peer)
				if sx == px || sy == py || sg == pg {
					if p.Cells[peer].Remove(p.Cells[pos].Cands[0]) {
						updated = true
					}
				}
			}
		}
	}
	return updated
}

func (p *Board) UpdateSub2() bool {
	updated := false
	for y := 0; y < 9; y++ {
		for x := 0; x < 9; x++ {
			pos := y * 9 + x
			if len(p.Cells[pos].Cands) == 1 {
				continue
			}
/*
			for _, cand := range p.Cells[pos].Cands {
			}
*/
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
					if p.Cells[y*9+x].HasCandOf(j*3 + i + 1) {
						s += strconv.Itoa(j*3 + i + 1)
					} else {
						if len(p.Cells[y*9+x].Cands) == 1 {
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
	return s
}
