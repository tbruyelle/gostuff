package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsPlain(t *testing.T) {
	lvl := LoadLevel(1)

	assert.True(t, lvl.IsPlain(0))
	assert.False(t, lvl.IsPlain(1))
	assert.True(t, lvl.IsPlain(2))
}

func TestIsPlainRotated(t *testing.T) {
	lvl := LoadLevel(1)

	lvl.RotateSwitch(lvl.switches[1])
	lvl.RotateSwitch(lvl.switches[1])
	lvl.RotateSwitch(lvl.switches[2])
	lvl.RotateSwitch(lvl.switches[2])
	lvl.RotateSwitch(lvl.switches[0])

	assert.False(t, lvl.IsPlain(0))
}
