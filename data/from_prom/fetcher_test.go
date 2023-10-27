package from_prom

import (
	"testing"
)

func TestFetch(t *testing.T) {
	fetcher := &PrometheusFetcher{
		URL:    "http://192.168.33.7:30092",
		PromQL: "1 - (avg(irate(node_cpu_seconds_total{mode=\"idle\"}[5m])) by (instance))",
	}
	data, err := fetcher.Fetch()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(data.GroupByLabels([]string{}))
}
