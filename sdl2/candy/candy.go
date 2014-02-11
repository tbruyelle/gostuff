package main

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	BlockSize      = 32
	NbBlockWidth   = 32
	NbBlockHeight  = 32
	DashboardWidth = 76
	WindowHeight   = BlockSize * NbBlockHeight
	WindowWidth    = DashboardWidth + BlockSize*NbBlockWidth
	Match3         = 3
	Match4         = 4
	Match5         = 5
)

type CandyType int

const (
	EmptyCandy CandyType = iota
	RedCandy
	GreenCandy
	BlueCandy
	YellowCandy
)

type Game struct {
	Candys [NbBlockWidth][NbBlockHeight]CandyType
	random *rand.Rand
}

type Match struct {
	start  int
	length int
}

var NoMatch = Match{}

func checkLine(line []CandyType) Match {
	var start, length int
	for i := 1; i < len(line); i++ {
		if line[start] == line[i] {
			length++
		} else {
			if length >= Match3 {
				return Match{start: start, length: length}
			}
			length = 0
			start = i
		}
	}
	fmt.Println("no match")
	return NoMatch
}

func NewGame() *Game {
	g := &Game{}
	g.random = rand.New(rand.NewSource(time.Now().Unix()))
	return g
}

func (g *Game) NewCandy() CandyType {
	var c int
	for c == 0 {
		c = g.random.Intn(4)
	}
	return CandyType(c)
}

func loopRowColumn(content func(i, j int)) {
	for i := 0; i < NbBlockWidth; i++ {
		for j := 0; j < NbBlockHeight; j++ {
			content(i, j)
		}
	}
}
