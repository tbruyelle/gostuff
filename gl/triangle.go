package main

import (
	"github.com/go-gl/gl"
)

type Triangle struct {
	buffer       gl.Buffer
	vertices     []Vertex
	sizeVertices int
	prg          gl.Program
	posLoc       gl.AttribLocation
	colLoc       gl.AttribLocation
	vao          gl.VertexArray
}

func NewTriangle(vertices []Vertex) *Triangle {
	t := &Triangle{}
	t.vertices = vertices
	t.sizeVertices = len(t.vertices) * sizeVertex

	vshader := loadShader(gl.VERTEX_SHADER, "shaders/colorv.vert")
	fshader := loadShader(gl.FRAGMENT_SHADER, "shaders/colorv.frag")
	t.prg = NewProgram(vshader, fshader)
	t.posLoc = gl.AttribLocation(0)
	t.colLoc = gl.AttribLocation(1)

	t.buffer = gl.GenBuffer()
	t.buffer.Bind(gl.ARRAY_BUFFER)

	gl.BufferData(gl.ARRAY_BUFFER, t.sizeVertices, nil, gl.STATIC_DRAW)
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, t.sizeVertices, t.vertices)

	t.buffer.Unbind(gl.ARRAY_BUFFER)

	t.vao = gl.GenVertexArray()
	t.vao.Bind()
	return t
}

func (t *Triangle) Draw() {
	t.prg.Use()
	t.buffer.Bind(gl.ARRAY_BUFFER)

	t.posLoc.EnableArray()
	t.posLoc.AttribPointer(4, gl.FLOAT, false, sizeVertex, uintptr(0))
	t.colLoc.EnableArray()
	t.colLoc.AttribPointer(4, gl.FLOAT, false, sizeVertex, uintptr(sizeCoords))

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
