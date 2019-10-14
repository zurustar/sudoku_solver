package main

import (
	"os"
)

func main() {
	for i, filename := range os.Args {
		if i != 0 {
			b := NewBoard()
			b.Load(filename)
		}
	}
}
