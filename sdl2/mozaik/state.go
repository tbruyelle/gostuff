package main

type State interface {
	Enter(sw *Switch)
	Exit(sw *Switch)
	Update(g *Game, sw *Switch) bool
}

type IdleState struct {
}

func (s *IdleState) Enter(sw *Switch) {}

func (s *IdleState) Exit(sw *Switch) {}

func (s *IdleState) Update(g *Game, sw *Switch) bool {
	return true
}

// RotateState performs a 90d rotation
type RotateState struct {
	IdleState
}

const rotateTicks = 20

func (s *RotateState) Update() bool {
return true
}
