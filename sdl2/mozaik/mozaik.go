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
	blocks    []*Block
	switches  []*Switch
	listening bool
}

func NewGame() *Game {
	return &Game{}
}

func (g *Game) Start() {
	// Load first level
	g.blocks = []*Block{
		&Block{Red, XMin, YMin}, &Block{Blue, XMin + BlockSize, YMin},
		&Block{Yellow, XMin, YMin + BlockSize}, &Block{Green, XMin + BlockSize, YMin + BlockSize},
	}
	g.addSwitch(XMin, YMin)
	g.listening = true
}

// addSwitch appends a new switch at the bottom right
// of the coordinates in parameters.
func (g *Game) addSwitch(x, y int) {
	g.switches = append(g.switches, &Switch{X: x + BlockSize - SwitchSize/2, Y: y + BlockSize - SwitchSize/2})
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

func (g *Game) Update() {}

func (g *Game) Reset() {}
