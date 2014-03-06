package main

import (
	"fmt"
)

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
	g.level.rotating = sw
	sw.rotate = 0
}

func (s *RotateState) Exit(g *Game, sw *Switch) {
	// Swap bocks according to the 90d rotation
	l, c := sw.line, sw.col
	fmt.Println("Swap from", l, c)
	b := g.level.blocks[l][c]
	g.level.blocks[l][c] = g.level.blocks[l+1][c]
	g.level.blocks[l+1][c] = g.level.blocks[l+1][c+1]
	g.level.blocks[l+1][c+1] = g.level.blocks[l][c+1]
	g.level.blocks[l][c+1] = b

	g.level.rotating = nil
}

func (s *RotateState) Update(g *Game, sw *Switch) {
	sw.rotate += rotatePerTick
	if sw.rotate >= rotateDegree {
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

// rotateStateReverse is used to cancel a previous rotate
type RotateStateReverse struct {
	RotateState
}

func NewRotateStateReverse() State {
	return &RotateStateReverse{}
}

func (s *RotateStateReverse) Exit(g *Game, sw *Switch) {
	// Swap bocks according to the -90d rotation
	l, c := sw.line, sw.col
	fmt.Println("Reverse swap from", l, c)
	b := g.level.blocks[l][c]
	g.level.blocks[l][c] = g.level.blocks[l][c+1]
	g.level.blocks[l][c+1] = g.level.blocks[l+1][c+1]
	g.level.blocks[l+1][c+1] = g.level.blocks[l+1][c]
	g.level.blocks[l+1][c] = b

	g.level.rotating = nil
}

func (s *RotateStateReverse) Update(g *Game, sw *Switch) {
	sw.rotate -= rotatePerTick
	if sw.rotate <= -rotateDegree {
		sw.ChangeState(NewIdleState())
	}
}
