package main

import (
	"fmt"
	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
	"runtime"
	"time"
)

const FRAME_RATE = time.Second / 40

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

func errorCallback(err glfw.ErrorCode, desc string) {
	fmt.Printf("%v: %v\n", err, desc)
}

var (
	window *glfw.Window
	err    error
	g      *Game
)

func main() {
	_ = fmt.Sprint()

	if !glfw.Init() {
		panic("Can't init glfw!")
	}
	defer glfw.Terminate()

	window, err = glfw.CreateWindow(WindowWidth, WindowHeight, "Mozaik", nil, nil)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	// Ensure thread context
	window.MakeContextCurrent()

	// TODO WHat fort ?
	//glfw.SwapInterval(1)

	window.SetKeyCallback(keyCb)
	window.SetMouseButtonCallback(mouseCb)

	gl.Init()
	gl.ClearColor(0.9, 0.9, 0.9, 0.0)

	// Use window coordinates
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Ortho(0, WindowWidth, WindowHeight, 0, 0, 1)

	g = NewGame()

	g.Start()
	go eventLoop(g)
	go renderLoop(g)
	Main()
	g.Stop()
}

func keyCb(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Press {
		switch key {
		case glfw.KeyEscape:
			close(mainfunc)
		}
	}
}

func mouseCb(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
	if action == glfw.Press {
		switch button {

		case glfw.MouseButtonLeft:
			x, y := window.GetCursorPosition()
			g.Click(int(x), int(y))
		}
	}
}

func eventLoop(g *Game) {
	defer close(mainfunc)
	for {
		do(func() {
			glfw.PollEvents()
		})
	}
}

func renderLoop(g *Game) {
	defer close(mainfunc)

	mainTicker := time.NewTicker(FRAME_RATE)

	for {
		select {
		case <-mainTicker.C:
			g.Update()
			do(func() {
				renderThings(g)
			})

		}
	}
}

func renderThings(g *Game) {
	gl.Clear(gl.COLOR_BUFFER_BIT)

	for _, b := range g.blocks {
		gl.Begin(gl.QUADS)
		switch b.Color {
		case Red:

			gl.Color3f(1.0, 0.0, 0.0)
		case Blue:
			gl.Color3f(0.0, 0.0, 1.0)

		case Green:
			gl.Color3f(0.0, 1.0, 0.0)

		case Pink:
			gl.Color3f(0.5, 0.0, 0.3)
		}
		gl.Vertex2i(b.X, b.Y)
		gl.Vertex2i(b.X+BlockSize, b.Y)
		gl.Vertex2i(b.X+BlockSize, b.Y+BlockSize)
		gl.Vertex2i(b.X, b.Y+BlockSize)
		gl.End()
	}
	window.SwapBuffers()
	//	//fmt.Println("rendering")
	//	renderer.Clear()
	//	// show dashboard
	//	renderer.SetDrawColor(50, 50, 50, 200)
	//	dashboard := sdl.Rect{0, 0, DashboardWidth, WindowHeight}
	//	renderer.FillRect(&dashboard)
	//
	//	// render blocks
	//	for _, b := range g.blocks {
	//		renderBlock(renderer, b, g)
	//	}
	//	// render switches
	//	for _, s := range g.switches {
	//		renderSwitch(renderer, s, g)
	//	}
	//	renderer.SetDrawColor(255, 255, 255, 255)
	//	renderer.Present()
}

//var block = sdl.Rect{W: BlockSize, H: BlockSize}

//var source = sdl.Rect{W: BlockSize, H: BlockSize}

// renderBlock renders a block
//func renderBlock(renderer *sdl.Renderer, b *Block, g *Game) {
//	//fmt.Printf("showCandy (%d,%d), %d\n", c.x, c.y, c._type)
//	block.X = int32(b.X)
//	block.Y = int32(b.Y)
//	alpha := uint8(255)
//
//	switch b.Color {
//	case Blue:
//		renderer.SetDrawColor(153, 50, 204, alpha)
//	case Yellow:
//		renderer.SetDrawColor(255, 215, 0, alpha)
//	case Green:
//		renderer.SetDrawColor(60, 179, 113, alpha)
//	case Red:
//		renderer.SetDrawColor(220, 20, 60, alpha)
//	case Pink:
//		renderer.SetDrawColor(255, 192, 203, alpha)
//
//	}
//	renderer.FillRect(&block)
//	renderer.SetDrawColor(255, 255, 255, 255)
//}
//
//var switch_ = sdl.Rect{W: SwitchSize, H: SwitchSize}
//
//// renderSwitch renders a switch
//func renderSwitch(renderer *sdl.Renderer, s *Switch, g *Game) {
//	switch_.X = int32(s.X)
//	switch_.Y = int32(s.Y)
//
//	renderer.FillRect(&switch_)
//}
