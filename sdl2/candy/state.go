package main

import (
	"fmt"
)

type StateType int

const (
	IdleStateType StateType = iota
	DyingStateType
	FallingStateType
	PermuteStateType
)

const (
	FallingInitSpeed = 8
	PermuteInitSpeed = 4
)

// State exposes the state methods
type State interface {
	Enter(c *Candy)
	Exit(c *Candy)
	// Update returns true when the state is in a finished state
	Update(g *Game, c *Candy) bool
	Type() StateType
}

type baseState struct{}

func (s *baseState) Enter(c *Candy) {
	c.sprite = NewSprite(CandySprite)
}

func (s *baseState) Exit(c *Candy) {
}

func (s *baseState) Update(g *Game, c *Candy) bool {
	return true
}

// idleState controls idle entities
type idleState struct {
	baseState
}

func NewIdleState() State {
	return &idleState{}
}

func (s *idleState) Enter(c *Candy) {
	c.crush = false
}

func (s *idleState) Type() StateType {
	return IdleStateType
}

// dyingState controls how the entity will die
type dyingState struct {
	baseState
	beforeDie int
}

func NewDyingState() State {
	return &dyingState{beforeDie: DyingFrames}
}

func NewDyingStateDelayed(delay int) State {
	return &dyingState{beforeDie: delay + DyingFrames}
}

func (s *dyingState) Enter(c *Candy) {
	c.crush = true
	c.sprite = NewSprite(DyingSprite)
}

func (s *dyingState) Update(g *Game, c *Candy) bool {
	if c.dead {
		return true
	}
	s.beforeDie--
	if s.beforeDie == 0 {
		c.dead = true
	} else if s.beforeDie <= c.sprite.nbframes {
		c.sprite.frame++
	}
	//fmt.Printf("Update dying state beforeDie=%d candy=%v\n", s.beforeDie, c)
	return false
}

func (s *dyingState) Type() StateType {
	return DyingStateType
}

// fallingState controls the fall of entities
type fallingState struct {
	baseState
}

func NewFallingState() State {
	return &fallingState{}
}

func (s *fallingState) Type() StateType {
	return FallingStateType
}

func (s *fallingState) Enter(c *Candy) {
	// initiaze gravity
	c.vy = FallingInitSpeed - 1
}

func (s *fallingState) Update(g *Game, c *Candy) bool {
	// increase gravity
	c.vy++
	// apply speed to coordinate
	c.y += c.vy
	// check collision
	if c.y < YMax && !g.collideColumn(c) {
		// no collision, still falling
		return false
	}
	// collision detected, adjust y position
	if c.y >= YMax {
		c.y = YMax
	} else {
		c.y--
		for g.collideColumn(c) {
			c.y--
		}
	}
	c.vy = 0
	c.ChangeState(NewIdleState())
	return true
}

// permuteState controls the candy translation
type permuteState struct {
	baseState
	// x,y represents the destination
	x, y int
}

func NewPermuteState(buddy *Candy) State {
	return &permuteState{x: buddy.x, y: buddy.y}
}

func (s *permuteState) Type() StateType {
	return PermuteStateType
}

func (s *permuteState) Enter(c *Candy) {
	if c.x > s.x {
		c.vx = -PermuteInitSpeed
	} else if c.x < s.x {
		c.vx = PermuteInitSpeed
	} else if c.y > s.y {
		c.vy = -PermuteInitSpeed
	} else {
		// c.y<s.y
		c.vy = PermuteInitSpeed
	}
	fmt.Printf("Enter permuteState c=%v, vx=%d, vy=%d, ox=%d, oy=%d\n", c, c.vx, c.vy, s.x, s.y)
}

func (s *permuteState) Update(g *Game, c *Candy) bool {
	fmt.Printf("Permuting %v\n", c)
	if c.vx != 0 {
		c.x += c.vx
	} else if c.vy != 0 {
		c.y += c.vy
	}
	if c.x == s.x && c.y == s.y {
		// Permutation done
		c.ChangeState(NewIdleState())
		return true
	}
	return false
}

func (s *permuteState) Exit(c *Candy) {
	c.vx = 0
	c.vy = 0
}
