package boid

import (
	"math"
	"math/rand"
	"sort"
	"sync"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/kctjohnson/bubble-boids/internal/dbscan"
	"github.com/kctjohnson/bubble-boids/internal/mathutil"
	"github.com/kctjohnson/bubble-boids/internal/quadtree"
)

type Flock struct {
	BoidSettings *BoidSettings
	Boids        []Boid
	QuadTree     quadtree.QuadTree[Boid]
	ScanInfo     dbscan.DBScanInfo[Boid]

	scatterCounter int // Starts at 0, when it hits ScatterCounterCap, all of the boids are scattered
}

func NewFlock(screenWidth float64, screenHeight float64) *Flock {
	// Initialize the settings of the flock
	newBoidSettings := NewBoidSettings()

	// Create the flock slice and initialize each boid
	newBoidSlice := make([]Boid, 0)
	numberOfBoids := BoidCount
	for i := 0; i < numberOfBoids; i++ {
		newBoidSlice = append(newBoidSlice, NewBoid(i, screenWidth, screenHeight, newBoidSettings))
	}

	return &Flock{
		Boids:          newBoidSlice,
		BoidSettings:   newBoidSettings,
		scatterCounter: 0,
	}
}

func (f *Flock) UpdateCluster(cluster []mathutil.Point[Boid], screenWidth float64, screenHeight float64, wg *sync.WaitGroup) {
	qtree := quadtree.NewQuadTree(mathutil.Rectangle[Boid]{X: 0, Y: 0, W: screenWidth, H: screenHeight}, 4)
	for _, bPoint := range cluster {
		qtree.Insert(bPoint)
	}

	for _, bPoint := range cluster {
		b := bPoint.UserData.Self()
		if f.scatterCounter >= ScatterCounterCap {
			// Randomize velocity and acceleration
			b.Velocity = mathutil.RandomVec2(-f.BoidSettings.MaxSpeed, f.BoidSettings.MaxSpeed)
			b.Acceleration = mathutil.RandomVec2(-f.BoidSettings.MaxSpeed, f.BoidSettings.MaxSpeed)
		} else {
			b.Edges(screenWidth, screenHeight)

			fPerc := float64(f.BoidSettings.Perception)
			inRangeOfBoid := qtree.Query(mathutil.Rectangle[Boid]{
				X: b.Position[0] - fPerc,
				Y: b.Position[1] - fPerc,
				W: fPerc * 2,
				H: fPerc * 2,
			})
			b.Flock(inRangeOfBoid)
		}
		b.Update()
		f.Boids[b.ID()] = b
	}

	wg.Done()
}

func (f *Flock) Update(screenWidth float64, screenHeight float64) {
	// Create the quadtree map of the boids
	positions := make([]mathutil.Point[Boid], len(f.Boids))
	qtree := quadtree.NewQuadTree(mathutil.Rectangle[Boid]{X: 0, Y: 0, W: screenWidth, H: screenHeight}, 4)
	for i, b := range f.Boids {
		positions[i] = mathutil.Point[Boid]{X: b.Position.X(), Y: b.Position.Y(), UserData: b}
		qtree.Insert(positions[i])
	}
	f.QuadTree = qtree

	// Scan the boids into groups using the quadtree
	var scanInfo dbscan.DBScanInfo[Boid] = dbscan.DBScan(qtree, positions, 4, float64(f.BoidSettings.Perception))
	f.ScanInfo = scanInfo

	// For each cluster start a new update, as well as one for the noise
	f.scatterCounter++
	var wg sync.WaitGroup
	wg.Add(1)
	go f.UpdateCluster(scanInfo.Noise, screenWidth, screenHeight, &wg)
	for _, c := range scanInfo.Clusters {
		wg.Add(1)
		go f.UpdateCluster(c, screenWidth, screenHeight, &wg)
	}

	wg.Wait()

	// for i := range f.Boids {
	// 	wg.Add(1)
	// 	go func(i int, b *Boid) {
	// 		if f.scatterCounter >= ScatterCounterCap {
	// 			// Randomize velocity and acceleration
	// 			b.Velocity = mathutil.RandomVec2(-f.BoidSettings.MaxSpeed, f.BoidSettings.MaxSpeed)
	// 			b.Acceleration = mathutil.RandomVec2(-f.BoidSettings.MaxSpeed, f.BoidSettings.MaxSpeed)
	// 		} else {
	// 			b.Edges(screenWidth, screenHeight)
	//
	// 			fPerc := float64(f.BoidSettings.Perception)
	// 			inRangeOfBoid := qtree.Query(mathutil.Rectangle[Boid]{
	// 				X: b.Position[0] - fPerc,
	// 				Y: b.Position[1] - fPerc,
	// 				W: fPerc * 2,
	// 				H: fPerc * 2,
	// 			})
	// 			b.Flock(inRangeOfBoid)
	// 		}
	// 		b.Update()
	// 		wg.Done()
	// 	}(i, &f.Boids[i])
	// }

	if f.scatterCounter >= ScatterCounterCap {
		f.scatterCounter = 0
	}

	wg.Wait()
}

func (f *Flock) Scatter() {
	f.scatterCounter = ScatterCounterCap
}

