package main

import (
	"fmt"
	"github.com/jackyb/go-sdl2/sdl"
	"github.com/jackyb/go-sdl2/sdl_image"
	"os"
)

const (
	BLOCK_SIZE      = 34
	NB_BLOCK_WIDTH  = 12
	NB_BLOCK_HEIGHT = 12
	WINDOW_HEIGHT   = BLOCK_SIZE * NB_BLOCK_HEIGHT
	WINDOW_WIDTH    = BLOCK_SIZE * NB_BLOCK_WIDTH
)

type Direction int
type Position sdl.Rect
type BlockType int
type Sprite *sdl.Texture
type Grid [NB_BLOCK_WIDTH][NB_BLOCK_HEIGHT]BlockType
type Things map[BlockType]*sdl.Texture
type Players [4]*sdl.Texture

type Game struct {
	grid          Grid
	mouseGrid     Grid
	things        Things
	players       Players
	playerPos     Position
	currentPlayer Sprite
	goals         int
}

const (
	UP Direction = iota
	DOWN
	LEFT
	RIGHT
)

const (
	EMPTY BlockType = iota
	WALL
	BOX
	GOAL
	GUS
	BOX_OK
)

func main() {
	window := sdl.CreateWindow("sokoban", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, WINDOW_WIDTH, WINDOW_HEIGHT, sdl.WINDOW_SHOWN)
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

	menu := img.Load("assets/menu.jpg")
	if menu == nil {
		fmt.Fprintf(os.Stderr, "Failed to open menu image %s", sdl.GetError())
		os.Exit(1)
	}
	defer menu.Free()
	texture := renderer.CreateTextureFromSurface(menu)
	defer texture.Destroy()


	running := true
	for running {
	renderer.Clear()
	renderer.Copy(texture, nil, nil)
	renderer.Present()
		var event sdl.Event
		sdl.WaitEvent(&event)
		switch t := event.(type) {
		case *sdl.QuitEvent:
			running = false
		case *sdl.KeyDownEvent:
			switch t.Keysym.Sym {
			case sdl.K_ESCAPE:
				running = false
			case sdl.K_1:
				play(renderer)
			case sdl.K_2:
				editor(renderer)
			}
		}

	}
}

func renderThings(renderer *sdl.Renderer, game *Game) {
	renderer.Clear()
	// show level
	loopGrid(func(i, j int) {
		block.X = int32(i * BLOCK_SIZE)
		block.Y = int32(j * BLOCK_SIZE)
		texture := game.things[game.grid[i][j]]
		if texture != nil {
			renderer.Copy(texture, nil, &block)
		}
	})
}

func loopGrid(content func(i, j int)) {
	for i := 0; i < NB_BLOCK_WIDTH; i++ {
		for j := 0; j < NB_BLOCK_HEIGHT; j++ {
			content(i, j)
		}
	}
}

func loadAsset(renderer *sdl.Renderer, path string) *sdl.Texture {
	asset := img.Load("assets/" + path)
	if asset == nil {
		fmt.Fprintf(os.Stderr, "Failed to create image %s : %s", path, sdl.GetError())
		os.Exit(1)
	}
	defer asset.Free()

	texture := renderer.CreateTextureFromSurface(asset)
	if texture == nil {
		fmt.Fprintf(os.Stderr, "Failed to create texture from %s : %s", path, sdl.GetError())
		os.Exit(1)
	}
	return texture
}

func NewGame(renderer *sdl.Renderer) *Game {
	g := Game{}
	g.things = make(Things)
	// init the game map
	loopGrid(func(i, j int) {
		g.grid[i][j] = 0
		g.mouseGrid[i][j] = 0
	})
	g.players[UP] = loadAsset(renderer, "mario_haut.gif")
	g.players[DOWN] = loadAsset(renderer, "mario_bas.gif")
	g.players[LEFT] = loadAsset(renderer, "mario_gauche.gif")
	g.players[RIGHT] = loadAsset(renderer, "mario_droite.gif")
	g.currentPlayer = g.players[DOWN]

	g.things[WALL] = loadAsset(renderer, "mur.jpg")
	g.things[BOX] = loadAsset(renderer, "caisse.jpg")
	g.things[BOX_OK] = loadAsset(renderer, "caisse_ok.jpg")
	g.things[GOAL] = loadAsset(renderer, "objectif.png")
	g.things[GUS] = loadAsset(renderer, "mario_bas.gif")

	return &g
}

func (g *Game) Destroy() {
	for _, t := range g.players {
		t.Destroy()
	}
	for _, t := range g.things {
		t.Destroy()
	}
}
