package main

import (
	"fmt"

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

type timerState int8

const (
	timerBegin timerState = iota
	timerPlay
	timerPause
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
	lblTime, lblTitle      *Label
	btnClock, btnTimer     *Button
	blinkTimer             *BlinkTimer
	analogClock            *AnalogClock
	timer                  *Timer
	sprites                []Sprite
	fnTime                 GetTime
}

func NewScreen(title string, window *sdl.Window, renderer *sdl.Renderer, width, height int32) *Screen {
	return &Screen{
		title:    title,
		window:   window,
		renderer: renderer,
		width:    width,
		height:   height,
		bg:       sdl.Color{192, 192, 192, 0},
		fg:       sdl.Color{0, 0, 0, 255},
		fnTime:   getTime,
	}
}

func (s *Screen) setup() {
	var err error
	fontSize := int(float64(s.height) * 0.03) // Главная константа перерисовки экрана
	s.font, err = ttf.OpenFont("assets/Roboto-Regular.ttf", fontSize)
	if err != nil {
		panic(err)
	}

	s.lblTitle = NewLabel(s.title, sdl.Point{0, 0}, s.fg, s.renderer, s.font)
	s.sprites = append(s.sprites, s.lblTitle)

	s.lblTime = NewLabel("--:--:--", sdl.Point{0, 0}, s.fg, s.renderer, s.font)
	lblRect := s.lblTime.GetSize()
	lblPos := sdl.Point{s.width/2 - lblRect.W/2, s.height - lblRect.H}
	s.lblTime.SetPos(lblPos)
	s.sprites = append(s.sprites, s.lblTime)

	s.btnClock = NewButton(s.renderer, "Clock", sdl.Rect{0, s.height - lblRect.H, lblRect.H * 3, lblRect.H}, s.fg, s.bg, s.font)
	s.sprites = append(s.sprites, s.btnClock)

	s.btnTimer = NewButton(s.renderer, "Timer", sdl.Rect{lblRect.H * 3, s.height - lblRect.H, lblRect.H * 3, lblRect.H}, s.fg, s.bg, s.font)
	s.sprites = append(s.sprites, s.btnTimer)

	s.analogClock = NewAnalogClock(s.renderer, sdl.Rect{(s.width - s.height) / 2, lblRect.H, s.height, s.height - lblRect.H*2}, s.fg, sdl.Color{255, 0, 0, 255}, sdl.Color{255, 255, 0, 255}, s.bg, s.blinkTimer, s.fnTime)
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
			s.setTimerStateBegin()
		}
		if t.Keysym.Sym == sdl.K_RETURN && t.State == sdl.RELEASED {
			if !s.timer.IsPaused() {
				s.setTimerStatePause()
			} else {
				s.setTimerStatePlay()
			}
		}
		if t.Keysym.Sym == sdl.K_SPACE && t.State == sdl.RELEASED {
		}
		if t.Keysym.Sym == sdl.K_F1 && t.State == sdl.RELEASED {
			s.selectClock()
		}
		if t.Keysym.Sym == sdl.K_F2 && t.State == sdl.RELEASED {
			s.selectTimer()
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
	if s.btnClock.IsReleased() {
		s.selectClock()
	}
	if s.btnTimer.IsReleased() {
		s.selectTimer()
	}
	for _, sprite := range s.sprites {
		sprite.Event(event)
	}
}

func (s *Screen) selectClock() {
	s.fnTime = getTime
	s.title = "Clock"
	s.Destroy()
	s.setup()
}

func (s *Screen) selectTimer() {
	s.fnTime = s.timer.GetTimer
	s.title = "Timer"
	s.Destroy()
	s.setup()
}

func (s *Screen) setTimerStateBegin() {
	s.timer.Reset()
}
func (s *Screen) setTimerStatePlay() {
	s.timer.Start()
}
func (s *Screen) setTimerStatePause() {
	s.timer.SetPause()
}

func (s *Screen) Update() {
	_, second, minute, hour := s.fnTime()
	lblStr := ""
	if s.blinkTimer.IsOn() {
		lblStr = fmt.Sprintf("%02d:%02d:%02d", hour, minute, second)
	} else {
		lblStr = fmt.Sprintf("%02d %02d %02d", hour, minute, second)
	}
	s.lblTime.SetText(lblStr)
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

func (s *Screen) Init() {
	s.timer = NewTimer()
	s.timer.Reset()
	go s.timer.Run()
	s.blinkTimer = &BlinkTimer{}
	go s.blinkTimer.Run()
	s.setup()
}

func (s *Screen) Run() {
	s.Init()
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
	s.blinkTimer.Stop()
	s.timer.Stop()
}

func (s *Screen) Destroy() {
	for _, sprite := range s.sprites {
		sprite.Destroy()
	}
	s.sprites = s.sprites[:0]
	s.font.Close()
}

func (s *Screen) quit() { s.running = false }
