package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIdleState(t *testing.T) {
	c := &Candy{}
	c.ChangeState(NewIdleState())

	r := c.Update()

	assert.True(t, r, "Update should return true")
	assert.Equal(t, c.sprite._type, CandySprite)
}

func TestDyingState(t *testing.T) {
	c := &Candy{}
	c.ChangeState(NewDyingState())

	r := c.Update()

	assert.False(t, r, "Update should return false")
	assert.Equal(t, c.sprite._type, DyingSprite)
}

func TestDyingStates(t *testing.T) {
	c := &Candy{}
	c.ChangeState(NewDyingState())

	for i := 0; i < DyingFrames; i++ {
		assert.False(t, c.Update(), "Until DyingFrames, Update should returns false (%d)", i)
	}
	r := c.Update()

	assert.True(t, r, "After being invoked DyingFrame times, Update should return true")
	assert.True(t, c.IsDead(), "Candy should be dead")
}

func TestDyingStatesDelayed(t *testing.T) {
	c := &Candy{}
	c.ChangeState(NewDyingStateDelayed(10))

	for i := 0; i < DyingFrames+10; i++ {
		assert.False(t, c.Update(), "Until DyingFrames+delay, Update should returns false (%d)", i)
	}
	r := c.Update()

	assert.True(t, r, "After being invoked DyingFrame+delay times, Update should return true")
	assert.True(t, c.IsDead(), "Candy should be dead")
}

func TestFallingState(t *testing.T) {
	setup()
	c:=&Candy{}
	c.ChangeState(NewFallingState())

	r:=c.Update()

	assert.False(t, r, "Update should return false")
}
