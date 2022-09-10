package quadtree

import (
	"github.com/kctjohnson/bubble-boids/internal/mathutil"
)

type QuadTree[T any] struct {
	Boundary mathutil.Rectangle[T]
	Points   []mathutil.Point[T]
	Capacity int
	Length   int
	Divided  bool

	Nodes []QuadTree[T]
}

func NewQuadTree[T any](boundary mathutil.Rectangle[T], capacity int) QuadTree[T] {
	return QuadTree[T]{
		Boundary: boundary,
		Points:   make([]mathutil.Point[T], 0, capacity),
		Capacity: capacity,
		Divided:  false,
		Length:   0,
		Nodes:    make([]QuadTree[T], 4, 4),
	}
}

func (qt *QuadTree[T]) Insert(p mathutil.Point[T]) bool {
	if !qt.Boundary.Contains(p) {
		return false
	}

	if !qt.Divided && qt.Length < qt.Capacity {
		qt.Points = append(qt.Points, p)
		qt.Length++
		return true
	} else {
		if !qt.Divided {
			qt.Subdivide()
		}

		for i := range qt.Nodes {
			if qt.Nodes[i].Insert(p) {
				return true
			}
		}
	}

	return false
}

func (qt *QuadTree[T]) Query(boundary mathutil.Rectangle[T]) []mathutil.Point[T] {
	pointsInRange := make([]mathutil.Point[T], 0)

	if !qt.Boundary.Intersects(boundary) {
		return pointsInRange
	}

	for _, p := range qt.Points {
		if boundary.Contains(p) {
			pointsInRange = append(pointsInRange, p)
		}
	}

	if qt.Divided {
		for i := range qt.Nodes {
			pointsInRange = append(pointsInRange, qt.Nodes[i].Query(boundary)...)
		}
	}

	return pointsInRange
}

func (qt *QuadTree[T]) Subdivide() {
	x := qt.Boundary.X
	y := qt.Boundary.Y
	w := qt.Boundary.W
	h := qt.Boundary.H
	subW := w / 2
	subH := h / 2

	qt.Nodes[0] = NewQuadTree(mathutil.Rectangle[T]{X: x + subW, Y: y + subH, W: subW, H: subH}, qt.Capacity)
	qt.Nodes[1] = NewQuadTree(mathutil.Rectangle[T]{X: x, Y: y + subH, W: subW, H: subH}, qt.Capacity)
	qt.Nodes[2] = NewQuadTree(mathutil.Rectangle[T]{X: x + subW, Y: y, W: subW, H: subH}, qt.Capacity)
	qt.Nodes[3] = NewQuadTree(mathutil.Rectangle[T]{X: x, Y: y, W: subW, H: subH}, qt.Capacity)

	qt.Divided = true
}
