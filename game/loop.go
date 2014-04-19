package game

import (
	"fmt"
	"log"

	"github.com/jackyb/go-sdl2/sdl"
	ttf "github.com/jackyb/go-sdl2/sdl_ttf"
	"github.com/jostly/arcadegame/gfx"
)

var running = true

var (
	RenderCallback  func(*sdl.Renderer) = func(_ *sdl.Renderer) {}
	MovementHandler func(float64)       = func(_ float64) {}
)

func MainLoop(renderer *sdl.Renderer, font *ttf.Font) {
	lastTick := sdl.GetTicks()
	lastSecond := lastTick
	frames := 0
	fpsText := gfx.CreateText(renderer, font, "FPS: ...", sdl.Color{255, 255, 255, 255})

	for running {

		tick := sdl.GetTicks()

		delta := float64(tick-lastTick) / 1000.0
		lastTick = tick

		if tick > (lastSecond + 1000) {
			lastSecond = tick
			fpsText.Destroy()
			fpsText = gfx.CreateText(renderer, font, fmt.Sprintf("FPS: %d", frames), sdl.Color{255, 255, 255, 255})
			frames = 0
		} else {
			frames++
		}

		for handleEvent() {
		}

		MovementHandler(delta)

		renderer.Clear()

		RenderCallback(renderer)

		gfx.RenderTexture(renderer, fpsText, 0, 0)

		renderer.Present()
	}

	fpsText.Destroy()
}

func handleEvent() bool {
	event := sdl.PollEvent()
	if event == nil {
		return false
	}
	switch event := event.(type) {
	case *sdl.QuitEvent:
		log.Println("Quit!")
		running = false
	case *sdl.KeyDownEvent:
		handleKey(event.Keysym)
	}
	return true
}

func handleKey(keysym sdl.Keysym) {
	switch keysym.Sym {
	case sdl.Keycode(sdl.K_ESCAPE):
		log.Println("KTHXBYE")
		running = false
	}
}
