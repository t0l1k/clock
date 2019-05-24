package main

import (
	"fmt"
	"time"
)

type Lap struct {
	split       int
	splitLenght time.Duration
	timerLenght time.Duration
}

func NewLap(split int, timerLenght, splitLenght time.Duration) *Lap {
	return &Lap{
		split:       split,
		splitLenght: splitLenght,
		timerLenght: timerLenght,
	}
}

func (s *Lap) String() string {
	return fmt.Sprintf("lap#%v:\tTimer:%v\tSplit:%v", s.split, s.timerLenght, s.splitLenght)
}
