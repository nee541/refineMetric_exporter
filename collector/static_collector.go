package collector

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/nee541/refineMetric_exporter/data"
	"github.com/nee541/refineMetric_exporter/data/configure"
	"github.com/nee541/refineMetric_exporter/data/from_prom"
	"github.com/prometheus/client_golang/prometheus"
)

type staticCollector struct {
	config *configure.Configure
	logger log.Logger
}

func NewStaticCollector(config *configure.Configure, logger log.Logger) *staticCollector {
	return &staticCollector{
		config: config,
		logger: logger,
	}
}

func (s *staticCollector) Update(ch chan<- prometheus.Metric, opts ...CollectorOption) error {
	wg := sync.WaitGroup{}
	wg.Add(len(s.config.Metrics))

	mutex := sync.Mutex{}
	hasError := false

	prometheusEndpoint, exists := s.config.Endpoints["prometheus"]
	if !exists {
		return fmt.Errorf("prometheus endpoint not found")
	}
	level.Debug(s.logger).Log("msg", "prometheus endpoint", "endpoint", prometheusEndpoint)

	for _, metric := range s.config.Metrics {
		fetcher := from_prom.NewPrometheusFetcher(prometheusEndpoint, metric.Query)
		labelNames := SliceDifference[string](metric.ExpectedLabels, metric.Aggregation.Labels)

		go func(metric configure.ConfiguredMetric) {
			defer wg.Done()
			err := updateMetric(ch, metric, fetcher, labelNames, s.logger)
			if err != nil {
				level.Error(s.logger).Log("msg", "failed to update metric", "name", metric.Name, "labels", metric.ExpectedLabels, "err", err)
				mutex.Lock()
				hasError = true
				mutex.Unlock()
			}
		}(metric)
	}
	wg.Wait()

	if hasError {
		return fmt.Errorf("failed to update metric")
	}
	return nil
}

func updateMetric(ch chan<- prometheus.Metric, metric configure.ConfiguredMetric, fetcher data.Fetcher, labelNames []string, logger log.Logger) error {
	tsdb, err := fetcher.Fetch()
	if err != nil {
		return err
	}
	grouped, err := tsdb.GroupByLabels(labelNames)
	if err != nil {
		return err
	}

	level.Debug(logger).Log("msg", "grouped", "name", metric.Name, "labels", metric.ExpectedLabels, "labelNames", labelNames, "grouped", grouped)

	for labelJson, values := range grouped {
		labelValues := []string{}
		err := json.Unmarshal([]byte(labelJson), &labelValues)
		if err != nil {
			return err
		}
		var aggregated float64
		switch configure.AggregateCollector(metric.Aggregation.Method) {
		case configure.AggregateCollectorSum:
			aggregated = Sum(values)
		case configure.AggregateCollectorAvg:
			aggregated = Avg(values)
		case configure.AggregateCollectorMin:
			aggregated = Min(values)
		case configure.AggregateCollectorMax:
			aggregated = Max(values)
		case configure.AggregateCollectorPencentile:
			if 0 < metric.Aggregation.Percentile && metric.Aggregation.Percentile < 100 {
				aggregated = Percentile(values, metric.Aggregation.Percentile)
			} else {
				level.Error(logger).Log("msg", "invalid pencentile", "name", metric.Name, "labels", labelJson, "value", aggregated)
				return fmt.Errorf("invalid pencentile: %d", metric.Aggregation.Percentile)
			}
		default:
			return fmt.Errorf("unknown aggregation method: %s", metric.Aggregation.Method)
		}
		level.Debug(logger).Log("msg", "metric", "name", metric.Name, "labels", labelJson, "value", aggregated)
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				metric.Name,
				metric.Description,
				labelNames,
				nil,
			),
			prometheus.GaugeValue,
			aggregated,
			labelValues...,
		)
	}

	return nil
}
