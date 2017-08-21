package gtfs

import (
	"time"
)

type BusDat struct {
	id          string
	route       string
	arrivaltime string
	dist        string
	goToStop    string
}

// NewBus create a new bus
func NewBus(id, route, arrtime, dist string) *BusDat {
	return &BusDat{id: id, route: route}
}

func (bus *BusDat) Id() string {
	return bus.id
}

func (bus *BusDat) Route() string {
	return bus.route
}

func (bus *BusDat) GoingToStop() string {
	return bus.goToStop
}

func (bus *BusDat) ArrivalTime() time.Duration {
	return time.Second * 0
}

func (bus *BusDat) DistToStop() float64 {
	return 0
}

func (bus *BusDat) SetArrivalTime(art string) {
	bus.arrivaltime = art
}

func (bus *BusDat) SetDistToStop(dst string) {
	bus.dist = dst
}

func (bus *BusDat) SetGoingToStop(stopid string) {
	bus.goToStop = stopid
}
