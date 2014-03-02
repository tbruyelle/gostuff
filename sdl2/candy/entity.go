package main

import (
	"fmt"
)

type CandyType int

const (
	EmptyCandy CandyType = iota
	RedCandy
	GreenCandy
	BlueCandy
	YellowCandy
	PinkCandy
	OrangeCandy
	RedHStripesCandy
	GreenHStripesCandy
	BlueHStripesCandy
	YellowHStripesCandy
	PinkHStripesCandy
	OrangeHStripesCandy
	RedVStripesCandy
	GreenVStripesCandy
	BlueVStripesCandy
	YellowVStripesCandy
	PinkVStripesCandy
	OrangeVStripesCandy
	RedPackedCandy
	GreenPackedCandy
	BluePackedCandy
	YellowPackedCandy
	PinkPackedCandy
	OrangePackedCandy
	BombCandy
)

type Candy struct {
	// state represents the current state of the Candy
	state                      State
	// sprite represents how the candu will be rendered
	sprite                     Sprite
	_type                      CandyType
	x, y, vx, vy, g            int
	visitedLine, visitedColumn bool
	// crush tells if the candy will be deleted on next Crush state
	crush bool
	// dead tells the candy can be removed from the game objects collection.
	dead bool
}

func NewCandy(_type CandyType) *Candy {
	c := &Candy{_type: _type}
	c.ChangeState(NewIdleState())
	return c
}

func (c *Candy) Update() {
	c.state.Update(c)
}

func (c *Candy) ChangeState(state State) {
	if c.state != nil {
		state.Exit(c)
	}
	c.state = state
	c.state.Enter(c)
}

func (c *Candy) String() string {
	return fmt.Sprintf("(%d,%d)t%d,%t-%t", c.x, c.y, c._type, c.crush, c.dead)
}

// IsNormal() returns true if the candy isn't special
func (c *Candy) IsNormal() bool {
	return c._type > 0 && c._type <= NbCandyType
}

func (c *Candy) IsStriped() bool {
	return c._type > NbCandyType && c._type <= NbCandyType*3
}

func (c *Candy) IsStripedH() bool {
	return c._type > NbCandyType && c._type <= NbCandyType*2
}

func (c *Candy) IsStripedV() bool {
	return c._type > NbCandyType*2 && c._type <= NbCandyType*3
}

func (c *Candy) IsPacked() bool {
	return c._type > NbCandyType*3 && c._type <= NbCandyType*4
}
