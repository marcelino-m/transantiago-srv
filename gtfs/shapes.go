package gtfs

import (
	"github.com/paulmach/go.geo"
)

type Shape struct {
	geo.Path
	lenShape float64
	dir      Direction
}

// NewShape ...
func NewShape() *Shape {
	return &Shape{
		Path:     *geo.NewPath(),
		lenShape: -1,
		dir:      -1,
	}
}

//  Measure length of shape
func (shape *Shape) Distance() float64 {
	if shape.lenShape == -1 {
		shape.lenShape = shape.Path.Distance()
		return shape.lenShape
	} else {
		return shape.lenShape
	}
}

func (shape *Shape) Direction() Direction {
	return shape.dir
}
