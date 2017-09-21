package main

import (
	"fmt"
	"net/http"

	"github.com/go-redis/redis"
	"github.com/marcelino-m/transantiago-srv/fetcher"
	"github.com/marcelino-m/transantiago-srv/gtfs"
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
				m := map[string]interface{}{"stop": s.Id(), "dist": b.DistToStop()}
				redisc.HMSet(busKey, m)
			}

			break
		}
	}
}
