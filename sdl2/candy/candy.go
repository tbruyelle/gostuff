package main

import (
	"fmt"
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
	Filling
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
	_type    CandyType
	x, y, v  int
	selected bool
}

type Column struct {
	candys []Candy
}

type Game struct {
	columns  []Column
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
	g.columns = make([]Column, NbBlockWidth)
	g.state = Filling
	return g
}

func (g *Game) Click(x, y int32) {
	if x <= DashboardWidth {
		// ignore click on dashboard for now
		return
	}
	// determine column
	col := &g.columns[determineColumn(int(x))]
	cy := determineYCandy(int(y))
	if c, found := findCandy(*col, cy); found {
		fmt.Printf("Found candy %d,%d\n", c.x, c.y)
		if c.selected {
			// already selected unselect
			c.selected = false
			g.selected = nil
		} else {
			// not selected
			c.selected = true
			g.selected = c
		}
	}
}

func findCandy(col Column, y int) (*Candy, bool) {
	for i := 0; i < len(col.candys); i++ {
		//fmt.Printf("%d finding candy %d current %d\n", i, col.candys[i].y, y)
		if col.candys[i].y == y {
			return &col.candys[i], true
		}
	}
	return nil, false
}

func determineColumn(x int) int {
	return (x - DashboardWidth) / BlockSize
}

func determineYCandy(y int) int {
	if y > BlockSize {
		return y - y%BlockSize
	}
	return 0
}

func (g *Game) Reset() {
	for i := 0; i < len(g.columns); i++ {
		g.columns[i].candys = nil
	}
	g.state = Filling
}

func (g *Game) Tick() bool {
	switch g.state {
	case Idle:
		return true
	case Crushing:
		g.applyVectors()
	case Filling:
		g.populateDropZone()
		g.applyVectors()
		if !g.move() {
			g.populateDropZone()
			fmt.Println("move->idle")
			g.state = Idle
		}
	}
	return false

}

func (g *Game) move() bool {
	moving := false
	for i := range g.columns {
		for j := range g.columns[i].candys {
			c := &g.columns[i].candys[j]
			if c.v > 0 {
				c.y += c.v
				if c.y < YMax && !collideColumnInd(j, g.columns[i]) {
					//fmt.Printf("moving %d -> %d\n", j, c.v)
					moving = true
				} else {
					// adjust y position according to the collision
					if c.y >= YMax {
						c.y = YMax
					} else {
						c.y--
						for collideColumnInd(j, g.columns[i]) {
							c.y--
						}
					}
					c.v = 0
				}

			}
		}
	}
	return moving
}

func (g *Game) populateDropZone() {
	for i := range g.columns {
		newc := g.newCandy()
		newc.x = DashboardWidth + BlockSize*i
		newc.y = 0
		col := &g.columns[i]
		if !collideColumn(newc, *col) {
			col.candys = append(col.candys, newc)
		}
	}
}

func (g *Game) applyVectors() {
	for i := 0; i < len(g.columns); i++ {
		col := &g.columns[i]
		for j := 0; j < len(col.candys); j++ {
			if col.candys[j].v == 0 {
				col.candys[j].v = i%2 + 1
			} else {
				col.candys[j].v++
			}
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

func (g *Game) newCandy() Candy {
	ct := g.random.Intn(NbCandyType) + 1
	return Candy{_type: CandyType(ct)}
}

func loopRowColumn(content func(i, j int)) {
	for i := 0; i < NbBlockWidth; i++ {
		for j := 0; j < NbBlockHeight; j++ {
			content(i, j)
		}
	}
}

func collide(c1, c2 Candy) bool {
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

func collideColumnInd(i int, col Column) bool {
	for j := 0; j < len(col.candys); j++ {
		if i != j && collide(col.candys[i], col.candys[j]) {
			//fmt.Printf("Collide between (%d,%d)/(%d,%d)\n", col.candys[i].x, col.candys[i].y, col.candys[j].x, col.candys[j].y)
			return true
		}
	}
	return false
}

func collideColumn(newc Candy, col Column) bool {
	for _, c := range col.candys {
		if collide(newc, c) {
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
