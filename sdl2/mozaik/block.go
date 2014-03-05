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
}

type Switch struct {
	state  State
	X, Y   int
	blocks [4]*Block
	rotate int
}

func (s *Switch) Rotate() {
	s.ChangeState(NewRotateState())
}

func (s *Switch) ChangeState(state State) {
	if s.state != nil {
		s.state.Exit(s)
	}
	s.state = state
	s.state.Enter(s)
}
