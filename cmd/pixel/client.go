package pixel

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/kctjohnson/bubble-boids/internal/boid"
	"github.com/kctjohnson/bubble-boids/internal/quadtree"
	"golang.org/x/image/colornames"
)

var Flock *boid.Flock
var ScreenWidth float64 = 1024
var ScreenHeight float64 = 768

func DrawQT(imd *imdraw.IMDraw, qt *quadtree.QuadTree[boid.Boid]) {
	imd.Color = colornames.Black
	imd.Rectangle(1)
	imd.Push(pixel.V(qt.Boundary.X, qt.Boundary.Y), pixel.V(qt.Boundary.X + qt.Boundary.W, qt.Boundary.Y + qt.Boundary.H))

	if qt.NE != nil {
		DrawQT(imd, qt.NE)
	}
	if qt.NW != nil {
		DrawQT(imd, qt.NW)
	}
	if qt.SE != nil {
		DrawQT(imd, qt.SE)
	}
	if qt.SW != nil {
		DrawQT(imd, qt.SW)
	}
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Pixel Rocks!",
		Bounds: pixel.R(0, 0, ScreenWidth, ScreenHeight),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	Flock = boid.NewFlock(ScreenWidth, ScreenHeight)

	imd := imdraw.New(nil)

	for !win.Closed() {
		imd.Clear()
		win.Clear(colornames.Aliceblue)

		Flock.Update(ScreenWidth, ScreenHeight)
		for _, b := range Flock.Boids {
			imd.Color = colornames.Brown
			imd.Circle(5, 0)
			imd.Push(pixel.V(b.Position.X(), b.Position.Y()))
		}

		// Draw the quadtree
		DrawQT(imd, Flock.QuadTree)


		imd.Draw(win)
		win.Update()
	}
}

func Execute() {
	pixelgl.Run(run)
}
