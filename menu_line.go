package main

import (
	"fmt"

	"github.com/t0l1k/sdl2/sdl2/ui"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type MenuLine struct {
	rect               sdl.Rect
	fg, bg             sdl.Color
	lblTitle, lblClock *ui.Label
	btnQuit            *ui.Button
}

func NewMenuLine(title string, rect sdl.Rect, fg, bg sdl.Color, renderer *sdl.Renderer, font *ttf.Font, fn func()) *MenuLine {
	lblTitle := ui.NewLabel(title, sdl.Point{rect.H + 3, rect.Y}, fg, renderer, font)
	clockStr := "00:00:00.000"
	w, _, _ := font.SizeUTF8(clockStr)
	lblClock := ui.NewLabel(clockStr, sdl.Point{rect.W - int32(w) - 3, rect.Y}, fg, renderer, font)
	btnQuit := ui.NewButton(renderer, "<-", sdl.Rect{rect.X, rect.Y, rect.H, rect.H}, fg, bg, font, fn)
	return &MenuLine{
		rect:     rect,
		fg:       fg,
		bg:       bg,
		lblTitle: lblTitle,
		lblClock: lblClock,
		btnQuit:  btnQuit,
	}
}

func (s *MenuLine) Render(renderer *sdl.Renderer) {
	ui.SetColor(renderer, s.bg)
	renderer.FillRect(&s.rect)
	s.lblTitle.Render(renderer)
	s.btnQuit.Render(renderer)
	s.lblClock.Render(renderer)
}

func (s *MenuLine) Update() {
	mSec, sec, min, hour := getTime()
	strClock := fmt.Sprintf("%02v:%02v:%02v.%03v", hour, min, sec, mSec)
	s.lblClock.SetText(strClock)
	s.lblClock.Update()
	s.lblTitle.Update()
	s.btnQuit.Update()
}
func (s *MenuLine) Event(e sdl.Event) {
	s.lblTitle.Event(e)
	s.btnQuit.Event(e)
	s.lblClock.Event(e)
}
func (s *MenuLine) Destroy() {
	s.lblTitle.Destroy()
	s.btnQuit.Destroy()
	s.lblClock.Destroy()
}
