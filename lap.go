package main

import (
	"fmt"
	"time"
)

type Lap struct {
	time        time.Time
	split       int
	splitLenght time.Duration
	timerLenght time.Duration
}

func NewLap(split int, timerLenght, splitLenght time.Duration) Lap {
	return Lap{
		time:        time.Now(),
		split:       split,
		splitLenght: splitLenght,
		timerLenght: timerLenght,
	}
}

func (s *Lap) String() string {
	return fmt.Sprintf("#%v SW:%v Split:%v", s.split, s.timerLenght, s.splitLenght)
}
