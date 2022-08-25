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

func (b Boid) Align(boids *[]*Boid) mgl64.Vec2 {
	steering := mgl64.Vec2{}
	total := 0
	for _, ob := range *boids {
		d := utils.Distance(b.Position, ob.Position)
		if ob.id != b.id && d < float64(utils.Perception) {
			steering = steering.Add(ob.Velocity)
			total++
		}
	}

	if total > 0 {
		steering = utils.Div(steering, float64(total))
		steering = utils.SetMag(steering, utils.MaxSpeed)
		steering = steering.Sub(b.Velocity)
		steering = utils.Limit(steering, utils.MaxForce)
	}

	return steering
}

func (b Boid) Cohesion(boids *[]*Boid) mgl64.Vec2 {
	steering := mgl64.Vec2{}
	total := 0
	for _, ob := range *boids {
		d := utils.Distance(b.Position, ob.Position)
		if ob.id != b.id && d < float64(utils.Perception) {
			steering = steering.Add(ob.Position)
			total++
		}
	}

	if total > 0 {
		steering = utils.Div(steering, float64(total))
		steering = steering.Sub(b.Position)
		steering = utils.SetMag(steering, utils.MaxSpeed)
		steering = steering.Sub(b.Velocity)
		steering = utils.Limit(steering, utils.MaxForce)
	}

	return steering
}

func (b Boid) Separation(boids *[]*Boid) mgl64.Vec2 {
	steering := mgl64.Vec2{}
	total := 0
	for _, ob := range *boids {
		d := utils.Distance(b.Position, ob.Position)
		if ob.id != b.id && d < float64(utils.Perception) {
			diff := b.Position.Sub(ob.Position)
			diff = utils.Div(diff, d)
			steering = steering.Add(diff)
			total++
		}
	}

	if total > 0 {
		steering = utils.Div(steering, float64(total))
		steering = utils.SetMag(steering, utils.MaxSpeed)
		steering = steering.Sub(b.Velocity)
		steering = utils.Limit(steering, utils.MaxForce)
	}

	return steering
}

func (b *Boid) Flock(boids *[]*Boid) {
	alignment := b.Align(boids)
	cohesion := b.Cohesion(boids)
	separation := b.Separation(boids)

	b.Acceleration = b.Acceleration.Add(alignment)
	b.Acceleration = b.Acceleration.Add(cohesion)
	b.Acceleration = b.Acceleration.Add(separation)
}

func (b *Boid) Update() {
	b.Position = b.Position.Add(b.Velocity)
	b.Velocity = b.Velocity.Add(b.Acceleration)
	b.Velocity = utils.Limit(b.Velocity, utils.MaxSpeed)
	b.Acceleration = b.Acceleration.Mul(0)
}
