package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

type BlinkTimer struct {
	DELAY   uint32
	blinkOn bool
	running bool
}

func (s *BlinkTimer) IsOn() bool { return s.blinkOn }
func (s *BlinkTimer) switchOn()  { s.blinkOn = !s.blinkOn }
func (s *BlinkTimer) Run() {
	s.DELAY = 1000 / 2
	s.running = true
	for s.running {
		sdl.Delay(s.DELAY)
		s.switchOn()
	}
}
func (s *BlinkTimer) Quit() { s.running = false }