type Boid struct {
	id           int
	Position     mgl64.Vec2
	Velocity     mgl64.Vec2
	Acceleration mgl64.Vec2

	// Owned by the flock structure, but the boids read these
	boidSettings *BoidSettings
}

func NewBoid(id int, screenWidth float64, screenHeight float64, boidSettings *BoidSettings) Boid {
	newBoid := Boid{}
	newBoid.id = id
	newBoid.boidSettings = boidSettings
	newBoid.Position = mgl64.Vec2{
		rand.Float64() * screenWidth,
		rand.Float64() * screenHeight,
	}
	newBoid.Velocity = mgl64.Vec2{rand.Float64(), rand.Float64()}
	newBoid.Velocity = mathutil.SetMag(newBoid.Velocity, mathutil.RandRange(-boidSettings.MaxSpeed, boidSettings.MaxSpeed))
	newBoid.Acceleration = mgl64.Vec2{0.0, 0.0}
	return newBoid
}

func (b Boid) ID() int {
	return b.id
}

func (b Boid) X() float64 {
	return b.Position[0]
}

func (b Boid) Y() float64 {
	return b.Position[1]
}

func (b Boid) Self() Boid {
	return b
}

// Makes sure the boid can't go outside the bounds of the screen
func (b *Boid) Edges(screenWidth float64, screenHeight float64) {
	if b.boidSettings.EdgeMode == EDGE_WARP {
		if b.Position.X() > screenWidth {
			b.Position[0] = 0
		} else if b.Position.X() < 0 {
			b.Position[0] = screenWidth
		}

		if b.Position.Y() > screenHeight {
			b.Position[1] = 0
		} else if b.Position.Y() < 0 {
			b.Position[1] = screenHeight
		}
	} else if b.boidSettings.EdgeMode == EDGE_AVOID {
		// Used to check which way the boid needs to move
		left := b.Position.X() < screenWidth/2
		top := b.Position.Y() < screenHeight/2

		// Get the heighest force being applies to the boid
		forces := []float64{
			b.boidSettings.MaxAlignmentForce,
			b.boidSettings.MaxCohesionForce,
			b.boidSettings.MaxSeparationForce,
		}
		sort.Slice(forces, func(i, j int) bool { return forces[i] > forces[j] })
		heighestForce := forces[0] * 10

		// Calculate the opposite force vector that will move the boid away from the walls
		var force mgl64.Vec2
		if left {
			force[0] = heighestForce / math.Max(b.Position.X(), 0.1)
		} else {
			force[0] = -(heighestForce / (screenWidth - math.Min(b.Position.X(), screenWidth-1)))
		}

		if top {
			force[1] = heighestForce / math.Max(b.Position.Y(), 0.1)
		} else {
			force[1] = -(heighestForce / (screenHeight - math.Min(b.Position.Y(), screenHeight-1)))
		}

		b.Acceleration = b.Acceleration.Add(force)
	}
}

func (b Boid) BoidLogic(boids []mathutil.Point[Boid]) mgl64.Vec2 {
	total := 0
	alignment := mgl64.Vec2{}
	cohesion := mgl64.Vec2{}
	separation := mgl64.Vec2{}
	for _, obi := range boids {
		ob := obi.UserData.Self()
		distance := mathutil.Distance(b.Position, ob.Position)
		if ob.id != b.id && distance < float64(b.boidSettings.Perception) {
			// Alignment
			alignment = alignment.Add(ob.Velocity)

			// Cohesion
			cohesion = cohesion.Add(ob.Position)

			// Separation
			diff := b.Position.Sub(ob.Position)
			diff = mathutil.Div(diff, distance)
			separation = separation.Add(diff)
			total++
		}
	}

	if total > 0 {
		// Alignment
		alignment = mathutil.Div(alignment, float64(total))
		alignment = mathutil.SetMag(alignment, b.boidSettings.MaxSpeed)
		alignment = alignment.Sub(b.Velocity)
		alignment = mathutil.Limit(alignment, b.boidSettings.MaxAlignmentForce)

		// Cohesion
		cohesion = mathutil.Div(cohesion, float64(total))
		cohesion = cohesion.Sub(b.Position)
		cohesion = mathutil.SetMag(cohesion, b.boidSettings.MaxSpeed)
		cohesion = cohesion.Sub(b.Velocity)
		cohesion = mathutil.Limit(cohesion, b.boidSettings.MaxCohesionForce)

		// Separation
		separation = mathutil.Div(separation, float64(total))
		separation = mathutil.SetMag(separation, b.boidSettings.MaxSpeed)
		separation = separation.Sub(b.Velocity)
		separation = mathutil.Limit(separation, b.boidSettings.MaxSeparationForce)
	}

	// Add them all up to get the final acceleration force
	force := alignment.Add(cohesion).Add(separation)
	return force
}

func (b *Boid) Flock(boids []mathutil.Point[Boid]) {
	force := b.BoidLogic(boids)
	b.Acceleration = b.Acceleration.Add(force)
}

func (b *Boid) Update() {
	b.Position = b.Position.Add(b.Velocity)
	b.Velocity = b.Velocity.Add(b.Acceleration)
	b.Velocity = mathutil.Limit(b.Velocity, b.boidSettings.MaxSpeed)
	b.Acceleration = b.Acceleration.Mul(0)
}
