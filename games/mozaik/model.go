package main

import (
	"fmt"
	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
)

type Model struct {
	buffer       gl.Buffer
	vertices     []Vertex
	sizeVertices int
	prg          gl.Program
	posLoc       gl.AttribLocation
	colLoc       gl.AttribLocation
	timeLoc      gl.UniformLocation
	vao          gl.VertexArray
}

func NewModel(vertices []Vertex, vshaderf, fshaderf string) *Model {
	t := &Model{}
	t.vertices = vertices
	fmt.Printf("newmodel with %d vertices", len(vertices))
	t.sizeVertices = len(t.vertices) * sizeVertex

	vshader := loadShader(gl.VERTEX_SHADER, vshaderf)
	fshader := loadShader(gl.FRAGMENT_SHADER, fshaderf)
	t.prg = NewProgram(vshader, fshader)
	t.posLoc = gl.AttribLocation(0)
	t.colLoc = gl.AttribLocation(1)
	t.timeLoc = t.prg.GetUniformLocation("time")
	loopLoc := t.prg.GetUniformLocation("loopDuration")
	fragLoopLoc := t.prg.GetUniformLocation("fragLoopDuration")
	perspectiveMatrix := t.prg.GetUniformLocation("perpectiveMatrix")

	// the projection matrix
	var frustrumScale, zNear, zFar float32 = 1.0, 0.5, 3.0
	var matrix [16]float32

	matrix[0] = frustrumScale
	matrix[5] = frustrumScale
	matrix[10] = (zFar + zNear) / (zNear - zFar)
	matrix[14] = (2 * zFar * zNear) / (zNear - zFar)
	matrix[11] = -1.0

	t.buffer = gl.GenBuffer()
	t.buffer.Bind(gl.ARRAY_BUFFER)

	gl.BufferData(gl.ARRAY_BUFFER, t.sizeVertices, nil, gl.STATIC_DRAW)
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, t.sizeVertices, t.vertices)

	t.buffer.Unbind(gl.ARRAY_BUFFER)

	t.vao = gl.GenVertexArray()
	t.vao.Bind()

	gl.Enable(gl.CULL_FACE)
	gl.CullFace(gl.BACK)
	gl.FrontFace(gl.CW)
	t.prg.Use()
	loopLoc.Uniform1f(5)
	fragLoopLoc.Uniform1f(10)
	perspectiveMatrix.UniformMatrix4fv(true, matrix)
	gl.ProgramUnuse()

	return t
}

func (t *Model) Draw() {
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

func (t *Model) Destroy() {
	t.buffer.Delete()
	t.vao.Delete()
}
