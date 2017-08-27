package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/go-redis/redis"
	Bi "github.com/marcelino-m/transantiago-srv/bi"
	"github.com/marcelino-m/transantiago-srv/fetcher"
	"github.com/marcelino-m/transantiago-srv/gtfs"
	"strings"
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
	count := 1

	for {
		for _, s := range stops {
			queue <- struct{}{}
			log.Printf("go routine %d", count)
			count++
			go worker(httpc, s, redisc, queue, 0)
		}
	}
}

func worker(httpc *http.Client, s *gtfs.Stop, redisc *redis.Client, q <-chan struct{}, retrytime int) {

	if retrytime > 10 {
		// TODO: Handle this
		log.Println("can't get data from stop %s, retry canceled after 10 times")
		<-q
		return
	}

	log.Printf("Proces stop %s", s.Id())
	buses, err := fetcher.FetchStopData(s.Id(), httpc)
	if err != nil {
		log.Printf("Erro geting data to stop %s, ", s.Id())
		log.Println(err)
		go worker(httpc, s, redisc, q, retrytime+1)
		return
	}

	<-q

	for _, b := range buses {
		stopKey := fmt.Sprintf("stop:%s", s.Id())
		busKey := fmt.Sprintf("bus:%s", b.Id())
		redisc.HSet(stopKey, b.Id(), b.DistToStop())
		redisc.HSet(busKey, s.Id(), b.DistToStop())
	}

	fmt.Printf("Done stop %s\n", s.Id())
}
