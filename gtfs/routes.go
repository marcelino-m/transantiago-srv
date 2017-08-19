package gtfs

type Route struct {
	id    string
	sname string
	lname string
}

//  Get id
func (r *Route) Id() string {
	return r.id
}

func NewRoute(id, sname, lname string) *Route {
	return &Route{
		id: id, sname: sname, lname: lname,
	}
}
