package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func setup() {
	g = NewGame()
}

func fill() {
	g.blocks = make([][]*Block, 2)
	g.blocks[0] = make([]*Block, 2)
	g.blocks[1] = make([]*Block, 2)
	g.blocks[0][0] = &Block{Color: Red}
	g.blocks[0][1] = &Block{Color: Blue}
	g.blocks[1][1] = &Block{Color: Pink}
	g.blocks[1][0] = &Block{Color: Yellow}
	g.addSwitch(0, 0)
}

func TestFindSwitch(t *testing.T) {
	setup()
	fill()

	s := g.findSwitch(XMin+BlockSize, YMin+BlockSize)

	assert.NotNil(t, s, "Should found a switch")
}

func TestLoadLevel(t *testing.T) {
	setup()

	lvl := `01
24

0,0

24`
	g.LoadLevelStr(lvl)

	assert.Equal(t, 1, len(g.switches))
	assert.Equal(t, 0, g.switches[0].line)
	assert.Equal(t, 0, g.switches[0].col)
	assert.Equal(t, 2, len(g.blocks))
	assert.Equal(t, 2, len(g.blocks[0]))
	assert.Equal(t, 2, len(g.blocks[1]))
	assert.Equal(t, 0, g.blocks[0][0].Color)
	assert.Equal(t, 1, g.blocks[0][1].Color)
	assert.Equal(t, 2, g.blocks[1][0].Color)
	assert.Equal(t, 4, g.blocks[1][1].Color)
	assert.Equal(t, "24\n", g.winSignature)
}

func TestRotateState(t *testing.T) {
	setup()
	fill()

	g.switches[0].ChangeState(NewRotateState())
	g.switches[0].ChangeState(NewIdleState())

	assert.Equal(t, 2, len(g.blocks))
	assert.Equal(t, 2, len(g.blocks[0]))
	assert.Equal(t, 2, len(g.blocks[1]))
	assert.Equal(t, Yellow, g.blocks[0][0].Color)
	assert.Equal(t, Red, g.blocks[0][1].Color)
	assert.Equal(t, Blue, g.blocks[1][1].Color)
	assert.Equal(t, Pink, g.blocks[1][0].Color)
}

func TestBlockSignature(t *testing.T) {
	setup()

	g.LoadLevelStr(`01
24

0,0`)

	signature := `01
24
`
	assert.Equal(t, signature, g.BlockSignature())
}
