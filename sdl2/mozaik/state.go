package main

type State interface {
	Enter(g *Game, sw *Switch)
	Exit(g *Game, sw *Switch)
	Update(g *Game, sw *Switch)
	AllowChange(state State) bool
}

type IdleState struct {
}

func NewIdleState() State {
	return &IdleState{}
}

func (s *IdleState) Enter(g *Game, sw *Switch) {
}

func (s *IdleState) Exit(g *Game, sw *Switch) {}

func (s *IdleState) Update(g *Game, sw *Switch) {
}

func (s *IdleState) AllowChange(state State) bool {
	switch state.(type) {
	case *IdleState:
		return false
	}
	return true
}

// RotateState performs a 90d rotation
type RotateState struct {
	IdleState
	startDegree int
}

func NewRotateState() State {
	return &RotateState{}
}

const (
	rotateTicks   = 15
	rotateDegree  = 90
	rotatePerTick = rotateDegree / rotateTicks
)

func (s *RotateState) Enter(g *Game, sw *Switch) {
	g.rotating = sw
	s.startDegree = sw.rotate
}

func (s *RotateState) Exit(g *Game, sw *Switch) {
	g.rotating = nil
}

func (s *RotateState) Update(g *Game, sw *Switch) {
	sw.rotate += rotatePerTick
	if sw.rotate-s.startDegree >= rotateDegree {
		sw.ChangeState(NewIdleState())
	}
}

func (s *RotateState) AllowChange(state State) bool {
	switch state.(type) {
	case *RotateState:
		return false
	}
	return true
}
