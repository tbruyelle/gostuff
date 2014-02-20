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
	Match3         = 3
	Match4         = 4
	Match5         = 5
	Speed          = 7
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
	RedStripesCandy
	GreenStripesCandy
	BlueStripesCandy
	YellowStripesCandy
	PinkStripesCandy
	OrangeStripesCandy
	RedPackedCandy
	GreenPackedCandy
	BluePackedCandy
	YellowPackedCandy
	PinkPackedCandy
	OrangePackedCandy
	BombCandy
	// CrushCandy is used to declare a candy crushable on next Crushing state
	CrushCandy
)

type Candy struct {
	_type                      CandyType
	x, y, vx, vy, g            int
	visitedLine, visitedColumn bool
	// crush indicates how the candy will be the next time
	// the application reaches the Crushing state.
	crush CandyType
}

func (c *Candy) String() string {
	return fmt.Sprintf("(%d,%d)t%dc%d", c.x, c.y, c._type, c.crush)
}

type Translation struct {
	c1, c2 *Candy
}

type Game struct {
	candys      []*Candy
	random      *rand.Rand
	state       State
	selected    *Candy
	translation *Translation
	flags       Flags
}

type Flags struct {
	keepUnmatchingTranslation bool
}

func NewGame() *Game {
	g := &Game{}
	g.random = rand.New(rand.NewSource(time.Now().Unix()))
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
			g.translation = nil
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

// remove crushed candys
func (g *Game) crushing() {
	var cds []*Candy
	for _, c := range g.candys {
		fmt.Printf("crushCandy %d t=%d\n", c.crush, c._type)
		if c.crush != CrushCandy {
			if c.crush != EmptyCandy {
				c._type = c.crush
				c.crush = EmptyCandy
			}
			cds = append(cds, c)
		}
	}
	fmt.Printf("Crushing %d candys\n", len(g.candys)-len(cds))
	g.candys = cds
	fmt.Printf("NOW %d candys\n", len(g.candys))
}

func withinLimits(x, y int) bool {
	return !(x < XMin || x > XMax+BlockSize || y < YMin || y > YMax+BlockSize)
}


func alligned(candys []*Candy) bool {
	xaligned := false
	for i := 1; i < len(candys); i++ {
		xaligned = candys[i-1].x == candys[i].x
		if !xaligned {
			break
		}
	}

	yaligned := false
	for i := 1; i < len(candys); i++ {
		yaligned = candys[i-1].y == candys[i].y
		if !yaligned {
			break
		}
	}
	return xaligned || yaligned
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
	ct := g.random.Intn(NbCandyType) + 1
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
