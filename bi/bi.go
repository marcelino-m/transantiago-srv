package bi

import (
	"errors"
	"fmt"

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

	if err != nil {
		return nil, err
	}

	bi := Bi{
		routecicle: make(map[string]*Cicle),
		allstops:   make(map[string]*Gtfs.Stop),
		allroutes:  make(map[string]*Gtfs.Route),
		gtfs:       &gtfs,
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

	gtfs := bi.gtfs
	var err error
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

//  Deduce position from bus metadata, asume disatance informed by
//  transantiago is along shape
func (bi *Bi) Position(bus *Gtfs.BusDat) (*geo.Point, error) {

	stop := bi.Stop(bus.GoingToStop())
	if stop == nil {
		return nil, errors.New(fmt.Sprintf("Stop code not found %s", bus.GoingToStop()))
	}

	route := bi.Route(bus.Route())
	if route == nil {
		return nil, errors.New(fmt.Sprintf("Route code not found %s", bus.Route()))
	}

	shape := bi.Shape(route, stop)
	if shape == nil {
		return nil, errors.New(fmt.Sprintf("Shape not found for stop %s and route %s", stop.Id(), route.Id()))
	}

	stopDist := shape.Measure(&stop.Point)
	rel := (stopDist - bus.DistToStop()) / shape.Distance()

	return shape.Interpolate(rel), nil

}
