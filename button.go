package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

// ver 2019.05.29
type Button struct {
	renderer                          *sdl.Renderer
	texFocus, texNotFocus, texPressed *sdl.Texture
	rect                              sdl.Rect
	fg, bg                            sdl.Color
	str                               string
	font                              *ttf.Font
	focus, pressed, show              bool
	fn                                func()
}

func NewButton(renderer *sdl.Renderer, str string, rect sdl.Rect, fg, bg sdl.Color, font *ttf.Font, fn func()) *Button {
	texFocus := newButtonTexture(renderer, str, rect, fg, bg, font, false)
	texNotFocus := newButtonTexture(renderer, str, rect, bg, fg, font, false)
	texPressed := newButtonTexture(renderer, str, rect, fg, bg, font, true)
	_, _, w, h, _ := texFocus.Query()
	return &Button{
		str:         str,
		renderer:    renderer,
		texFocus:    texFocus,
		texNotFocus: texNotFocus,
		texPressed:  texPressed,
		rect:        sdl.Rect{rect.X, rect.Y, w, h},
		fg:          fg,
		bg:          bg,
		font:        font,
		focus:       false,
		show:        true,
		fn:          fn,
	}
}

func newButtonTexture(renderer *sdl.Renderer, str string, rect sdl.Rect, fg, bg sdl.Color, font *ttf.Font, pressed bool) *sdl.Texture {
	labelTexture := newLabelTexture(str, fg, renderer, font)
	defer labelTexture.Destroy()
	buttonTexture, err := renderer.CreateTexture(sdl.PIXELFORMAT_ARGB8888, sdl.TEXTUREACCESS_TARGET, rect.W, rect.H)
	if err != nil {
		panic(err)
	}
	_, _, w, h, _ := labelTexture.Query()
	labelRect := sdl.Rect{(rect.W - w) / 2, (rect.H - h) / 2, w, h}
	renderer.SetRenderTarget(buttonTexture)
	renderer.SetDrawColor(bg.R, bg.G, bg.B, bg.A)
	renderer.Clear()
	renderer.SetDrawColor(fg.R, fg.G, fg.B, fg.A)
	renderer.Copy(labelTexture, nil, &labelRect)
	renderer.SetDrawColor(255, 0, 0, 255)
	renderer.DrawRect(&sdl.Rect{2, 2, rect.W - 4, rect.H - 4})
	if pressed {
		renderer.DrawRect(&sdl.Rect{4, 4, rect.W - 8, rect.H - 8})
	}
	renderer.SetRenderTarget(nil)
	return buttonTexture
}

func (s *Button) GetText() string {
	return s.str
}

func (s *Button) SetText(str string) {
	if s.show {
		s.Destroy()
		s.texFocus = newButtonTexture(s.renderer, str, s.rect, s.fg, s.bg, s.font, false)
		s.texNotFocus = newButtonTexture(s.renderer, str, s.rect, s.bg, s.fg, s.font, false)
		s.texPressed = newButtonTexture(s.renderer, str, s.rect, s.fg, s.bg, s.font, true)
	}
}

func (s *Button) GetSize() sdl.Rect {
	return s.rect
}

func (s *Button) SetSize(width, height int32) {
	if s.show {
		s.rect.W = width
		s.rect.H = height
		s.SetText(s.str)
	}
}

func (s *Button) GetShow() bool {
	return s.show
}

func (s *Button) SetShow(show bool) {
	if s.show {
		s.show = show
	}
}

func (s *Button) SetPos(pos sdl.Point) {
	s.rect.X = pos.X
	s.rect.Y = pos.Y
}

func (s *Button) IsPressed() bool {
	return s.pressed
}

func (s *Button) Render(renderer *sdl.Renderer) {
	if s.show {
		if s.focus && !s.pressed {
			if err := renderer.Copy(s.texFocus, nil, &s.rect); err != nil {
				panic(err)
			}
		} else if s.focus && s.pressed {
			if err := renderer.Copy(s.texPressed, nil, &s.rect); err != nil {
				panic(err)
			}
		} else if !s.focus {
			if err := renderer.Copy(s.texNotFocus, nil, &s.rect); err != nil {
				panic(err)
			}
		}
	}
}

func (s *Button) Update() {
	if s.show {
		x, y, state := sdl.GetMouseState()
		mousePoint := sdl.Point{x, y}
		if mousePoint.InRect(&s.rect) {
			s.focus = true
		} else {
			s.focus = false
		}
		if s.focus && state > 0 {
			s.pressed = true
		} else if s.focus && state == 0 && s.pressed {
			s.fn()
			s.pressed = false
		}
	}
}

func (s *Button) Event(e sdl.Event) {}

func (s *Button) Destroy() {
	s.texNotFocus.Destroy()
	s.texFocus.Destroy()
	s.texPressed.Destroy()
}

func (s *Button) String() string {
	return fmt.Sprintf("%v %v %v %v %v", s.rect, s.focus, s.texFocus, s.texNotFocus, s.texPressed)
}
