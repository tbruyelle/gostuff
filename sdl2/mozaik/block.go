package main

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
	X, Y  int
}

type Switch struct {
	X, Y   int
	blocks [4]*Block
}

func (s *Switch) Rotate() {

}
