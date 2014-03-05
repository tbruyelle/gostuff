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

	for _, s := range g.switches {
		renderSwitch(s)
	}
	// TODO What for?
	//gl.Flush()
	window.SwapBuffers()
}

func renderSwitch(s *Switch) {
	// TODO pb to call it on every block?
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
	// TODO constant
	v := SwitchSize / 2

	gl.Translatef(float32(s.X+v), float32(s.Y+v), 0)

	if s.rotate > 0 {
		gl.Rotatef(float32(s.rotate), 0, 0, 1)
	}

	// render block top left
	gl.Begin(gl.QUADS)
	setColor(s.blocks[0].Color)
	gl.Vertex2i(-BlockSize, -BlockSize)
	gl.Vertex2i(0, -BlockSize)
	gl.Vertex2i(0, 0)
	gl.Vertex2i(-BlockSize, 0)
	gl.End()

	// render block top right
	gl.Begin(gl.QUADS)
	setColor(s.blocks[1].Color)
	gl.Vertex2i(0, -BlockSize)
	gl.Vertex2i(BlockSize, -BlockSize)
	gl.Vertex2i(BlockSize, 0)
	gl.Vertex2i(0, 0)
	gl.End()

	// render block bottom right
	gl.Begin(gl.QUADS)
	setColor(s.blocks[2].Color)
	gl.Vertex2i(0, 0)
	gl.Vertex2i(BlockSize, 0)
	gl.Vertex2i(BlockSize, BlockSize)
	gl.Vertex2i(0, BlockSize)
	gl.End()

	// render block bottom left
	gl.Begin(gl.QUADS)
	setColor(s.blocks[3].Color)
	gl.Vertex2i(-BlockSize, 0)
	gl.Vertex2i(0, 0)
	gl.Vertex2i(0, BlockSize)
	gl.Vertex2i(-BlockSize, BlockSize)
	gl.End()

	// render the switch
	gl.Begin(gl.QUADS)
	gl.Color3f(1.0, 1.0, 1.0)
	renderSquare(v)
	gl.End()
}

func setColor(color ColorDef) {
	switch color {
	case Red:
		gl.Color3f(1.0, 0.0, 0.0)
	case Blue:
		gl.Color3f(0.0, 0.0, 1.0)

	case Green:
		gl.Color3f(0.0, 1.0, 0.0)

	case Yellow:
		gl.Color3f(0.5, 0.0, 0.3)
	}
}

func renderSquare(side int) {
	gl.Vertex2i(-side, side)
	gl.Vertex2i(side, side)
	gl.Vertex2i(side, -side)
	gl.Vertex2i(-side, -side)
}
