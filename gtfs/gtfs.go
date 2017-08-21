package gtfs

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/paulmach/go.geo"
)

type Gtfs struct {
	conn *sql.DB
}

// New ...
func Connect(dbpath string) (Gtfs, error) {
	gtfs := Gtfs{}
	db, err := sql.Open("sqlite3", dbpath)
	if err != nil {
		return gtfs, err
	}
	gtfs.conn = db
	return gtfs, nil

}

func (gtfd Gtfs) Close() {
	gtfd.conn.Close()
}

// Stops get stops of transantiago
func (gtfs Gtfs) AllStops() (map[string]*Stop, error) {
	rows, err := gtfs.conn.Query(
		"SELECT stop_id, stop_lat, stop_lon  FROM stops",
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	stops := make(map[string]*Stop)

	for rows.Next() {
		var id string
		var lon, lat float64

		err := rows.Scan(&id, &lat, &lon)
		if err != nil {
			return nil, err
		}

		stop := Stop{id: id}
		stop.SetLat(lat)
		stop.SetLng(lon)
		stop.Transform(geo.Mercator.Project)
		stops[stop.id] = &stop
	}

	return stops, nil
}

func (gtfs Gtfs) StopsByRoute(route *Route, dir Direction, astops map[string]*Stop) ([]*Stop, error) {
	rows, err := gtfs.conn.Query(
		`SELECT stop_id FROM stop_times
         WHERE trip_id =
          (SELECT trip_id FROM trips  WHERE direction_id = ? AND route_id = ? limit 1)
         ORDER BY
          stop_sequence ASC`,
		dir,
		route.Id(),
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	stops := []*Stop{}

	for rows.Next() {
		var stopid string
		err := rows.Scan(&stopid)

		if err != nil {
			return nil, err
		}

		stops = append(stops, astops[stopid])
	}

	return stops, nil
}

// Get a shape  by route and direction
func (gtfs Gtfs) Shape(route *Route, dir Direction) (*Shape, error) {

	rows, err := gtfs.conn.Query(
		`SELECT
           shape_pt_lat, shape_pt_lon FROM shapes WHERE shape_id =
            (SELECT  DISTINCT shape_id FROM trips WHERE direction_id = ?  and route_id = ? LIMIT 1)
           ORDER BY
             shape_pt_sequence ASC`,
		dir,
		route.Id(),
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	shape := Shape{}
	var lat, lon float64

	for rows.Next() {
		err := rows.Scan(&lat, &lon)
		if err != nil {
			rows.Close()
			return nil, err
		}

		p := geo.NewPoint(lon, lat)
		shape.Push(p)
	}

	shape.Transform(geo.Mercator.Project)

	return &shape, nil
}

//  Get all routes
func (gtfs Gtfs) Routes() (map[string]*Route, error) {
	rows, err := gtfs.conn.Query(
		"SELECT route_id, route_short_name, route_long_name  FROM routes",
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	routes := make(map[string]*Route)

	for rows.Next() {
		r := Route{}
		err := rows.Scan(&r.id, &r.sname, &r.lname)
		if err != nil {
			return nil, err
		}

		routes[r.Id()] = &r
	}

	return routes, nil
}
