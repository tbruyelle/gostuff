package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

const (
	WindowWidth  = 800
	WindowHeight = 600
	BlockSize    = 128
	SwitchSize   = 32
	XMin         = 32
	YMin         = 32
	XMax         = WindowHeight - 32
	YMax         = WindowWidth - 32
)

type Game struct {
	blocks   []*Block
	switches []*Switch
	// rotating represents a rotate which
	// is currently rotating
	rotating *Switch
}

func NewGame() *Game {
	return &Game{}
}

func (g *Game) Start() {
	// Load first level
	g.LoadLevel(1)
}

func (g *Game) addBlock(x, y int, color ColorDef) {
	b := &Block{X: x, Y: y, Color: color}
	g.blocks = append(g.blocks, b)
}

// addSwitch appends a new switch at the bottom right
// of the coordinates in parameters.
func (g *Game) addSwitch(b1, b2, b3, b4 int) {

	s := &Switch{
		X: g.blocks[b1].X + BlockSize - SwitchSize/2,
		Y: g.blocks[b1].Y + BlockSize - SwitchSize/2,
	}
	s.blocks[0] = g.blocks[b1]
	s.blocks[1] = g.blocks[b2]
	s.blocks[2] = g.blocks[b3]
	s.blocks[3] = g.blocks[b4]
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

func (g *Game) Reset() {}

func (g *Game) LoadLevel(lvl int) {
	b, err := ioutil.ReadFile(fmt.Sprintf("./levels/%d", lvl))
	if err != nil {
		panic(err)
	}
	g.LoadLevelStr(string(b))
}

func (g *Game) LoadLevelStr(str string) {
	lines := strings.Split(str, "\n")
	handleBlocks := true
	cx, cy := XMin, YMin
	for i := 0; i < len(lines); i++ {
		if len(lines[i]) == 0 {
			handleBlocks = false
			continue
		}
		tokens := strings.Split(lines[i], ",")
		if handleBlocks {
			for j := 0; j < len(tokens); j++ {
				g.addBlock(cx, cy, atoc(tokens[j]))
				cx += BlockSize
			}
			cx = XMin
			cy += BlockSize
		} else {
			g.addSwitch(atoi(tokens[0]), atoi(tokens[1]),
				atoi(tokens[2]), atoi(tokens[3]))
		}
	}
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
