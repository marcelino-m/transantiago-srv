package main

import (
	// "net/http"

	"fmt"
	// "github.com/marcelino-m/transantiago-srv/fetcher"
	// "time"
	"github.com/marcelino-m/transantiago-srv/gtfs"
)

func main() {

	// tr := &http.Transport{
	// 	MaxIdleConns:       10,
	// 	IdleConnTimeout:    30 * time.Second,
	// 	DisableCompression: true,
	// }

	// client := &http.Client{Transport: tr}
	// buses, err := fetcher.FetchStopData("PA11", client)

	// if err != nil {
	// 	fmt.Printf("%+v\n", err)
	// }

	// for _, bus := range buses {
	// 	fmt.Printf("%+v\n", bus)
	// }

	con, err := gtfs.Connect("/home/marcelo/lab/tase/gtfs/gtfs.db")
	if err != nil {
		fmt.Printf("%+v\n", err)
		return
	}

	// stops, _ := con.Stops()
	// for _, s := range stops {
	// 	fmt.Printf("%+v\n", s)
	// }

	route := gtfs.NewRoute("G08", "S. name", "L. name")

	sh, _ := con.Shape(route, gtfs.Regress)
	length := sh.Length()
	fmt.Printf("%+v\n", length)

}
