package main

type State interface {
	Enter(sw *Switch)
	Exit(sw *Switch)
	Update(g *Game, sw *Switch)
	AllowChange(state State) bool
}

type IdleState struct {
}

func NewIdleState() State {
	return &IdleState{}
}

func (s *IdleState) Enter(sw *Switch) {
}

func (s *IdleState) Exit(sw *Switch) {}

func (s *IdleState) Update(g *Game, sw *Switch) {
}

func (s *IdleState) AllowChange(state State) bool{
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

func (s *RotateState) Enter(sw *Switch) {
	s.startDegree = sw.rotate
}

func (s *RotateState) Exit(sw *Switch) {
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
