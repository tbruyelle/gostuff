package main

import (
	"math/rand"
	"time"
)

const (
	BlockSize      = 32
	NbBlockWidth   = 16
	NbBlockHeight  = 16
	DashboardWidth = 128
	WindowHeight   = BlockSize * NbBlockHeight
	WindowWidth    = DashboardWidth + BlockSize*NbBlockWidth
	Match3         = 3
	Match4         = 4
	Match5         = 5
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
)

type Candy struct {
	_type   CandyType
	x, y, v int
}

type Column struct {
	candys []Candy
}

type Game struct {
	columns []Column
	random  *rand.Rand
	state   State
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
	for _, col := range g.columns {
		col.candys = make([]Candy, NbBlockHeight)
	}
	return g
}

func (g *Game) Tick() bool {
	switch g.state {
	case Idle:
		return true
	case Crushing:
	case Filling:
		g.applyVectors()

	}
	return false

}

func (g *Game) populateDropZone() {
	for _, col := range g.columns {
		col.candys[0] = g.newCandy()
	}
}

func applyVector(col *Column) {
	if len(col.candys) == 0 {
		return
	}
	c := &col.candys[0]
	for i := 1; i < len(col.candys); i++ {
		if col.candys[i]._type == EmptyCandy {
			c.v += BlockSize
		} else {
			c = &col.candys[i]
		}
	}
}

func (g *Game) applyVectors() {
	for _, col := range g.columns {
		applyVector(&col)
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
	var c int
	for c == 0 {
		c = g.random.Intn(4)
	}
	return Candy{_type: CandyType(c)}
}

func loopRowColumn(content func(i, j int)) {
	for i := 0; i < NbBlockWidth; i++ {
		for j := 0; j < NbBlockHeight; j++ {
			content(i, j)
		}
	}
}

func (g *Game) Destroy() {

}

func (g *Game) Start() {
}

func (g *Game) Stop() {
}
