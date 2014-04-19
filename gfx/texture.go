package gfx

import "github.com/jackyb/go-sdl2/sdl"

func RenderTexture(renderer *sdl.Renderer, texture *sdl.Texture, x, y int) {
	var w, h int
	sdl.QueryTexture(texture, nil, nil, &w, &h)
	dst := sdl.Rect{int32(x), int32(y), int32(w), int32(h)}
	renderer.Copy(texture, nil, &dst)
}
