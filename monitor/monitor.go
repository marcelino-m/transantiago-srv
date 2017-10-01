package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-redis/redis"
	"github.com/marcelino-m/transantiago-srv/fetcher"
	"github.com/marcelino-m/transantiago-srv/gtfs"
	"github.com/paulmach/go.geo"
	"strings"
)

func monitor(s *gtfs.Stop, route string, httpc *http.Client, redisc *redis.Client, q chan struct{}) {

	for {
		q <- struct{}{}
		for {
			buses, err := fetcher.FetchStopData(s.Id(), httpc)
			if err != nil {
				log.Print("Error:", err)
				continue
			}

			<-q

			all := make(map[string]interface{})

			for _, b := range buses {
				if route != "" {
					if !strings.EqualFold(b.Route(), route) {
						continue
					}
				}
				if b.DistToStop() > 1650 {
					continue
				}
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

				all[b.Id()] = fmt.Sprintf(
					"{\"pos\":[%v,%v], \"id\": \"%v\"}",
					ptoll.Lng(),
					ptoll.Lat(),
					busKey,
				)

				_, err = redisc.HMSet(busKey, m).Result()
				if err != nil {
					log.Print("Error inserting bus in redis:", err)
				}

			}

			if len(all) != 0 {
				_, err = redisc.HMSet("curr:buses", all).Result()
				if err != nil {
					log.Print("Erro inserting in redis:", err)
				}
			}

			break
		}
	}
}
