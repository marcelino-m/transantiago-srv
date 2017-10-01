package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/go-redis/redis"
	Bi "github.com/marcelino-m/transantiago-srv/bi"
	"strings"
)

var bi *Bi.Bi

// Stops flags
type StopFlag []string

//  ...
func (s *StopFlag) String() string {
	return fmt.Sprintf("%v", *s)
}

//  ...
func (s *StopFlag) Set(value string) error {
	*s = append(*s, strings.ToUpper(value))
	return nil
}

func main() {
	stops := StopFlag{}

	flag.Var(&stops, "stop", "Stop to momintor, you can pas multiple `-stop'")
	route := flag.String("route", "", "Route to monitor")
	nthread := flag.Int("nth", 200, "Number of concurrent request")

	flag.Parse()

	var err error
	if *route != "" {
		bi, err = Bi.NewBiFromRoute("/home/marcelo/lab/tase/gtfs/gtfs.db", *route)
	} else {
		bi, err = Bi.NewBiFromStops("/home/marcelo/lab/tase/gtfs/gtfs.db", stops...)
	}

	if err != nil {
		if err != nil {
			log.Fatal("Fail to initialize deamon\n Error:", err)
		}
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
	log.Println("Initialize bi...")
	err = bi.InitializeBi()
	if err != nil {
		log.Fatal("Fail to initialize bi!.\n: ", err)
	} else {
		log.Println("Bi redy is ready")
	}

	wg.Add(1)

	for _, s := range bi.AllStop() {
		go monitor(s, *route, httpc, redisc, queue)
	}

	wg.Wait()
}
