package main

import (
	"fmt"
)

type ColorDef int

const (
	Red ColorDef = iota
	Yellow
	Blue
	Green
	Pink
)

type Block struct {
	Color ColorDef
}

type Switch struct {
	X, Y   int
	blocks [4]*Block
}

func (s *Switch) Rotate() {
	fmt.Println("Rotate", s)
}
