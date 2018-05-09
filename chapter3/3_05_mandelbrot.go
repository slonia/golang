package main

import (
	"image"
	"image/color"
	"image/png"
	"math/cmplx"
	"math/rand"
	"os"
	"time"
)

var colors []color.RGBA

func main() {
	const (
		xmin, ymin, xmax, ymax = -2, -2, +2, +2
		width, height          = 1024, 1024
	)
	generateColors(200)
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for py := 0; py < height; py++ {
		y := float64(py)/height*(ymax-ymin) + ymin
		for px := 0; px < width; px++ {
			x := float64(px)/width*(xmax-xmin) + xmin
			z := complex(x, y)
			img.Set(px, py, mandelbrot(z))
		}
	}
	png.Encode(os.Stdout, img)
}

func generateColors(n int) {
	rand.Seed(time.Now().UTC().UnixNano())
	colors = make([]color.RGBA, 0)
	for i := 0; i < n; i++ {
		r := uint8(rand.Float64() * 255.0)
		g := uint8(rand.Float64() * 255.0)
		b := uint8(rand.Float64() * 255.0)
		colors = append(colors, color.RGBA{r, g, b, 255})
	}
}

func mandelbrot(z complex128) color.Color {
	const iterations = 200
	// const contrast = 15
	var v complex128
	for n := uint8(0); n < iterations; n++ {
		v = v*v + z
		if cmplx.Abs(v) > 2 {
			return colors[n]
		}
	}
	return color.Black
}
