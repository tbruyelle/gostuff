package main

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
)

const (
	WINDOW_WIDTH  = 640
	WINDOW_HEIGHT = 480
	TITLE         = "Fuck yeah opengl"
)

func errorCallback(err glfw.ErrorCode, desc string) {
	fmt.Printf("Error callback %v: %v\n", err, desc)
}

var vertices = [...]float32{-0.5, 0.0, 0.5, 0.5, 0.5, -0.5} // triange 1
//	-0.8, -0.8, -0.3, -0.8, -0.8, -0.3} // Triangle 2
var colors = [...]float32{1.0, 0, 0}

var triangle *Triangle

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
	initScene()
	defer destroyScene()

	for !window.ShouldClose() {
		drawScene()
		window.SwapBuffers()
		<-time.After(time.Second * 2)
		break
	}
}

func showVersion() {
	//maj, min, v := glfw.GetVersion()
	//fmt.Println("version=", maj, min, v)
	fmt.Println("version=", glfw.GetVersionString())
}

func initScene() {
	triangle = NewTriangle(vertices, colors)
	triangle.Load()
}

func drawScene() {
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	triangle.Draw()
}

func loadShader(type_ gl.GLenum, file string) gl.Shader {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	source:=string(b)
	shader := gl.CreateShader(type_)
	shader.Source(source)
	shader.Compile()
	if shader.Get(gl.COMPILE_STATUS) != gl.TRUE {
		panic("fragment error for source " + source + "\n" + shader.GetInfoLog())
	}
	return shader
}

func attachShaders(prg gl.Program, shaders ...gl.Shader) {
	for _, shader := range shaders {
		prg.AttachShader(shader)
	}
	prg.Link()
	if prg.Get(gl.LINK_STATUS) != gl.TRUE {
		panic("linker error: " + prg.GetInfoLog())
	}
	prg.Validate()
	for _, shader := range shaders {
		prg.DetachShader(shader)
		shader.Delete()
	}
}

func destroyScene() {
	triangle.Destroy()
}
