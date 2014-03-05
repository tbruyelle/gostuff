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
	g.addBlock(XMin, YMin, Red)
	g.addBlock(XMin+BlockSize, YMin, Red)
	g.addBlock(XMin, YMin+BlockSize, Red)
	g.addBlock(XMin+BlockSize, YMin+BlockSize, Red)
	g.addSwitch(0, 1, 2, 3)

	s := g.findSwitch(XMin+BlockSize, YMin+BlockSize)

	assert.NotNil(t, s, "Should found a switch")
}

func TestLoadLevel(t *testing.T) {
	setup()

	g.LoadLevelStr(`0,1
2,4

0,1,2,3`)

	assert.Equal(t, 1, len(g.switches))
	assert.Equal(t, 4, len(g.blocks))
	assert.Equal(t, 0, g.blocks[0].Color)
	assert.Equal(t, 1, g.blocks[1].Color)
	assert.Equal(t, 2, g.blocks[2].Color)
	assert.Equal(t, 4, g.blocks[3].Color)


}
