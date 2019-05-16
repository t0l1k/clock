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

func FillCircle(renderer *sdl.Renderer, x0, y0, radius int32, color sdl.Color) {
	setColor(renderer, color)
	for y := -radius; y <= radius; y++ {
		for x := -radius; x <= radius; x++ {
			if x*x+y*y <= radius*radius {
				renderer.DrawPoint(x0+x, y0+y)
			}
		}
	}
}

// func getPixel(texture *sdl.Texture, x0, y0 int32) (color sdl.Color) {
// 	var surface *sdl.Surface
// 	format, _, width, height, _ := texture.Query()
// 	surface = sdl.CreateRGBSurface(format,width,height)
// 	return color
// }

// func FloodFill(renderer *sdl.Renderer, x0, y0 int32, newColor, oldColor sdl.Color) {
// }
