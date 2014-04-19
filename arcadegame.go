package main

import (
	"fmt"
	"log"

	"github.com/jackyb/go-sdl2/sdl"
	"github.com/jackyb/go-sdl2/sdl_image"
	ttf "github.com/jackyb/go-sdl2/sdl_ttf"
	"github.com/jostly/arcadegame/game"
)

const (
	SCREEN_WIDTH  = 640
	SCREEN_HEIGHT = 480
)

var (
	x = float64(SCREEN_WIDTH / 4)
	y = float64(SCREEN_HEIGHT / 2)
)

func error(message string) {
	log.Printf(message+": %v\n", sdl.GetError())
}

func main() {
	log.Println("SDL2 Tutorial #1")

	if sdl.Init(sdl.INIT_EVERYTHING) != 0 {
		error("Error initializing SDL")
		return
	}

	defer sdl.Quit()

	ttf.Init()

	window := sdl.CreateWindow("Hello World!", 0, 0, SCREEN_WIDTH, SCREEN_HEIGHT,
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

	font, error := ttf.OpenFont("fonts/ComicNeue-Regular.ttf", 18)

	if error != nil {
		log.Printf("Error when loading font: %v", error)
		return
	}

	defer font.Close()

	game.RenderCallback = func(r *sdl.Renderer) {

		drawShip(r)
		drawMissiles(r)

	}

	game.UpdateCallback = func(delta float64) {
		moveSpeed := delta * 100

		keystate := sdl.GetKeyboardState()

		updateMissiles(moveSpeed * 3)

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
		if keystate[sdl.SCANCODE_RSHIFT] != 0 {
			currentTick := sdl.GetTicks()
			if currentTick > lastFire+300 {
				lastFire = currentTick
				missiles = append(missiles, FloatPoint{x + 32, y})
			}
		}
	}

	game.StatusCallback = func() string {
		return fmt.Sprintf("Active missiles: %d", len(missiles))
	}

	game.MainLoop(renderer, font)

}

type FloatPoint struct {
	X, Y float64
}

func updateMissiles(moveSpeed float64) {
	newMissiles := make([]FloatPoint, 0, len(missiles))
	for _, p := range missiles {
		p.X += moveSpeed
		if p.X <= SCREEN_WIDTH {
			newMissiles = append(newMissiles, p)
		}
	}
	missiles = newMissiles
}

var lastFire = sdl.GetTicks()

var shipPoints = [...]sdl.Point{sdl.Point{30, 0}, sdl.Point{-20, -15}, sdl.Point{-20, 15}}

func drawShip(r *sdl.Renderer) {
	ship := make([]sdl.Point, 4)
	for i := 0; i < 4; i++ {
		p := shipPoints[i%3]
		ship[i] = sdl.Point{p.X + int32(x), p.Y + int32(y)}
	}

	r.SetDrawColor(255, 255, 255, 255)
	err := r.DrawLines(ship)
	if err != 0 {
		error("Draw error")
	}
}

var missiles = []FloatPoint{}

func drawMissiles(r *sdl.Renderer) {
	if len(missiles) > 0 {
		for _, p := range missiles {
			r.DrawPoint(int(p.X), int(p.Y))
		}
	}

}
