package main

import (
	"fmt"
	"unsafe"

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

var vertices = []float32{-0.5, 0.0, 0.5, 0.5, 0.5, -0.5, // triange 1
	-0.8, -0.8, -0.3, -0.8, -0.8, -0.3} // Triangle 2
var colors = []float32{1.0, 0, 0}

var prg gl.Program
var buffer gl.Buffer

func main() {
	glfw.SetErrorCallback(errorCallback)

	if !glfw.Init() {
		panic("Can't init glfw!")
	}
	defer glfw.Terminate()

	//showVersion()
	//glfw.WindowHint(glfw.ContextVersionMajor, 3)
	//glfw.WindowHint(glfw.ContextVersionMinor, 3)
	//glfw.WindowHint(glfw.OpenglForwardCompatible, glfw.True)
	//glfw.WindowHint(glfw.OpenglProfile, glfw.OpenglCoreProfile)
	//showVersion()

	window, err := glfw.CreateWindow(WINDOW_WIDTH, WINDOW_HEIGHT, TITLE, nil, nil)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()
	window.MakeContextCurrent()

	glfw.SwapInterval(1)

	gl.Init()
	//initScene()
	defer destroyScene()

	for !window.ShouldClose() {
		drawScene()
		window.SwapBuffers()
	}
}

func showVersion() {
	//maj, min, v := glfw.GetVersion()
	//fmt.Println("version=", maj, min, v)
	fmt.Println("version=", glfw.GetVersionString())
}

func initScene() error {
	buffer = gl.GenBuffer()
	buffer.Bind(gl.ARRAY_BUFFER)
	var f float32
	floatSize := int(unsafe.Sizeof(f))
	sizeVert := len(vertices) * floatSize
	sizeColors := len(colors) * floatSize

	gl.BufferData(gl.ARRAY_BUFFER, sizeVert+sizeColors, nil, gl.STATIC_DRAW)
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, sizeVert, vertices)
	gl.BufferSubData(gl.ARRAY_BUFFER, sizeVert, sizeColors, vertices)
	buffer.Unbind(gl.ARRAY_BUFFER)

	//gl.Viewport(0, 0, WINDOW_WIDTH, WINDOW_HEIGHT)
	//gl.MatrixMode(gl.PROJECTION)
	vs := `
	#version 330
	in vec3 in_vertex;

		void main() {
			gl_Position = vec4(in_vertex, 0.0);
				}
				`
	fs := `
	#version 330
	out vec4 out_color;

		void main() {
			out_color= vec4(1.0,1.0,1.0,1.0);
				}
				`
	vshader := gl.CreateShader(gl.VERTEX_SHADER)
	vshader.Source(vs)
	vshader.Compile()
	if vshader.Get(gl.COMPILE_STATUS) != gl.TRUE {
		panic("vertex shader error :" + vshader.GetInfoLog())
	}
	fshader := gl.CreateShader(gl.FRAGMENT_SHADER)
	fshader.Source(fs)
	fshader.Compile()
	if fshader.Get(gl.COMPILE_STATUS) != gl.TRUE {
		panic("fragment shader error:" + fshader.GetInfoLog())
	}

	prg = gl.CreateProgram()
	//prg.AttachShader(vshader)
	//prg.AttachShader(fshader)
	//prg.BindAttribLocation(0, "in_vertex")
	//prg.BindAttribLocation(1, "out_color")
	prg.Link()
	if prg.Get(gl.LINK_STATUS) != gl.TRUE {
		panic("linker error: " + prg.GetInfoLog())
	}

	return nil
}

func destroyScene() {
	prg.Delete()
}

func drawScene() {
	gl.Clear(gl.COLOR_BUFFER_BIT)
	//prg.Use()
	gl.EnableClientState(gl.VERTEX_ARRAY)
	gl.VertexPointer(2, gl.FLOAT, 0, vertices)
	gl.DrawArrays(gl.TRIANGLES, 0, 6)
	gl.DisableClientState(gl.VERTEX_ARRAY)
}
