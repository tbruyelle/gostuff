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

	window := sdl.CreateWindow("Candy Crush Saga", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, WindowWidth, WindowHeight, sdl.WINDOW_SHOWN)
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

	game := NewGame()
	defer game.Destroy()

	game.Start()
	loop(game, renderer)
	game.Stop()
}

func loop(game *Game, renderer *sdl.Renderer) {
	mainTicker := time.NewTicker(FRAME_RATE)
	for {
		select {
		case <-mainTicker.C:
			wait:=game.Tick()
			renderThings(renderer, game)
			if wait{
			event := sdl.WaitEvent()
			switch t := event.(type) {
			case *sdl.QuitEvent:
				return
			case *sdl.KeyDownEvent:
				switch t.Keysym.Sym {
				case sdl.K_ESCAPE:
					return
				}
			}
		}
		}
	}
}

func renderThings(renderer *sdl.Renderer, game *Game) {
	renderer.Clear()
	// show dashboard
	renderer.SetDrawColor(50, 50, 50, 200)
	dashboard := sdl.Rect{0, 0, DashboardWidth, WindowHeight}
	renderer.FillRect(&dashboard)

	// show candys

	renderer.SetDrawColor(255, 255, 255, 255)
	renderer.Present()
}

func showCandy(renderer *sdl.Renderer, x, y int, candy CandyType, game *Game) {
	if candy == EmptyCandy {
		return
	}
	block := sdl.Rect{}
	block.X = int32(x * BlockSize)
	block.Y = int32(y * BlockSize)
	switch candy {
	case BlueCandy:
		renderer.SetDrawColor(0, 0, 255, 255)
	case YellowCandy:
		renderer.SetDrawColor(0, 0, 0, 255)
	case GreenCandy:
		renderer.SetDrawColor(100, 0, 0, 255)
	case RedCandy:
		renderer.SetDrawColor(0, 255, 0, 255)
	}
	renderer.FillRect(&block)
	renderer.SetDrawColor(255, 255, 255, 255)
}
