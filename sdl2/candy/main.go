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
	_ = fmt.Sprint()
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
			wait := game.Tick()
			renderThings(renderer, game)
			if wait {
				event := sdl.WaitEvent()
				switch t := event.(type) {
				case *sdl.QuitEvent:
					return
				case *sdl.KeyDownEvent:
					switch t.Keysym.Sym {
					case sdl.K_ESCAPE:
						return
					case sdl.K_r:
						game.Reset()
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
	for _, col := range game.columns {
		for _, c := range col.candys {
			showCandy(renderer, c, game)
		}
	}
	renderer.SetDrawColor(255, 255, 255, 255)
	renderer.Present()
}

var block = sdl.Rect{W: BlockSize - 2, H: BlockSize - 2}

func showCandy(renderer *sdl.Renderer, c Candy, game *Game) {
	if c._type == EmptyCandy {
		return
	}
	//fmt.Printf("showCandy (%d,%d), %d\n", c.x, c.y, c._type)
	block.X = int32(c.x + 1)
	block.Y = int32(c.y + 1)
	switch c._type {
	case BlueCandy:
		renderer.SetDrawColor(153, 50, 204, 255)
	case YellowCandy:
		renderer.SetDrawColor(255, 215, 0, 255)
	case GreenCandy:
		renderer.SetDrawColor(60, 179, 113, 255)
	case RedCandy:
		renderer.SetDrawColor(220, 20, 60, 255)
	case PinkCandy:
		renderer.SetDrawColor(255, 192, 203, 255)

	}
	renderer.FillRect(&block)
	renderer.SetDrawColor(255, 255, 255, 255)
}
