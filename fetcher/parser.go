package fetcher

import (
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/marcelino-m/transantiago-srv/gtfs"
)

func Parser(doc *goquery.Document) []*gtfs.BusDat {

	buses := []*gtfs.BusDat{}

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

			bus := gtfs.NewBus(id, currService, arrt, dist)
			buses = append(buses, bus)

		case 3:
			id := re.ReplaceAllString(strings.Trim(nodestd.Eq(0).Text(), "\n\r "), "")
			arrt := re.ReplaceAllString(strings.Trim(nodestd.Eq(1).Text(), "\n\r "), "")
			dist := re.ReplaceAllString(strings.Trim(nodestd.Eq(2).Text(), "\n\r "), "")

			bus := gtfs.NewBus(id, currService, arrt, dist)
			buses = append(buses, bus)

		case 2:
			// Bus out of service
		case 1:
			// end of data

		}
	})

	return buses
}
