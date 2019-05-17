package main

import (
	"math"
	"time"
)

type Timer struct {
	nowTick, lastTick, mSecond, second int
	running, pause                     bool
}

func NewTimer() *Timer {
	return &Timer{}
}

func (s *Timer) Reset() {
	s.running = true
	s.pause = true
	s.mSecond = 0
	s.second = 0
	s.lastTick = s.update()
}

func (s *Timer) Start() {
	s.lastTick = s.update()
	s.nowTick = 0
	s.pause = false
}

func (s *Timer) IsPaused() bool {
	return s.pause
}

func (s *Timer) SetPause() {
	s.pause = true
}

func (s *Timer) Stop() {
	if s.running {
		s.running = false
	}
}

func (s *Timer) update() int {
	return time.Now().Nanosecond() / 1000000
	// return sdl.GetTicks()
}

func (s *Timer) Run() {
	for s.running {
		if s.running && !s.pause {
			var diff int
			s.nowTick = s.update()
			if s.nowTick >= s.lastTick {
				diff = s.nowTick - s.lastTick
			} else {
				diff = int(math.Abs(float64(s.lastTick - s.nowTick - s.lastTick)))
			}
			s.mSecond += diff
			if s.mSecond >= 1000 {
				s.mSecond -= 1000
				s.second++
			}
			s.lastTick = s.nowTick
		}
	}
}

func (s *Timer) GetTimer() (int, int, int, int) {
	second := s.second % 60
	minute := s.second % 3600 / 60
	hour := s.second % 86400 / 3600
	return int(s.mSecond), int(second), int(minute), int(hour)
}