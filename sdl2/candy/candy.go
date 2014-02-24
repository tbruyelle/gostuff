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

type State int

const (
	Idle State = iota
	Matching
	Crushing
	Falling
	Translating
)

type CandyType int

const (
	EmptyCandy CandyType = iota
	RedCandy
	GreenCandy
	BlueCandy
	YellowCandy
	PinkCandy
	OrangeCandy
	RedHStripesCandy
	GreenHStripesCandy
	BlueHStripesCandy
	YellowHStripesCandy
	PinkHStripesCandy
	OrangeHStripesCandy
	RedVStripesCandy
	GreenVStripesCandy
	BlueVStripesCandy
	YellowVStripesCandy
	PinkVStripesCandy
	OrangeVStripesCandy
	RedPackedCandy
	GreenPackedCandy
	BluePackedCandy
	YellowPackedCandy
	PinkPackedCandy
	OrangePackedCandy
	BombCandy
)

type Candy struct {
	_type                      CandyType
	x, y, vx, vy, g            int
	visitedLine, visitedColumn bool
	// crush tells if the candy will be deleted on next Crush state
	crush bool
}

func (c *Candy) String() string {
	return fmt.Sprintf("(%d,%d)t%d,%t", c.x, c.y, c._type, c.crush)
}

func (c *Candy) isStriped() bool {
	return c._type > NbCandyType && c._type <= NbCandyType*3
}

func (c *Candy) isStripedH() bool {
	return c._type > NbCandyType && c._type <= NbCandyType*2
}

func (c *Candy) isStripedV() bool {
	return c._type > NbCandyType*2 && c._type <= NbCandyType*3
}

type Translation struct {
	c1, c2 *Candy
}

type Direction int

const (
	Left Direction = iota
	Top
	Right
	Bottom
)

type Game struct {
	candys          []*Candy
	state           State
	selected        *Candy
	translation     *Translation
	flags           Flags
	randomCandyType RandomCandyType
}

type Flags struct {
	keepUnmatchingTranslation bool
}

type RandomCandyType func() CandyType

type CandyTypeRandomizer struct {
	random *rand.Rand
}

func (c CandyTypeRandomizer) RandomCandyType() CandyType {
	return CandyType(c.random.Intn(NbCandyType) + 1)
}

func NewGame() *Game {
	g := &Game{}
	g.randomCandyType = CandyTypeRandomizer{rand.New(rand.NewSource(time.Now().Unix()))}
	g.state = Falling
	return g
}

func (g *Game) ToggleKeepUnmatchingTranslation() {
	g.flags.keepUnmatchingTranslation = !g.flags.keepUnmatchingTranslation
}

func (g *Game) Tick() bool {
	switch g.state {
	case Idle:
		fmt.Println("Idle")
		return true
	case Matching:
		fmt.Println("Matching")
		if g.matching() {
			g.state = Crushing
		} else {
			// no match
			if !g.flags.keepUnmatchingTranslation && g.translation != nil {
				// cancel previous translation
				g.permute(g.translation.c2, g.translation.c1)
			} else {
				g.state = Idle
			}
			g.translation = nil
		}

	case Crushing:
		fmt.Println("Crushing")
		g.crushing()
		g.translation = nil
		// trigger the fall of new candys
		g.state = Falling

	case Falling:
		fmt.Println("Falling")
		g.populateDropZone()
		g.applyGravity()
		if !g.fall() {
			g.populateDropZone()
			g.state = Matching
		}
	case Translating:
		fmt.Println("Translating")
		if !g.translate() {
			g.unselectAll()
			g.state = Matching
		}
	}

	return false
}

func withinLimits(x, y int) bool {
	return !(x < XMin || x > XMax+BlockSize || y < YMin || y > YMax+BlockSize)
}

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
					g.permute(c, g.selected)
				} else {
					// remove previous selection
					g.selected = nil
				}
			}
			g.selected = c
		}
	}
}

func (g *Game) permute(c1, c2 *Candy) {
	g.state = Translating
	if c1.x > c2.x {
		c1.vx = -BlockSize
		c2.vx = BlockSize
		return
	}
	if c1.x < c2.x {
		c1.vx = BlockSize
		c2.vx = -BlockSize
		return
	}
	if c1.y > c2.y {
		c1.vy = -BlockSize
		c2.vy = BlockSize
		return
	}
	// c1.y<c2.y
	c1.vy = BlockSize
	c2.vy = -BlockSize
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

var tSpeed = 4

func (g *Game) translate() bool {
	moving := false
	for _, c := range g.candys {
		if c.vx > 0 {
			c.x += tSpeed
			c.vx -= tSpeed
			moving = true
		} else if c.vx < 0 {
			c.x -= tSpeed
			c.vx += tSpeed
			moving = true
		}
		if c.vy > 0 {
			c.y += tSpeed
			c.vy -= tSpeed
			moving = true
		} else if c.vy < 0 {
			c.y -= tSpeed
			c.vy += tSpeed
			moving = true
		}
	}
	return moving

}

func (g *Game) fall() bool {
	falling := false
	for i, c := range g.candys {
		if c.g > 0 {
			c.y += c.g
			if c.y < YMax && !g.collideColumn(c, i) {
				//fmt.Printf("moving %d -> %d\n", j, c.g)
				falling = true
			} else {
				// adjust y position according to the collision
				if c.y >= YMax {
					c.y = YMax
				} else {
					c.y--
					for g.collideColumn(c, i) {
						c.y--
					}
				}
				c.g = 0
			}

		}
	}
	return falling
}

func (g *Game) populateDropZone() {
	for i := 0; i < NbBlockWidth; i++ {
		newc := g.newCandy()
		newc.x = XMin + BlockSize*i
		newc.y = 0
		if !g.collideColumn(newc, -1) {
			g.candys = append(g.candys, newc)
		}
	}
}

func (g *Game) applyGravity() {
	for _, c := range g.candys {
		if c.g == 0 {
			// init gravity
			c.g = 1 + (c.x / BlockSize % 2)
		} else {
			// increase gravity
			c.g++
		}
	}
}

func (g *Game) newCandy() *Candy {
	ct := g.randomCandyType()
	return &Candy{_type: CandyType(ct)}
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

func (g *Game) collideColumn(newc *Candy, ind int) bool {
	for i, c := range g.candys {
		if i != ind && c.x == newc.x && collide(newc, c) {
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
