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
	Color    ColorDef
	Rendered bool
}

type Switch struct {
	state     State
	line, col int
	X, Y      int
	rotate    int
}

func (s *Switch) Rotate() {
	s.ChangeState(NewRotateState())
}

func (s *Switch) ChangeState(state State) {
	if s.state != nil {
		s.state.Exit(g, s)
		if !s.state.AllowChange(state) {
			fmt.Println("Change state not allowed")
			return
		}
	}
	s.state = state
	s.state.Enter(g, s)
}
