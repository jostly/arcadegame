package main

import (
	"log"

	"github.com/jackyb/go-sdl2/sdl"
	"github.com/jackyb/go-sdl2/sdl_image"
	ttf "github.com/jackyb/go-sdl2/sdl_ttf"
	"github.com/jostly/arcadegame/game"
	"github.com/jostly/arcadegame/gfx"
)

const (
	SCREEN_WIDTH  = 320
	SCREEN_HEIGHT = 320
)

var (
	x = 0.0
	y = 0.0
)

func error(message string) {
	log.Printf(message+": %v\n", sdl.GetError())
}

func renderTexture(renderer *sdl.Renderer, texture *sdl.Texture, x, y int) {
	var w, h int
	sdl.QueryTexture(texture, nil, nil, &w, &h)
	dst := sdl.Rect{int32(x), int32(y), int32(w), int32(h)}
	renderer.Copy(texture, nil, &dst)
}

func createText(renderer *sdl.Renderer, font *ttf.Font, message string, color sdl.Color) *sdl.Texture {
	surf := font.RenderText_Blended(message, color)
	if surf == nil {
		error("Error when rendering text")
		return nil
	}
	defer surf.Free()

	return renderer.CreateTextureFromSurface(surf)
}

var running = true

func main() {
	log.Println("SDL2 Tutorial #1")

	if sdl.Init(sdl.INIT_EVERYTHING) != 0 {
		error("Error initializing SDL")
		return
	}

	defer sdl.Quit()

	ttf.Init()

	window := sdl.CreateWindow("Hello World!", 100, 100, SCREEN_WIDTH, SCREEN_HEIGHT,
		sdl.WINDOW_SHOWN)

	if window == nil {
		error("Error opening window")
		return
	}

	defer window.Destroy()

	renderer := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)

	if renderer == nil {
		error("Error creating renderer")
		return
	}

	defer renderer.Destroy()

	texture := img.LoadTexture(renderer, "images/exampletexture_2.png")

	if texture == nil {
		error("Error creating texture")
		return
	}

	defer texture.Destroy()

	//font, error := ttf.OpenFont("/Users/johan/Library/Fonts/Anonymous Pro Minus B.ttf", 20)
	font, error := ttf.OpenFont("fonts/ComicNeue-Regular.ttf", 30)

	if error != nil {
		log.Printf("Error when loading font: %v", error)
		return
	}

	defer font.Close()

	game.RenderCallback = func(r *sdl.Renderer) {
		gfx.RenderTexture(r, texture, int(x), int(y))
	}

	game.MovementHandler = func(delta float64) {
		moveSpeed := delta * 100

		keystate := sdl.GetKeyboardState()

		if keystate[sdl.SCANCODE_W] != 0 {
			y -= moveSpeed
		}
		if keystate[sdl.SCANCODE_S] != 0 {
			y += moveSpeed
		}
		if keystate[sdl.SCANCODE_A] != 0 {
			x -= moveSpeed
		}
		if keystate[sdl.SCANCODE_D] != 0 {
			x += moveSpeed
		}
	}

	game.MainLoop(renderer, font)

}
