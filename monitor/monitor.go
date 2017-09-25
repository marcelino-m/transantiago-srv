package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-redis/redis"
	"github.com/marcelino-m/transantiago-srv/fetcher"
	"github.com/marcelino-m/transantiago-srv/gtfs"
	"github.com/paulmach/go.geo"
)

func monitor(httpc *http.Client, s *gtfs.Stop, redisc *redis.Client, q chan struct{}) {

	for {
		q <- struct{}{}
		for {
			buses, err := fetcher.FetchStopData(s.Id(), httpc)
			if err != nil {
				continue
			}

			<-q

			for _, b := range buses {
				busKey := fmt.Sprintf("bus:%s", b.Id())
				pto, err := bi.Position(b)

				if err != nil {
					log.Println(err)
					continue
				}

				var ptoll *geo.Point

				if pto != nil {
					ptoll = pto.Clone()
					ptoll.Transform(geo.Mercator.Inverse)
				} else {
					ptoll = geo.NewPoint(0, 0)
				}

				m := map[string]interface{}{
					"stop":  s.Id(),
					"dist":  b.DistToStop(),
					"route": b.Route(),
					"lon":   ptoll.Lng(),
					"lat":   ptoll.Lat(),
				}

				redisc.HMSet(busKey, m)
			}

			break
		}
	}
}
