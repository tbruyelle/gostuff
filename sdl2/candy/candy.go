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
	YMax           = WindowHeight - BlockSize
	NbCandyType    = 5
)

type State int

const (
	Idle State = iota
	Crushing
	Falling
	Translate
)

type CandyType int

const (
	EmptyCandy CandyType = iota
	RedCandy
	GreenCandy
	BlueCandy
	YellowCandy
	PinkCandy
)

type Candy struct {
	_type           CandyType
	x, y, vx, vy, g int
	selected        bool
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

func (g *Game) Click(x, y int32) {
	if x <= DashboardWidth {
		// ignore click on dashboard for now
		return
	}
	cx := determineXCandy(int(x))
	cy := determineYCandy(int(y))
	if c, found := g.findCandy(cx, cy); found {
		//fmt.Printf("Found candy %d,%d\n", c.x, c.y)
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
	g.state = Translate
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

func (g *Game) findCandy(x, y int) (*Candy, bool) {
	for _, c := range g.candys {
		//fmt.Printf("%d finding candy %d current %d\n", i, col.candys[i].y, y)
		if c.y == y && c.x == x {
			return c, true
		}
	}
	return nil, false
}

func determineXCandy(x int) int {
	return determineYCandy(x-DashboardWidth) + DashboardWidth
}

func determineYCandy(y int) int {
	if y > BlockSize {
		return y - y%BlockSize
	}
	return 0
}

func (g *Game) Reset() {
	g.candys = nil
	g.state = Falling
}

func (g *Game) Tick() bool {
	switch g.state {
	case Idle:
		return true
	case Crushing:
		g.applyGravity()
	case Falling:
		g.populateDropZone()
		g.applyGravity()
		if !g.fall() {
			g.populateDropZone()
			fmt.Println("move->idle")
			g.state = Idle
		}
	case Translate:
		if !g.translate() {
			g.unselectAll()
			g.state = Idle
		}
	}

	return false

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
		newc.x = DashboardWidth + BlockSize*i
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
