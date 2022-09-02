package mathutil

import (
	"math"
	"math/rand"

	"github.com/go-gl/mathgl/mgl64"
)

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
