package main

import (
	"unsafe"

	"github.com/go-gl/gl"
)

type Triangle struct {
	buffer     gl.Buffer
	vertexData []float32
	sizeData   int
	prg        gl.Program
	posLoc     gl.AttribLocation
	colLoc     gl.AttribLocation
	vao        gl.VertexArray
}

func NewTriangle(vertexData []float32) *Triangle {
	t := &Triangle{}
	t.vertexData = vertexData
	t.sizeData = len(t.vertexData) * int(unsafe.Sizeof(float32(0)))

	vshader := loadShader(gl.VERTEX_SHADER, "shaders/colorv.vert")
	fshader := loadShader(gl.FRAGMENT_SHADER, "shaders/colorv.frag")
	t.prg = NewProgram(vshader, fshader)
	t.posLoc = gl.AttribLocation(0)
	t.colLoc = gl.AttribLocation(1)

	t.buffer = gl.GenBuffer()
	t.buffer.Bind(gl.ARRAY_BUFFER)

	gl.BufferData(gl.ARRAY_BUFFER, t.sizeData, nil, gl.STATIC_DRAW)
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, t.sizeData, t.vertexData)

	t.buffer.Unbind(gl.ARRAY_BUFFER)

	t.vao = gl.GenVertexArray()
	t.vao.Bind()
	return t
}

func (t *Triangle) Draw() {
	t.prg.Use()
	t.buffer.Bind(gl.ARRAY_BUFFER)

	t.posLoc.EnableArray()
	t.posLoc.AttribPointer(4, gl.FLOAT, false, 0, uintptr(0))
	t.colLoc.EnableArray()
	t.colLoc.AttribPointer(4, gl.FLOAT, false, 0, uintptr(48))

	//t.colLoc.AttribPointer(4, gl.FLOAT, false, 0, uintptr(t.sizeData))
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
