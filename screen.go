package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type Screen struct {
	title                  string
	window                 *sdl.Window
	renderer               *sdl.Renderer
	width, height          int32
	font                   *ttf.Font
	flags                  uint32
	running                bool
	bg, fg                 sdl.Color
	fpsCountTime, fpsCount uint32
	lblTime                *Label
	blinkTimer             *BlinkTimer
	analogClock            *AnalogClock
}

func NewScreen(title string, window *sdl.Window, renderer *sdl.Renderer, width, height int32, font *ttf.Font) *Screen {
	return &Screen{
		title:    title,
		window:   window,
		renderer: renderer,
		width:    width,
		height:   height,
		font:     font,
		bg:       sdl.Color{0, 64, 0, 0},
		fg:       sdl.Color{255, 0, 255, 255},
	}
}

func (s *Screen) setup() {
	s.lblTime = NewLabel("--:--", sdl.Point{0, 0}, s.fg, s.renderer, s.font)
	lblRect := s.lblTime.GetSize()
	s.lblTime.SetPos(sdl.Point{s.width/2 - lblRect.W/2, s.height - lblRect.H})
	s.analogClock = NewAnalogClock(s.renderer, sdl.Rect{50, 50, 300, 300}, s.fg, s.bg)
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
	case *sdl.WindowEvent:
		switch t.Event {
		case sdl.WINDOWEVENT_RESIZED:
			s.width, s.height = t.Data1, t.Data2
			s.Destroy()
			s.setup()
			fmt.Println("window resized", s.width, s.height)
		case sdl.WINDOWEVENT_FOCUS_GAINED:
			fmt.Println("window focus gained", s.width, s.height)
		case sdl.WINDOWEVENT_FOCUS_LOST:
			fmt.Println("window focus lost", s.width, s.height)
		case sdl.WINDOW_MINIMIZED:
			s.Destroy()
		case sdl.WINDOWEVENT_RESTORED:
			s.setup()
		}
	}
}
func (s *Screen) Update() {
	_, _, minute, hour := getTime()
	lblStr := ""
	if s.blinkTimer.IsOn() {
		lblStr = fmt.Sprintf("%02d:%02d", hour, minute)
	} else {
		lblStr = fmt.Sprintf("%02d %02d", hour, minute)
	}
	s.lblTime.SetText(lblStr)

	if sdl.GetTicks()-s.fpsCountTime > 999 {
		s.window.SetTitle(fmt.Sprintf("%s fps:%v", s.title, s.fpsCount))
		s.fpsCount = 0
		s.fpsCountTime = sdl.GetTicks()
	}
}
func (s *Screen) Render() {
	setColor(s.renderer, s.bg)
	s.renderer.Clear()
	setColor(s.renderer, s.fg)
	s.lblTime.Render(s.renderer)
	s.analogClock.Render(s.renderer)

	s.renderer.Present()
	s.fpsCount++
}
func (s *Screen) quit() { s.running = false }
func (s *Screen) Run() {
	s.setup()
	s.blinkTimer = &BlinkTimer{}
	go s.blinkTimer.Run()
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
	s.blinkTimer.Quit()
}
func (s *Screen) Destroy() {}
