package main

type bounds struct {
	x0 float64
	x1 float64
	y0 float64
	y1 float64
}

const (
	WIDTH  = 1280
	HEIGHT = 800

	MAXITERATIONS = 1024
)

var (
	ribnds = bounds{
		x0: -2,
		x1: 1,
		y0: -1,
		y1: 1,
	}
)
