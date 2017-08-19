package gtfs

import (
	"github.com/paulmach/go.geo"
)

const (
	tbStops     = "stops"
	tbTrips     = "trips"
	tbStopTimes = "stoptimes"
)

type Stop struct {
	id string
	geo.Point
}

func (stop *Stop) Id() string {
	return stop.id
}
