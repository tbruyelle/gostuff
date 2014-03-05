package main

import (
	"fmt"
)

const (
	WindowWidth    = 800
	WindowHeight   = 600
	DashboardWidth = 256
	BlockSize      = 128
	SwitchSize     = 32
	XMin           = DashboardWidth + 32
	YMin           = 32
	XMax           = WindowHeight - 32
	YMax           = WindowWidth - 32
)

type Game struct {
	switches  []*Switch
	listening bool
}

func NewGame() *Game {
	return &Game{}
}

func (g *Game) Start() {
	// Load first level
	g.addSwitch(XMin, YMin, Red, Blue, Yellow, Green)
	g.listening = true
}

// addSwitch appends a new switch at the bottom right
// of the coordinates in parameters.
func (g *Game) addSwitch(x, y int, c1, c2, c3, c4 ColorDef) {
	s := &Switch{X: x + BlockSize - SwitchSize/2, Y: y + BlockSize - SwitchSize/2}
	s.blocks[0] = &Block{c1}
	s.blocks[1] = &Block{c2}
	s.blocks[2] = &Block{c3}
	s.blocks[3] = &Block{c4}
	s.ChangeState(NewIdleState())
	g.switches = append(g.switches, s)
}

func (g *Game) findSwitch(x, y int) *Switch {
	for _, s := range g.switches {
		if x >= s.X && x <= s.X+SwitchSize && y >= s.Y && y <= s.Y+SwitchSize {
			return s
		}
	}
	return nil
}

func (g *Game) Stop() {
}

func (g *Game) Click(x, y int) {
	fmt.Println("click", x, y)
	if g.listening {
		if s := g.findSwitch(x, y); s != nil {
			s.Rotate()
		}
	}
}

func (g *Game) Update() {
	for _, s := range g.switches {
		s.state.Update(g, s)
	}
}

func (g *Game) Reset() {}
