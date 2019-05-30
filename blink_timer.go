package main

import "time"

// ver 2019.05.29
type BlinkTimer struct {
	delay   time.Duration
	blinkOn bool
	running bool
}

func NewBlinkTimer(delay time.Duration) *BlinkTimer {
	return &BlinkTimer{
		delay: delay,
	}
}

func (s *BlinkTimer) IsOn() bool { return s.blinkOn }

func (s *BlinkTimer) switchOn() { s.blinkOn = !s.blinkOn }

func (s *BlinkTimer) Run() {
	s.running = true
	for s.running {
		time.Sleep(s.delay)
		s.switchOn()
	}
}

func (s *BlinkTimer) Start() {
	s.running = true
}

func (s *BlinkTimer) Stop() {
	s.blinkOn = true
	s.running = false
}
