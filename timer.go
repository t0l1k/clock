package main

import (
	"math"
	"time"
)

// Timer умеет засекать время.
type Timer struct {
	nowTick, lastTick, mSecond, second int
	running, pause                     bool
}

// NewTimer создает экземпляр
func NewTimer() *Timer {
	return &Timer{}
}

// Reset обнулить текущее время таймера
func (s *Timer) Reset() {
	s.running = true
	s.pause = true
	s.mSecond = 0
	s.second = 0
	s.lastTick = s.update()
}

// Start старт таймера
func (s *Timer) Start() {
	s.lastTick = s.update()
	s.nowTick = 0
	s.pause = false
}

// IsPaused проверить установлена ли пауза
func (s *Timer) IsPaused() bool {
	return s.pause
}

// SetPause поставить таймер на паузу
func (s *Timer) SetPause() {
	s.pause = true
}

// Stop остановить таймер
func (s *Timer) Stop() {
	if s.running {
		s.running = false
	}
}

func (s *Timer) update() int {
	return time.Now().Nanosecond() / 1000000
}

// Run запуск экземпляра таймера в отдельной горутине
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
		time.Sleep(10 * time.Millisecond) // задержка для предотвращения троутлинга
	}
}

// GetTimer передать сколько текущее время таймера
func (s *Timer) GetTimer() (int, int, int, int) {
	second := s.second % 60
	minute := s.second % 3600 / 60
	hour := s.second % 86400 / 3600
	return s.mSecond, second, minute, hour
}
