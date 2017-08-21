package bi

import (
	"math"

	Gtfs "github.com/marcelino-m/transantiago-srv/gtfs"
	"github.com/paulmach/go.geo"
	"log"
)

type Cicle struct {
	going      *Gtfs.Shape
	regres     *Gtfs.Shape
	stop2shape map[string]*Gtfs.Shape
}

type RouteCicle map[string]*Cicle

type Bi struct {
	routecicle RouteCicle
	allstops   map[string]*Gtfs.Stop
	allroutes  map[string]*Gtfs.Route
}

// NewBi ...
func NewBi(gtfsdb string) (*Bi, error) {
	gtfs, err := Gtfs.Connect(gtfsdb)
	defer gtfs.Close()

	if err != nil {
		return nil, err
	}

	bi := Bi{}
	bi.allstops, err = gtfs.AllStops()
	if err != nil {
		return nil, err
	}

	bi.allroutes, err = gtfs.Routes()
	if err != nil {
		return nil, err
	}

	for _, r := range bi.allroutes {
		cicle := Cicle{}

		cicle.going, err = gtfs.Shape(r, Gtfs.Going)
		if err != nil {
			return nil, err
		}

		cicle.regres, err = gtfs.Shape(r, Gtfs.Regress)
		if err != nil {
			return nil, err
		}

		stopg, err := gtfs.StopsByRoute(r, Gtfs.Going, bi.allstops)
		if err != nil {
			return nil, err
		}

		for _, s := range stopg {
			cicle.stop2shape[s.Id()] = cicle.going
		}

		stopr, err := gtfs.StopsByRoute(r, Gtfs.Regress, bi.allstops)
		if err != nil {
			return nil, err
		}

		for _, s := range stopr {
			cicle.stop2shape[s.Id()] = cicle.regres
		}

		bi.routecicle[r.Id()] = &cicle
	}

	return &bi, nil
}

//  ...
func (bi *Bi) GetShape(route *Gtfs.Route, stop *Gtfs.Stop) *Gtfs.Shape {
	shape := bi.routecicle[route.Id()].stop2shape[stop.Id()]
	return shape
}

//  Get Stop from stopid
func (bi *Bi) GetStop(stopid string) *Gtfs.Stop {
	stop, ok := bi.allstops[stopid]
	if !ok {
		return nil
	}

	return stop
}

//  Get Route from routeid
func (bi *Bi) GetRoute(routeid string) *Gtfs.Route {
	r, ok := bi.allroutes[routeid]
	if !ok {
		return nil
	}

	return r
}

//  Deduce position from bus metadata
func (bi *Bi) GetPosition(bus *Gtfs.BusDat) *geo.Point {
	stop := bi.GetStop(bus.Id())
	route := bi.GetRoute(bus.Route())
	if stop == nil || route == nil {
		return nil
	}

	shape := bi.GetShape(route, stop)
	ts := shape.Project(&stop.Point)
	ps := shape.Interpolate(ts)

	var eps float64 = 10

	dis2stop := bus.DistToStop()
	delta := ts / 2
	ft := delta
	count := 100
	for {

		if count > 100 {
			log.Printf("Can't locate bus %s\n,stop: %s, dist to stop %f", bus.Id(), stop.Id(), dis2stop)
			return nil
		}

		bus := shape.Interpolate(ft)
		dis := bus.DistanceFrom(ps)
		if math.Abs(dis-dis2stop) <= eps {
			return bus
		} else if dis < dis2stop {
			delta := delta / 2
			ft = ft + delta
		} else {
			delta := delta / 2
			ft = ft - delta
		}
	}

}
