package gtfs

const (
	tbStops     = "stops"
	tbTrips     = "trips"
	tbStopTimes = "stoptimes"
)

type Stop struct {
	id  string
	lon float64
	lat float64
}

func (stop *Stop) Id() string {
	return stop.id
}

func (stop *Stop) Lon() float64 {
	return stop.lon
}

func (stop *Stop) Lat() float64 {
	return stop.lat
}
