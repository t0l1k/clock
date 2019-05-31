package main

// ver 2019.05.28
import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type Label struct {
	str      string
	renderer *sdl.Renderer
	texture  *sdl.Texture
	rect     sdl.Rect
	color    sdl.Color
	font     *ttf.Font
}

func NewLabel(str string, point sdl.Point, color sdl.Color, renderer *sdl.Renderer, font *ttf.Font) *Label {
	texture := newLabelTexture(str, color, renderer, font)
	_, _, w, h, _ := texture.Query()
	return &Label{
		str:      str,
		rect:     sdl.Rect{point.X, point.Y, w, h},
		color:    color,
		renderer: renderer,
		font:     font,
		texture:  texture,
	}
}

func newLabelTexture(str string, color sdl.Color, renderer *sdl.Renderer, font *ttf.Font) *sdl.Texture {
	var (
		surface *sdl.Surface
		texture *sdl.Texture
		err     error
	)
	if surface, err = font.RenderUTF8Blended(str, color); err != nil {
		panic(err)
	}
	defer surface.Free()
	if texture, err = renderer.CreateTextureFromSurface(surface); err != nil {
		panic(err)
	}
	return texture
}

func (s *Label) GetTexture() *sdl.Texture {
	return s.texture
}

func (s *Label) SetText(str string) {
	s.texture.Destroy()
	s.str = str
	s.texture = newLabelTexture(s.str, s.color, s.renderer, s.font)
	_, _, w, h, _ := s.texture.Query()
	s.SetSize(w, h)
}

func (s *Label) GetSize() sdl.Rect {
	return s.rect
}

func (s *Label) SetSize(w, h int32) {
	s.rect.W = w
	s.rect.H = h
}

func (s *Label) SetPos(pos sdl.Point) {
	s.rect.X = pos.X
	s.rect.Y = pos.Y
}

func (s *Label) Render(renderer *sdl.Renderer) {
	if err := renderer.Copy(s.texture, nil, &s.rect); err != nil {
		panic(err)
	}
}

func (s *Label) Update()         {}
func (s *Label) Event(sdl.Event) {}
func (s *Label) Destroy()        { s.texture.Destroy() }
func (s *Label) String()         { fmt.Sprintf("label:%v %v %v", s.str, s.color, s.rect) }
