package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type AnalogClock struct {
	renderer                                   *sdl.Renderer
	rect                                       sdl.Rect
	texFace                                    *sdl.Texture
	fg, bg, secHandColor, tweentyPointColor    sdl.Color
	font                                       *ttf.Font
	hourHand, minuteHand, secondHand, mSecHand *ClockHand
	drawMsec                                   bool
	tweentyBlinkTimer                          *BlinkTimer
	tipTweentyPoint                            sdl.Point
	fnAnalog, fnDigit                          GetTime
	lblDigTime                                 *Label
}

func NewAnalogClock(renderer *sdl.Renderer, rect sdl.Rect, fg, secHandColor, tweentyPointColor, bg sdl.Color, font *ttf.Font, blinkTimer *BlinkTimer, fnAnalog, fnDigit GetTime) *AnalogClock {
	texFace := NewClockFace(renderer, rect, fg, bg)
	rectWidth, rectHeight := int32(float64(rect.H)*0.470), int32(float64(rect.H)*0.02)
	mSecHand := NewSmallHand(renderer, rect.W, rect.H, sdl.Rect{rect.X, rect.Y, int32(float64(rectWidth) * 1), rectHeight / 2}, sdl.Point{int32(float64(rectHeight) * 0.2), rectHeight / 4}, secHandColor, bg)
	secondHand := NewSmallHandRounded(renderer, rect.W, rect.H, sdl.Rect{rect.X, rect.Y, int32(float64(rectWidth) * 1.20), rectHeight}, sdl.Point{int32(float64(rectWidth) * 0.2), rectHeight / 2}, secHandColor, bg)
	minuteHand := NewBigHand(renderer, rect.W, rect.H, sdl.Rect{rect.X, rect.Y, int32(float64(rectWidth) * 0.9), rectHeight * 2}, sdl.Point{rectHeight * 2, rectHeight / 2 * 2}, fg, bg)
	hourHand := NewBigHand(renderer, rect.W, rect.H, sdl.Rect{rect.X, rect.Y, int32(float64(rectWidth) * 0.7), rectHeight * 2}, sdl.Point{rectHeight * 2, rectHeight / 2 * 2}, fg, bg)
	tipTweentyPoint := getTip(sdl.Point{rect.W / 2, rect.H / 2}, 0/60, float64(rect.H/2-(rect.H/90)*3), 0, 0)

	lblDigTime := NewLabel("00:00:00", sdl.Point{0, 0}, fg, renderer, font)
	lblRect := lblDigTime.GetSize()
	lblPos := sdl.Point{rect.X + (rect.W-lblRect.W)/2, rect.Y + int32(float64(rect.H)*0.25)}
	lblDigTime.SetPos(lblPos)

	return &AnalogClock{
		renderer:          renderer,
		rect:              rect,
		texFace:           texFace,
		font:              font,
		fg:                fg,
		bg:                bg,
		tweentyPointColor: tweentyPointColor,
		hourHand:          hourHand,
		minuteHand:        minuteHand,
		secondHand:        secondHand,
		mSecHand:          mSecHand,
		tweentyBlinkTimer: blinkTimer,
		tipTweentyPoint:   tipTweentyPoint,
		fnAnalog:          fnAnalog,
		fnDigit:           fnDigit,
		lblDigTime:        lblDigTime,
	}
}

func (s *AnalogClock) Render(renderer *sdl.Renderer) {
	if err := renderer.Copy(s.texFace, nil, &s.rect); err != nil {
		panic(err)
	}
	if s.tweentyBlinkTimer.IsOn() {
		FillCircle(s.renderer, s.rect.X+s.tipTweentyPoint.X, s.rect.Y+s.tipTweentyPoint.Y, s.rect.H/200, s.tweentyPointColor)
	}
	s.lblDigTime.Render(s.renderer)
	s.hourHand.Render(s.renderer)
	s.minuteHand.Render(s.renderer)
	s.secondHand.Render(s.renderer)
	if s.drawMsec {
		s.mSecHand.Render(s.renderer)
	}
}

func (s *AnalogClock) Update() {
	mSec, second, minute, hour := s.fnAnalog()
	s.mSecHand.Update(float64(mSec) / 1000.0)
	s.secondHand.Update((float64(second) + s.mSecHand.GetFraction()) / 60.0)
	s.minuteHand.Update((float64(minute) + s.secondHand.GetFraction()) / 60.0)
	s.hourHand.Update((float64(hour) + s.minuteHand.GetFraction()) / 12.0)

	mSec, second, minute, hour = s.fnDigit()
	lblStr := ""
	if s.tweentyBlinkTimer.IsOn() {
		lblStr = fmt.Sprintf("%02d:%02d:%02d", hour, minute, second)
	} else {
		lblStr = fmt.Sprintf("%02d %02d %02d", hour, minute, second)
	}
	s.lblDigTime.SetText(lblStr)
}

func (s *AnalogClock) Event(sdl.Event) {}
func (s *AnalogClock) Destroy() {
	s.lblDigTime.Destroy()
	s.texFace.Destroy()
	s.hourHand.Destroy()
	s.minuteHand.Destroy()
	s.secondHand.Destroy()
	s.mSecHand.Destroy()
}
func (s *AnalogClock) String() { fmt.Sprintf("AnalogClock:%v %v %v", s.fg, s.bg, s.rect) }

func NewClockFace(renderer *sdl.Renderer, rect sdl.Rect, fg, bg sdl.Color) (texClockFace *sdl.Texture) {
	var err error
	if texClockFace, err = renderer.CreateTexture(sdl.PIXELFORMAT_ARGB8888, sdl.TEXTUREACCESS_TARGET, rect.W, rect.H); err != nil {
		panic(err)
	}
	center := sdl.Point{rect.W / 2, rect.H / 2}
	margin := rect.H / 90
	renderer.SetRenderTarget(texClockFace)
	texClockFace.SetBlendMode(sdl.BLENDMODE_BLEND)
	setColor(renderer, bg)
	renderer.Clear()
	setColor(renderer, sdl.Color{128, 255, 255, 255})
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
		FillCircle(renderer, tip.X, tip.Y, radius, bg)
		FillCircle(renderer, tip.X, tip.Y, radius/2, fg)
	}
	renderer.SetRenderTarget(nil)
	return texClockFace
}
