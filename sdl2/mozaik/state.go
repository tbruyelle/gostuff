package main

type State interface {
	Enter(s *Switch)
	Exit(s *Switch)
	Update(g *Game, s *Switch) bool
}

type IdleState struct {
}

func (s *IdleState) Enter(s *Switch) {}

func (s *IdleState) Exit(s *Switch) {}

func (s *IdleState) Update(g *Game, s *Switch) bool {
	return true
}

// RotateState performs a 90d rotation
type RotateState struct {
	IdleState
}

const rotateTicks=20

func (s *RotateState) Update() bool {

}
