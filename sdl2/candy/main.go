package main

import (
	"fmt"
	"github.com/jackyb/go-sdl2/sdl"
	"os"
	"runtime"
	"time"
)

const FRAME_RATE = time.Second / 40

var (
	window          *sdl.Window
	tileset         *sdl.Texture
	tilesetSelected *sdl.Texture
)

func main() {
	_ = fmt.Sprint()
	runtime.LockOSThread()

	window = sdl.CreateWindow("Candy Crush Saga", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, WindowWidth, WindowHeight, sdl.WINDOW_SHOWN)
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
	renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND)

	tilesetFile := os.Getenv("GOPATH") + "/src/github.com/tbruyelle/gostuff/sdl2/candy/assets/tileset.bmp"
	tilesetSurface := sdl.LoadBMP(tilesetFile)
	if tilesetSurface == nil {
		fmt.Fprintf(os.Stderr, "Failed to load bitmap %s", tilesetFile)
		os.Exit(1)
	}
	tileset = renderer.CreateTextureFromSurface(tilesetSurface)
	tilesetSurface.SetAlphaMod(190)
	tilesetSelected = renderer.CreateTextureFromSurface(tilesetSurface)

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
			for wait {
				event := sdl.PollEvent()
				switch t := event.(type) {
				case *sdl.QuitEvent:
					return
				case *sdl.KeyDownEvent:
					switch t.Keysym.Sym {
					case sdl.K_ESCAPE:
						return
					case sdl.K_r:
						game.Reset()
						wait = false
					}
				case *sdl.MouseButtonEvent:
					if t.State != 0 {
						//fmt.Println("Click", t.X, t.Y)
						game.Click(int(t.X), int(t.Y))
						wait = false
					}
				}

			}
		}
	}
}

func renderThings(renderer *sdl.Renderer, game *Game) {
	//fmt.Println("rendering")
	renderer.Clear()
	// show dashboard
	renderer.SetDrawColor(50, 50, 50, 200)
	dashboard := sdl.Rect{0, 0, DashboardWidth, WindowHeight}
	renderer.FillRect(&dashboard)

	// show candys
	for _, c := range game.candys {
		showCandy(renderer, *c, game)
	}
	renderer.SetDrawColor(255, 255, 255, 255)
	renderer.Present()
}

var block = sdl.Rect{W: BlockSize, H: BlockSize}
var source = sdl.Rect{W: BlockSize, H: BlockSize}

func showCandy(renderer *sdl.Renderer, c Candy, game *Game) {
	if c._type == EmptyCandy {
		return
	}
	//fmt.Printf("showCandy (%d,%d), %d\n", c.x, c.y, c._type)
	block.X = int32(c.x)
	block.Y = int32(c.y)
	//alpha := uint8(255)
	//if c.selected {
	//	alpha = 150
	//}
	switch c._type {
	case BlueCandy:
		source.X = BlockSize
	case YellowCandy:
		source.X = BlockSize * 4
	case GreenCandy:
		source.X = BlockSize * 3
	case RedCandy:
		source.X = BlockSize * 5
	case PinkCandy:
		source.X = BlockSize * 2
	case OrangeCandy:
		source.X = 0
	}
	if c.selected {
		renderer.Copy(tilesetSelected, &source, &block)
	} else {
		renderer.Copy(tileset, &source, &block)
	}
}

func showCandySquare(renderer *sdl.Renderer, c Candy, game *Game) {
	if c._type == EmptyCandy {
		return
	}
	//fmt.Printf("showCandy (%d,%d), %d\n", c.x, c.y, c._type)
	block.X = int32(c.x + 1)
	block.Y = int32(c.y + 1)
	alpha := uint8(255)
	if c.selected {
		alpha = 150
	}
	switch c._type {
	case BlueCandy:
		renderer.SetDrawColor(153, 50, 204, alpha)
	case YellowCandy:
		renderer.SetDrawColor(255, 215, 0, alpha)
	case GreenCandy:
		renderer.SetDrawColor(60, 179, 113, alpha)
	case RedCandy:
		renderer.SetDrawColor(220, 20, 60, alpha)
	case PinkCandy:
		renderer.SetDrawColor(255, 192, 203, alpha)

	}
	renderer.FillRect(&block)
	renderer.SetDrawColor(255, 255, 255, 255)
}
