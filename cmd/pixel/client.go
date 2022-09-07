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

func DrawQT(imd *imdraw.IMDraw, qt quadtree.QuadTree[boid.Boid]) {
	imd.Color = colornames.Black
	imd.Push(pixel.V(qt.Boundary.X, qt.Boundary.Y), pixel.V(qt.Boundary.X + qt.Boundary.W, qt.Boundary.Y + qt.Boundary.H))
	imd.Rectangle(1)

	for _, node := range(qt.Nodes) {
		DrawQT(imd, node)
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
			imd.Color = colornames.Blue
			imd.Push(pixel.V(b.Position.X(), b.Position.Y()))
			imd.Circle(5, 0)

			// imd.Color = colornames.Red
			// viewAngle := mathutil.GetVecAngle(b.Velocity)
			// viewPoint := mathutil.GetPointFromAngle(b.Position, float64(Flock.BoidSettings.Perception), viewAngle)
			// leftViewpoint := mathutil.GetPointFromAngle(b.Position, float64(Flock.BoidSettings.Perception), -boid.FOV / 2 + viewAngle)
			// rightViewpoint := mathutil.GetPointFromAngle(b.Position, float64(Flock.BoidSettings.Perception), boid.FOV / 2 + viewAngle)
   //
			// imd.Push(pixel.V(b.Position.X(), b.Position.Y()), pixel.V(viewPoint[0], viewPoint[1]))
			// imd.Push(pixel.V(b.Position.X(), b.Position.Y()), pixel.V(leftViewpoint[0], leftViewpoint[1]))
			// imd.Push(pixel.V(b.Position.X(), b.Position.Y()), pixel.V(rightViewpoint[0], rightViewpoint[1]))
			// imd.Line(1)

			// imd.Color = colornames.Blue
			// imd.Circle(float64(Flock.BoidSettings.Perception), 1)
			// imd.Push(pixel.V(b.Position.X(), b.Position.Y()))
		}

		// Angle between first and second boid
		// angleBetween := mathutil.GetAngleBetween(Flock.Boids[0].Position, Flock.Boids[1].Position)
		// perspectivePoint := mathutil.GetPointFromAngle(Flock.Boids[0].Position, 50, angleBetween)
		// imd.Push(pixel.V(Flock.Boids[0].Position.X(), Flock.Boids[0].Position.Y()), pixel.V(perspectivePoint[0], perspectivePoint[1]))
		// imd.Line(1)


		// Draw the quadtree
		//DrawQT(imd, Flock.QuadTree)


		imd.Draw(win)
		win.Update()
	}
}

func Execute() {
	pixelgl.Run(run)
}
