package main

import (
	"fmt"
	"github.com/jackyb/go-sdl2/sdl"
	"os"
	"runtime"
	"time"
)

const FRAME_RATE = time.Second / 50

func main() {
	runtime.LockOSThread()

	window, err := sdl.CreateWindow("snake", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, WINDOW_WIDTH, WINDOW_HEIGHT, sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create window %s\n", sdl.GetError())
		os.Exit(1)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create renderer %s\n", sdl.GetError())
		os.Exit(1)
	}
	defer renderer.Destroy()
	renderer.SetDrawColor(255, 255, 255, 255)

	game := NewGame(renderer)
	defer game.Destroy()

	game.Start()
	loop(game, renderer)
	game.Stop()
}

func loop(game *Game, renderer *sdl.Renderer) {
	mainTicker := time.NewTicker(FRAME_RATE)
	for {
		select {
		case <-game.EndLoop:
			fmt.Println("endloop")
			return
		case <-game.SnakeTick():
			fmt.Println("snake")
			game.Tick()
		case <-game.AppleTick():
			fmt.Println("apple")
			game.NewApple()
		case <-mainTicker.C:
			event := sdl.PollEvent()
			switch t := event.(type) {
			case *sdl.QuitEvent:
				game.StopLoop()
			case *sdl.KeyDownEvent:
				switch t.Keysym.Sym {
				case sdl.K_ESCAPE:
					game.StopLoop()
				case sdl.K_UP:
					game.Command(UP)
				case sdl.K_DOWN:
					game.Command(DOWN)
				case sdl.K_LEFT:
					game.Command(LEFT)
				case sdl.K_RIGHT:
					game.Command(RIGHT)
				}
			}
			renderThings(renderer, game)
		}
	}
}

func renderThings(renderer *sdl.Renderer, game *Game) {
	renderer.Clear()
	// show level
	loopGrid(func(i, j int) {
		b := game.Grid[i][j]
		if b != SNAKE && b != SNAKE_HEAD {
			show(renderer, i, j, b, game)
		}
	})
	// show snake
	snakeType := SNAKE_HEAD
	for _, sp := range game.Snake {
		show(renderer, sp.Pos.X, sp.Pos.Y, snakeType, game)
		if snakeType == SNAKE_HEAD {
			snakeType = SNAKE
		}
	}
	renderer.Present()
}

func show(renderer *sdl.Renderer, x, y int, thing BlockType, game *Game) {
	if thing == EMPTY {
		return
	}
	block.X = int32(x * BLOCK_SIZE)
	block.Y = int32(y * BLOCK_SIZE)
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
