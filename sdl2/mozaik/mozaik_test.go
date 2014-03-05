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
	g.blocks = make([][]*Block, 2)
	g.blocks[0] = make([]*Block, 2)
	g.blocks[1] = make([]*Block, 2)
	g.blocks[0][0] = &Block{Color: Red}
	g.blocks[0][1] = &Block{Color: Red}
	g.blocks[1][0] = &Block{Color: Red}
	g.blocks[1][1] = &Block{Color: Red}
	g.addSwitch(0, 0)

	s := g.findSwitch(XMin+BlockSize, YMin+BlockSize)

	assert.NotNil(t, s, "Should found a switch")
}

func TestLoadLevel(t *testing.T) {
	setup()

	g.LoadLevelStr(`0,1
2,4

0,0`)

	assert.Equal(t, 1, len(g.switches))
	assert.Equal(t, 0, g.switches[0].bx)
	assert.Equal(t, 0, g.switches[0].by)
	assert.Equal(t, 2, len(g.blocks))
	assert.Equal(t, 2, len(g.blocks[0]))
	assert.Equal(t, 2, len(g.blocks[1]))
	assert.Equal(t, 0, g.blocks[0][0].Color)
	assert.Equal(t, 1, g.blocks[0][1].Color)
	assert.Equal(t, 2, g.blocks[1][0].Color)
	assert.Equal(t, 4, g.blocks[1][1].Color)
}
