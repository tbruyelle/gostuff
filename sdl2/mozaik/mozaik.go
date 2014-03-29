package main

import (
	"fmt"
)

const (
	WindowWidth        = 800
	WindowHeight       = 800
	BlockSize          = 128
	SwitchSize         = 48
	DashboardHeight    = 128
	XMin               = 32
	YMin               = 32
	XMax               = WindowHeight - 32
	YMax               = WindowWidth - 32 - DashboardHeight
	SignatureBlockSize = 32
)

type Game struct {
	currentLevel int
	level        Level
	listen       bool
}

func NewGame() *Game {
	return &Game{currentLevel: 1, listen: true}
}

func (g *Game) Start() {
	// Load first level
	g.level = LoadLevel(g.currentLevel)
}

func (g *Game) Stop() {
}

func (g *Game) Click(x, y int) {
	if g.listen {
		g.level.PressSwitch(x, y)
	}
}

func (g *Game) Update() {
	for _, s := range g.level.switches {
		if StateName(s.state) != "main.IdleState" {
			fmt.Println("state", s.name, StateName(s.state))
		}
		s.state.Update(g, s)
	}
}

func (g *Game) Continue() {
	if g.level.Win() {
		g.Warp()
	}
}
func (g *Game) Warp() {
	if g.listen {
		// Next level
		g.currentLevel++
		g.level = LoadLevel(g.currentLevel)
	}
}

func (g *Game) UndoLastMove() {
	if g.listen {
		g.level.UndoLastMove()
	}
}

func (g *Game) Reset() {
	sw := g.level.PopLastRotated()
	if sw != nil {
		g.listen = false
		sw.ChangeState(NewResetState())
	}
}
