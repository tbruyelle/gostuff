package main

import (
	"fmt"
	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
	"io/ioutil"
	"math"
	"strconv"
	"strings"
	"unsafe"
)

type Vertex struct {
	Coords Coords
	Color  Color
}

func NewVertex(X, Y, Z float32, color Color) Vertex {
	return Vertex{Coords: Coords{X, Y, Z, 1.0}, Color: color}
}

var (
	RedColor   = Color{1.0, 0.0, 0.0, 1.0}
	GreenColor = Color{0.0, 1.0, 0.0, 1.0}
	BlueColor  = Color{0.0, 0.0, 1.0, 1.0}
	BgColor    = Color{1.0, 0.85, 0.23, 1.0}
)

type Coords struct{ X, Y, Z, W float32 }
type Color struct{ R, G, B, A float32 }

func Sequence(seqSize, ind int) int {
	r := ind / seqSize
	for r >= seqSize {
		r -= seqSize
	}
	return r

}

func readVertexFile(file string) []Vertex {
	vertexes := make([]Vertex, 0)
	b, err := ioutil.ReadFile(file + ".coords")
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(b), "\n")
	for _, line := range lines {
		coords := strings.Split(line, ",")
		if len(coords) >= 4 {
			v := Vertex{}
			v.Coords.X = atof(coords[0])
			v.Coords.Y = atof(coords[1])
			v.Coords.Z = atof(coords[2])
			v.Coords.W = atof(coords[3])
			vertexes = append(vertexes, v)
		}
	}
	b, err = ioutil.ReadFile(file + ".colors")
	if err != nil {
		panic(err)
	}
	vind := 0
	lines = strings.Split(string(b), "\n")
	for _, line := range lines {
		colors := strings.Split(line, ",")
		if len(colors) >= 4 {
			v := &vertexes[vind]
			v.Color.R = atof(colors[0])
			v.Color.G = atof(colors[1])
			v.Color.B = atof(colors[2])
			v.Color.A = atof(colors[3])
			vind++
		}
	}
	return vertexes
}

func atof(s string) float32 {
	f, err := strconv.ParseFloat(strings.TrimSpace(s), 10)
	if err != nil {
		panic(err)
	}
	return float32(f)
}

var (
	sizeFloat  = int(unsafe.Sizeof(float32(0)))
	sizeCoords = sizeFloat * 4
	sizeVertex = int(unsafe.Sizeof(Vertex{}))
)

func showVersion() {
	//maj, min, v := glfw.GetVersion()
	//fmt.Println("version=", maj, min, v)
	fmt.Println("version=", glfw.GetVersionString())
}
func NewProgram(shaders ...gl.Shader) gl.Program {
	prg := gl.CreateProgram()
	for _, shader := range shaders {
		prg.AttachShader(shader)
	}
	prg.Link()
	if prg.Get(gl.LINK_STATUS) != gl.TRUE {
		panic("linker error: " + prg.GetInfoLog())
	}
	prg.Validate()
	for _, shader := range shaders {
		prg.DetachShader(shader)
		shader.Delete()
	}
	return prg
}

func loadShader(type_ gl.GLenum, file string) gl.Shader {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	source := string(b)
	shader := gl.CreateShader(type_)
	shader.Source(source)
	shader.Compile()
	if shader.Get(gl.COMPILE_STATUS) != gl.TRUE {
		panic("fragment error for source " + source + "\n" + shader.GetInfoLog())
	}
	return shader
}

func rotateOffsets(x, y float32) (float32, float32) {
	const fLoopDuration = float64(5)
	const fScale = 3.14159 * 2 / fLoopDuration

	fElapsedTime := glfw.GetTime()
	fCurrTimeThroughLoop := math.Mod(fElapsedTime, fLoopDuration)
	return float32(math.Cos(fCurrTimeThroughLoop*fScale) * 0.5),
		float32(math.Sin(fCurrTimeThroughLoop*fScale) * 0.5)
}
