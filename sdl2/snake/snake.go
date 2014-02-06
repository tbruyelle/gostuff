package main

import (
	"fmt"
	"github.com/jackyb/go-sdl2/sdl"
	"github.com/jackyb/go-sdl2/sdl_image"
	"math/rand"
	"os"
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
type Position sdl.Rect
type Sprite *sdl.Texture
type Grid [NB_BLOCK_WIDTH][NB_BLOCK_HEIGHT]BlockType

//type Things map[BlockType]*sdl.Texture
type Players [4]*sdl.Texture

type Game struct {
	grid    Grid
	snake   Snake
	dir     Direction
	invDir  map[Direction]Direction
	tickers []*time.Ticker
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

func main() {
	window := sdl.CreateWindow("snake", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, WINDOW_WIDTH, WINDOW_HEIGHT, sdl.WINDOW_SHOWN)
	if window == nil {
		fmt.Fprintf(os.Stderr, "failed to create window %s\n", sdl.GetError())
		os.Exit(1)
	}
	defer window.Destroy()

	renderer := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if renderer == nil {
		fmt.Fprintf(os.Stderr, "Failed to create renderer %s\n", sdl.GetError())
		os.Exit(1)
	}
	defer renderer.Destroy()
	renderer.SetDrawColor(255, 255, 255, 255)

	game := NewGame(renderer)
	defer game.Destroy()

	running := true
	prevTs := uint32(0)
	game.Start()
	for running {
		event := sdl.PollEvent()
		switch t := event.(type) {
		case *sdl.QuitEvent:
			running = false
		case *sdl.KeyDownEvent:
			switch t.Keysym.Sym {
			case sdl.K_ESCAPE:
				running = false
			case sdl.K_UP:
				nextDir(game, UP)
			case sdl.K_DOWN:
				nextDir(game, DOWN)
			case sdl.K_LEFT:
				nextDir(game, LEFT)
			case sdl.K_RIGHT:
				nextDir(game, RIGHT)
			}
		}
		ts := sdl.GetTicks()
		if ts-prevTs > TICK {
			prevTs = ts
			moveSnake(game)
		} else {
			sdl.Delay(TICK - (ts - prevTs))
		}
		renderThings(renderer, game)
	}
	game.Stop()
}

func thingPoper(g *Game, ticker *time.Ticker, thing BlockType) {
	for _ = range ticker.C {
		newThing(g, thing)
	}
}

func newThing(g *Game, thing BlockType) {
	// dtermine random coordinates
	var pos Position
	pos.X = int32(rand.Intn(NB_BLOCK_HEIGHT))
	pos.Y = int32(rand.Intn(NB_BLOCK_WIDTH))

	g.grid[pos.X][pos.Y] = thing
}

func nextDir(game *Game, dir Direction) {
	game.snake[0].nextDir = dir
}

func beyondLimits(pos Position) bool {
	return pos.X >= 0 && pos.X < NB_BLOCK_HEIGHT && pos.Y >= 0 && pos.Y < NB_BLOCK_WIDTH
}

func moveSnake(g *Game) {
	head := &g.snake[0]
	movePos(head.nextDir, &head.pos)
	dir := head.nextDir
	for i := 1; i < len(g.snake); i++ {
		movePos(g.snake[i].nextDir, &g.snake[i].pos)
		dir, g.snake[i].nextDir = g.snake[i].nextDir, dir
	}
	// eat apple ?
	if beyondLimits(head.pos)&&g.grid[head.pos.X][head.pos.Y] == APPLE {
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
		pos.Y--
	case DOWN:
		pos.Y++
	case LEFT:
		pos.X--
	case RIGHT:
		pos.X++
	}
}

func renderThings(renderer *sdl.Renderer, game *Game) {
	renderer.Clear()
	// show level
	loopGrid(func(i, j int) {
		b := game.grid[i][j]
		if b != SNAKE && b != SNAKE_HEAD {
			show(renderer, int32(i), int32(j), b, game)
		}
	})
	// show snake
	snakeType := SNAKE_HEAD
	for _, sp := range game.snake {
		show(renderer, sp.pos.X, sp.pos.Y, snakeType, game)
		if snakeType == SNAKE_HEAD {
			snakeType = SNAKE
		}
	}
	renderer.Present()
}

func show(renderer *sdl.Renderer, x, y int32, thing BlockType, game *Game) {
	if thing == EMPTY {
		return
	}
	block.X = x * BLOCK_SIZE
	block.Y = y * BLOCK_SIZE
	switch thing {
	case APPLE:
		renderer.SetDrawColor(0, 255, 0, 255)
	case SNAKE_HEAD:
		renderer.SetDrawColor(0, 0, 0, 255)
	case SNAKE:
		renderer.SetDrawColor(100, 0, 0, 255)
	}
	renderer.FillRect(&block)
	renderer.SetDrawColor(255, 255, 255, 255)
}

func loopGrid(content func(i, j int)) {
	for i := 0; i < NB_BLOCK_WIDTH; i++ {
		for j := 0; j < NB_BLOCK_HEIGHT; j++ {
			content(i, j)
		}
	}
}

func loadAsset(renderer *sdl.Renderer, path string) *sdl.Texture {
	asset := img.LoadTexture(renderer, "assets/"+path)
	if asset == nil {
		fmt.Fprintf(os.Stderr, "Failed to create image %s : %s", path, sdl.GetError())
		os.Exit(1)
	}
	return asset
}

func NewGame(renderer *sdl.Renderer) *Game {
	g := Game{}
	g.dir = START_DIR
	g.invDir = make(map[Direction]Direction)
	g.invDir[UP] = DOWN
	g.invDir[DOWN] = UP
	g.invDir[LEFT] = RIGHT
	g.invDir[RIGHT] = LEFT

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

	//g.players[UP] = loadAsset(renderer, "mario_haut.gif")
	//g.players[DOWN] = loadAsset(renderer, "mario_bas.gif")
	//g.players[LEFT] = loadAsset(renderer, "mario_gauche.gif")
	//g.players[RIGHT] = loadAsset(renderer, "mario_droite.gif")
	//g.currentPlayer = g.players[DOWN]

	//g.things = make(Things)
	//g.things[SNAKE] = renderer.CreateTextureFromSurface(surface)
	//g.things[SNAKE_HEAD] = renderer.CreateTextureFromSurface(surface)

	return &g
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
