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
	RenderCallback func(*sdl.Renderer) = func(_ *sdl.Renderer) {}
	UpdateCallback func(float64)       = func(_ float64) {}
	StatusCallback func() string       = func() string { return "" }
)

func MainLoop(renderer *sdl.Renderer, font *ttf.Font) {
	lastTick := sdl.GetTicks()
	lastSecond := lastTick
	frames := 0
	fpsText := gfx.CreateText(renderer, font, "FPS: ... "+StatusCallback(), sdl.Color{255, 255, 255, 255})
	framesPerSecond := 0

	for running {

		tick := sdl.GetTicks()
		if tick < lastTick+10 {
			sdl.Delay(1)

			continue
		}

		delta := float64(tick-lastTick) / 1000.0
		lastTick = tick
		fpsText.Destroy()
		fpsText = gfx.CreateText(renderer, font, fmt.Sprintf("FPS: %3d ", framesPerSecond)+StatusCallback(), sdl.Color{255, 255, 255, 255})

		if tick > (lastSecond + 1000) {
			lastSecond = tick
			framesPerSecond = frames
			frames = 0
		} else {
			frames++
		}

		handleEvents()

		UpdateCallback(delta)

		renderer.SetDrawColor(0, 0, 0, 0)
		renderer.Clear()

		RenderCallback(renderer)

		gfx.RenderTexture(renderer, fpsText, 0, 0)

		renderer.Present()

	}

	fpsText.Destroy()
}

func handleEvents() {
	for {
		event := sdl.PollEvent()
		if event == nil {
			return
		}
		switch event := event.(type) {
		case *sdl.QuitEvent:
			log.Println("Quit!")
			running = false
		case *sdl.KeyDownEvent:
			handleKey(event.Keysym)
		}

	}
}

func handleKey(keysym sdl.Keysym) {
	switch keysym.Sym {
	case sdl.Keycode(sdl.K_ESCAPE):
		log.Println("KTHXBYE")
		running = false
	}
}
