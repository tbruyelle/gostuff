package main

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
	blocks   []*Block
	switches []*Switch
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
	g.switches = []*Switch{
		&Switch{XMin + BlockSize - SwitchSize/2, YMin + BlockSize - SwitchSize/2},
	}
}

func (g *Game) Stop() {
}

func (g *Game) Click(x, y int) {}

func (g *Game) Update() {}

func (g *Game) Reset() {}
