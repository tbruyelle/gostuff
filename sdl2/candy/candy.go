package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

const (
	BlockSize      = 64
	NbBlockWidth   = 8
	NbBlockHeight  = 8
	DashboardWidth = 128
	WindowHeight   = BlockSize * NbBlockHeight
	WindowWidth    = DashboardWidth + BlockSize*NbBlockWidth
	YMin           = 0
	YMax           = WindowHeight - BlockSize
	XMin           = DashboardWidth
	XMax           = WindowWidth - BlockSize
	NbCandyType    = 6
)

type GameState int

const (
	Idle GameState = iota
	Matching
	Crushing
	Falling
	Translating
)

type Translation struct {
	c1, c2 *Candy
}

type Direction int

const (
	Left Direction = iota
	Top
	Right
	Bottom
	TopLeft
	TopRight
	BottomLeft
	BottomRight
	NbDirections
)

type Game struct {
	candys       []*Candy
	state        GameState
	selected     *Candy
	translation  *Translation
	flags        Flags
	candyTypeGen CandyTypeGenerator
}

type Flags struct {
	keepUnmatchingTranslation bool
}

type CandyTypeGenerator interface {
	NewCandyType() CandyType
}

type CandyTypeRandomizer struct {
	random *rand.Rand
}

func (c CandyTypeRandomizer) NewCandyType() CandyType {
	return CandyType(c.random.Intn(NbCandyType) + 1)
}

func NewGame() *Game {
	g := &Game{}
	g.candyTypeGen = CandyTypeRandomizer{rand.New(rand.NewSource(time.Now().Unix()))}
	g.state = Falling
	return g
}

func (g *Game) ToggleKeepUnmatchingTranslation() {
	g.flags.keepUnmatchingTranslation = !g.flags.keepUnmatchingTranslation
}

func (g *Game) Update() bool {
	allUpdated := true
	// Update all candys until all of them are
	// in a finished state.
	for _, c := range g.candys {
		if !c.Update(g) {
			allUpdated = false
		}
	}
	return allUpdated
}

func (g *Game) Tick() {
	if g.state == Falling {
		// Ensure new candys will appear on top
		// during the Falling state
		g.populateDropZone()
	}
	allUpdated := g.Update()
	//fmt.Printf("%d allupdated=%t\n", g.state, allUpdated)
	if allUpdated {
		// All candys updates, compute whats next
		switch g.state {
		case Crushing:
			// End of crushing, time to remove dead candys
			fmt.Println("Remove dead candys")
			var kept []*Candy
			for _, c := range g.candys {
				if !c.dead {
					kept = append(kept, c)
					c.ChangeState(NewFallingState())
				}
			}
			g.candys = kept
			g.state = Falling

			// TODO merge translating and falling
		case Falling:
			// End of falling, check if there is some match
			if g.matching() {
				// TODO put state=crusing in g.crushing()
				g.crushing()
				g.state = Crushing
			}

		case Translating:
			// End of translating, check if there a match
			fmt.Println("Translating")
			g.unselectAll()
			if g.matching() {
				g.crushing()
				g.state = Crushing
			} else {
				if !g.flags.keepUnmatchingTranslation {
					// revert translation
					fmt.Println("Revert permutation")
					g.translation.c1.ChangeState(NewPermuteState(g.translation.c2))
					g.translation.c2.ChangeState(NewPermuteState(g.translation.c1))
				}
				g.state = Idle
			}
			g.translation = nil
		}
	}
}

func withinLimits(x, y int) bool {
	return !(x < XMin || x > XMax+BlockSize || y < YMin || y > YMax+BlockSize)
}

// TODO Register clicks only in Idle state
func (g *Game) Click(x, y int) {
	if !withinLimits(x, y) {
		//fmt.Printf("Out of limits %d,%d\n", x, y)
		return
	}
	cx := determineXCandy(int(x))
	cy := determineYCandy(int(y))
	if c, found := findCandy(g.candys, cx, cy); found {
		//fmt.Printf("Found candy %d,%d, selected=%t\n", c.x, c.y, c.selected)
		if c == g.selected {
			// already selected unselect
			g.selected = nil
		} else {
			// not selected
			if g.selected != nil {
				// check if previous selection is near the new one
				if near(c, g.selected) {
					// init permutation
					g.translation = &Translation{c, g.selected}
					c.ChangeState(NewPermuteState(g.selected))
					g.selected.ChangeState(NewPermuteState(c))
					g.state = Translating
				} else {
					// remove previous selection
					g.selected = nil
				}
			}
			g.selected = c
		}
	}
}

