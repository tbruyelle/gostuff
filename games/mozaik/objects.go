package main

import (
	"github.com/go-gl/gl"
	"github.com/remogatto/mathgl"
	"math"
)

type BlockModel struct {
	ModelGroup
	block *Block
}

func getBlockColor(b *Block) Color {
	switch b.Color {
	case Red:
		return RedColor
	case Blue:
		return BlueColor
	case LightBlue:
		return LightBlueColor
	case Orange:
		return OrangeColor
	case Green:
		return GreenColor
	case Pink:
		return PinkColor
	case Yellow:
		return YellowColor
	}
	return WhiteColor
}

func NewBlockModel(b *Block) *BlockModel {
	model := &BlockModel{block: b}

	c := getBlockColor(b)
	vs := []Vertex{
		NewVertex(0, 0, 0, c),
		NewVertex(0, BlockSize, 0, c),
		NewVertex(BlockSize, 0, 0, c),
		NewVertex(BlockSize, BlockSize, 0, c),
	}
	sub:=&ModelBase{}
	model.Add(sub)
	sub.Init(gl.TRIANGLE_STRIP, vs, "shaders/basic.vert", "shaders/basic.frag")

	return model
}

type SwitchModel struct {
	ModelBase
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
	model.Init(gl.TRIANGLE_FAN, vs, "shaders/basic.vert", "shaders/basic.frag")

	v := SwitchSize / 2
	model.modelView = mathgl.Ortho2D(0, WindowWidth, WindowHeight, 0).Mul4(mathgl.Translate3D(float32(sw.X+v), float32(sw.Y+v), 0))
	return model
}

var (
	topLeftModelView     = mathgl.Translate3D(-BlockSize, -BlockSize, 0)
	topRightModelView    = mathgl.Translate3D(0, -BlockSize, 0)
	bottomRightModelView = mathgl.Ident4f()
	bottomLeftModelView  = mathgl.Translate3D(-BlockSize, 0, 0)
)

// TODO the switch number
func (t *SwitchModel) Draw() {
	modelViewBackup := t.modelView
	s := t.sw
	if s.rotate != 0 {
		t.modelView = t.modelView.Mul4(mathgl.HomogRotate3D(t.sw.rotate, [3]float32{0, 0, 1}))
	}
	scale := mathgl.Scale3D(s.scale, s.scale, 0)

	// Draw the associated blocks
	// top left block
	t.drawBlock(g.level.blocks[s.line][s.col], scale.Mul4(topLeftModelView))
	// top right block
	t.drawBlock(g.level.blocks[s.line][s.col+1], scale.Mul4(topRightModelView))
	// bottom right block
	t.drawBlock(g.level.blocks[s.line+1][s.col+1], scale.Mul4(bottomRightModelView))
	// bottom left block
	t.drawBlock(g.level.blocks[s.line+1][s.col], scale.Mul4(bottomLeftModelView))

	t.ModelBase.Draw()

	t.modelView = modelViewBackup
}

func (t *SwitchModel) drawBlock(b *Block, modelView mathgl.Mat4f) {
	if !b.Rendered {
		b.Rendered = true
		bm := g.world.blocks[b]
		bm.modelView = t.modelView.Mul4(modelView)
		bm.Draw()
	}
}

type Background struct {
	ModelBase
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
	model.Init(gl.TRIANGLES, vs, "shaders/basic.vert", "shaders/basic.frag")
	return model
}

func (t *Background) Draw() {
	if t.angle > math.Pi {
		t.angle = t.angle - math.Pi
	} else {
		t.angle += 0.03
	}
	modelViewBackup := t.modelView
	t.modelView = t.modelView.Mul4(mathgl.HomogRotate3D(-t.angle, [3]float32{0, 0, 1}))

	t.ModelBase.Draw()

	t.modelView = modelViewBackup
}
