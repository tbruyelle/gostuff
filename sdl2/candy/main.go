package main

import (
	"fmt"
	"git.tideland.biz/goas/loop"
	"github.com/jackyb/go-sdl2/sdl"
	"os"
	"runtime"
	"time"
)

const FRAME_RATE = time.Second / 30

var (
	window          *sdl.Window
	tileset         *sdl.Texture
	tilesetSelected *sdl.Texture
)

func main() {
	_ = fmt.Sprint()
	runtime.LockOSThread()

	window = sdl.CreateWindow("Candy Crush Saga", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, WindowWidth, WindowHeight, sdl.WINDOW_SHOWN)
	if window == nil {
		fmt.Fprintf(os.Stderr, "failed to create window %s\n", sdl.GetError())
		os.Exit(1)
	}
	defer window.Destroy()

	renderer := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if renderer == nil {
		fmt.Fprintf(os.Stderr, "Failed to create renderer %s\n", sdl.GetError())
		os.Exit(1)
	}
	defer renderer.Destroy()
	renderer.SetDrawColor(255, 255, 255, 255)
	renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND)

	tilesetFile := os.Getenv("GOPATH") + "/src/github.com/tbruyelle/gostuff/sdl2/candy/assets/tileset.bmp"
	tilesetSurface := sdl.LoadBMP(tilesetFile)
	if tilesetSurface == nil {
		fmt.Fprintf(os.Stderr, "Failed to load bitmap %s", tilesetFile)
		os.Exit(1)
	}
	tileset = renderer.CreateTextureFromSurface(tilesetSurface)
	tilesetSurface.SetAlphaMod(190)
	tilesetSelected = renderer.CreateTextureFromSurface(tilesetSurface)

	game := NewGame()
	defer game.Destroy()

	// render loop
	loop.GoRecoverable(
		func(loop loop.Loop) error {
			game.Start()

			ticker := time.NewTicker(FRAME_RATE)

			for {
				select {
				case <-ticker.C:
					game.Tick()
					renderThings(renderer, game)

				case <-loop.ShallStop():
					ticker.Stop()
					game.Stop()
					return nil

				}
			}
		},
		func(rs loop.Recoverings) (loop.Recoverings, error) {
			for _, r := range rs {
				fmt.Printf("%s\n%s", r.Reason)
			}
			return rs, fmt.Errorf("Unrecoverable loop\n")
		},
	)

	// event loop
	eventChan := make(chan sdl.Event)
	loop.GoRecoverable(
		func(loop loop.Loop) error {
			var evt sdl.Event
			for {
				evt = <-eventChan
				switch evt := evt.(type) {

				case *sdl.QuitEvent:
					//return nil
					break
				case *sdl.KeyDownEvent:
					switch evt.Keysym.Sym {
					case sdl.K_ESCAPE:
						//return nil
						break
					case sdl.K_r:
						game.Reset()
					case sdl.K_k:
						game.ToggleKeepUnmatchingTranslation()
					}
				case *sdl.MouseButtonEvent:
					if evt.State != 0 {
						game.Click(int(evt.X), int(evt.Y))
					}
				}

			}
		},
		func(rs loop.Recoverings) (loop.Recoverings, error) {
			for _, r := range rs {
				fmt.Printf("%s\n%s", r.Reason)
			}
			return rs, fmt.Errorf("Unrecoverable loop\n")
		},
	)

	for {
		eventChan <- sdl.WaitEvent()
	}

}

func renderThings(renderer *sdl.Renderer, game *Game) {
	fmt.Println("rendering")
	renderer.Clear()
	fmt.Println("rendering2")
	// show dashboard
	renderer.SetDrawColor(50, 50, 50, 200)
	dashboard := sdl.Rect{0, 0, DashboardWidth, WindowHeight}
	renderer.FillRect(&dashboard)

	// show candys
	for _, c := range game.candys {
		showCandy(renderer, c, game)
	}
	renderer.SetDrawColor(255, 255, 255, 255)
	fmt.Println("rendering3")
	renderer.Present()
	fmt.Println("rendering4")
}

var block = sdl.Rect{W: BlockSize, H: BlockSize}
var source = sdl.Rect{W: BlockSize, H: BlockSize}

