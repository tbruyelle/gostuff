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
		renderBlock(b)
	}
	for _, s := range g.switches {
		renderSwitch(s)
	}
	window.SwapBuffers()
}

func renderBlock(b *Block) {
	gl.Begin(gl.QUADS)
	switch b.Color {
	case Red:
		gl.Color3f(1.0, 0.0, 0.0)
	case Blue:
		gl.Color3f(0.0, 0.0, 1.0)

	case Green:
		gl.Color3f(0.0, 1.0, 0.0)

	case Yellow:
		gl.Color3f(0.5, 0.0, 0.3)
	}
	gl.Vertex2i(b.X, b.Y)
	gl.Vertex2i(b.X+BlockSize, b.Y)
	gl.Vertex2i(b.X+BlockSize, b.Y+BlockSize)
	gl.Vertex2i(b.X, b.Y+BlockSize)
	gl.End()
}

func renderSwitch(s *Switch) {
	gl.Begin(gl.QUADS)
	gl.Color3f(1.0, 1.0, 1.0)
	gl.Vertex2i(s.X, s.Y)
	gl.Vertex2i(s.X+SwitchSize, s.Y)
	gl.Vertex2i(s.X+SwitchSize, s.Y+SwitchSize)
	gl.Vertex2i(s.X, s.Y+SwitchSize)
	gl.End()
}

//
//// renderSwitch renders a switch
//func renderSwitch(renderer *sdl.Renderer, s *Switch, g *Game) {
//	switch_.X = int32(s.X)
//	switch_.Y = int32(s.Y)
//
//	renderer.FillRect(&switch_)
//}
