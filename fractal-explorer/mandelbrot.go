package main

func MandelBrot(x0, y0 float64) int {
	var x = 0.0
	var y = 0.0

	iteration := 0

	for x*x+y*y <= 4 && iteration < MAXITERATIONS {
		xtmp := x*x - y*y + x0
		y = 2*x*y + y0
		x = xtmp
		iteration++
	}

	return int(iteration)
}
