package main

import (
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"net/url"
	"strconv"
)

const (
	width, height = 600, 320
	xyrange       = 30.0
	xyscale       = width / 2 / xyrange
	zscale        = height * 0.4
	angle         = math.Pi / 6
)

var cells float64
var sin30, cos30 = math.Sin(angle), math.Cos(angle)

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func corner(i, j int) (float64, float64, float64) {
	x := xyrange * (float64(i)/cells - 0.5)
	y := xyrange * (float64(j)/cells - 0.5)
	z := f(x, y)
	sx := width/2 + (x-y)*cos30*xyscale
	sy := height/2 + (x+y)*sin30*xyscale - z*zscale
	return sx, sy, z
}

func f(x, y float64) float64 {
	r := math.Hypot(x, y)
	val := math.Sin(r) / r
	if math.IsNaN(val) {
		return 0.0
	}
	return val
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/svg+xml")
	graph(w, r)
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

func graph(out io.Writer, r *http.Request) {
	cells = paramOrDefault(r.URL.Query(), "cells", 100.0)
	fmt.Fprintf(out, "<svg xmlns='http://www.w3.org/2000/svg' "+
		"style='stroke: grey; fill: white; stroke-width: 0.7' "+
		"width='%d' height='%d'>\n", width, height)
	for i := 0; i < int(cells); i++ {
		for j := 0; j < int(cells); j++ {
			ax, ay, az := corner(i+1, j)
			bx, by, bz := corner(i, j)
			cx, cy, cz := corner(i, j+1)
			dx, dy, dz := corner(i+1, j+1)
			avgz := (az + bz + cz + dz) / 4.0
			avgz = (avgz + 0.2/1.2) + 0.3
			r := int(255 * (1.0 - avgz))
			g := int(255 * avgz)
			fmt.Fprintf(out, "<polygon points='%g,%g %g,%g %g,%g %g,%g' style='fill: rgb(%d,%d,128);stroke: rgb(%d,%d,128)'/>\n", ax, ay, bx, by, cx, cy, dx, dy, r, g, r, g)
		}
	}
	fmt.Fprintf(out, "</svg>")
}
