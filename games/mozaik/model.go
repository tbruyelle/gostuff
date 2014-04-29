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
	vshader, fshader                        gl.Shader
}

func (t *Model) Init(vertices []Vertex, vshaderf, fshaderf string) {
	t.vertices = vertices
	t.sizeVertices = len(t.vertices) * sizeVertex

	// Shaders
	t.vshader = loadShader(gl.VERTEX_SHADER, vshaderf)
	t.fshader = loadShader(gl.FRAGMENT_SHADER, fshaderf)
	t.prg = NewProgram(t.vshader, t.fshader)
	t.posLoc = gl.AttribLocation(0)
	t.colLoc = gl.AttribLocation(1)
	t.uniformProjectionView = t.prg.GetUniformLocation("projectionView")
	t.uniformModelView = t.prg.GetUniformLocation("modelView")

	// the projection matrix
	t.projectionView = mathgl.Ident4f()

	// the model view
	t.modelView = mathgl.Ident4f()

	// Create VBO
	t.buffer = gl.GenBuffer()
	t.buffer.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, t.sizeVertices, nil, gl.STATIC_DRAW)
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, t.sizeVertices, t.vertices)
	t.buffer.Unbind(gl.ARRAY_BUFFER)

	// Create VAO
	t.vao = gl.GenVertexArray()
	t.vao.Bind()
	t.buffer.Bind(gl.ARRAY_BUFFER)

	// Attrib vertex data to VAO
	t.posLoc.AttribPointer(4, gl.FLOAT, false, sizeVertex, uintptr(0))
	t.posLoc.EnableArray()
	t.colLoc.AttribPointer(4, gl.FLOAT, false, sizeVertex, uintptr(sizeCoords))
	t.colLoc.EnableArray()

	t.buffer.Unbind(gl.ARRAY_BUFFER)
	(gl.VertexArray(0)).Bind()
}

func (t *Model) Destroy() {
	t.buffer.Delete()
	t.vao.Delete()
	t.vshader.Delete()
	t.fshader.Delete()
	t.prg.Delete()
}
