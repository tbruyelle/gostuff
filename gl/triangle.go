package main

import (
	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
)

type Triangle struct {
	buffer       gl.Buffer
	vertices     []Vertex
	sizeVertices int
	prg          gl.Program
	posLoc       gl.AttribLocation
	colLoc       gl.AttribLocation
	timeLoc      gl.UniformLocation
	vao          gl.VertexArray
}

func NewTriangle(vertices []Vertex) *Triangle {
	t := &Triangle{}
	t.vertices = vertices
	t.sizeVertices = len(t.vertices) * sizeVertex

	vshader := loadShader(gl.VERTEX_SHADER, "shaders/rotateOffset.vert")
	fshader := loadShader(gl.FRAGMENT_SHADER, "shaders/offset.frag")
	t.prg = NewProgram(vshader, fshader)
	t.posLoc = gl.AttribLocation(0)
	t.colLoc = gl.AttribLocation(1)
	t.timeLoc = t.prg.GetUniformLocation("time")
	loopLoc := t.prg.GetUniformLocation("loopDuration")
	fragLoopLoc := t.prg.GetUniformLocation("fragLoopDuration")

	t.buffer = gl.GenBuffer()
	t.buffer.Bind(gl.ARRAY_BUFFER)

	gl.BufferData(gl.ARRAY_BUFFER, t.sizeVertices, nil, gl.STATIC_DRAW)
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, t.sizeVertices, t.vertices)

	t.buffer.Unbind(gl.ARRAY_BUFFER)

	t.vao = gl.GenVertexArray()
	t.vao.Bind()

	t.prg.Use()
	loopLoc.Uniform1f(5)
	fragLoopLoc.Uniform1f(10)
	gl.ProgramUnuse()

	return t
}

func (t *Triangle) Draw() {
	t.prg.Use()

	t.timeLoc.Uniform1f(float32(glfw.GetTime()))

	t.buffer.Bind(gl.ARRAY_BUFFER)

	t.posLoc.EnableArray()
	t.posLoc.AttribPointer(4, gl.FLOAT, false, sizeVertex, uintptr(0))
	t.colLoc.EnableArray()
	t.colLoc.AttribPointer(4, gl.FLOAT, false, sizeVertex, uintptr(sizeCoords))

	//for i := float64(0); i < 30; i++ {
		gl.DrawArrays(gl.TRIANGLES, 0, len(t.vertices))
	//	t.timeLoc.Uniform1f(float32(glfw.GetTime() + i/20))
	//}

	t.posLoc.DisableArray()
	t.buffer.Unbind(gl.ARRAY_BUFFER)
	gl.ProgramUnuse()
}

func (t *Triangle) Destroy() {
	t.buffer.Delete()
	t.vao.Delete()
}
