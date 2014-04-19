package gfx

import (
	"log"

	"github.com/jackyb/go-sdl2/sdl"
	ttf "github.com/jackyb/go-sdl2/sdl_ttf"
)

func CreateText(renderer *sdl.Renderer, font *ttf.Font, message string, color sdl.Color) *sdl.Texture {
	surf := font.RenderText_Blended(message, color)
	if surf == nil {
		log.Println("Error when rendering text")
		return nil
	}
	defer surf.Free()

	return renderer.CreateTextureFromSurface(surf)
}
