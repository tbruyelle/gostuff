package main

import (
	"fmt"

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

var prg gl.Program

func main() {
	glfw.SetErrorCallback(errorCallback)

	if !glfw.Init() {
		panic("Can't init glfw!")
	}
	defer glfw.Terminate()

	//showVersion()
	//glfw.WindowHint(glfw.ContextVersionMajor, 3)
	//glfw.WindowHint(glfw.ContextVersionMinor, 2)
	//glfw.WindowHint(glfw.OpenglForwardCompatible, glfw.True)
	//glfw.WindowHint(glfw.OpenglProfile, glfw.OpenglCoreProfile)
	showVersion()
	window, err := glfw.CreateWindow(WINDOW_WIDTH, WINDOW_HEIGHT, TITLE, nil, nil)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()
	window.MakeContextCurrent()

	glfw.SwapInterval(1)

	gl.Init()
	if err := initScene(); err != nil {
		panic(err)
	}
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
	//gl.Viewport(0, 0, WINDOW_WIDTH, WINDOW_HEIGHT)
	//gl.MatrixMode(gl.PROJECTION)
	vs := `
uniform vec2 offset;
	attribute vec4 vPosition;

		void main() {
				gl_Position = vec4(vPosition.xy+offset, vPosition.zw);
				}
				`
	fs := `
//precision mediump float;
	uniform vec3 color;

		void main() {
				gl_FragColor = vec4(color.xyz, 1.0);
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
	prg.AttachShader(vshader)
	prg.AttachShader(fshader)
	prg.BindAttribLocation(0, "vPosition")
	prg.BindAttribLocation(1, "offset")
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
