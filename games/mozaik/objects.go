package main

import (
	"github.com/go-gl/gl"
	"github.com/remogatto/mathgl"
	"math"
)

type BlockModel struct {
	Model
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
	model := &BlockModel{}

	c := getBlockColor(b)
	vs := []Vertex{
		NewVertex(0, 0, 0, c),
		NewVertex(0, BlockSize, 0, c),
		NewVertex(BlockSize, 0, 0, c),
		NewVertex(BlockSize, BlockSize, 0, c),
	}
	model.Init(vs, "shaders/basic.vert", "shaders/basic.frag")
	return model
}

func (t *BlockModel) Draw() {
	t.prg.Use()

	t.vao.Bind()

	t.sendMVP()

	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, len(t.vertices))

	t.vao.Unbind()
	t.prg.Unuse()
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
	model.modelView = mathgl.Ortho2D(0, WindowWidth, WindowHeight, 0).Mul4(mathgl.Translate3D(float32(sw.X+v), float32(sw.Y+v), 0))

	return model
}

// TODO the switch number
func (t *SwitchModel) Draw() {

	modelViewBackup := t.modelView
	if t.sw.rotate != 0 {
		t.modelView = t.modelView.Mul4(mathgl.HomogRotate3D(t.sw.rotate, [3]float32{0, 0, 1}))
	}

	// Draw the blocks
	var b *Block
	s := t.sw
	bsf := float32(BlockSize - s.Z)
	// top left block
	b = t.lvl.blocks[s.line][s.col]
	if !b.Rendered {
		drawBlock(b, t.modelView.Mul4(mathgl.Translate3D(-bsf, -bsf, 0)))
	}
	// top right block
	b = t.lvl.blocks[s.line][s.col+1]
	if !b.Rendered {
		drawBlock(b, t.modelView.Mul4(mathgl.Translate3D(0, -bsf, 0)))
	}
	// bottom right block
	b = t.lvl.blocks[s.line+1][s.col+1]
	if !b.Rendered {
		drawBlock(b, t.modelView)
	}
	// bottom left block
	b = t.lvl.blocks[s.line+1][s.col]
	if !b.Rendered {
		drawBlock(b, t.modelView.Mul4(mathgl.Translate3D(-bsf, 0, 0)))
	}

	// Draw the switch
	t.prg.Use()

	t.vao.Bind()

	t.sendMVP()

	gl.DrawArrays(gl.TRIANGLE_FAN, 0, len(t.vertices))

	t.vao.Unbind()
	t.prg.Unuse()

	t.modelView = modelViewBackup
}

func drawBlock(b *Block, modelView mathgl.Mat4f) {
	bl := NewBlockModel(b)
	bl.modelView = modelView
	bl.Draw()
	b.Rendered = true
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
	if t.angle > math.Pi {
		t.angle = t.angle - math.Pi
	} else {
		t.angle += 0.03
	}
	modelViewBackup := t.modelView
	t.modelView = t.modelView.Mul4(mathgl.HomogRotate3D(-t.angle, [3]float32{0, 0, 1}))

	t.prg.Use()

	t.vao.Bind()

	t.sendMVP()


	gl.DrawArrays(gl.TRIANGLES, 0, len(t.vertices))

	t.vao.Unbind()
	t.prg.Unuse()

	t.modelView = modelViewBackup
}
