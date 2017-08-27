package bi

import (
	"math"

	Gtfs "github.com/marcelino-m/transantiago-srv/gtfs"
	"github.com/paulmach/go.geo"
)

type Cicle struct {
	going      *Gtfs.Shape
	regres     *Gtfs.Shape
	stop2shape map[string]*Gtfs.Shape
}

type Bi struct {
	routecicle map[string]*Cicle
	allstops   map[string]*Gtfs.Stop
	allroutes  map[string]*Gtfs.Route
	gtfs       *Gtfs.Gtfs
}

// NewCicle ...
func NewCicle() *Cicle {
	c := &Cicle{
		stop2shape: make(map[string]*Gtfs.Shape),
	}
	return c
}

// NewBi ...
func NewBi(gtfsdb string) (*Bi, error) {
	gtfs, err := Gtfs.Connect(gtfsdb)
	defer gtfs.Close()

	if err != nil {
		return nil, err
	}

	bi := Bi{
		routecicle: make(map[string]*Cicle),
		allstops:   make(map[string]*Gtfs.Stop),
		allroutes:  make(map[string]*Gtfs.Route),
	}

	bi.allstops, err = gtfs.AllStops()
	if err != nil {
		return nil, err
	}

	bi.allroutes, err = gtfs.Routes()
	if err != nil {
		return nil, err
	}

	return &bi, nil

}

//  InitializeBi Initialize internal state of Bi, this operation
//  take a lot of  time to finish, for this reason is splited away
//  so we can postpone until the last minute
func (bi *Bi) InitializeBi() error {

	var err error = nil
	gtfs := bi.gtfs

	for _, r := range bi.allroutes {
		cicle := NewCicle()

		cicle.going, err = gtfs.Shape(r, Gtfs.Going)
		if err != nil {
			return err
		}

		cicle.regres, err = gtfs.Shape(r, Gtfs.Regress)
		if err != nil {
			return err
		}

		stopg, err := gtfs.StopsByRoute(r, Gtfs.Going, bi.allstops)
		if err != nil {
			return err
		}

		for _, s := range stopg {
			cicle.stop2shape[s.Id()] = cicle.going
		}

		stopr, err := gtfs.StopsByRoute(r, Gtfs.Regress, bi.allstops)
		if err != nil {
			return err
		}

		for _, s := range stopr {
			cicle.stop2shape[s.Id()] = cicle.regres
		}

		bi.routecicle[r.Id()] = cicle
	}

	return nil
}

//  ...
func (bi *Bi) Shape(route *Gtfs.Route, stop *Gtfs.Stop) *Gtfs.Shape {
	shape := bi.routecicle[route.Id()].stop2shape[stop.Id()]
	return shape
}

//  Get Stop from stopid
func (bi *Bi) Stop(stopid string) *Gtfs.Stop {
	stop, ok := bi.allstops[stopid]
	if !ok {
		return nil
	}

	return stop
}

//  Get all Stop from stopid
func (bi *Bi) AllStop() map[string]*Gtfs.Stop {
	return bi.allstops
}

//  Get Route from routeid
func (bi *Bi) Route(routeid string) *Gtfs.Route {
	r, ok := bi.allroutes[routeid]
	if !ok {
		return nil
	}

	return r
}

//  Deduce position from bus metadata
func (bi *Bi) Position(bus *Gtfs.BusDat) *geo.Point {
	stop := bi.Stop(bus.Id())
	route := bi.Route(bus.Route())
	if stop == nil || route == nil {
		return nil
	}

	shape := bi.Shape(route, stop)
	dis2stop := bus.DistToStop()
	length := shape.Length()
	eps := 10.0

	for i := 0; i < length; i++ {
		dis := stop.DistanceFrom(shape.GetAt(i))
		delta := dis2stop - dis
		if math.Abs(delta) <= eps {
			return shape.GetAt(i)
		} else if delta < 0 {
			continue
		}

		line := geo.NewLine(shape.GetAt(i-1), shape.GetAt(i))
		step := 0.5
		t := step
		for {
			// binary search
			p := line.Interpolate(t)
			dis = stop.DistanceFrom(p)
			delta = dis2stop - dis
			if math.Abs(delta) <= eps {
				return p
			} else if delta < 0 {
				step /= 2
				t += step
			} else {
				step /= 2
				t -= step
			}
		}

	}

	return nil
}
