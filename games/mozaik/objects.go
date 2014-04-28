package main

import (
	"math"
)

func NewBackground() *Model {
	vs := []Vertex{}

	for i := float64(0); i <= BgSegments; i++ {
		if math.Mod(i, 2) == 0 {
			vs = append(vs, NewVertex(0, 0, 0, BgColor))
		}
		a := 2 * math.Pi * i / BgSegments
		vs = append(vs, NewVertex(float32(math.Sin(a)*windowRadius), float32(math.Cos(a)*windowRadius), 0, BgColor))
	}
	m := NewModel(vs, "shaders/rotate.vert", "shaders/colorv.frag")
	return m
}

func NewCube2() *Model {
	return NewModel(readVertexFile("data/cube"), "shaders/rorateCube.vert", "shaders/cube.frag")
}

func NewTriangle() *Model {
	return NewModel([]Vertex{
		NewVertex(0.0, 0.5, 0.0, RedColor),
		NewVertex(0.5, -0.366, 0.0, GreenColor),
		NewVertex(-0.5, -0.366, 0.0, BlueColor),
	}, "shaders/rotateOffset.vert", "shaders/offset.frag")
}

func NewCube() *Model {
	return NewModel([]Vertex{
		NewVertex(0.25, 0.25, 0.75, RedColor),
		NewVertex(0.25, -0.25, 0.75, RedColor),
		NewVertex(-0.25, 0.25, 0.75, RedColor),

		NewVertex(0.25, -0.25, 0.75, RedColor),
		NewVertex(-0.25, -0.25, 0.75, RedColor),
		NewVertex(-0.25, 0.25, 0.75, RedColor),

		NewVertex(0.25, 0.25, -0.75, GreenColor),
		NewVertex(-0.25, 0.25, -0.75, GreenColor),
		NewVertex(0.25, -0.25, -0.75, GreenColor),

		NewVertex(0.25, -0.25, -0.75, GreenColor),
		NewVertex(-0.25, 0.25, -0.75, GreenColor),
		NewVertex(-0.25, -0.25, -0.75, GreenColor),

		NewVertex(-0.25, 0.25, 0.75, BlueColor),
		NewVertex(-0.25, -0.25, 0.75, BlueColor),
		NewVertex(-0.25, -0.25, -0.75, BlueColor),

		NewVertex(-0.25, 0.25, 0.75, BlueColor),
		NewVertex(-0.25, -0.25, -0.75, BlueColor),
		NewVertex(-0.25, 0.25, -0.75, BlueColor),

		NewVertex(0.25, 0.25, 0.75, RedColor),
		NewVertex(0.25, -0.25, -0.75, RedColor),
		NewVertex(0.25, -0.25, 0.75, RedColor),

		NewVertex(0.25, 0.25, 0.75, RedColor),
		NewVertex(0.25, 0.25, -0.75, RedColor),
		NewVertex(0.25, -0.25, -0.75, RedColor),

		NewVertex(0.25, 0.25, -0.75, RedColor),
		NewVertex(0.25, 0.25, 0.75, RedColor),
		NewVertex(-0.25, 0.25, 0.75, RedColor),

		NewVertex(0.25, 0.25, -0.75, RedColor),
		NewVertex(-0.25, 0.25, 0.75, RedColor),
		NewVertex(-0.25, 0.25, -0.75, RedColor),

		NewVertex(0.25, -0.25, -0.75, RedColor),
		NewVertex(-0.25, -0.25, 0.75, RedColor),
		NewVertex(0.25, -0.25, 0.75, RedColor),

		NewVertex(0.25, -0.25, -0.75, RedColor),
		NewVertex(-0.25, -0.25, -0.75, RedColor),
		NewVertex(-0.25, -0.25, 0.75, RedColor),
	}, "shaders/cube.vert", "shaders/cube.frag")
}
