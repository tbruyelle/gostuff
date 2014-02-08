package main

import (
	"github.com/jackyb/go-sdl2/sdl"
	"math/rand"
	"time"
)

const (
	BLOCK_SIZE      = 12
	NB_BLOCK_WIDTH  = 32
	NB_BLOCK_HEIGHT = 32
	WINDOW_HEIGHT   = BLOCK_SIZE * NB_BLOCK_HEIGHT
	WINDOW_WIDTH    = BLOCK_SIZE * NB_BLOCK_WIDTH
	START_LENGTH    = 4
	START_X         = NB_BLOCK_WIDTH / 2
	START_Y         = NB_BLOCK_HEIGHT / 2
	TICK            = 60
	START_DIR       = RIGHT
	APPLE_TICK      = 2000
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
	grid    Grid
	snake   Snake
	dir     Direction
	invDir  map[Direction]Direction
	tickers []*time.Ticker
	running bool
}

type Snake []SnakePart

type SnakePart struct {
	pos     Position
	nextDir Direction
}

const (
	UP Direction = iota
	DOWN
	LEFT
	RIGHT
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
	g.dir = START_DIR
	g.invDir = make(map[Direction]Direction)
	g.invDir[UP] = DOWN
	g.invDir[DOWN] = UP
	g.invDir[LEFT] = RIGHT
	g.invDir[RIGHT] = LEFT
	g.running = true

	// init the game map
	loopGrid(func(i, j int) {
		g.grid[i][j] = EMPTY
	})

	// snake init
	var pos Position
	pos.X = START_X
	pos.Y = START_Y
	for i := 0; i < START_LENGTH; i++ {
		g.snake = append(g.snake, SnakePart{pos: Position{X: pos.X, Y: pos.Y}, nextDir: START_DIR})
		movePos(g.invDir[START_DIR], &pos)
	}
	return &g
}

func (g *Game) Loop() bool {
	return g.running
}

func (g *Game) StopLoop() {
	g.running = false
}

func thingPoper(g *Game, ticker *time.Ticker, thing BlockType) {
	for _ = range ticker.C {
		newThing(g, thing)
	}
}

func newThing(g *Game, thing BlockType) {
	// determine random coordinates
	var pos Position
	pos.X = r.Intn(NB_BLOCK_HEIGHT)
	pos.Y = r.Intn(NB_BLOCK_WIDTH)

	g.grid[pos.X][pos.Y] = thing
}

func (g *Game) Command(dir Direction) {
	g.snake[0].nextDir = dir
}

func beyondLimits(pos Position) bool {
	return pos.X >= 0 && pos.X < NB_BLOCK_HEIGHT && pos.Y >= 0 && pos.Y < NB_BLOCK_WIDTH
}

func (g *Game) Tick() {
	head := &g.snake[0]
	movePos(head.nextDir, &head.pos)
	dir := head.nextDir
	for i := 1; i < len(g.snake); i++ {
		movePos(g.snake[i].nextDir, &g.snake[i].pos)
		dir, g.snake[i].nextDir = g.snake[i].nextDir, dir
	}
	// eat apple ?
	if beyondLimits(head.pos) && g.grid[head.pos.X][head.pos.Y] == APPLE {
		g.grid[head.pos.X][head.pos.Y] = EMPTY
		queue := g.snake[len(g.snake)-1]
		var pos Position
		pos.X = queue.pos.X
		pos.Y = queue.pos.Y
		movePos(g.invDir[queue.nextDir], &pos)
		// increase snake length
		g.snake = append(g.snake, SnakePart{pos, queue.nextDir})
	}
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

func (g *Game) Start() {
	g.NewTicker(APPLE_TICK, APPLE)
}

func (g *Game) NewTicker(tick time.Duration, thing BlockType) {
	ticker := time.NewTicker(time.Millisecond * tick)
	go thingPoper(g, ticker, thing)
	g.tickers = append(g.tickers, ticker)
}

func (g *Game) Stop() {
	for _, t := range g.tickers {
		t.Stop()
	}
}

func (g *Game) Destroy() {
	//for _, t := range g.players {
	//	t.Destroy()
	//}
	//for _, t := range g.things {
	//	t.Destroy()
	//}
}
