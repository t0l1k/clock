package main

import (
	"github.com/t0l1k/sdl2/sdl2/ui"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type Frame struct {
	renderer *sdl.Renderer
	texture  *sdl.Texture
	rect     sdl.Rect
	fg, bg   sdl.Color
	font     *ttf.Font
	show     bool
}

func NewFrame(arr []string, point sdl.Point, fg, bg sdl.Color, renderer *sdl.Renderer, font *ttf.Font) *Frame {
	texture := newFrameTexture(arr, fg, bg, renderer, font)
	_, _, w, h, _ := texture.Query()
	return &Frame{
		rect:     sdl.Rect{point.X, point.Y, w, h},
		fg:       fg,
		bg:       bg,
		renderer: renderer,
		font:     font,
		texture:  texture,
		show:     true,
	}
}

func newFrameTexture(arr []string, fg, bg sdl.Color, renderer *sdl.Renderer, font *ttf.Font) *sdl.Texture {
	var w, h, max int
	if len(arr) > 0 {
		for i, str := range arr {
			if len(str) > max {
				w, h, _ = font.SizeUTF8(arr[i])
				max = len(str)
			}
		}
	}
	texture, err := renderer.CreateTexture(sdl.PIXELFORMAT_ARGB8888, sdl.TEXTUREACCESS_TARGET, int32(w), int32(h*len(arr)))
	texture.SetBlendMode(sdl.BLENDMODE_BLEND)
	if err != nil {
		panic(err)
	}
	renderer.SetRenderTarget(texture)
	renderer.SetDrawColor(bg.R, bg.G, bg.B, bg.A)
	renderer.Clear()
	renderer.SetDrawColor(fg.R, fg.G, fg.B, fg.A)
	for i := range arr {
		renderer.SetDrawColor(255, 0, 0, 255)
		label := ui.NewLabel(arr[len(arr)-i-1], sdl.Point{0, int32(i * h)}, fg, renderer, font)
		label.Render(renderer)
		label.Destroy()
	}
	renderer.SetRenderTarget(nil)

	return texture
}

func (s *Frame) SetText(arr []string) {
	if s.show {
		s.Destroy()
		s.texture = newFrameTexture(arr, s.fg, s.bg, s.renderer, s.font)
		_, _, w, h, _ := s.texture.Query()
		s.rect = sdl.Rect{s.rect.X, s.rect.Y, w, h}
	}
}

func (s *Frame) Render(renderer *sdl.Renderer) {
	if s.show {
		if err := renderer.Copy(s.texture, nil, &s.rect); err != nil {
			panic(err)
		}
	}
}

func (s *Frame) Update()         {}
func (s *Frame) Event(sdl.Event) {}
func (s *Frame) Destroy()        { s.texture.Destroy() }
