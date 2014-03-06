package main

import (
	"fmt"
	"github.com/andrebq/gas"
	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
	"github.com/go-gl/gltext"
	"math"
	"os"
	"runtime"
	"time"
)

const FRAME_RATE = time.Second / 40

// Arrange that main.main runs on main thread.
func init() {
	runtime.LockOSThread()
}

// Queue of work to run in main thread.
var mainfunc = make(chan func())

// Run all the functions that need to run in the main thread.
func Main() {
	var f func()
	for f = range mainfunc {
		f()
	}
}

// Put the function f on the main thread function queue.
func do(f func()) {
	done := make(chan bool, 1)
	mainfunc <- func() {
		f()
		done <- true
	}
	<-done
}

func errorCallback(err glfw.ErrorCode, desc string) {
	fmt.Printf("%v: %v\n", err, desc)
}

var (
	window  *glfw.Window
	err     error
	g       *Game
	fonts   []*gltext.Font
	fontInd int
)

func main() {
	_ = fmt.Sprint()

	if !glfw.Init() {
		panic("Can't init glfw!")
	}
	defer glfw.Terminate()

	window, err = glfw.CreateWindow(WindowWidth, WindowHeight, "Mozaik", nil, nil)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	// Ensure thread context
	window.MakeContextCurrent()

	// TODO WHat fort ?
	//glfw.SwapInterval(1)

	window.SetKeyCallback(keyCb)
	window.SetMouseButtonCallback(mouseCb)

	gl.Init()
	gl.ClearColor(0.9, 0.9, 0.9, 0.0)
	// useless in 2D
	gl.Disable(gl.DEPTH_TEST)

	for i := int32(64); i < 72; i++ {
		font := loadFonts(i)
		defer font.Release()
		fonts = append(fonts, font)
	}

	// Use window coordinates
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Ortho(0, WindowWidth, WindowHeight, 0, 0, 1)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()

	g = NewGame()

	g.Start()
	go eventLoop(g)
	go renderLoop(g)
	Main()
	g.Stop()
}

// Load the font
func loadFonts(scale int32) *gltext.Font {
	file, err := gas.Abs("code.google.com/p/freetype-go/testdata/luxisr.ttf")
	if err != nil {
		panic(err)
	}
	fd, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer fd.Close()

	font, err := gltext.LoadTruetype(fd, scale, 32, 127, gltext.LeftToRight)
	if err != nil {
		panic(err)
	}
	return font
}

func keyCb(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Press {
		switch key {
		case glfw.KeyEscape:
			close(mainfunc)
		case glfw.KeyR:
			g.Reset()
		case glfw.KeyW:
			g.Warp()
		case glfw.KeyC:
			g.Cancel()
		case glfw.KeySpace:
			g.Continue()
		}
	}
}

func mouseCb(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
	if action == glfw.Press {
		switch button {

		case glfw.MouseButtonLeft:
			x, y := window.GetCursorPosition()
			g.Click(int(x), int(y))
		}
	}
}

func eventLoop(g *Game) {
	defer close(mainfunc)
	for {
		do(func() {
			glfw.PollEvents()
		})
	}
}

func renderLoop(g *Game) {
	defer close(mainfunc)

	mainTicker := time.NewTicker(FRAME_RATE)

	for {
		select {
		case <-mainTicker.C:
			g.Update()
			do(func() {
				renderThings(g)
			})

		}
	}
}

func renderThings(g *Game) {
	gl.Clear(gl.COLOR_BUFFER_BIT)

	// reinit blocks as not renderered
	for i := 0; i < len(g.blocks); i++ {
		for j := 0; j < len(g.blocks[i]); j++ {
			if g.blocks[i][j] != nil {
				g.blocks[i][j].Rendered = false
			}
		}
	}

	// render first the rotating switch if there's one
	if g.rotating != nil {
		renderRotatingSwitch(g.rotating)
	}

	for i := 0; i < len(g.blocks); i++ {
		gl.LoadIdentity()
		gl.Translatef(float32(XMin), float32(YMin+BlockSize*i), 0)
		for j := 0; j < len(g.blocks[i]); j++ {
			if g.blocks[i][j] != nil && !g.blocks[i][j].Rendered {
				gl.Begin(gl.QUADS)
				setColor(g.blocks[i][j].Color)
				gl.Vertex2i(0, 0)
				gl.Vertex2i(BlockSize, 0)
				gl.Vertex2i(BlockSize, BlockSize)
				gl.Vertex2i(0, BlockSize)
				gl.End()
			}
			gl.Translatef(float32(BlockSize), 0, 0)
		}
	}

	// render the switches
	// TODO can we use z to make them upper?
	for _, s := range g.switches {
		renderSwitch(s)
	}

	renderDashboard()

	if g.Win() {
		fontInd++
		if fontInd >= len(fonts) {
			fontInd = 0
		}
		gl.LoadIdentity()
		gl.Color3ub(107, 142, 35)
		fonts[fontInd].Printf(float32(WindowWidth/2)-100, float32(WindowHeight/2), "GAGNE !!")
	}

	// TODO What for?
	//gl.Flush()
	window.SwapBuffers()
}

