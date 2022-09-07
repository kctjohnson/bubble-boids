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

func GetVecAngle(vec mgl64.Vec2) float64 {
	return math.Atan2(vec[1], vec[0]) * (180 / math.Pi)
}

func GetAngleBetween(v1 mgl64.Vec2, v2 mgl64.Vec2) float64 {
	delta := v2.Sub(v1)
	angleBetween := math.Atan2(delta.Y(), delta.X())
	if angleBetween < 0 {
		angleBetween += 2 * math.Pi
	}
	return angleBetween * (180 / math.Pi)
}

func GetPointFromAngle(offset mgl64.Vec2, length float64, angle float64) mgl64.Vec2 {
	radAngle := angle * math.Pi / 180

	var temp mgl64.Vec2
	temp[0] = length*math.Cos(radAngle) + offset[0]
	temp[1] = length*math.Sin(radAngle) + offset[1]
	return temp
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
