package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

type AnalogClock struct {
	renderer *sdl.Renderer
	rect     sdl.Rect
	texFace  *sdl.Texture
	fg, bg   sdl.Color
}

func NewAnalogClock(renderer *sdl.Renderer, rect sdl.Rect, fg, bg sdl.Color) *AnalogClock {
	texFace := NewClockFace(renderer, rect, fg, bg)
	return &AnalogClock{
		renderer: renderer,
		rect:     rect,
		texFace:  texFace,
		fg:       fg,
		bg:       bg,
	}
}

func (s *AnalogClock) Render(renderer *sdl.Renderer) {
	if err := renderer.Copy(s.texFace, nil, &s.rect); err != nil {
		panic(err)
	}
}

func (s *AnalogClock) Update()         {}
func (s *AnalogClock) Event(sdl.Event) {}
func (s *AnalogClock) Destroy()        { s.texFace.Destroy() }
func (s *AnalogClock) String()         { fmt.Sprintln("AnalogClock:%v %v %v", s.fg, s.bg, s.rect) }

func NewClockFace(renderer *sdl.Renderer, rect sdl.Rect, fg, bg sdl.Color) (texClockFace *sdl.Texture) {
	var err error
	if texClockFace, err = renderer.CreateTexture(sdl.PIXELFORMAT_RGBA8888, sdl.TEXTUREACCESS_TARGET, rect.W, rect.H); err != nil {
		panic(err)
	}
	center := sdl.Point{rect.W / 2, rect.H / 2}
	margin := rect.H / 90
	renderer.SetRenderTarget(texClockFace)
	texClockFace.SetBlendMode(sdl.BLENDMODE_BLEND)
	setColor(renderer, fg)
	renderer.Clear()
	FillCircle(renderer, center.X, center.Y, rect.H/2, bg)
	renderer.DrawLine(0, 0, rect.W, rect.H)
	renderer.DrawLine(0, rect.H, rect.W, 0)
	renderer.DrawLine(rect.W/2, 0, rect.W/2, 0)
	renderer.DrawLine(0, rect.H/2, rect.W, rect.H/2)
	var x, y int32
	for y = 0; y < rect.H; y += 5 {
		for x = 0; x < rect.W; x += 5 {
			renderer.DrawLine(x, 0, x, rect.H)
			renderer.DrawLine(0, y, rect.W, y)
		}
	}

	for i := 0; i < 60; i++ {
		var (
			tip    sdl.Point
			radius int32
		)
		if i%5 == 0 {
			radius = margin * 2
		} else {
			radius = margin
		}
		tip = getTip(center, float64(i)/60.0, float64(center.Y-margin*3), 0, 0)
		FillCircle(renderer, tip.X, tip.Y, radius, fg)
	}
	renderer.SetRenderTarget(nil)
	return texClockFace
}

func FillCircle(renderer *sdl.Renderer, x0, y0, radius int32, color sdl.Color) {
	setColor(renderer, color)
	for y := -radius; y <= radius; y++ {
		for x := -radius; x <= radius; x++ {
			if x*x+y*y <= radius*radius {
				renderer.DrawPoint(x0+x, y0+y)
			}
		}
	}
}
