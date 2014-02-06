package main

import (
	"fmt"

	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
)

const (
	WINDOW_WIDTH  = 500
	WINDOW_HEIGHT = 500
	TITLE         = "Fuck yeah opengl"
)

func errorCallback(err glfw.ErrorCode, desc string) {
	fmt.Printf("Error callback %v: %v\n", err, desc)
}

type World struct {
	triangle *Triangle
}

var world = &World{}

func main() {
	glfw.SetErrorCallback(errorCallback)

	if !glfw.Init() {
		panic("Can't init glfw!")
	}
	defer glfw.Terminate()

	showVersion()
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenglForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.OpenglProfile, glfw.OpenglCoreProfile)

	window, err := glfw.CreateWindow(WINDOW_WIDTH, WINDOW_HEIGHT, TITLE, nil, nil)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()
	window.MakeContextCurrent()

	glfw.SwapInterval(1)

	gl.Init()
	world.initScene()
	defer world.destroyScene()

	for !window.ShouldClose() {
		world.drawScene()
		window.SwapBuffers()
	}
}

func (w *World) initScene() {
	w.triangle = NewTriangle([]Vertex{
		NewVertex(0.0, 0.5, 0.0, Red),
		NewVertex(0.5, -0.366, 0.0, Green),
		NewVertex(-0.5, -0.366, 0.0, Blue),
	})
}

func (w *World) drawScene() {
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	w.triangle.Draw()
}

func (w *World) destroyScene() {
	w.triangle.Destroy()
}
