package boid

import (
	"math/rand"

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
	newBoid.Velocity = mathutil.SetMag(newBoid.Velocity, rand.Float64()*8-4) // 2 to 4
	newBoid.Acceleration = mgl64.Vec2{0.0, 0.0}
	return newBoid
}

// Makes sure the boid can't go outside the bounds of the screen
func (b *Boid) Edges(screenWidth float64, screenHeight float64) {
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
		alignment = mathutil.Limit(alignment, MaxForce)

		// Cohesion
		cohesion = mathutil.Div(cohesion, float64(total))
		cohesion = cohesion.Sub(b.Position)
		cohesion = mathutil.SetMag(cohesion, MaxSpeed)
		cohesion = cohesion.Sub(b.Velocity)
		cohesion = mathutil.Limit(cohesion, MaxForce)

		// Separation
		separation = mathutil.Div(separation, float64(total))
		separation = mathutil.SetMag(separation, MaxSpeed)
		separation = separation.Sub(b.Velocity)
		separation = mathutil.Limit(separation, MaxForce)
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
