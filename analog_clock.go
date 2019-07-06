package main

import (
	"fmt"
	"time"

	"github.com/t0l1k/sdl2/sdl2/ui"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type GetTime func() (int, int, int, int)

type stopWatchState int8

const (
	stopWatchBegin stopWatchState = iota
	stopWatchPlay
	stopWatchPause
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
	stopWatch                                  *StopWatch
	prevStopWatch                              StopWatch
	stopWatchFrame                             *Frame
	laps                                       []string
	lapCount                                   int
	tipTweentyPoint                            sdl.Point
	fnAnalog, fnDigit                          GetTime
	lblDigTime                                 *ui.Label
}

func NewAnalogClock(renderer *sdl.Renderer, rect sdl.Rect, fg, secHandColor, tweentyPointColor, bg sdl.Color, font *ttf.Font, stopWatch *StopWatch, blinkTimer *BlinkTimer) *AnalogClock {
	texFace := NewClockFace(renderer, rect, fg, bg)
	rectWidth, rectHeight := int32(float64(rect.H)*0.470), int32(float64(rect.H)*0.02)
	mSecHand := NewSmallHand(renderer, rect.W, rect.H, sdl.Rect{rect.X, rect.Y, int32(float64(rectWidth) * 1), rectHeight / 2}, sdl.Point{int32(float64(rectHeight) * 0.2), rectHeight / 4}, secHandColor, bg)
	secondHand := NewSmallHandRounded(renderer, rect.W, rect.H, sdl.Rect{rect.X, rect.Y, int32(float64(rectWidth) * 1.235), rectHeight}, sdl.Point{int32(float64(rectWidth) * 0.2), rectHeight / 2}, secHandColor, bg)
	minuteHand := NewBigHand(renderer, rect.W, rect.H, sdl.Rect{rect.X, rect.Y, int32(float64(rectWidth) * 0.9), rectHeight * 2}, sdl.Point{rectHeight * 2, rectHeight / 2 * 2}, fg, bg)
	hourHand := NewBigHand(renderer, rect.W, rect.H, sdl.Rect{rect.X, rect.Y, int32(float64(rectWidth) * 0.7), rectHeight * 2}, sdl.Point{rectHeight * 2, rectHeight / 2 * 2}, fg, bg)
	tipTweentyPoint := getTip(sdl.Point{rect.W / 2, rect.H / 2}, 0/60, float64(rect.H/2-(rect.H/90)*3), 0, 0)

	clockStr := "00:00:00.000"
	w, _, _ := font.SizeUTF8(clockStr)
	lblClock := ui.NewLabel(clockStr, sdl.Point{rect.W - int32(w) - 3, rect.Y}, fg, renderer, font)
	stopWatchFrame := NewFrame([]string{" "}, sdl.Point{rect.X, rect.Y}, fg, bg, renderer, font)

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
		stopWatch:         stopWatch,
		tweentyBlinkTimer: blinkTimer,
		tipTweentyPoint:   tipTweentyPoint,
		fnAnalog:          getTime,
		fnDigit:           stopWatch.GetStopWatch,
		lblDigTime:        lblClock,
		stopWatchFrame:    stopWatchFrame,
	}
}

func (s *AnalogClock) Render(renderer *sdl.Renderer) {
	if err := renderer.Copy(s.texFace, nil, &s.rect); err != nil {
		panic(err)
	}
	if s.tweentyBlinkTimer.IsOn() {
		ui.FillCircle(s.renderer, s.rect.X+s.tipTweentyPoint.X, s.rect.Y+s.tipTweentyPoint.Y, s.rect.H/200, s.tweentyPointColor)
	}
	s.lblDigTime.Render(s.renderer)
	s.hourHand.Render(s.renderer)
	s.minuteHand.Render(s.renderer)
	s.secondHand.Render(s.renderer)
	if s.drawMsec {
		s.mSecHand.Render(s.renderer)
	}
	s.stopWatchFrame.Render(renderer)
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
		lblStr = fmt.Sprintf("%02d:%02d %02d.%03d", hour, minute, second, mSec)
	} else {
		lblStr = fmt.Sprintf("%02d %02d:%02d %03d", hour, minute, second, mSec)
	}
	s.lblDigTime.SetText(lblStr)
}

func (s *AnalogClock) Event(event sdl.Event) {
	switch t := event.(type) {
	case *sdl.KeyboardEvent:
		if t.Keysym.Sym == sdl.K_r && t.State == sdl.RELEASED {
			s.setStopWatchStateBegin()
		}
		if t.Keysym.Sym == sdl.K_RETURN && t.State == sdl.RELEASED {
			if !s.stopWatch.IsPaused() {
				s.setStopWatchStatePause()
			} else {
				s.setStopWatchStatePlay()
			}
		}
		if t.Keysym.Sym == sdl.K_SPACE && t.State == sdl.RELEASED {
			s.setStopWatchLap()
		}
	}
}

func (s *AnalogClock) setStopWatchStateBegin() {
	s.stopWatch.Reset()
	s.lapCount = 0
	s.prevStopWatch = *s.stopWatch
	s.laps = s.laps[:0]
	s.stopWatchFrame.SetText([]string{" "})
}

func (s *AnalogClock) setStopWatchStatePlay() {
	if !s.stopWatch.IsPaused() {
		s.prevStopWatch = *s.stopWatch
	}
	s.stopWatch.Start()
}
func (s *AnalogClock) setStopWatchStatePause() {
	s.stopWatch.SetPause()
}

func (s *AnalogClock) setStopWatchLap() {
	if !s.stopWatch.IsPaused() {
		s.lapCount++
		dur, _ := time.ParseDuration(s.stopWatch.String())
		lap := NewLap(s.lapCount, dur, s.stopWatch.Sub(s.prevStopWatch).Round(time.Millisecond))
		s.laps = append(s.laps, lap.String())
		fmt.Println(lap.String())
		s.stopWatchFrame.SetText(s.laps)
		s.prevStopWatch = *s.stopWatch
	}
}

func (s *AnalogClock) Destroy() {
	s.lblDigTime.Destroy()
	s.texFace.Destroy()
	s.hourHand.Destroy()
	s.minuteHand.Destroy()
	s.secondHand.Destroy()
	s.mSecHand.Destroy()
}

func NewClockFace(renderer *sdl.Renderer, rect sdl.Rect, fg, bg sdl.Color) (texClockFace *sdl.Texture) {
	var err error
	if texClockFace, err = renderer.CreateTexture(sdl.PIXELFORMAT_ARGB8888, sdl.TEXTUREACCESS_TARGET, rect.W, rect.H); err != nil {
		panic(err)
	}
	center := sdl.Point{rect.W / 2, rect.H / 2}
	margin := rect.H / 90
	renderer.SetRenderTarget(texClockFace)
	texClockFace.SetBlendMode(sdl.BLENDMODE_BLEND)
	ui.SetColor(renderer, bg)
	renderer.Clear()
	ui.SetColor(renderer, sdl.Color{255, 255, 128, 255})
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
		ui.FillCircle(renderer, tip.X, tip.Y, radius, bg)
		ui.FillCircle(renderer, tip.X, tip.Y, radius/2, fg)
	}
	renderer.SetRenderTarget(nil)
	return texClockFace
}
