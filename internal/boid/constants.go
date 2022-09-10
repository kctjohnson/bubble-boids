package boid

const (
	BoidCount         = 500
	ScatterCounterCap = 600
)

type EdgeMode int

const (
	EDGE_AVOID EdgeMode = iota
	EDGE_WARP
)

type BoidSettings struct {
	EdgeMode           EdgeMode
	MaxAlignmentForce  float64
	MaxCohesionForce   float64
	MaxSeparationForce float64
	MaxSpeed           float64
	Perception         int
}

func NewBoidSettings() *BoidSettings {
	return &BoidSettings{
		EdgeMode:           EDGE_AVOID,
		MaxAlignmentForce:  0.3,
		MaxCohesionForce:   0.25,
		MaxSeparationForce: 0.35,
		MaxSpeed:           3.5,
		Perception:         20,
	}
}
