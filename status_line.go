package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type StatusLine struct {
	rect             sdl.Rect
	fg, bg           sdl.Color
	sprites          []Sprite
	fnClock, fnTimer func()
}

func NewStatusLine(rect sdl.Rect, fg, bg sdl.Color, renderer *sdl.Renderer, font *ttf.Font, fnClock, fnTimer func()) *StatusLine {
	var sprites []Sprite
	btnClock := NewButton(renderer, "Clock", sdl.Rect{rect.X, rect.Y, rect.H * 3, rect.H}, fg, bg, font, fnClock)
	sprites = append(sprites, btnClock)

	btnTimer := NewButton(renderer, "Timer", sdl.Rect{rect.X + rect.H*3, rect.Y, rect.H * 3, rect.H}, fg, bg, font, fnTimer)
	sprites = append(sprites, btnTimer)

	return &StatusLine{
		rect:    rect,
		fg:      fg,
		bg:      bg,
		sprites: sprites,
		fnClock: fnClock,
		fnTimer: fnTimer,
	}
}

func (s *StatusLine) Render(renderer *sdl.Renderer) {
	setColor(renderer, s.bg)
	renderer.FillRect(&s.rect)
	for _, sprite := range s.sprites {
		sprite.Render(renderer)
	}
}

func (s *StatusLine) Update() {
	for _, sprite := range s.sprites {
		sprite.Update()
	}
}

func (s *StatusLine) Event(e sdl.Event) {
	switch t := e.(type) {
	case *sdl.KeyboardEvent:

		if t.Keysym.Sym == sdl.K_F1 && t.State == sdl.RELEASED {
			s.fnClock()
		}
		if t.Keysym.Sym == sdl.K_F2 && t.State == sdl.RELEASED {
			s.fnTimer()
		}
	}
	for _, sprite := range s.sprites {
		sprite.Event(e)
	}
}

func (s *StatusLine) Destroy() {
	for _, sprite := range s.sprites {
		sprite.Destroy()
	}
}
