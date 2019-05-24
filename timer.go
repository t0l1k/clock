package main

import (
	"fmt"
	"time"
)

// Timer умеет засекать время.
type Timer struct {
	startTick       time.Time
	nSecond, second int64
	running, pause  bool
}

// NewTimer создает экземпляр
func NewTimer() *Timer {
	return &Timer{}
}

// Reset обнулить таймер
func (s *Timer) Reset() {
	s.running = true
	s.pause = true
	s.nSecond = 0
	s.second = 0
	s.startTick = s.update()
}

// Start старт таймера
func (s *Timer) Start() {
	s.startTick = s.update()
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
	s.running = false
}

func (s *Timer) update() time.Time {
	return time.Now()
}

// Run запуск экземпляра таймера в отдельной горутине
func (s *Timer) Run() {
	for s.running {
		if !s.IsPaused() {
			nowTick := s.update()
			diff := nowTick.Sub(s.startTick).Nanoseconds()
			s.nSecond += diff
			if s.nSecond > 1000000000 {
				s.nSecond -= 1000000000
				s.second++
			}
			// fmt.Println("now:", s.nowTick, s.lastTick, diff, s.nSecond, s.second)
			s.startTick = nowTick
			time.Sleep(1 * time.Millisecond) // задержка для предотвращения троутлинга
		}
	}
}

func (s *Timer) Sub(u Timer) time.Duration {
	return time.Duration(s.second-u.second)*time.Second + time.Duration(s.nSecond-u.nSecond)
}

// GetTimer передать текущее время таймера
func (s *Timer) GetTimer() (int, int, int, int) {
	second := s.second % 60
	minute := s.second % 3600 / 60
	hour := s.second % 86400 / 3600
	return int(s.nSecond / 1000000), int(second), int(minute), int(hour)
}

func (s *Timer) String() (str string) {
	mS, sec, m, h := s.GetTimer()
	if h > 0 {
		str = fmt.Sprintf("%02vh%02vm%02vs%03vms", h, m, sec, mS)
	} else if h == 0 && m > 0 {
		str = fmt.Sprintf("%02vm%02vs%03vms", m, sec, mS)
	} else if h == 0 && m == 0 && sec > 0 {
		str = fmt.Sprintf("%02vs%03vms", sec, mS)
	} else if h == 0 && m == 0 && sec == 0 {
		str = fmt.Sprintf("%03vms", mS)
	}
	return str
}
