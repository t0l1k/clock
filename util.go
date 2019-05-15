package main

import (
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

func setColor(renderer *sdl.Renderer, color sdl.Color) {
	renderer.SetDrawColor(color.R, color.G, color.B, color.A)
}

func getTime() (int, int, int, int) {
	mSec := time.Now().Nanosecond() / 1000000
	sec := time.Now().Second()
	minute := time.Now().Minute()
	hour := time.Now().Hour()
	return mSec, sec, minute, hour
}
