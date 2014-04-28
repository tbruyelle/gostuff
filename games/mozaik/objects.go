package main

import (
	"github.com/go-gl/gl"
	"github.com/remogatto/mathgl"
	"math"
)

type SwitchModel struct {
	Model
	sw *Switch
}

func NewSwitchModel(sw *Switch) *SwitchModel {
	model := &SwitchModel{sw: sw}

	vs := []Vertex{NewVertex(0, 0, 0, WhiteColor)}
	vv := float64(SwitchSize / 2)
	for i := float64(0); i <= SwitchSegments; i++ {
		a := 2 * math.Pi * i / SwitchSegments
		vs = append(vs, NewVertex(float32(math.Sin(a)*vv), float32(math.Cos(a)*vv), 0, WhiteColor))
	}
	model.Init(vs, "shaders/basic.vert", "shaders/basic.frag")
	return model
}

// TODO the switch number
func (t *SwitchModel) Draw() {
	v := SwitchSize / 2
	t.modelView = mathgl.Ident4f()
	t.projectionView = mathgl.Ortho2D(0, WindowWidth, WindowHeight, 0).Mul4(mathgl.Translate3D(float32(t.sw.X+v), float32(t.sw.Y+v), 0))

	t.prg.Use()

	t.buffer.Bind(gl.ARRAY_BUFFER)

	t.posLoc.EnableArray()
	t.posLoc.AttribPointer(4, gl.FLOAT, false, sizeVertex, uintptr(0))
	t.colLoc.EnableArray()
	t.colLoc.AttribPointer(4, gl.FLOAT, false, sizeVertex, uintptr(sizeCoords))

	t.uniformModelView.UniformMatrix4f(false, (*[16]float32)(&t.modelView))
	t.uniformProjectionView.UniformMatrix4f(false, (*[16]float32)(&t.projectionView))

	gl.DrawArrays(gl.TRIANGLE_FAN, 0, len(t.vertices))

	t.posLoc.DisableArray()
	t.buffer.Unbind(gl.ARRAY_BUFFER)
	gl.ProgramUnuse()
}

type Background struct {
	Model
	angle float32
}

func NewBackground() *Background {
	model := &Background{}
	vs := []Vertex{}

	for i := float64(0); i <= BgSegments; i++ {
		if math.Mod(i, 2) == 0 {
			vs = append(vs, NewVertex(0, 0, 0, BgColor))
		}
		a := 2 * math.Pi * i / BgSegments
		vs = append(vs, NewVertex(float32(math.Sin(a)*windowRadius), float32(math.Cos(a)*windowRadius), 0, BgColor))
	}
	model.Init(vs, "shaders/basic.vert", "shaders/basic.frag")
	return model
}

func (t *Background) Draw() {
	t.angle += 0.05
	t.modelView = mathgl.HomogRotate3D(t.angle, [3]float32{0, 0, 1})

	t.prg.Use()

	t.buffer.Bind(gl.ARRAY_BUFFER)

	t.posLoc.EnableArray()
	t.posLoc.AttribPointer(4, gl.FLOAT, false, sizeVertex, uintptr(0))
	t.colLoc.EnableArray()
	t.colLoc.AttribPointer(4, gl.FLOAT, false, sizeVertex, uintptr(sizeCoords))

	t.uniformModelView.UniformMatrix4f(false, (*[16]float32)(&t.modelView))
	t.uniformProjectionView.UniformMatrix4f(false, (*[16]float32)(&t.projectionView))

	gl.DrawArrays(gl.TRIANGLES, 0, len(t.vertices))

	t.posLoc.DisableArray()
	t.buffer.Unbind(gl.ARRAY_BUFFER)
	gl.ProgramUnuse()
}
