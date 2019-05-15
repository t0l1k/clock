package main

import (
	"math"
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

func getTip(center sdl.Point, percent, lenght, width, height float64) (tip sdl.Point) {
	angle := (0.5 - percent) * (2.0 * math.Pi)
	sine := math.Sin(angle)
	cosine := math.Cos(angle)
	tip.X = center.X + int32(lenght*sine-width)
	tip.Y = center.Y + int32(lenght*cosine-height)
	return tip
}
