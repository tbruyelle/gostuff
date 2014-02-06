package main

import (
	"fmt"
	"github.com/go-gl/gl"
glfw 	"github.com/go-gl/glfw3"
	"io/ioutil"
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
