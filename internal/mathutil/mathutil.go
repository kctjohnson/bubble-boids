package mathutil

import (
	"math"
	"math/rand"

	"github.com/go-gl/mathgl/mgl64"
)

type Position[T any] interface {
	ID() int
	X() float64
	Y() float64
	Self() T
}

type Point[T any] struct {
	X, Y     float64
	UserData Position[T]
}

type Rectangle[T any] struct {
	X, Y, W, H float64
}

func (r Rectangle[T]) Contains(p Point[T]) bool {
	return (p.X >= r.X && p.X <= r.X+r.W && p.Y >= r.Y && p.Y <= r.Y+r.H)
}

func (r Rectangle[T]) Intersects(other Rectangle[T]) bool {
	return (r.X < other.X+other.W && r.X+r.W > other.X && r.Y < other.Y+other.H && r.Y+r.H > other.Y)
}

func RandomVec2(min float64, max float64) mgl64.Vec2 {
	x := RandRange(min, max)
	y := RandRange(min, max)
	return mgl64.Vec2{x, y}
}

func MagSq(vec mgl64.Vec2) float64 {
	return vec[0]*vec[0] + vec[1]*vec[1]
}

func Mag(vec mgl64.Vec2) float64 {
	return math.Sqrt(MagSq(vec))
}

func SetMag(vec mgl64.Vec2, mag float64) mgl64.Vec2 {
	return vec.Normalize().Mul(mag)
}

func Distance(vec1 mgl64.Vec2, vec2 mgl64.Vec2) float64 {
	return Mag(vec2.Sub(vec1))
}

func Div(vec mgl64.Vec2, scalar float64) mgl64.Vec2 {
	if scalar == 0 {
		panic("Divide by zero error")
	}

	vec[0] /= scalar
	vec[1] /= scalar

	return vec
}

func Limit(vec mgl64.Vec2, max float64) mgl64.Vec2 {
	mSq := MagSq(vec)
	if mSq > max*max {
		vec = Div(vec, math.Sqrt(mSq))
		vec = vec.Mul(max)
	}
	return vec
}

func RandRange(min float64, max float64) float64 {
	return rand.Float64()*(max-min) + min
}