func renderRotatingSwitch(s *Switch) {
	gl.LoadIdentity()
	// TODO constant
	v := SwitchSize / 2

	gl.Translatef(float32(s.X+v), float32(s.Y+v), 0)
	gl.Rotatef(float32(s.rotate), 0, 0, 1)

	var b *Block
	// render block top left
	b = g.blocks[s.line][s.col]
	gl.Begin(gl.QUADS)
	setColor(b.Color)
	gl.Vertex2i(-BlockSize, -BlockSize)
	gl.Vertex2i(0, -BlockSize)
	gl.Vertex2i(0, 0)
	gl.Vertex2i(-BlockSize, 0)
	gl.End()
	b.Rendered = true

	// render block top right
	b = g.blocks[s.line][s.col+1]
	gl.Begin(gl.QUADS)
	setColor(b.Color)
	gl.Vertex2i(0, -BlockSize)
	gl.Vertex2i(BlockSize, -BlockSize)
	gl.Vertex2i(BlockSize, 0)
	gl.Vertex2i(0, 0)
	gl.End()
	b.Rendered = true

	// render block bottom right
	b = g.blocks[s.line+1][s.col+1]
	gl.Begin(gl.QUADS)
	setColor(b.Color)
	gl.Vertex2i(0, 0)
	gl.Vertex2i(BlockSize, 0)
	gl.Vertex2i(BlockSize, BlockSize)
	gl.Vertex2i(0, BlockSize)
	gl.End()
	b.Rendered = true

	// render block bottom left
	b = g.blocks[s.line+1][s.col]
	gl.Begin(gl.QUADS)
	setColor(b.Color)
	gl.Vertex2i(-BlockSize, 0)
	gl.Vertex2i(0, 0)
	gl.Vertex2i(0, BlockSize)
	gl.Vertex2i(-BlockSize, BlockSize)
	gl.End()
	b.Rendered = true
}

func renderSwitch(s *Switch) {
	gl.LoadIdentity()
	// TODO constant
	v := SwitchSize / 2

	gl.Translatef(float32(s.X+v), float32(s.Y+v), 0)
	// render the switch
	gl.Begin(gl.TRIANGLE_FAN)
	gl.Color3f(1.0, 1.0, 1.0)
	vv := float64(v)
	for a := float64(0); a < 360; a += 5 {
		gl.Vertex2d(math.Sin(a)*vv, math.Cos(a)*vv)
	}
	//gl.Vertex2i(-v, v)
	//gl.Vertex2i(v, v)
	//gl.Vertex2i(v, -v)
	//gl.Vertex2i(-v, -v)
	gl.End()
}

func setColor(color ColorDef) {
	switch color {
	case Red:
		gl.Color3ub(255, 51, 51)
	case Yellow:
		gl.Color3ub(255, 215, 0)
	case Blue:
		gl.Color3ub(100, 149, 237)
	case Green:
		gl.Color3ub(102, 204, 0)
	case Pink:
		gl.Color3ub(255, 104, 255)

	}
}

func renderDashboard() {
	gl.LoadIdentity()
	setColor(Blue)
	fonts[0].Printf(float32(XMin), float32(WindowHeight-DashboardHeight),
		"Level %d", g.currentLevel)
	gl.Translatef(float32(XMin)+300, float32(WindowHeight-DashboardHeight), 0)
	line := 0
	for _, c := range g.winSignature {
		if c == '\n' {
			line++
			gl.LoadIdentity()
			gl.Translatef(float32(XMin)+300, float32(WindowHeight-DashboardHeight+SignatureBlockSize*line), 0)
			continue
		}
		if c != '-' {
			gl.Begin(gl.QUADS)
			setColor(atoc(string(c)))
			gl.Vertex2i(0, 0)
			gl.Vertex2i(SignatureBlockSize, 0)
			gl.Vertex2i(SignatureBlockSize, SignatureBlockSize)
			gl.Vertex2i(0, SignatureBlockSize)
			gl.End()
		}
		gl.Translated(SignatureBlockSize, 0, 0)
	}
}
