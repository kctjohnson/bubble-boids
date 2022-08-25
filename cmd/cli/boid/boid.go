package boid

import (
	"github.com/kctjohnson/bubble-boids/cmd/cli/utils"
	"math/rand"

	"github.com/go-gl/mathgl/mgl64"
)

type Boid struct {
	id           int
	Position     mgl64.Vec2
	Velocity     mgl64.Vec2
	Acceleration mgl64.Vec2
}

func NewBoid(id int, screenWidth int, screenHeight int) *Boid {
	newBoid := new(Boid)
	newBoid.id = id
	newBoid.Position = mgl64.Vec2{
		rand.Float64() * float64(screenWidth),
		rand.Float64() * float64(screenHeight),
	}
	newBoid.Velocity = mgl64.Vec2{rand.Float64(), rand.Float64()}
	newBoid.Velocity = utils.SetMag(newBoid.Velocity, rand.Float64()*8-4) // 2 to 4
	newBoid.Acceleration = mgl64.Vec2{0.0, 0.0}
	return newBoid
}

// Makes sure the boid can't go outside the bounds of the screen
func (b *Boid) Edges(screenWidth int, screenHeight int) {
	if b.Position.X() > float64(screenWidth) {
		b.Position[0] = 0
	} else if b.Position.X() < 0 {
		b.Position[0] = float64(screenWidth)
	}

	if b.Position.Y() > float64(screenHeight) {
		b.Position[1] = 0
	} else if b.Position.Y() < 0 {
		b.Position[1] = float64(screenHeight)
	}
}

func (b Boid) BoidLogic(boids *[]*Boid) mgl64.Vec2 {
	total := 0
	alignment := mgl64.Vec2{}
	cohesion := mgl64.Vec2{}
	separation := mgl64.Vec2{}
	for _, ob := range *boids {
		distance := utils.Distance(b.Position, ob.Position)
		if ob.id != b.id && distance < float64(utils.Perception) {
			// Alignment
			alignment = alignment.Add(ob.Velocity)
			
			// Cohesion
			cohesion = cohesion.Add(ob.Position)

			// Separation
			diff := b.Position.Sub(ob.Position)
			diff = utils.Div(diff, distance)
			separation = separation.Add(diff)
			total++
		}
	}

	if total > 0 {
		// Alignment
		alignment = utils.Div(alignment, float64(total))
		alignment = utils.SetMag(alignment, utils.MaxSpeed)
		alignment = alignment.Sub(b.Velocity)
		alignment = utils.Limit(alignment, utils.MaxForce)

		// Cohesion
		cohesion = utils.Div(cohesion, float64(total))
		cohesion = cohesion.Sub(b.Position)
		cohesion = utils.SetMag(cohesion, utils.MaxSpeed)
		cohesion = cohesion.Sub(b.Velocity)
		cohesion = utils.Limit(cohesion, utils.MaxForce)

		// Separation
		separation = utils.Div(separation, float64(total))
		separation = utils.SetMag(separation, utils.MaxSpeed)
		separation = separation.Sub(b.Velocity)
		separation = utils.Limit(separation, utils.MaxForce)
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
	b.Velocity = utils.Limit(b.Velocity, utils.MaxSpeed)
	b.Acceleration = b.Acceleration.Mul(0)
}
