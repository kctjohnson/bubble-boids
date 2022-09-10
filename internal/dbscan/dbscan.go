package dbscan

import (
	"github.com/kctjohnson/bubble-boids/internal/mathutil"
	"github.com/kctjohnson/bubble-boids/internal/quadtree"
)

type PointState int

const (
	Undefined PointState = iota
	Core
	Border
	Noise
)

type DBScanInfo[T any] struct {
	Clusters map[int][]mathutil.Point[T]
	Noise    []mathutil.Point[T]
}

func DBScan[T any](qtree quadtree.QuadTree[T], points []mathutil.Point[T], minDensity int, perception float64) DBScanInfo[T] {
	visited := make(map[int]PointState, len(points))
	markedForNoise := make([]mathutil.Point[T], 0, len(points))
	info := DBScanInfo[T]{
		Clusters: make(map[int][]mathutil.Point[T]),
		Noise:    make([]mathutil.Point[T], 0, len(points)),
	}

	c := 0
	for _, p := range points {
		// We've already visited this point, skip it
		if visited[p.UserData.ID()] != Undefined {
			continue
		}

		// Get the neighbors for this undefined point and determine if it's core or not
		neighbors := qtree.Query(mathutil.Rectangle[T]{
			X: p.X - perception,
			Y: p.Y - perception,
			W: perception * 2,
			H: perception * 2,
		})
		if len(neighbors) < minDensity {
			visited[p.UserData.ID()] = Noise
			markedForNoise = append(markedForNoise, p)
			continue
		}

		// The last point was core, we're going to expand now into a new cluster
		c++
		info.Clusters[c] = make([]mathutil.Point[T], 0)
		info.Clusters[c] = append(info.Clusters[c], p)
		visited[p.UserData.ID()] = Core
		for i := 0; i < len(neighbors); i++ {
			// Skip if it's the same point
			if neighbors[i].UserData.ID() == p.UserData.ID() {
				continue
			}

			// Change noise to border
			if visited[neighbors[i].UserData.ID()] == Noise {
				visited[neighbors[i].UserData.ID()] = Border
			}

			// Skip if this point is already defined
			if visited[neighbors[i].UserData.ID()] != Undefined {
				continue
			}

			// Add to the current cluster and label it as core
			visited[neighbors[i].UserData.ID()] = Core
			info.Clusters[c] = append(info.Clusters[c], neighbors[i])

			// Get the next neighbors to search
			nextNeighbors := qtree.Query(mathutil.Rectangle[T]{
				X: neighbors[i].X - perception,
				Y: neighbors[i].Y - perception,
				W: perception * 2,
				H: perception * 2,
			})
			for _, n := range nextNeighbors {
				neighbors = append(neighbors, n)
			}
		}

	}

	// Add all noise points to the noise group
	for _, m := range markedForNoise {
		if visited[m.UserData.ID()] == Noise {
			info.Noise = append(info.Noise, m)
		}
	}

	return info
}
