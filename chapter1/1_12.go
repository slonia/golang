package main

import (
	"image"
	"image/color"
	"image/gif"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var palette = []color.Color{color.Black, color.RGBA{0x00, 0xFF, 0x00, 0xFF}}

const (
	whiteIndex = 0 // first color in palette
	blackIndex = 1 // next color in palette
)

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	lissajous(w, r)
}

func paramOrDefault(r url.Values, name string, def float64) float64 {
	keys, ok := r[name]

	if !ok || len(keys) < 1 {
		log.Printf("Url Param '%s' is missing, will use %v\n", name, def)
		return def
	}
	val, _ := strconv.Atoi(keys[0])

	log.Printf("Url Param '%s' found, using %v\n", name, val)
	return float64(val)
}

func lissajous(out io.Writer, r *http.Request) {
	rand.Seed(time.Now().UTC().UnixNano())
	cycles := paramOrDefault(r.URL.Query(), "cycles", 5.0)
	res := paramOrDefault(r.URL.Query(), "res", 0.001)
	size := paramOrDefault(r.URL.Query(), "size", 100.0)
	delay := paramOrDefault(r.URL.Query(), "nframes", 64.0)
	nframes := paramOrDefault(r.URL.Query(), "delay", 8.0)
	freq := rand.Float64() * 3.0 // relative frequency of y oscillator
	anim := gif.GIF{LoopCount: int(nframes)}
	phase := 0.0 // phase difference
	for i := 0; i < int(nframes); i++ {
		rect := image.Rect(0, 0, int(2*size+1), int(2*size+1))
		img := image.NewPaletted(rect, palette)
		for t := 0.0; t < cycles*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			img.SetColorIndex(int(size)+int(x*float64(size)+0.5), int(size)+int(y*float64(size)+0.5),
				blackIndex)
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, int(delay))
		anim.Image = append(anim.Image, img)
	}
	gif.EncodeAll(out, &anim)
}
