package main



type Cell struct {
	cands []int
	x int
	y int
}

func NewCell(pos int) *Cell {
	p := new(Cell)
	p.cands = []int{1,2,3,4,5,6,7,8,9}
	p.x = pos % 9
	p.y = (pos - p.x) / 9
	return p
}

func (p *Cell)Set(v int) {
	if (0 < v) && (v<=9) {
		p.cands = []int{v}
	}
}


