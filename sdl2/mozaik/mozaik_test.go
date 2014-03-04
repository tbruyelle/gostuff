package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func setup() {
	g = NewGame()
}

func TestFindSwitch(t *testing.T) {
	setup()
	g.addSwitch(XMin, YMin,Red,Red,Red,Red)

	s := g.findSwitch(XMin+BlockSize, YMin+BlockSize)

	assert.NotNil(t, s, "Should found a switch")
}
