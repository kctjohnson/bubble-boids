package boid

const (
	BoidCount         = 300
	ScatterCounterCap = 600
)

type BoidSettings struct {
	MaxAlignmentForce  float64
	MaxCohesionForce   float64
	MaxSeparationForce float64
	MaxSpeed           float64
	Perception         int
}

func NewBoidSettings() *BoidSettings {
	return &BoidSettings{
		MaxAlignmentForce:  0.3,
		MaxCohesionForce:   0.3,
		MaxSeparationForce: 0.3,
		MaxSpeed:           2.5,
		Perception:         30,
	}
}
