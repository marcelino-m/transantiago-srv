package fetcher

import (
	"fmt"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/marcelino-m/transantiago-srv/gtfs"
)

var (
	fmturl = "http://m.ibus.cl/Servlet?paradero=%s&servicio=&button=Consulta+Paradero"
)

func FetchStopDataRaw(stopcode string, client *http.Client) (*goquery.Document, error) {

	if client == nil {
		tr := &http.Transport{
			MaxIdleConns:       12000,
			IdleConnTimeout:    1 * time.Second,
			DisableCompression: true,
		}

		client = &http.Client{Transport: tr}
	}

	req, err := http.NewRequest("GET", makUrl(stopcode), nil)
	if nil != err {
		return nil, err
	}

	req.Header.Set("User-Agent", randMovilUa())
	resp, err := client.Do(req)
	if nil != err {
		return nil, err
	}

	defer resp.Body.Close()

	return goquery.NewDocumentFromResponse(resp)
}

func FetchStopData(stopcode string, client *http.Client) ([]*gtfs.Bus, error) {

	doc, err := FetchStopDataRaw(stopcode, client)
	if err != nil {
		return nil, err
	}

	return Parser(doc), nil
}

func makUrl(stopcode string) string {
	return fmt.Sprintf(fmturl, stopcode)
}
