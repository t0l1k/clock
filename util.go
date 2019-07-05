package main

import (
	"math"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

func getTime() (int, int, int, int) {
	mSec := time.Now().Nanosecond() / 1000000
	sec := time.Now().Second()
	minute := time.Now().Minute()
	hour := time.Now().Hour()
	return mSec, sec, minute, hour
}

func getTip(center sdl.Point, percent, lenght, width, height float64) (tip sdl.Point) {
	radians := (0.5 - percent) * (2.0 * math.Pi)
	sine := math.Sin(radians)
	cosine := math.Cos(radians)
	tip.X = center.X + int32(lenght*sine-width)
	tip.Y = center.Y + int32(lenght*cosine-height)
	return tip
}

func getAngle(percent float64) float64 {
	radians := (0.5 - percent) * (2.0 * math.Pi)
	angle := (radians * -180 / math.Pi) + 90
	return angle
}
