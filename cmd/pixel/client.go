package pixel

import (
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/kctjohnson/bubble-boids/internal/boid"
	"github.com/kctjohnson/bubble-boids/internal/dbscan"
	"github.com/kctjohnson/bubble-boids/internal/mathutil"
	"github.com/kctjohnson/bubble-boids/internal/quadtree"
	"golang.org/x/image/colornames"
)

var Flock *boid.Flock
var ScreenWidth float64 = 1920
var ScreenHeight float64 = 1080

var ClusterColors = []color.RGBA{
	colornames.Red,
	colornames.Green,
	colornames.Blue,
	colornames.Purple,
	colornames.Brown,
	colornames.Lightblue,
}

func DrawQT(imd *imdraw.IMDraw, qt quadtree.QuadTree[boid.Boid]) {
	imd.Color = colornames.Black
	imd.Rectangle(1)
	imd.Push(pixel.V(qt.Boundary.X, qt.Boundary.Y), pixel.V(qt.Boundary.X+qt.Boundary.W, qt.Boundary.Y+qt.Boundary.H))

	for _, node := range qt.Nodes {
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

		qtree := quadtree.NewQuadTree(mathutil.Rectangle[boid.Boid]{X: 0, Y: 0, W: ScreenWidth, H: ScreenHeight}, 10)
		positions := make([]mathutil.Point[boid.Boid], len(Flock.Boids))
		for i, b := range Flock.Boids {
			positions[i] = mathutil.Point[boid.Boid]{
				X:        b.Position[0],
				Y:        b.Position[1],
				UserData: b,
			}
			qtree.Insert(positions[i])
		}
		scanInfo := dbscan.DBScan(qtree, positions, 4, float64(Flock.BoidSettings.Perception))

		for i, cluster := range scanInfo.Clusters {
			imd.Color = ClusterColors[i%len(ClusterColors)]
			for _, b := range cluster {
				imd.Circle(5, 0)
				imd.Push(pixel.V(b.UserData.X(), b.UserData.Y()))
			}
		}

		for _, b := range scanInfo.Noise {
			imd.Color = colornames.Aqua
			imd.Circle(5, 0)
			imd.Push(pixel.V(b.UserData.X(), b.UserData.Y()))
		}

		for _, b := range Flock.Boids {
			imd.Color = colornames.Blue
			imd.Circle(2, 0)
			imd.Push(pixel.V(b.Position.X(), b.Position.Y()))

			// imd.Color = colornames.Blue
			// imd.Circle(float64(Flock.BoidSettings.Perception), 1)
			// imd.Push(pixel.V(b.Position.X(), b.Position.Y()))
		}

		// Draw the quadtree
		DrawQT(imd, qtree)

		imd.Draw(win)
		win.Update()
	}
}

func Execute() {
	pixelgl.Run(run)
}
