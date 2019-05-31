package main

import (
	"fmt"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type Sprite interface {
	Render(*sdl.Renderer)
	Update()
	Event(sdl.Event)
	Destroy()
}

type GetTime func() (int, int, int, int)

type stopWatchState int8

const (
	stopWatchBegin stopWatchState = iota
	stopWatchPlay
	stopWatchPause
)

type Screen struct {
	title                  string
	window                 *sdl.Window
	renderer               *sdl.Renderer
	width, height          int32
	font                   *ttf.Font
	flags                  uint32
	running                bool
	bg, fg, lineBg         sdl.Color
	fpsCountTime, fpsCount uint32
	blinkTimer             *BlinkTimer
	stopWatch              *StopWatch
	prevStopWatch          StopWatch
	sprites                []Sprite
	fnAnalog, fnDigit      GetTime
	laps                   []Lap
	lapCount               int
	menuLine               *MenuLine
	analogClock            *AnalogClock
	statusLine             *StatusLine
}

func NewScreen(title string, window *sdl.Window, renderer *sdl.Renderer, width, height int32) *Screen {
	stopWatch := NewStopWatch()
	stopWatch.Reset()
	go stopWatch.Run()
	blinkTimer := NewBlinkTimer(1000 / 2)
	go blinkTimer.Run()
	return &Screen{
		title:      title,
		window:     window,
		renderer:   renderer,
		width:      width,
		height:     height,
		bg:         sdl.Color{192, 192, 192, 0},
		fg:         sdl.Color{0, 0, 0, 255},
		lineBg:     sdl.Color{128, 128, 128, 0},
		fnAnalog:   getTime,
		fnDigit:    stopWatch.GetStopWatch,
		stopWatch:  stopWatch,
		blinkTimer: blinkTimer,
	}
}

func (s *Screen) setup() {
	var err error
	fontSize := int(float64(s.height) * 0.03) // Главная константа перерисовки экрана
	s.font, err = ttf.OpenFont("assets/Roboto-Regular.ttf", fontSize)
	if err != nil {
		panic(err)
	}
	lineHeight := int32(float64(fontSize) * 1.5)
	s.menuLine = NewMenuLine(
		s.title,
		sdl.Rect{0, 0, s.width, lineHeight},
		s.fg,
		s.lineBg,
		s.renderer,
		s.font,
		func() { s.quit() })
	s.sprites = append(s.sprites, s.menuLine)

	s.statusLine = NewStatusLine(
		sdl.Rect{0, s.height - lineHeight, s.width, lineHeight},
		s.fg,
		s.lineBg,
		s.renderer,
		s.font, s.selectClock, s.selectStopWatch)
	s.sprites = append(s.sprites, s.statusLine)

	s.analogClock = NewAnalogClock(
		s.renderer,
		sdl.Rect{(s.width - s.height) / 2, lineHeight, s.height, s.height - lineHeight*2},
		s.fg,
		sdl.Color{255, 0, 0, 255},
		s.bg,
		s.bg,
		s.font,
		s.blinkTimer,
		s.fnAnalog,
		s.fnDigit)
	s.sprites = append(s.sprites, s.analogClock)
}

func (s *Screen) setMode() {
	if s.flags == 0 {
		s.flags = sdl.WINDOW_FULLSCREEN_DESKTOP
		mode, err := sdl.GetCurrentDisplayMode(0)
		if err != nil {
			panic(err)
		}
		s.width, s.height = mode.W, mode.H
	} else {
		s.flags = 0
		s.width, s.height = 800, 600
	}
	s.window.SetFullscreen(s.flags)
	s.window.SetSize(s.width, s.height)
	s.Destroy()
	s.setup()
}

func (s *Screen) Event() {
	event := sdl.WaitEventTimeout(3)
	switch t := event.(type) {
	case *sdl.QuitEvent:
		s.quit()
	case *sdl.KeyboardEvent:
		if t.Keysym.Sym == sdl.K_ESCAPE && t.State == sdl.RELEASED {
			s.quit()
		}
		if t.Keysym.Sym == sdl.K_F11 && t.State == sdl.RELEASED {
			s.setMode()
		}
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
	case *sdl.WindowEvent:
		switch t.Event {
		case sdl.WINDOWEVENT_RESIZED:
			s.width, s.height = t.Data1, t.Data2
			s.Destroy()
			s.setup()
			// fmt.Println("window resized", s.width, s.height)
		case sdl.WINDOWEVENT_FOCUS_GAINED:
			// fmt.Println("window focus gained", s.width, s.height)
		case sdl.WINDOWEVENT_FOCUS_LOST:
			// fmt.Println("window focus lost", s.width, s.height)
		case sdl.WINDOW_MINIMIZED:
			s.Destroy()
		case sdl.WINDOWEVENT_RESTORED:
			s.setup()
		}
	}
	for _, sprite := range s.sprites {
		sprite.Event(event)
	}
}

func (s *Screen) selectClock() {
	s.fnAnalog = getTime
	s.fnDigit = s.stopWatch.GetStopWatch
	s.title = "Clock"
	s.Destroy()
	s.setup()
}

func (s *Screen) selectStopWatch() {
	s.fnAnalog = s.stopWatch.GetStopWatch
	s.fnDigit = getTime
	s.title = "StopWatch"
	s.Destroy()
	s.setup()
}

func (s *Screen) setStopWatchStateBegin() {
	s.stopWatch.Reset()
	s.lapCount = 0
	s.prevStopWatch = *s.stopWatch
}

func (s *Screen) setStopWatchStatePlay() {
	if !s.stopWatch.IsPaused() {
		s.prevStopWatch = *s.stopWatch
	}
	s.stopWatch.Start()
}
func (s *Screen) setStopWatchStatePause() {
	s.stopWatch.SetPause()
}

func (s *Screen) setStopWatchLap() {
	if !s.stopWatch.IsPaused() {
		s.lapCount++
		dur, _ := time.ParseDuration(s.stopWatch.String())
		lap := NewLap(s.lapCount, dur, s.stopWatch.Sub(s.prevStopWatch).Round(time.Millisecond))
		fmt.Println(lap)
		s.prevStopWatch = *s.stopWatch
	}
}

func (s *Screen) Update() {
	for _, sprite := range s.sprites {
		sprite.Update()
	}
	if sdl.GetTicks()-s.fpsCountTime > 999 {
		s.window.SetTitle(fmt.Sprintf("%s fps:%v", s.title, s.fpsCount))
		s.fpsCount = 0
		s.fpsCountTime = sdl.GetTicks()
	}
}

func (s *Screen) Render() {
	setColor(s.renderer, s.bg)
	s.renderer.Clear()
	for _, sprite := range s.sprites {
		sprite.Render(s.renderer)
	}
	s.renderer.Present()
	s.fpsCount++
}

func (s *Screen) Run() {
	s.setup()
	frameRate := uint32(1000 / 60)
	lastTime := sdl.GetTicks()
	s.running = true
	for s.running {
		now := sdl.GetTicks()
		if now >= lastTime {
			i := 0
			for {
				s.Event()
				s.Update()
				lastTime += frameRate
				now = sdl.GetTicks()
				if lastTime > now {
					break
				}
				i++
				if i >= 3 {
					lastTime = now + frameRate
					break
				}
			}
			s.Render()
		} else {
			sdl.Delay(lastTime - now)
		}
	}
}

func (s *Screen) Destroy() {
	for _, sprite := range s.sprites {
		sprite.Destroy()
	}
	s.sprites = s.sprites[:0]
	s.font.Close()
}

func (s *Screen) quit() {
	s.blinkTimer.Stop()
	s.stopWatch.Stop()
	s.running = false
}
