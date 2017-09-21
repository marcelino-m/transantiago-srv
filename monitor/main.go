package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis"
	Bi "github.com/marcelino-m/transantiago-srv/bi"
	"github.com/marcelino-m/transantiago-srv/gtfs"
)

func main() {

	stopmon := flag.String("stop", "all", "Stops to momintor")
	nthread := flag.Int("nthread", 200, "Number of concurrent request")
	flag.Parse()

	bi, err := Bi.NewBi("/home/marcelo/lab/tase/gtfs/gtfs.db")

	if err != nil {
		log.Fatalln("Fail to connect to gtfs DB (sqlite3 backend)")
	}

	var stops []*gtfs.Stop

	if *stopmon == "all" {
		tmp := bi.AllStop()
		for _, s := range tmp {
			stops = append(stops, s)
		}
	} else {
		bs := strings.Split(*stopmon, ",")
		for _, b := range bs {
			if b == "" {
				continue
			}
			stop := bi.Stop(strings.ToUpper(b))
			if stop == nil {
				log.Fatalf("No exit stops %s\n", b)
			}
			stops = append(stops, stop)
		}
	}

	if err != nil {
		log.Fatalln("Fail to get stops data")
	}

	redisc := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err = redisc.Ping().Result()
	if err != nil {
		log.Fatal("Fail to conect to redis server.")
	}

	maxidleconns := *nthread / 4
	if maxidleconns < 2 {
		maxidleconns = 2
	}

	tr := &http.Transport{
		MaxIdleConns:    maxidleconns,
		IdleConnTimeout: 1 * time.Second,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
	}

	httpc := &http.Client{Transport: tr}

	queue := make(chan struct{}, *nthread)
	wg := sync.WaitGroup{}
	wg.Add(1)

	for _, s := range stops {
		go monitor(httpc, s, redisc, queue)
	}

	wg.Wait()
}