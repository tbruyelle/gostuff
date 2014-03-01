package main

import (
	"fmt"
)

type StateType int

const (
	IdleStateType StateType = iota
	DyingStateType
)

// State exposes the state methods
type State interface {
	Enter(c *Candy)
	Exit(c *Candy)
	Update(c *Candy)
	Type() StateType
}

type baseState struct{}

func (s *baseState) Enter(c *Candy) {
	c.sprite = NewSprite(CandySprite)
}

func (s *baseState) Exit(c *Candy) {
}

func (s *baseState) Update(c *Candy) {
}

type idleState struct {
	baseState
}

func (s *idleState) Type() StateType {
	return IdleStateType
}

// dyingState controls how the candy will die
type dyingState struct {
	baseState
	beforeDie int
}

func (s *dyingState) Enter(c *Candy) {
	c.sprite = NewSprite(DyingSprite)
}

func (s *dyingState) Update(c *Candy) {
	s.beforeDie--
	if s.beforeDie == 0 {
		c.dead = true
	} else if s.beforeDie <= c.sprite.nbframes {
		c.sprite.frame++
	}
	fmt.Printf("Update dying state beforeDie=%d candy=%v\n", s.beforeDie, c)
}

func (s *dyingState) Type() StateType {
	return DyingStateType
}

func NewIdleState() State {
	return &idleState{}
}

func NewDyingState() State {
	return &dyingState{beforeDie: DyingFrames}
}

func NewDyingStateDelayed(delay int) State {
	return &dyingState{beforeDie: delay + DyingFrames}
}
