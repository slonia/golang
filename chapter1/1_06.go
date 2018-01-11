// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// Run with "web" command-line argument for web server.
// See page 13.
//!+main

// Lissajous generates GIF animations of random Lissajous figures.
package main

import (
	"image"
	"image/color"
	"image/gif"
	"io"
	"math"
	"math/rand"
	"os"
	"time"
)

//!-main
// Packages not needed by version in book.

var palette = []color.Color{}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	for i := 0; i <= 16; i++ {
		colors := make([]int, 3)
		for i, _ := range colors {
			colors[i] = randInInt(0, 255)
		}
		palette = append(palette, color.RGBA{uint8(colors[0]), uint8(colors[1]), uint8(colors[2]), 255})
	}
	lissajous(os.Stdout)
}

func randInInt(from int, to int) int {
	f := rand.Float32()
	diff := float32(to - from)
	return int(float32(from) + f*diff)
}
func lissajous(out io.Writer) {
	const (
		cycles  = 5     // number of complete x oscillator revolutions
		res     = 0.001 // angular resolution
		size    = 100   // image canvas covers [-size..+size]
		nframes = 64    // number of animation frames
		delay   = 8     // delay between frames in 10ms units
	)
	freq := rand.Float64() * 3.0 // relative frequency of y oscillator
	anim := gif.GIF{LoopCount: nframes}
	phase := 0.0 // phase difference
	var imCol uint8
	for i := 0; i < nframes; i++ {
		imCol = uint8(i)/4 + 1
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, palette)
		for t := 0.0; t < cycles*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			img.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.5),
				imCol)
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	gif.EncodeAll(out, &anim)
}
