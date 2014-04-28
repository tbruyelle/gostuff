package main

import (
	"github.com/go-gl/gl"
	"github.com/remogatto/mathgl"
)

type Model struct {
	buffer                                  gl.Buffer
	vertices                                []Vertex
	sizeVertices                            int
	prg                                     gl.Program
	posLoc                                  gl.AttribLocation
	colLoc                                  gl.AttribLocation
	vao                                     gl.VertexArray
	uniformModelView, uniformProjectionView gl.UniformLocation
	modelView, projectionView               mathgl.Mat4f
}

func (t *Model) Init(vertices []Vertex, vshaderf, fshaderf string) {
	t.vertices = vertices
	t.sizeVertices = len(t.vertices) * sizeVertex

	// Shaders
	vshader := loadShader(gl.VERTEX_SHADER, vshaderf)
	fshader := loadShader(gl.FRAGMENT_SHADER, fshaderf)
	t.prg = NewProgram(vshader, fshader)
	t.posLoc = gl.AttribLocation(0)
	t.colLoc = gl.AttribLocation(1)
	t.uniformProjectionView = t.prg.GetUniformLocation("projectionView")
	t.uniformModelView = t.prg.GetUniformLocation("modelView")

	// the projection matrix
	t.projectionView = mathgl.Ident4f()

	// the model view
	t.modelView = mathgl.Ident4f()

	t.buffer = gl.GenBuffer()
	t.buffer.Bind(gl.ARRAY_BUFFER)

	gl.BufferData(gl.ARRAY_BUFFER, t.sizeVertices, nil, gl.STATIC_DRAW)
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, t.sizeVertices, t.vertices)

	t.buffer.Unbind(gl.ARRAY_BUFFER)

	t.vao = gl.GenVertexArray()
	t.vao.Bind()
}

func (t *Model) Destroy() {
	t.buffer.Delete()
	t.vao.Delete()
}
