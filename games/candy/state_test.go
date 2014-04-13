package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIdleState(t *testing.T) {
	setup()
	c := &Candy{}
	c.ChangeState(NewIdleState())

	r := c.Update(g)

	assert.True(t, r, "Update should return true")
	assert.Equal(t, c.sprite._type, CandySprite)
}

func TestDyingState(t *testing.T) {
	c := &Candy{}
	c.ChangeState(NewDyingState())

	r := c.Update(g)

	assert.False(t, r, "Update should return false")
	assert.Equal(t, c.sprite._type, DyingSprite)
}

func TestDyingStates(t *testing.T) {
	setup()
	c := &Candy{}
	c.ChangeState(NewDyingState())

	for i := 0; i < DyingFrames; i++ {
		assert.False(t, c.Update(g), "Until DyingFrames, Update should returns false (%d)", i)
	}
	r := c.Update(g)

	assert.True(t, r, "After being invoked DyingFrame times, Update should return true")
	assert.True(t, c.IsDead(), "Candy should be dead")
}

func TestDyingStatesDelayed(t *testing.T) {
	setup()
	c := &Candy{}
	c.ChangeState(NewDyingStateDelayed(10))

	for i := 0; i < DyingFrames+10; i++ {
		assert.False(t, c.Update(g), "Until DyingFrames+delay, Update should returns false (%d)", i)
	}
	r := c.Update(g)

	assert.True(t, r, "After being invoked DyingFrame+delay times, Update should return true")
	assert.True(t, c.IsDead(), "Candy should be dead")
}

func TestFallingState(t *testing.T) {
	setup()
	c := &Candy{}
	c.ChangeState(NewFallingState())

	r := c.Update(g)

	assert.False(t, r, "Update should return false")
}

func TestFallingStates(t *testing.T) {
	setup()
	c := &Candy{x: XMin, y: YMin}
	c.ChangeState(NewFallingState())

	for !c.Update(g) {
	}

	assert.Equal(t, c.x, XMin, "x should not change during fall")
	assert.Equal(t, c.y, YMax, "candy should fall until the bottom")
}

func TestFallingStatesMulti(t *testing.T) {
	setup()
	c1 := &Candy{x: XMin, y: YMin}
	c1.ChangeState(NewFallingState())
	c2 := &Candy{x: XMin + BlockSize, y: YMin}
	c2.ChangeState(NewFallingState())
	// add candy to create collisions
	g.candys = append(g.candys, &Candy{x: c1.x, y: YMin + 2*BlockSize})
	g.candys = append(g.candys, &Candy{x: c2.x, y: YMin + 5*BlockSize})

	for !c1.Update(g) || !c2.Update(g) {
	}

	assert.Equal(t, c1.y, YMin+BlockSize, "candy should fall until the bottom")
	assert.Equal(t, c2.y, YMin+4*BlockSize, "candy should fall until the bottom")
}

func TestPermuteState(t *testing.T) {
	c := &Candy{x: XMin, y: YMin}
	c2 := &Candy{x: XMin + BlockSize, y: YMin}
	c.ChangeState(NewPermuteState(c2))

	for !c.Update(g) {
	}

	assert.Equal(t, c.x, c2.x, "Candy should permute to other candy")
}

func TestPermuteStateExit(t *testing.T) {
	c := &Candy{x: XMin, y: YMin}
	c2 := &Candy{x: XMin + BlockSize, y: YMin}
	c.ChangeState(NewPermuteState(c2))

	c.ChangeState(NewIdleState())

	assert.Equal(t, c.vx, 0, "Quit permute state should reset vx")
	assert.Equal(t, c.vy, 0, "Quit permute state should reset vy")
}

func TestMorphState(t *testing.T) {
	setup()
	c := &Candy{}
	c.ChangeState(NewIdleState())
	c.ChangeState(NewMorphState(RedPackedCandy))

	for !c.Update(g) {
	}

	assert.Equal(t, c._type, RedPackedCandy, "Candy should have morph to RedPackedCandy")
}
