package main

import (
	"fmt"
	"github.com/jackyb/go-sdl2/sdl"
)

var currentThing BlockType
var mousePos sdl.Rect

func editor(renderer *sdl.Renderer) {
	game := NewGame(renderer)
	defer game.Destroy()

	currentThing = WALL
	fmt.Println("Editor ready")

	//main loop
	running := true
	btnPressed := uint8(0)
	for running {
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
				currentThing = WALL
			case sdl.K_2:
				currentThing = BOX
			case sdl.K_3:
				currentThing = GOAL
			case sdl.K_4:
				currentThing = GUS
			}
			mark(mousePos.X, mousePos.Y, currentThing, &(game.mouseGrid))

		case *sdl.MouseButtonEvent:
			if t.State == 0 {
				btnPressed = 0
			} else {
				btnPressed = t.Button
				markBtn(btnPressed, t.X, t.Y, game)
			}
		case *sdl.MouseMotionEvent:
			mark(mousePos.X, mousePos.Y, EMPTY, &(game.mouseGrid))
			mousePos.X = t.X
			mousePos.Y = t.Y
			mark(mousePos.X, mousePos.Y, currentThing, &(game.mouseGrid))
			if btnPressed != 0 {
				markBtn(btnPressed, mousePos.X, mousePos.Y, game)
			}
		}
		renderEditor(renderer, game)
	}
}

func markBtn(button uint8, x, y int32, game *Game) {
	switch button {
	case 1:
		mark(x, y, currentThing, &(game.grid))
	case 3:
		mark(x, y, EMPTY, &(game.grid))
	}
}

func mark(x, y int32, blockt BlockType, grid *Grid) BlockType {
	bx := x / BLOCK_SIZE
	by := y / BLOCK_SIZE
	previousBlock := grid[bx][by]
	grid[bx][by] = blockt
	return previousBlock
}

func renderEditor(renderer *sdl.Renderer, g *Game) {
	renderThings(renderer, g)
	// show mouse cursor
	loopGrid(func(i, j int) {
		block.X = int32(i * BLOCK_SIZE)
		block.Y = int32(j * BLOCK_SIZE)
		texture := g.things[g.mouseGrid[i][j]]
		if texture != nil {
			texture.SetBlendMode(sdl.BLENDMODE_BLEND)
			texture.SetAlphaMod(160)
			renderer.Copy(texture, nil, &block)
			texture.SetAlphaMod(255)
		}
	})

	renderer.Present()
}