func near(c1, c2 *Candy) bool {
	if c1.x == c2.x && math.Abs(float64(c1.y-c2.y)) == BlockSize {
		return true
	}
	if c1.y == c2.y && math.Abs(float64(c1.x-c2.x)) == BlockSize {
		return true
	}
	return false
}

func findCandy(candys []*Candy, x, y int) (*Candy, bool) {
	if !withinLimits(x, y) {
		return nil, false
	}
	for _, c := range candys {
		if c.y == y && c.x == x {
			return c, true
		}
	}
	return nil, false
}

func findCandyInDir(cs []*Candy, c *Candy, dir Direction) *Candy {
	switch dir {
	case Left:
		return leftCandy(cs, c)

	case Right:
		return rightCandy(cs, c)

	case Top:
		return topCandy(cs, c)

	case Bottom:
		return bottomCandy(cs, c)

	case TopLeft:
		return topLeftCandy(cs, c)

	case TopRight:
		return topRightCandy(cs, c)

	case BottomLeft:
		return bottomLeftCandy(cs, c)

	case BottomRight:
		return bottomRightCandy(cs, c)
	}
	return nil
}

func topCandy(candys []*Candy, c *Candy) *Candy {
	found, _ := findCandy(candys, c.x, c.y-BlockSize)
	return found
}

func bottomCandy(candys []*Candy, c *Candy) *Candy {
	found, _ := findCandy(candys, c.x, c.y+BlockSize)
	return found
}

func leftCandy(candys []*Candy, c *Candy) *Candy {
	found, _ := findCandy(candys, c.x-BlockSize, c.y)
	return found
}

func rightCandy(candys []*Candy, c *Candy) *Candy {
	found, _ := findCandy(candys, c.x+BlockSize, c.y)
	return found
}

func topLeftCandy(candys []*Candy, c *Candy) *Candy {
	found, _ := findCandy(candys, c.x-BlockSize, c.y-BlockSize)
	return found
}

func bottomRightCandy(candys []*Candy, c *Candy) *Candy {
	found, _ := findCandy(candys, c.x+BlockSize, c.y+BlockSize)
	return found
}

func bottomLeftCandy(candys []*Candy, c *Candy) *Candy {
	found, _ := findCandy(candys, c.x-BlockSize, c.y+BlockSize)
	return found
}

func topRightCandy(candys []*Candy, c *Candy) *Candy {
	found, _ := findCandy(candys, c.x+BlockSize, c.y-BlockSize)
	return found
}
func determineXCandy(x int) int {
	return determineYCandy(x-XMin) + XMin
}

func determineYCandy(y int) int {
	return determineCoord(y-YMin) + YMin
}

func determineCoord(c int) int {
	if c > BlockSize {
		return c - c%BlockSize
	}
	return 0
}

func (g *Game) Reset() {
	g.candys = nil
	g.state = Falling
}

func (g *Game) unselectAll() {
	g.selected = nil
}

// populateDropZone populates the top of the grid
// with new randomly generated candys.
func (g *Game) populateDropZone() {
	for i := 0; i < NbBlockWidth; i++ {
		c := g.newCandy()
		c.x = XMin + BlockSize*i
		// Apply -BlockSize to y to see the candy fall
		// from the top of the screen
		c.y = -BlockSize
		if !g.collideColumn(c) {
			// pop the candy
			c.ChangeState(NewFallingState())
			g.candys = append(g.candys, c)
		}
	}
}

func (g *Game) newCandy() *Candy {
	ct := g.candyTypeGen.NewCandyType()
	return NewCandy(CandyType(ct))
}

func collide(c1, c2 *Candy) bool {
	if c1.x+BlockSize <= c2.x {
		return false
	}
	if c1.y+BlockSize <= c2.y {
		return false
	}
	if c2.x+BlockSize <= c1.x {
		return false
	}
	if c2.y+BlockSize <= c1.y {
		return false
	}
	return true
}

// collideColumn returns true if the candy in parameter
// collides with another candy.
func (g *Game) collideColumn(newc *Candy) bool {
	for _, c := range g.candys {
		if c != newc && c.x == newc.x && collide(newc, c) {
			return true
		}
	}
	return false
}

func (g *Game) Destroy() {

}

func (g *Game) Start() {
}

func (g *Game) Stop() {
}
