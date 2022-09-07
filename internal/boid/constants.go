package boid

const (
	BoidCount         = 300
	ScatterCounterCap = 600
	QuadCap           = 10
	FOV               = 90
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
		MaxAlignmentForce:  0.1,
		MaxCohesionForce:   0.1,
		MaxSeparationForce: 0.1,
		MaxSpeed:           5.5,
		Perception:         100,
	}
}
