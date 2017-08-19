package gtfs

import (
	"time"
)

type Bus struct {
	id          string
	route       string
	arrivaltime string
	dist        string
}

// NewBus create a new bus
func NewBus(id, route, arrtime, dist string) *Bus {
	return &Bus{id: id, route: route}
}

func (bus *Bus) Id() string {
	return bus.id
}

func (bus *Bus) Route() string {
	return bus.route
}

func (bus *Bus) ArrivalTime() time.Duration {
	return time.Second * 0
}

func (bus *Bus) DistToStop() float32 {
	return 0
}

func (bus *Bus) SetArrivalTime(art string) {
	bus.arrivaltime = art
}

func (bus *Bus) SetDistToStop(dst string) {
	bus.dist = dst
}
