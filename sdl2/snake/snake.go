package main

import (
	"fmt"
	"github.com/jackyb/go-sdl2/sdl"
	"math/rand"
	"time"
)

const (
	BLOCK_SIZE           = 12
	NB_BLOCK_WIDTH       = 32
	NB_BLOCK_HEIGHT      = 32
	WINDOW_HEIGHT        = BLOCK_SIZE * NB_BLOCK_HEIGHT
	WINDOW_WIDTH         = BLOCK_SIZE * NB_BLOCK_WIDTH
	START_LENGTH         = 4
	START_X              = NB_BLOCK_WIDTH / 2
	START_Y              = NB_BLOCK_HEIGHT / 2
	START_DIR            = RIGHT
	START_APPLE_POP_RATE = 0.5
	START_SNAKE_RATE     = 3
)

type Direction int
type BlockType int

//type Position sdl.Rect
type Position struct {
	X, Y int
}
type Grid [NB_BLOCK_WIDTH][NB_BLOCK_HEIGHT]BlockType

//type Things map[BlockType]*sdl.Texture
type Players [4]*sdl.Texture

type Game struct {
	Grid         Grid
	Snake        Snake
	currentDir   Direction
	loose        bool
	snakeRate    float32
	applePopRate float32
	EndLoop      chan bool
	snakeTicker  *time.Ticker
	appleTicker  *time.Ticker
}

type Snake []SnakePart

type SnakePart struct {
	Pos     Position
	nextDir Direction
}

const (
	UP    Direction = 1
	DOWN  Direction = -1
	LEFT  Direction = 2
	RIGHT Direction = -2
)

const (
	EMPTY BlockType = iota
	SNAKE_HEAD
	SNAKE
	APPLE
)

var block = sdl.Rect{W: BLOCK_SIZE, H: BLOCK_SIZE}
var r = rand.New(rand.NewSource(time.Now().Unix()))

func NewGame(renderer *sdl.Renderer) *Game {
	g := Game{}
	g.currentDir = START_DIR

	// init the game map
	loopGrid(func(i, j int) {
		g.Grid[i][j] = EMPTY
	})

	// snake init
	pos := Position{START_X, START_Y}
	for i := 0; i < START_LENGTH; i++ {
		g.Snake = append(g.Snake, SnakePart{Pos: Position{X: pos.X, Y: pos.Y}, nextDir: START_DIR})
		movePos(-START_DIR, &pos)
	}
	g.snakeRate = START_SNAKE_RATE
	g.applePopRate = START_APPLE_POP_RATE
	g.generateSnakeTicker()
	g.generateAppleTicker()
	g.EndLoop = make(chan bool, 1)
	return &g
}

func (g *Game) NewApple() {
	g.newThing(APPLE)
}

func generateDuration(rate float32) time.Duration {
	return time.Duration(float32(time.Second) / rate)
}

func (g *Game) generateAppleTicker() {
	g.appleTicker = time.NewTicker(generateDuration(g.applePopRate))
}

func (g *Game) generateSnakeTicker() {
	g.snakeTicker = time.NewTicker(generateDuration(g.snakeRate))
}

func (g *Game) AppleTick() <-chan time.Time {
	return g.appleTicker.C
}

func (g *Game) SnakeTick() <-chan time.Time {
	return g.snakeTicker.C
}

func (g *Game) newThing(thing BlockType) {
	// determine random coordinates
	var pos Position
	pos.X = r.Intn(NB_BLOCK_HEIGHT)
	pos.Y = r.Intn(NB_BLOCK_WIDTH)

	g.Grid[pos.X][pos.Y] = thing
}

func (g *Game) Command(dir Direction) {
	// ignore command if direction its the inverse of current direction
	if dir != -g.currentDir {
		g.Snake[0].nextDir = dir
	}
}

func (g *Game) Tick() {
	head := &g.Snake[0]
	movePos(head.nextDir, &head.Pos)
	g.currentDir = head.nextDir
	dir := g.currentDir
	for i := 1; i < len(g.Snake); i++ {
		movePos(g.Snake[i].nextDir, &g.Snake[i].Pos)
		dir, g.Snake[i].nextDir = g.Snake[i].nextDir, dir
	}
	// snake collision?
	for _, part := range g.Snake[1:] {
		if part.Pos.X == head.Pos.X && part.Pos.Y == head.Pos.Y {
			fmt.Println("Loose", head.Pos.X, head.Pos.Y)
			g.loose = true
			g.EndLoop <- true
			return
		}
	}
	// eat apple ?
	if g.Grid[head.Pos.X][head.Pos.Y] == APPLE {
		g.Grid[head.Pos.X][head.Pos.Y] = EMPTY
		g.grow()
	}
}

func (g *Game) grow() {
	queue := g.Snake[len(g.Snake)-1]
	pos := queue.Pos
	movePos(-queue.nextDir, &pos)
	g.Snake = append(g.Snake, SnakePart{pos, queue.nextDir})
}

func movePos(dir Direction, pos *Position) {
	switch dir {
	case UP:
		if pos.Y == 0 {
			pos.Y = NB_BLOCK_HEIGHT - 1
		} else {
			pos.Y--
		}
	case DOWN:
		if pos.Y == NB_BLOCK_HEIGHT-1 {
			pos.Y = 0
		} else {
			pos.Y++
		}
	case LEFT:
		if pos.X == 0 {
			pos.X = NB_BLOCK_WIDTH - 1
		} else {
			pos.X--
		}
	case RIGHT:
		if pos.X == NB_BLOCK_WIDTH-1 {
			pos.X = 0
		} else {
			pos.X++
		}
	}
}

func loopGrid(content func(i, j int)) {
	for i := 0; i < NB_BLOCK_WIDTH; i++ {
		for j := 0; j < NB_BLOCK_HEIGHT; j++ {
			content(i, j)
		}
	}
}

func (g *Game) StopLoop() {
	fmt.Println("endLoop")
	g.EndLoop <- true
	fmt.Println("endLoop2")
}

func (g *Game) Start() {
}

func (g *Game) Stop() {
}

func (g *Game) Destroy() {
	//for _, t := range g.players {
	//	t.Destroy()
	//}
	//for _, t := range g.things {
	//	t.Destroy()
	//}
}
