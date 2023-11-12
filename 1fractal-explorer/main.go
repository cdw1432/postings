package main

import (
	"image"
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

func main() {
	FractalWindow()
}

func run() {

	cfg := pixelgl.WindowConfig{
		Title:  "Fractals",
		Bounds: pixel.R(0, 0, WIDTH, HEIGHT),
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	buffer := image.NewRGBA(image.Rect(0, 0, WIDTH, HEIGHT))
	canvas := pixelgl.NewCanvas(win.Bounds())

	go FractalDraw(buffer)
	for !win.Closed() {
		win.Clear(color.Black)

		canvas.SetPixels(buffer.Pix)
		canvas.Draw(win, pixel.IM.Moved(win.Bounds().Center()))
		win.Update()
	}
}

type imageIterator struct {
	image.Point
	bounds image.Rectangle
}

func FractalDraw(buffer *image.RGBA) {

	b := &imageIterator{bounds: buffer.Bounds()}
	for b.check() {
		r := ribnds.x0 + (float64(b.X)/WIDTH)*(ribnds.x1-ribnds.x0)
		i := ribnds.y0 + (float64(b.Y)/HEIGHT)*(ribnds.y1-ribnds.y0)
		iter := MandelBrot(r, i)

		if iter < MAXITERATIONS-1 {
			buffer.Set(b.X, b.Y, color.RGBA{255, 255, 255, 255})
		} else {
			buffer.Set(b.X, b.Y, color.RGBA{0, 0, 0, 255})
		}
	}

}

func (i *imageIterator) check() bool {
	if i.X < i.bounds.Max.X {
		i.X++
		return true
	}

	if i.Y < i.bounds.Max.Y {
		i.X = 0
		i.Y++
		return true
	}
	return false
}
func FractalWindow() {
	pixelgl.Run(run)
}
