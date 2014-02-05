package main

import (
	"fmt"
	"unsafe"

	"github.com/go-gl/gl"
)

const shaderVert = `
#version 330

layout(location = 0) in vec4 vert;

//uniform mat4 projection;
//uniform mat4 view;
//uniform mat4 model;

void main()
{
    //gl_Position = projection * view * model * vert;
    gl_Position = vert;
}`

const shaderFrag = `
#version 330

out vec4 fragColor;

void main() {
  fragColor = vec4(1.0f, 0.0f, 0.0f, 1.0f);
}`

type Triangle struct {
	buffer       gl.Buffer
	vertices     []float32
	colors       []float32
	sizeVertices int
	sizeColors   int
	prg          gl.Program
	posLoc       gl.AttribLocation
	colLoc       gl.AttribLocation
	vao          gl.VertexArray
}

func NewTriangle(vertices [6]float32, colors [3]float32) *Triangle {
	t := &Triangle{}
	t.vertices = []float32{
		0.75, 0.75, 0.0, 1.0,
		0.75, -0.75, 0.0, 1.0,
		-0.75, -0.75, 0.0, 1.0,
	}
	t.colors = []float32{1.0, 0, 0, 1.0}

	var f float32
	floatSize := int(unsafe.Sizeof(f))
	t.sizeVertices = len(t.vertices) * floatSize
	t.sizeColors = len(t.colors) * floatSize

	vshader := loadShader(gl.VERTEX_SHADER, shaderVert)
	fshader := loadShader(gl.FRAGMENT_SHADER, shaderFrag)

	t.prg = gl.CreateProgram()
	attachShaders(t.prg, vshader, fshader)
	t.posLoc = t.prg.GetAttribLocation("vert")
	//t.colLoc = t.prg.GetAttribLocation("outColor")
	fmt.Println("postLoc=", t.posLoc)
	//fmt.Println("colLoc=", t.colLoc)
	return t
}

func (t *Triangle) Load() {

	//t.buffer.Delete()
	t.buffer = gl.GenBuffer()
	t.buffer.Bind(gl.ARRAY_BUFFER)

	gl.BufferData(gl.ARRAY_BUFFER, t.sizeVertices, t.vertices, gl.STATIC_DRAW)
	//gl.BufferSubData(gl.ARRAY_BUFFER, 0, t.sizeVertices, t.vertices)
	//gl.BufferSubData(gl.ARRAY_BUFFER, t.sizeVertices, t.sizeColors, t.colors)

	t.buffer.Unbind(gl.ARRAY_BUFFER)

	t.vao = gl.GenVertexArray()
	t.vao.Bind()

}

func (t *Triangle) Draw() {
	t.prg.Use()
	t.buffer.Bind(gl.ARRAY_BUFFER)
	//t.posLoc.EnableArray()
	//t.colLoc.EnableArray()

	//t.colLoc.EnableArray()

	t.posLoc.AttribPointer(4, gl.FLOAT, false, 0, uintptr(0))
	t.posLoc.EnableArray()

	//t.colLoc.AttribPointer(4, gl.FLOAT, false, 0, uintptr(t.sizeVertices))
	gl.DrawArrays(gl.TRIANGLES, 0, 3)

	t.posLoc.DisableArray()
	//t.colLoc.DisableArray()
	t.buffer.Unbind(gl.ARRAY_BUFFER)
	gl.ProgramUnuse()
}

func (t *Triangle) Destroy() {
	t.buffer.Delete()
	t.vao.Delete()
}
