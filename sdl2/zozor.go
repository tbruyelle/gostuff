package main

// #include <SDL2/SDL_keyboard.h>
import "C"
import "fmt"
import "github.com/jackyb/go-sdl2/sdl"
import "os"

const imageName = "zozor.bmp"

func main() {
	window := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		800, 600, sdl.WINDOW_SHOWN)
	if window == nil {
		fmt.Fprintf(os.Stderr, "Failed to create window: %s\n", sdl.GetError())
		os.Exit(1)
	}
	defer window.Destroy()

	renderer := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if renderer == nil {
		fmt.Fprintf(os.Stderr, "Failed to create renderer: %s\n", sdl.GetError())
		os.Exit(2)
	}
	defer renderer.Destroy()
	renderer.SetDrawColor(255, 255, 255, 255)
	//surface := window.GetSurface()
	//if surface == nil {
	//	fmt.Fprintf(os.Stderr, "Failed to create surface: %s\n", sdl.GetError())
	//	os.Exit(1)
	//}

	var image *sdl.Surface = sdl.LoadBMP(imageName)
	if image == nil {
		fmt.Fprintf(os.Stderr, "Failed to load BMP: %s", sdl.GetError())
		os.Exit(3)
	}
	defer image.Free()
	texture := renderer.CreateTextureFromSurface(image)
	defer texture.Destroy()

	//src := sdl.Rect{10, 10, 64, 64}
	dst := sdl.Rect{100, 50, 64, 64}

	renderer.Clear()
	//renderer.Copy(texture, nil, &sdl.Rect{X: 20, Y: 20})
	renderer.Copy(texture, nil, &dst)
	renderer.Present()

	//surface.FillRect(nil, 0xffffffff)
	//surface.Blit(nil, image, &sdl.Rect{X: 200, Y: 200})
	//rect := sdl.Rect{0, 0, 200, 200}
	//window.UpdateSurface()
	running := true
	for running {
		//var event sdl.Event
		//sdl.WaitEvent(&event)
		event := sdl.PollEvent()
		switch t := event.(type) {
		case *sdl.QuitEvent:
			running = false
		case *sdl.KeyDownEvent:
			switch t.Keysym.Sym {
			case C.SDLK_UP:
				dst.Y--
			case C.SDLK_DOWN:
				dst.Y++
			case C.SDLK_LEFT:
				dst.X--
			case C.SDLK_RIGHT:
				dst.X++
			case C.SDLK_ESCAPE:
				running = false
			}
		case *sdl.MouseButtonEvent:
			dst.X = t.X
			dst.Y = t.Y
		}
		renderer.Clear()
		//renderer.Copy(texture, nil, &sdl.Rect{X: 20, Y: 20})
		renderer.Copy(texture, nil, &dst)
		renderer.Present()

	}
}
