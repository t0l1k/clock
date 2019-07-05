package main

import (
	"fmt"

	"github.com/t0l1k/sdl2/sdl2/ui"
	"github.com/veandco/go-sdl2/sdl"
)

type ClockHand struct {
	renderer        *sdl.Renderer
	texture         *sdl.Texture
	rect, paintRect sdl.Rect
	handCenter      sdl.Point
	width, height   int32
	fg, bg          sdl.Color
	angle, fraction float64
}

func (s *ClockHand) GetFraction() float64 { return s.fraction }

func (s *ClockHand) Render(renderer *sdl.Renderer) {
	if err := renderer.CopyEx(s.texture, nil, &s.paintRect, s.angle, &s.handCenter, sdl.FLIP_NONE); err != nil {
		panic(err)
	}
}

func (s *ClockHand) Update(percent float64) {
	s.fraction = percent
	s.angle = getAngle(percent)
}
func (s *ClockHand) Event(sdl.Event) {}
func (s *ClockHand) Destroy()        { s.texture.Destroy() }
func (s *ClockHand) String()         { fmt.Sprintf("ClockHand:%v %v %v", s.fg, s.bg, s.rect) }

func NewSmallHand(renderer *sdl.Renderer, width, height int32, rect sdl.Rect, center sdl.Point, fg, bg sdl.Color) *ClockHand {
	texHand, err := renderer.CreateTexture(sdl.PIXELFORMAT_ARGB8888, sdl.TEXTUREACCESS_TARGET, rect.W, rect.H)
	if err != nil {
		panic(err)
	}
	renderer.SetRenderTarget(texHand)
	texHand.SetBlendMode(sdl.BLENDMODE_BLEND)
	ui.SetColor(renderer, bg)
	renderer.Clear()
	ui.SetColor(renderer, fg)
	renderer.FillRect(&sdl.Rect{0, 0, center.X - center.X/4, rect.H})
	renderer.FillRect(&sdl.Rect{0, rect.H / 3, rect.W, rect.H - rect.H/3*2})
	ui.FillCircle(renderer, center.X, center.Y, rect.H/5, bg)
	renderer.SetRenderTarget(nil)
	paintRect := sdl.Rect{rect.X + width/2 - center.X, rect.Y + height/2 - rect.H/2, rect.W, rect.H}

	return &ClockHand{
		renderer:   renderer,
		texture:    texHand,
		rect:       rect,
		paintRect:  paintRect,
		handCenter: center,
		fg:         fg,
		bg:         bg,
	}
}

func NewSmallHandRounded(renderer *sdl.Renderer, width, height int32, rect sdl.Rect, center sdl.Point, fg, bg sdl.Color) *ClockHand {
	texHand, err := renderer.CreateTexture(sdl.PIXELFORMAT_ARGB8888, sdl.TEXTUREACCESS_TARGET, rect.W, rect.H)
	if err != nil {
		panic(err)
	}
	renderer.SetRenderTarget(texHand)
	texHand.SetBlendMode(sdl.BLENDMODE_BLEND)
	ui.SetColor(renderer, bg)
	renderer.Clear()
	ui.SetColor(renderer, fg)
	ui.FillCircle(renderer, rect.H-rect.H/2, rect.H/2, rect.H/2, fg)
	h := rect.H / 4
	if h < 1 {
		h = 1
	}
	renderer.FillRect(&sdl.Rect{rect.H / 3, (rect.H - h) / 2, rect.W - rect.H, h})
	ui.FillCircle(renderer, int32(float64(rect.W)*0.97), int32(float64(rect.H)/2), int32(float64(rect.H)/3), fg)
	ui.FillCircle(renderer, int32(float64(rect.W)*0.92), int32(float64(rect.H)/2), int32(float64(rect.H)/2.5), fg)
	ui.FillCircle(renderer, int32(float64(rect.W)*0.87), int32(float64(rect.H)/2), int32(float64(rect.H)/2.25), fg)
	ui.FillCircle(renderer, int32(float64(rect.W)*0.82), int32(float64(rect.H)/2), int32(float64(rect.H)/2), fg)
	ui.FillCircle(renderer, center.X, center.Y, rect.H/5, bg)
	renderer.SetRenderTarget(nil)
	paintRect := sdl.Rect{rect.X + width/2 - center.X, rect.Y + height/2 - rect.H/2, rect.W, rect.H}

	return &ClockHand{
		renderer:   renderer,
		texture:    texHand,
		rect:       rect,
		paintRect:  paintRect,
		handCenter: center,
		fg:         fg,
		bg:         bg,
	}
}

func NewBigHand(renderer *sdl.Renderer, width, height int32, rect sdl.Rect, center sdl.Point, fg, bg sdl.Color) *ClockHand {
	texHand, err := renderer.CreateTexture(sdl.PIXELFORMAT_ARGB8888, sdl.TEXTUREACCESS_TARGET, rect.W, rect.H)
	if err != nil {
		panic(err)
	}
	renderer.SetRenderTarget(texHand)
	texHand.SetBlendMode(sdl.BLENDMODE_BLEND)
	ui.SetColor(renderer, bg)
	renderer.Clear()
	ui.SetColor(renderer, fg)
	renderer.DrawLine(0, 0, rect.W-rect.H/2, rect.H/4)
	renderer.DrawLine(rect.W-rect.H/2, rect.H/4, rect.W, rect.H-rect.H/2)
	renderer.DrawLine(rect.W-rect.H/2, (rect.H-1)-rect.H/4, rect.W, rect.H-rect.H/2)
	renderer.DrawLine(0, rect.H-1, rect.W-rect.H/2, (rect.H-1)-rect.H/4)
	renderer.DrawLine(0, 0, 0, rect.H)
	ui.FillCircle(renderer, center.X, center.Y, rect.H/5, fg)
	renderer.SetRenderTarget(nil)
	paintRect := sdl.Rect{rect.X + width/2 - center.X, rect.Y + height/2 - rect.H/2, rect.W, rect.H}

	return &ClockHand{
		renderer:   renderer,
		texture:    texHand,
		rect:       rect,
		paintRect:  paintRect,
		handCenter: center,
		fg:         fg,
		bg:         bg,
	}
}
