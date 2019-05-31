package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type StatusLine struct {
	rect       sdl.Rect
	fg, bg     sdl.Color
	lblMessage *Label
}

func NewStatusLine(message string, rect sdl.Rect, fg, bg sdl.Color, renderer *sdl.Renderer, font *ttf.Font) *StatusLine {
	lblMessage := NewLabel(message, sdl.Point{rect.X + 5, rect.Y}, fg, renderer, font)
	return &StatusLine{
		rect:       rect,
		fg:         fg,
		bg:         bg,
		lblMessage: lblMessage,
	}
}

func (s *StatusLine) SetMessage(str string) {
	s.lblMessage.SetText(str)
}

func (s *StatusLine) Render(renderer *sdl.Renderer) {
	setColor(renderer, s.bg)
	renderer.FillRect(&s.rect)
	s.lblMessage.Render(renderer)
}

func (s *StatusLine) Update() {
	s.lblMessage.Update()
}

func (s *StatusLine) Event(e sdl.Event) {
	s.lblMessage.Event(e)
}

func (s *StatusLine) Destroy() {
	s.lblMessage.Destroy()
}
