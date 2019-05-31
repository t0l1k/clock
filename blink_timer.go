package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

type BlinkTimer struct {
	delay   uint32
	blinkOn bool
	running bool
}

func NewBlinkTimer(delay uint32) *BlinkTimer {
	return &BlinkTimer{
		delay:   delay,
		running: true,
	}
}

func (s *BlinkTimer) IsOn() bool { return s.blinkOn }

func (s *BlinkTimer) switchOn() { s.blinkOn = !s.blinkOn }

func (s *BlinkTimer) Run() {
	for s.running {
		sdl.Delay(s.delay)
		s.switchOn()
	}
}

func (s *BlinkTimer) Stop() {
	s.running = false
}
