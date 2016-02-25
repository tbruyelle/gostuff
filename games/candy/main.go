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
	window  *sdl.Window
	tileset *sdl.Texture
)

// Arrange that main.main runs on main thread.
func init() {
	runtime.LockOSThread()
}

// Queue of work to run in main thread.
var mainfunc = make(chan func())

// Run all the functions that need to run in the main thread.
func Main() {
	var f func()
	for f = range mainfunc {
		f()
	}
}

// Put the function f on the main thread function queue.
func do(f func()) {
	done := make(chan bool, 1)
	mainfunc <- func() {
		f()
		done <- true
	}
	<-done
}

func main() {
	_ = fmt.Sprint()
	defer sdl.Quit()
	var err error

	window, err = sdl.CreateWindow("Candy Crush Saga", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, WindowWidth, WindowHeight, sdl.WINDOW_SHOWN)
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
	renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND)

	tilesetFile := os.Getenv("GOPATH") + "/src/github.com/tbruyelle/gostuff/games/candy/assets/tileset.bmp"
	tilesetSurface, err := sdl.LoadBMP(tilesetFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load bitmap %s", tilesetFile)
		os.Exit(1)
	}
	tileset, err = renderer.CreateTextureFromSurface(tilesetSurface)
	if err != nil {
		panic(err)
	}

	game := NewGame()
	defer game.Destroy()

	game.Start()
	go eventLoop(game)
	go renderLoop(game, renderer)
	Main()
	game.Stop()
}

func eventLoop(game *Game) {
	defer close(mainfunc)
	var evt sdl.Event
	for {
		do(func() {
			evt = sdl.PollEvent()
		})
		switch t := evt.(type) {
		case *sdl.QuitEvent:
			return
		case *sdl.KeyDownEvent:
			switch t.Keysym.Sym {
			case sdl.K_ESCAPE:
				return
			case sdl.K_r:
				game.Reset()
			case sdl.K_k:
				game.ToggleKeepUnmatchingTranslation()
			}
		case *sdl.MouseButtonEvent:
			if t.State != 0 {
				game.Click(int(t.X), int(t.Y))
			}
		}
	}
}

func renderLoop(game *Game, renderer *sdl.Renderer) {
	defer close(mainfunc)

	mainTicker := time.NewTicker(FRAME_RATE)

	for {
		select {
		case <-mainTicker.C:
			game.Tick()
			do(func() {
				renderThings(renderer, game)
			})

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
		showCandy(renderer, c, game)
	}
	renderer.SetDrawColor(255, 255, 255, 255)
	renderer.Present()
}

var block = sdl.Rect{W: BlockSize, H: BlockSize}
var source = sdl.Rect{W: BlockSize, H: BlockSize}

// showCandy shows the candy according to a tileset.
func showCandy(renderer *sdl.Renderer, c *Candy, game *Game) {
	if c._type == EmptyCandy {
		return
	}
	//fmt.Printf("showCandy (%d,%d), %d\n", c.x, c.y, c._type)
	block.X = int32(c.x)
	block.Y = int32(c.y)
	alpha := uint8(255)
	switch c._type {
	case BlueCandy:
		source.X = BlockSize
		source.Y = 0
	case YellowCandy:
		source.X = BlockSize * 4
		source.Y = 0
	case GreenCandy:
		source.X = BlockSize * 3
		source.Y = 0
	case RedCandy:
		source.X = BlockSize * 5
		source.Y = 0
	case PinkCandy:
		source.X = BlockSize * 2
		source.Y = 0
	case OrangeCandy:
		source.X = 0
		source.Y = 0
	case BlueHStripesCandy:
		source.X = BlockSize
		source.Y = BlockSize
	case YellowHStripesCandy:
		source.X = BlockSize * 4
		source.Y = BlockSize
	case GreenHStripesCandy:
		source.X = BlockSize * 3
		source.Y = BlockSize
	case RedHStripesCandy:
		source.X = BlockSize * 5
		source.Y = BlockSize
	case PinkHStripesCandy:
		source.X = BlockSize * 2
		source.Y = BlockSize
	case OrangeHStripesCandy:
		source.X = 0
		source.Y = BlockSize
	case BlueVStripesCandy:
		source.X = BlockSize
		source.Y = BlockSize * 2
	case YellowVStripesCandy:
		source.X = BlockSize * 4
		source.Y = BlockSize * 2
	case GreenVStripesCandy:
		source.X = BlockSize * 3
		source.Y = BlockSize * 2
	case RedVStripesCandy:
		source.X = BlockSize * 5
		source.Y = BlockSize * 2
	case PinkVStripesCandy:
		source.X = BlockSize * 2
		source.Y = BlockSize * 2
	case OrangeVStripesCandy:
		source.X = 0
		source.Y = BlockSize * 2
	case BluePackedCandy:
		source.X = BlockSize
		source.Y = BlockSize * 3
	case YellowPackedCandy:
		source.X = BlockSize * 4
		source.Y = BlockSize * 3
	case GreenPackedCandy:
		source.X = BlockSize * 3
		source.Y = BlockSize * 3
	case RedPackedCandy:
		source.X = BlockSize * 5
		source.Y = BlockSize * 3
	case PinkPackedCandy:
		source.X = BlockSize * 2
		source.Y = BlockSize * 3
	case OrangePackedCandy:
		source.X = 0
		source.Y = BlockSize * 3
	case BombCandy:
		source.X = 0
		source.Y = BlockSize * 4
	}

	if c == game.selected {
		alpha = uint8(150)
	} else {
		alpha = c.sprite.alpha
	}

	tileset.SetAlphaMod(alpha)
	renderer.Copy(tileset, &source, &block)
}

// showCandySquare is a deprecated method which shows candys as
// simples colored squares.
func showCandySquare(renderer *sdl.Renderer, c *Candy, game *Game) {
	if c._type == EmptyCandy {
		return
	}
	//fmt.Printf("showCandy (%d,%d), %d\n", c.x, c.y, c._type)
	block.X = int32(c.x + 1)
	block.Y = int32(c.y + 1)
	alpha := uint8(255)
	if c == game.selected {
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
