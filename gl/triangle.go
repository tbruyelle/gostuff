package main

import (
	"unsafe"

	"github.com/go-gl/gl"
)

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

func NewTriangle(vertices []float32) *Triangle {
	t := &Triangle{}
	t.vertices = vertices

	floatSize := int(unsafe.Sizeof(0))
	t.sizeVertices = len(t.vertices) * floatSize
	t.sizeColors = len(t.colors) * floatSize

	vshader := loadShader(gl.VERTEX_SHADER, "shaders/basic.vert")
	fshader := loadShader(gl.FRAGMENT_SHADER, "shaders/colory.frag")

	t.prg = NewProgram(vshader, fshader)
	t.posLoc = gl.AttribLocation(0)

	//t.buffer.Delete()
	t.buffer = gl.GenBuffer()
	t.buffer.Bind(gl.ARRAY_BUFFER)

	gl.BufferData(gl.ARRAY_BUFFER, t.sizeVertices, nil, gl.STATIC_DRAW)
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, t.sizeVertices, t.vertices)
	//gl.BufferSubData(gl.ARRAY_BUFFER, t.sizeVertices, t.sizeColors, t.colors)

	t.buffer.Unbind(gl.ARRAY_BUFFER)

	t.vao = gl.GenVertexArray()
	t.vao.Bind()
	return t
}

func (t *Triangle) Draw() {
	t.prg.Use()
	t.buffer.Bind(gl.ARRAY_BUFFER)

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
