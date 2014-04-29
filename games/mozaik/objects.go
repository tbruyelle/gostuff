package main

import (
	"github.com/go-gl/gl"
	"github.com/remogatto/mathgl"
	"math"
)

type BlockModel struct {
	Model
}

func NewBlockModel(bl *Block) *BlockModel {
	model := &BlockModel{}

	vs := []Vertex{
		NewVertex(0, 0, 0, BlueColor),
		NewVertex(0, BlockSize, 0, BlueColor),
		NewVertex(BlockSize, 0, 0, BlueColor),
		NewVertex(BlockSize, BlockSize, 0, BlueColor),
	}
	model.Init(vs, "shaders/basic.vert", "shaders/basic.frag")
	return model
}

func (t *BlockModel) Draw() {
	t.prg.Use()

	t.buffer.Bind(gl.ARRAY_BUFFER)

	t.posLoc.EnableArray()
	t.posLoc.AttribPointer(4, gl.FLOAT, false, sizeVertex, uintptr(0))
	t.colLoc.EnableArray()
	t.colLoc.AttribPointer(4, gl.FLOAT, false, sizeVertex, uintptr(sizeCoords))

	t.uniformModelView.UniformMatrix4f(false, (*[16]float32)(&t.modelView))
	t.uniformProjectionView.UniformMatrix4f(false, (*[16]float32)(&t.projectionView))

	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, len(t.vertices))

	t.posLoc.DisableArray()
	t.colLoc.DisableArray()
	t.buffer.Unbind(gl.ARRAY_BUFFER)
	gl.ProgramUnuse()
}

type SwitchModel struct {
	Model
	sw  *Switch
	lvl *Level
}

func NewSwitchModel(sw *Switch, lvl *Level) *SwitchModel {
	model := &SwitchModel{sw: sw, lvl: lvl}

	vs := []Vertex{NewVertex(0, 0, 0, WhiteColor)}
	vv := float64(SwitchSize / 2)
	for i := float64(0); i <= SwitchSegments; i++ {
		a := 2 * math.Pi * i / SwitchSegments
		vs = append(vs, NewVertex(float32(math.Sin(a)*vv), float32(math.Cos(a)*vv), 0, WhiteColor))
	}
	model.Init(vs, "shaders/basic.vert", "shaders/basic.frag")

	v := SwitchSize / 2
	model.projectionView = mathgl.Ortho2D(0, WindowWidth, WindowHeight, 0).Mul4(mathgl.Translate3D(float32(sw.X+v), float32(sw.Y+v), 0))

	return model
}

// TODO the switch number
func (t *SwitchModel) Draw() {
	// Draw the blocks
	var b *Block
	s := t.sw
	bsf := float32(BlockSize - s.Z)
	// top left block
	b = t.lvl.blocks[s.line][s.col]
	if !b.Rendered {
		bl := NewBlockModel(b)
		bl.projectionView = t.projectionView.Mul4(mathgl.Translate3D(-bsf, -bsf, 0))
		bl.Draw()
		b.Rendered = true
	}
	// top right block
	b = t.lvl.blocks[s.line][s.col+1]
	if !b.Rendered {
		bl := NewBlockModel(b)
		bl.projectionView = t.projectionView.Mul4(mathgl.Translate3D(0, -bsf, 0))
		bl.Draw()
		b.Rendered = true
	}
	// bottom right block
	b = t.lvl.blocks[s.line+1][s.col+1]
	if !b.Rendered {
		bl := NewBlockModel(b)
		bl.projectionView = t.projectionView
		bl.Draw()
		b.Rendered = true
	}
	// bottom left block
	b = t.lvl.blocks[s.line+1][s.col]
	if !b.Rendered {
		bl := NewBlockModel(b)
		bl.projectionView = t.projectionView.Mul4(mathgl.Translate3D(-bsf, 0, 0))
		bl.Draw()
		b.Rendered = true
	}

	// Draw the switch
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
	t.colLoc.DisableArray()
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
	t.colLoc.DisableArray()
	t.buffer.Unbind(gl.ARRAY_BUFFER)
	gl.ProgramUnuse()
}
