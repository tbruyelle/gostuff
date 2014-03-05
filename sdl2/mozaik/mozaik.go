package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

const (
	WindowWidth        = 800
	WindowHeight       = 800
	BlockSize          = 128
	SwitchSize         = 48
	DashboardHeight    = 256
	XMin               = 32
	YMin               = 32
	XMax               = WindowHeight - 32
	YMax               = WindowWidth - 32 - DashboardHeight
	SignatureBlockSize = 48
)

type Game struct {
	blocks   [][]*Block
	switches []*Switch
	// rotating represents a rotate which
	// is currently rotating
	rotating     *Switch
	winSignature string
	currentLevel int
}

func NewGame() *Game {
	return &Game{currentLevel: 1}
}

func (g *Game) Start() {
	// Load first level
	g.LoadLevel()
}

// addSwitch appends a new switch at the bottom right
// of the coordinates in parameters.
func (g *Game) addSwitch(line, col int) {

	s := &Switch{
		line: line, col: col,
		X: XMin + col*BlockSize + BlockSize - SwitchSize/2,
		Y: YMin + line*BlockSize + BlockSize - SwitchSize/2,
	}
	s.ChangeState(NewIdleState())
	g.switches = append(g.switches, s)
	fmt.Println("Switch added", s.X, s.Y)
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
	// Handle click only when no switch are rotating
	if g.rotating == nil {
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

func (g *Game) Continue() {
	if g.Win() {
		// Next level
		g.currentLevel++
		g.LoadLevel()
	}
}

func (g *Game) Reset() {
	g.Stop()
	g.Start()
}

func (g *Game) LoadLevel() {
	b, err := ioutil.ReadFile(fmt.Sprintf("./levels/%d", g.currentLevel))
	if err != nil {
		panic(err)
	}
	g.LoadLevelStr(string(b))
}

func (g *Game) LoadLevelStr(str string) {
	lines := strings.Split(str, "\n")
	step := 0
	g.blocks = nil
	g.switches = nil
	g.winSignature=""

	for i := 0; i < len(lines); i++ {
		if len(lines[i]) == 0 {
			step++
			continue
		}
		tokens := strings.Split(lines[i], ",")
		switch step {
		case 0:
			// read block colors
			bline := make([]*Block, len(tokens))
			g.blocks = append(g.blocks, bline)
			for j := 0; j < len(tokens); j++ {
				bline[j] = &Block{Color: atoc(tokens[j])}
			}
		case 1:
			// reead switch locations
			g.addSwitch(atoi(tokens[0]), atoi(tokens[1]))
		case 2:
			//read win
			g.winSignature += lines[i] + "\n"
		}
	}
	fmt.Printf("Level loaded blocks=%d, swicthes=%d\n", len(g.blocks), len(g.switches))

	for i := 0; i < len(g.blocks); i++ {
		fmt.Printf("line %d blocks %d\n", i, len(g.blocks[i]))
	}
	fmt.Printf("winSignature\n%s\n---\n", g.winSignature)
}

func (g *Game) Win() bool {
	return g.winSignature == g.BlockSignature()
}

func (g *Game) BlockSignature() string {
	var signature string
	for i := 0; i < len(g.blocks); i++ {
		for j := 0; j < len(g.blocks[i]); j++ {
			signature += ctoa(g.blocks[i][j].Color)
		}
		signature += "\n"
	}
	return signature
}

func atoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}

func atoc(s string) ColorDef {
	return ColorDef(atoi(s))
}

func ctoa(c ColorDef) string {
	return fmt.Sprintf("%d", c)
}
