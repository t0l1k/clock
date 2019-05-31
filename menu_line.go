package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type MenuLine struct {
	rect     sdl.Rect
	fg, bg   sdl.Color
	lblTitle *Label
	btnQuit  *Button
}

func NewMenuLine(title string, rect sdl.Rect, fg, bg sdl.Color, renderer *sdl.Renderer, font *ttf.Font, fn func()) *MenuLine {
	lblTitle := NewLabel(title, sdl.Point{rect.H + 3, rect.Y}, fg, renderer, font)
	btnQuit := NewButton(renderer, "<-", sdl.Rect{rect.X, rect.Y, rect.H, rect.H}, fg, bg, font, fn)
	return &MenuLine{
		rect:     rect,
		fg:       fg,
		bg:       bg,
		lblTitle: lblTitle,
		btnQuit:  btnQuit,
	}
}

func (s *MenuLine) Render(renderer *sdl.Renderer) {
	setColor(renderer, s.bg)
	renderer.FillRect(&s.rect)
	s.lblTitle.Render(renderer)
	s.btnQuit.Render(renderer)
}

func (s *MenuLine) Update() {
	s.lblTitle.Update()
	s.btnQuit.Update()
}
func (s *MenuLine) Event(e sdl.Event) {
	s.lblTitle.Event(e)
	s.btnQuit.Event(e)
}
func (s *MenuLine) Destroy() {
	s.lblTitle.Destroy()
	s.btnQuit.Destroy()
}
