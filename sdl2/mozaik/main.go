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

	window = sdl.CreateWindow("Mozaik", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, WindowWidth, WindowHeight, sdl.WINDOW_SHOWN)
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

	//tilesetFile := os.Getenv("GOPATH") + "/src/github.com/tbruyelle/gostuff/sdl2/candy/assets/tileset.bmp"
	//tilesetSurface := sdl.LoadBMP(tilesetFile)
	//if tilesetSurface == nil {
	//	fmt.Fprintf(os.Stderr, "Failed to load bitmap %s", tilesetFile)
	//	os.Exit(1)
	//}
	//tileset = renderer.CreateTextureFromSurface(tilesetSurface)

	g := NewGame()

	g.Start()
	go eventLoop(g)
	go renderLoop(g, renderer)
	Main()
	g.Stop()
}

func eventLoop(g *Game) {
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
				g.Reset()
			}
		case *sdl.MouseButtonEvent:
			if t.State != 0 {
				g.Click(int(t.X), int(t.Y))
			}
		}
	}
}

func renderLoop(g *Game, renderer *sdl.Renderer) {
	defer close(mainfunc)

	mainTicker := time.NewTicker(FRAME_RATE)

	for {
		select {
		case <-mainTicker.C:
			g.Update()
			do(func() {
				renderThings(renderer, g)
			})

		}
	}
}

func renderThings(renderer *sdl.Renderer, g *Game) {
	//fmt.Println("rendering")
	renderer.Clear()
	// show dashboard
	renderer.SetDrawColor(50, 50, 50, 200)
	dashboard := sdl.Rect{0, 0, DashboardWidth, WindowHeight}
	renderer.FillRect(&dashboard)

	// render blocks
	for _, b := range g.blocks {
		renderBlock(renderer, b, g)
	}
	// render switches
	for _, s := range g.switches {
		renderSwitch(renderer, s, g)
	}
	renderer.SetDrawColor(255, 255, 255, 255)
	renderer.Present()
}

var block = sdl.Rect{W: BlockSize, H: BlockSize}

//var source = sdl.Rect{W: BlockSize, H: BlockSize}

// renderBlock renders a block
func renderBlock(renderer *sdl.Renderer, b *Block, g *Game) {
	//fmt.Printf("showCandy (%d,%d), %d\n", c.x, c.y, c._type)
	block.X = int32(b.X)
	block.Y = int32(b.Y)
	alpha := uint8(255)

	switch b.Color {
	case Blue:
		renderer.SetDrawColor(153, 50, 204, alpha)
	case Yellow:
		renderer.SetDrawColor(255, 215, 0, alpha)
	case Green:
		renderer.SetDrawColor(60, 179, 113, alpha)
	case Red:
		renderer.SetDrawColor(220, 20, 60, alpha)
	case Pink:
		renderer.SetDrawColor(255, 192, 203, alpha)

	}
	renderer.FillRect(&block)
	renderer.SetDrawColor(255, 255, 255, 255)
}

var switch_ = sdl.Rect{W: SwitchSize, H: SwitchSize}

// renderSwitch renders a switch
func renderSwitch(renderer *sdl.Renderer, s *Switch, g *Game) {
	switch_.X = int32(s.X)
	switch_.Y = int32(s.Y)

	renderer.FillRect(&switch_)
}
