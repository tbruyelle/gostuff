package main

func NewCube2() *Model {
	return NewModel(readVertexFile("data/cube"), "shaders/rorateCube.vert", "shaders/cube.frag")
}
func NewTriangle() *Model {
	return NewModel([]Vertex{
		NewVertex(0.0, 0.5, 0.0, Red),
		NewVertex(0.5, -0.366, 0.0, Green),
		NewVertex(-0.5, -0.366, 0.0, Blue),
	}, "shaders/rotateOffset.vert", "shaders/offset.frag")
}

func NewCube() *Model {
	return NewModel([]Vertex{
		NewVertex(0.25, 0.25, 0.75, Red),
		NewVertex(0.25, -0.25, 0.75, Red),
		NewVertex(-0.25, 0.25, 0.75, Red),

		NewVertex(0.25, -0.25, 0.75, Red),
		NewVertex(-0.25, -0.25, 0.75, Red),
		NewVertex(-0.25, 0.25, 0.75, Red),

		NewVertex(0.25, 0.25, -0.75, Green),
		NewVertex(-0.25, 0.25, -0.75, Green),
		NewVertex(0.25, -0.25, -0.75, Green),

		NewVertex(0.25, -0.25, -0.75, Green),
		NewVertex(-0.25, 0.25, -0.75, Green),
		NewVertex(-0.25, -0.25, -0.75, Green),

		NewVertex(-0.25, 0.25, 0.75, Blue),
		NewVertex(-0.25, -0.25, 0.75, Blue),
		NewVertex(-0.25, -0.25, -0.75, Blue),

		NewVertex(-0.25, 0.25, 0.75, Blue),
		NewVertex(-0.25, -0.25, -0.75, Blue),
		NewVertex(-0.25, 0.25, -0.75, Blue),

		NewVertex(0.25, 0.25, 0.75, Red),
		NewVertex(0.25, -0.25, -0.75, Red),
		NewVertex(0.25, -0.25, 0.75, Red),

		NewVertex(0.25, 0.25, 0.75, Red),
		NewVertex(0.25, 0.25, -0.75, Red),
		NewVertex(0.25, -0.25, -0.75, Red),

		NewVertex(0.25, 0.25, -0.75, Red),
		NewVertex(0.25, 0.25, 0.75, Red),
		NewVertex(-0.25, 0.25, 0.75, Red),

		NewVertex(0.25, 0.25, -0.75, Red),
		NewVertex(-0.25, 0.25, 0.75, Red),
		NewVertex(-0.25, 0.25, -0.75, Red),

		NewVertex(0.25, -0.25, -0.75, Red),
		NewVertex(-0.25, -0.25, 0.75, Red),
		NewVertex(0.25, -0.25, 0.75, Red),

		NewVertex(0.25, -0.25, -0.75, Red),
		NewVertex(-0.25, -0.25, -0.75, Red),
		NewVertex(-0.25, -0.25, 0.75, Red),
	}, "shaders/cube.vert", "shaders/cube.frag")
}
