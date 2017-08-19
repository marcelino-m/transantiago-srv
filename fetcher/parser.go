package fetcher

import (
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Bus struct {
	id          string
	service     string
	arrivaltime string
	dist        string
}

func (bus *Bus) Id() string {
	return bus.id
}

func (bus *Bus) Service() string {
	return bus.service
}

func (bus *Bus) ArrivalTime() time.Duration {
	return time.Second * 0
}

func (bus *Bus) DistToStop() float32 {
	return 0
}

func Parser(doc *goquery.Document) []Bus {

	buses := []Bus{}
	var currService string
	var re = regexp.MustCompile(`(\r\n|\n|\r)`)
	table := doc.Find("div > div > table").Not(".cabecera4")
	table.Find("tr").Each(func(itr int, node *goquery.Selection) {
		// skip header
		if itr == 0 {
			return
		}

		nodestd := node.Find("td")
		// Parse table
		switch ntd := nodestd.Length(); ntd {
		case 4:
			currService = re.ReplaceAllString(strings.Trim(nodestd.Eq(0).Text(), "\n\r "), "")
			id := re.ReplaceAllString(strings.Trim(nodestd.Eq(1).Text(), "\n\r "), "")
			arrt := re.ReplaceAllString(strings.Trim(nodestd.Eq(2).Text(), "\n\r "), "")
			dist := re.ReplaceAllString(strings.Trim(nodestd.Eq(3).Text(), "\n\r "), "")

			bus := Bus{
				id:          id,
				arrivaltime: arrt,
				dist:        dist,
				service:     currService,
			}
			buses = append(buses, bus)
		case 3:
			id := re.ReplaceAllString(strings.Trim(nodestd.Eq(0).Text(), "\n\r "), "")
			arrt := re.ReplaceAllString(strings.Trim(nodestd.Eq(1).Text(), "\n\r "), "")
			dist := re.ReplaceAllString(strings.Trim(nodestd.Eq(2).Text(), "\n\r "), "")

			bus := Bus{
				id:          id,
				arrivaltime: arrt,
				dist:        dist,
				service:     currService,
			}

			buses = append(buses, bus)
		case 2:
			// Bus out of service
		case 1:
			// end of data

		}
	})

	return buses
}
