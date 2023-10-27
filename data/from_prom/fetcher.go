package from_prom

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/nee541/refineMetric_exporter/data"
)

type PrometheusFetcher struct {
	URL    string
	PromQL string
}

func NewPrometheusFetcher(url string, promql string) *PrometheusFetcher {
	return &PrometheusFetcher{
		URL:    url,
		PromQL: promql,
	}
}

func (p *PrometheusFetcher) Fetch() (data.TSBD, error) {
	// url encode promql
	requestURL := fmt.Sprintf("%s/api/v1/query?query=%s", p.URL, url.QueryEscape(p.PromQL))
	// fetch data from http
	res, err := http.Get(requestURL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	// read data
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	bodyObject := &Response{}
	err = json.Unmarshal(bodyBytes, bodyObject)
	if err != nil || bodyObject.Status != "success" {
		return nil, err
	}
	if bodyObject.Status != "success" {
		return nil, fmt.Errorf("status is not success")
	}
	return &bodyObject.Data, nil
}
