package main

import (
	"fmt"

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
func (s *ClockHand) String()         { fmt.Sprintln("ClockHand:%v %v %v", s.fg, s.bg, s.rect) }

func NewSmallHand(renderer *sdl.Renderer, width, height int32, rect sdl.Rect, center sdl.Point, fg, bg sdl.Color) *ClockHand {
	texHand, err := renderer.CreateTexture(sdl.PIXELFORMAT_ARGB8888, sdl.TEXTUREACCESS_TARGET, rect.W, rect.H)
	if err != nil {
		panic(err)
	}
	renderer.SetRenderTarget(texHand)
	texHand.SetBlendMode(sdl.BLENDMODE_BLEND)
	setColor(renderer, bg)
	renderer.Clear()
	setColor(renderer, fg)
	renderer.FillRect(&sdl.Rect{0, 0, center.X - center.X/4, rect.H})
	renderer.FillRect(&sdl.Rect{0, rect.H / 3, rect.W, rect.H - rect.H/3*2})
	FillCircle(renderer, center.X, center.Y, rect.H/5, bg)
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
	setColor(renderer, bg)
	renderer.Clear()
	setColor(renderer, fg)
	renderer.DrawLine(0, 0, rect.W-rect.H/2, rect.H/4)
	renderer.DrawLine(rect.W-rect.H/2, rect.H/4, rect.W, rect.H-rect.H/2)
	renderer.DrawLine(rect.W-rect.H/2, (rect.H-1)-rect.H/4, rect.W, rect.H-rect.H/2)
	renderer.DrawLine(0, rect.H-1, rect.W-rect.H/2, (rect.H-1)-rect.H/4)
	renderer.DrawLine(0, 0, 0, rect.H)
	FillCircle(renderer, center.X, center.Y, rect.H/5, fg)
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
