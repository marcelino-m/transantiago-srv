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

				var lon, lat float64 = 0, 0
				if pto != nil {
					c := pto.Clone()
					c.Transform(geo.Mercator.Inverse)
					lon = c.Lng()
					lat = c.Lat()
				}

				m := map[string]interface{}{
					"stop":  s.Id(),
					"dist":  b.DistToStop(),
					"lon":   lon,
					"lat":   lat,
					"route": b.Route(),
				}

				redisc.HMSet(busKey, m)
			}

			break
		}
	}
}
