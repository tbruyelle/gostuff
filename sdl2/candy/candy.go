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

	NbCandyType = 6
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
)

type Candy struct {
	_type                      CandyType
	x, y, vx, vy, g            int
	selected                   bool
	visitedLine, visitedColumn bool
	crushme                    int
}

type Game struct {
	candys   []*Candy
	random   *rand.Rand
	state    State
	selected *Candy
}

type Match struct {
	start  int
	length int
}

var NoMatch = Match{}

func NewGame() *Game {
	g := &Game{}
	g.random = rand.New(rand.NewSource(time.Now().Unix()))
	g.state = Falling
	return g
}

func (g *Game) Tick() bool {
	switch g.state {
	case Idle:
		fmt.Println("Idle")
		return true
	case Matching:
		if g.matching() {
			g.state = Crushing
		} else {
			g.state = Idle
		}

	case Crushing:
		// remove crushed candys
		var cds []*Candy
		for _, c := range g.candys {
			if c.crushme == 0 {
				cds = append(cds, c)
			}
		}
		fmt.Printf("Crushing %d candys\n", len(g.candys)-len(cds))
		g.candys = cds
		// trigger the fall of new candys
		g.state = Falling

	case Falling:
		g.populateDropZone()
		g.applyGravity()
		if !g.fall() {
			g.populateDropZone()
			g.state = Matching
		}
	case Translating:
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

func (g *Game) matching() bool {
	match := false
	// remove selection
	for _, c := range g.candys {
		c.visitedColumn = false
		c.visitedLine = false
	}
	fmt.Println("check lines")
	for _, c := range g.candys {
		lines := g.findInLine(c, c._type)
		match = checkRegion(lines)||match
	}
	fmt.Println("check columns")
	for _, c := range g.candys {
		columns := g.findInColumn(c, c._type)
		if len(columns) > 1 {
			fmt.Printf("match coliu %v\n", columns)
		}
		match = checkRegion(columns)||match
	}
	return match
}

func checkRegion(region Region) bool {
	nbMatch := len(region) - 2
	if nbMatch > 0 {
		fmt.Printf("match region %v\n", region)
		for _, c := range region {
			c.crushme += nbMatch
		}
		return true
	}
	return false
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
		if c.selected {
			// already selected unselect
			c.selected = false
			g.selected = nil
		} else {
			// not selected
			if g.selected != nil {
				// check if previous selection is near the new one
				if near(c, g.selected) {
					// init permutation
					g.permute(c, g.selected)
				} else {
					// remove previous selection
					g.selected.selected = false
				}
			}
			c.selected = true
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
	for _, c := range g.candys {
		c.selected = false
	}
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

func matchType(t1, t2 CandyType) bool {
	return t1 == t2
}

func checkLine(line []CandyType) Match {
	var start, length int
	length = 1
	for i := 1; i < len(line); i++ {
		if line[start] == line[i] {
			length++
		} else {
			if length >= Match3 {
				return Match{start: start, length: length}
			}
			length = 1
			start = i
		}
	}
	return NoMatch
}

func checkGrid(candys [][]CandyType) []Match {
	matches := make([]Match, 0)
	// check lines
	for i := 0; i < len(candys); i++ {
		m := checkLine(candys[i])
		if m.length > 0 {
			matches = append(matches, m)
		}
	}
	return matches
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
