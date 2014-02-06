package main

import (
	"fmt"
	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
	"io/ioutil"
	"unsafe"
)

type Vertex struct {
	Coords Coords
	Color  Color
}

func NewVertex(X, Y, Z float32, color Color) Vertex {
	return Vertex{Coords: Coords{X, Y, Z, 1.0}, Color: color}
}

var (
	Red   = Color{1.0, 0.0, 0.0, 1.0}
	Green = Color{0.0, 1.0, 0.0, 1.0}
	Blue  = Color{0.0, 0.0, 1.0, 1.0}
)

type Coords struct{ X, Y, Z, W float32 }
type Color struct{ R, G, B, A float32 }

var (
	sizeFloat  = int(unsafe.Sizeof(float32(0)))
	sizeCoords = sizeFloat * 4
	sizeVertex = int(unsafe.Sizeof(Vertex{}))
)

func showVersion() {
	//maj, min, v := glfw.GetVersion()
	//fmt.Println("version=", maj, min, v)
	fmt.Println("version=", glfw.GetVersionString())
}
func NewProgram(shaders ...gl.Shader) gl.Program {
	prg := gl.CreateProgram()
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
	return prg
}

func loadShader(type_ gl.GLenum, file string) gl.Shader {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	source := string(b)
	shader := gl.CreateShader(type_)
	shader.Source(source)
	shader.Compile()
	if shader.Get(gl.COMPILE_STATUS) != gl.TRUE {
		panic("fragment error for source " + source + "\n" + shader.GetInfoLog())
	}
	return shader
}
