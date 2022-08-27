package boid

import (
	"math"
	"math/rand"
	"sort"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/kctjohnson/bubble-boids/internal/mathutil"
)

type Boid struct {
	id           int
	Position     mgl64.Vec2
	Velocity     mgl64.Vec2
	Acceleration mgl64.Vec2
}

func NewBoid(id int, screenWidth float64, screenHeight float64) *Boid {
	newBoid := new(Boid)
	newBoid.id = id
	newBoid.Position = mgl64.Vec2{
		rand.Float64() * screenWidth,
		rand.Float64() * screenHeight,
	}
	newBoid.Velocity = mgl64.Vec2{rand.Float64(), rand.Float64()}
	newBoid.Velocity = mathutil.SetMag(newBoid.Velocity, mathutil.RandRange(-MaxSpeed, MaxSpeed))
	newBoid.Acceleration = mgl64.Vec2{0.0, 0.0}
	return newBoid
}

// Makes sure the boid can't go outside the bounds of the screen
func (b *Boid) Edges(screenWidth float64, screenHeight float64) {
	// Used to check which way the boid needs to move
	left := b.Position.X() < screenWidth/2
	top := b.Position.Y() < screenHeight/2

	// Get the heighest force being applies to the boid
	forces := []float64{MaxAlignmentForce, MaxCohesionForce, MaxSeparationForce, MaxSpeed}
	sort.Slice(forces, func(i, j int) bool { return forces[i] < forces[j] })
	heighestForce := forces[0] * 10

	// Calculate the opposite force vector that will move the boid away from the walls
	var force mgl64.Vec2
	if left {
		force[0] = heighestForce / math.Max(b.Position.X(), 0.1)
	} else {
		force[0] = -(heighestForce / (screenWidth - math.Min(b.Position.X(), screenWidth - 1)))
	}

	if top {
		force[1] = heighestForce / math.Max(b.Position.Y(), 0.1)
	} else {
		force[1] = -(heighestForce / (screenHeight - math.Min(b.Position.Y(), screenHeight - 1)))
	}

	b.Acceleration = b.Acceleration.Add(force)
}

func (b Boid) BoidLogic(boids *[]*Boid) mgl64.Vec2 {
	total := 0
	alignment := mgl64.Vec2{}
	cohesion := mgl64.Vec2{}
	separation := mgl64.Vec2{}
	for _, ob := range *boids {
		distance := mathutil.Distance(b.Position, ob.Position)
		if ob.id != b.id && distance < float64(Perception) {
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
		alignment = mathutil.SetMag(alignment, MaxSpeed)
		alignment = alignment.Sub(b.Velocity)
		alignment = mathutil.Limit(alignment, MaxAlignmentForce)

		// Cohesion
		cohesion = mathutil.Div(cohesion, float64(total))
		cohesion = cohesion.Sub(b.Position)
		cohesion = mathutil.SetMag(cohesion, MaxSpeed)
		cohesion = cohesion.Sub(b.Velocity)
		cohesion = mathutil.Limit(cohesion, MaxCohesionForce)

		// Separation
		separation = mathutil.Div(separation, float64(total))
		separation = mathutil.SetMag(separation, MaxSpeed)
		separation = separation.Sub(b.Velocity)
		separation = mathutil.Limit(separation, MaxSeparationForce)
	}

	// Add them all up to get the final acceleration force
	force := alignment.Add(cohesion).Add(separation)
	return force
}

func (b *Boid) Flock(boids *[]*Boid) {
	force := b.BoidLogic(boids)
	b.Acceleration = b.Acceleration.Add(force)
}

func (b *Boid) Update() {
	b.Position = b.Position.Add(b.Velocity)
	b.Velocity = b.Velocity.Add(b.Acceleration)
	b.Velocity = mathutil.Limit(b.Velocity, MaxSpeed)
	b.Acceleration = b.Acceleration.Mul(0)
}
