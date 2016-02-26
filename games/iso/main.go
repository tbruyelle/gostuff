package main

import (
	"fmt"
	"github.com/jackyb/go-sdl2/sdl"
	"log"
	"os"
	"runtime"
	"time"
)

const (
	FRAME_RATE       = time.Second / 40
	BlockWidth       = 110
	BlockFloorHeight = 65
	BlockHeight      = 128
	WindowWidth      = BlockWidth * 6
	WindowHeight     = BlockHeight * 4
)

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

var game *Game

func main() {
	_ = fmt.Sprint()
	defer sdl.Quit()
	var err error

	window, err = sdl.CreateWindow("iso", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, WindowWidth, WindowHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		log.Fatalf("failed to create window %s", sdl.GetError())
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		log.Fatalf("Failed to create renderer %s", sdl.GetError())
	}
	defer renderer.Destroy()
	renderer.SetDrawColor(255, 255, 255, 255)
	renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND)

	tilesetFile := os.Getenv("GOPATH") + "/src/github.com/tbruyelle/gostuff/games/iso/assets/tileset.bmp"
	tilesetSurface, err := sdl.LoadBMP(tilesetFile)
	if err != nil {
		log.Fatalf("Failed to load tileset %s, %v", tilesetFile, err)
	}
	tileset, err = renderer.CreateTextureFromSurface(tilesetSurface)
	if err != nil {
		log.Fatal(err)
	}

	game = NewGame()
	defer game.Destroy()

	game.Start()
	go eventLoop()
	go renderLoop(renderer)
	Main()
	game.Stop()
}

func eventLoop() {
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
			}
		case *sdl.MouseButtonEvent:
			if t.State != 0 {
				game.Click(int(t.X), int(t.Y))
			}
		}
	}
}

func renderLoop(renderer *sdl.Renderer) {
	defer close(mainfunc)

	mainTicker := time.NewTicker(FRAME_RATE)

	for {
		select {
		case <-mainTicker.C:
			game.Tick()
			do(func() {
				renderThings(renderer)
			})
		}
	}
}

func renderThings(renderer *sdl.Renderer) {
	//fmt.Println("rendering")
	renderer.Clear()

	// show blocks
	var offsetx int
	for i := 0; i < len(game.Board); i++ {
		if i%2 != 0 {
			offsetx = BlockWidth / 2
		} else {
			offsetx = 0
		}
		for j := 0; j < len(game.Board[i]); j++ {
			x := j*BlockWidth + offsetx
			y := i * BlockFloorHeight / 2
			showBlock(renderer, x, y, game.Board[i][j])
		}
	}
	renderer.SetDrawColor(255, 255, 255, 255)
	renderer.Present()
}

var (
	block  = sdl.Rect{W: BlockWidth, H: BlockHeight}
	source = sdl.Rect{W: BlockWidth, H: BlockHeight}
	grid   = sdl.Rect{W: BlockWidth, H: BlockHeight, X: (BlockWidth + 1) * 10, Y: BlockHeight * 5}
)

func showBlock(renderer *sdl.Renderer, x, y, t int) {
	block.X = int32(x)
	block.Y = int32(y)
	if t != 0 {
		switch t {
		case 1:
			source.X = 0
			source.Y = 0
		case 2:
			source.X = BlockWidth + 1
			source.Y = 0
		}

		alpha := uint8(255)
		/*if c == game.selected {
			alpha = uint8(150)
		} else {
			alpha = c.sprite.alpha
		}*/

		tileset.SetAlphaMod(alpha)
		renderer.Copy(tileset, &source, &block)
	}
	if game.ShowGrid {
		renderer.Copy(tileset, &grid, &block)
	}
}