// showCandy shows the candy according to a tileset.
func showCandy(renderer *sdl.Renderer, c *Candy, game *Game) {
	if c._type == EmptyCandy {
		return
	}
	//fmt.Printf("showCandy (%d,%d), %d\n", c.x, c.y, c._type)
	block.X = int32(c.x)
	block.Y = int32(c.y)
	alpha := uint8(255)
	switch c.sprite._type {
	case BlueSprite:
		source.X = BlockSize
		source.Y = 0
	case YellowSprite:
		source.X = BlockSize * 4
		source.Y = 0
	case GreenSprite:
		source.X = BlockSize * 3
		source.Y = 0
	case RedSprite:
		source.X = BlockSize * 5
		source.Y = 0
	case PinkSprite:
		source.X = BlockSize * 2
		source.Y = 0
	case OrangeSprite:
		source.X = 0
		source.Y = 0
	case BlueHStripesSprite:
		source.X = BlockSize
		source.Y = BlockSize
	case YellowHStripesSprite:
		source.X = BlockSize * 4
		source.Y = BlockSize
	case GreenHStripesSprite:
		source.X = BlockSize * 3
		source.Y = BlockSize
	case RedHStripesSprite:
		source.X = BlockSize * 5
		source.Y = BlockSize
	case PinkHStripesSprite:
		source.X = BlockSize * 2
		source.Y = BlockSize
	case OrangeHStripesSprite:
		source.X = 0
		source.Y = BlockSize
	case BlueVStripesSprite:
		source.X = BlockSize
		source.Y = BlockSize * 2
	case YellowVStripesSprite:
		source.X = BlockSize * 4
		source.Y = BlockSize * 2
	case GreenVStripesSprite:
		source.X = BlockSize * 3
		source.Y = BlockSize * 2
	case RedVStripesSprite:
		source.X = BlockSize * 5
		source.Y = BlockSize * 2
	case PinkVStripesSprite:
		source.X = BlockSize * 2
		source.Y = BlockSize * 2
	case OrangeVStripesSprite:
		source.X = 0
		source.Y = BlockSize * 2
	case BluePackedSprite:
		source.X = BlockSize
		source.Y = BlockSize * 3
	case YellowPackedSprite:
		source.X = BlockSize * 4
		source.Y = BlockSize * 3
	case GreenPackedSprite:
		source.X = BlockSize * 3
		source.Y = BlockSize * 3
	case RedPackedSprite:
		source.X = BlockSize * 5
		source.Y = BlockSize * 3
	case PinkPackedSprite:
		source.X = BlockSize * 2
		source.Y = BlockSize * 3
	case OrangePackedSprite:
		source.X = 0
		source.Y = BlockSize * 3
	case BombSprite:
		source.X = 0
		source.Y = BlockSize * 4
	case DyingSprite:
		source.X = 0
		source.Y = BlockSize * 4
		alpha = alpha / uint8(c.sprite.frame+1)

	}
	tileset.SetAlphaMod(alpha)
	if c == game.selected {
		renderer.Copy(tilesetSelected, &source, &block)
	} else {
		renderer.Copy(tileset, &source, &block)
	}
}

// showCandySquare is a deprecated method which shows candys as
// simples colored squares.
func showCandySquare(renderer *sdl.Renderer, c *Candy, game *Game) {
	if c._type == EmptyCandy {
		return
	}
	//fmt.Printf("showCandy (%d,%d), %d\n", c.x, c.y, c._type)
	block.X = int32(c.x + 1)
	block.Y = int32(c.y + 1)
	alpha := uint8(255)
	if c == game.selected {
		alpha = 150
	}
	switch c._type {
	case BlueCandy:
		renderer.SetDrawColor(153, 50, 204, alpha)
	case YellowCandy:
		renderer.SetDrawColor(255, 215, 0, alpha)
	case GreenCandy:
		renderer.SetDrawColor(60, 179, 113, alpha)
	case RedCandy:
		renderer.SetDrawColor(220, 20, 60, alpha)
	case PinkCandy:
		renderer.SetDrawColor(255, 192, 203, alpha)

	}
	renderer.FillRect(&block)
	renderer.SetDrawColor(255, 255, 255, 255)
}
