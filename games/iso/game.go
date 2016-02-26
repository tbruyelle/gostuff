package main

type Game struct {
	Board [][]int
}

type Block struct {
	Type int
	X, Y int
}

func NewGame() *Game {
	g := &Game{}
	g.Board = [][]int{
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		{0, 1, 1, 0},
		{0, 1, 1, 0},
		{2, 2, 0, 0},
		{0, 0, 0, 0},
	}
	return g
}

func (g *Game) Click(x, y int) {}
func (g *Game) Tick()          {}
func (g *Game) Start()         {}
func (g *Game) Stop()          {}
func (g *Game) Destroy()       {}
