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
	FRAME_RATE   = time.Second / 40
	BlockWidth   = 110
	BlockHeight  = 128
	WindowWidth  = BlockWidth * 4
	WindowHeight = BlockHeight * 4
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
	for i := range game.Board {
		for j := range game.Board[i] {
			x, y := i*BlockWidth, j*BlockHeight
			showBlock(renderer, x, y, game.Board[i][j])
		}
	}
	renderer.SetDrawColor(255, 255, 255, 255)
	renderer.Present()
}

var block = sdl.Rect{W: BlockWidth, H: BlockHeight}
var source = sdl.Rect{W: BlockWidth, H: BlockHeight}

func showBlock(renderer *sdl.Renderer, x, y, t int) {
	if t == 0 {
		return
	}
	//fmt.Printf("showBlock (%d,%d), %d\n", x, y, t)
	block.X = int32(x)
	block.Y = int32(y)
	alpha := uint8(255)
	switch t {
	case 1:
		source.X = 0
		source.Y = 0
	}

	/*if c == game.selected {
		alpha = uint8(150)
	} else {
		alpha = c.sprite.alpha
	}*/

	tileset.SetAlphaMod(alpha)
	renderer.Copy(tileset, &source, &block)
}
