package configure

import "github.com/prometheus/client_golang/prometheus"

type AggregateCollector string

const (
	AggregateCollectorSum AggregateCollector = "sum"
	AggregateCollectorAvg AggregateCollector = "avg"
	AggregateCollectorMin AggregateCollector = "min"
	AggregateCollectorMax AggregateCollector = "max"
)

type DataSourceType string

const (
	DataSourceTypePromQL DataSourceType = "promql"
)

// Supported data source types
var supportedTypes = map[DataSourceType]struct{}{
	DataSourceTypePromQL: {},
}

// Supported aggregation methods
var supportedMethods = map[AggregateCollector]struct{}{
	AggregateCollectorSum: {},
	AggregateCollectorAvg: {},
	AggregateCollectorMin: {},
	AggregateCollectorMax: {},
}

type Configure struct {
	Endpoints map[string]string  `yaml:"endpoints"`
	Metrics   []ConfiguredMetric `yaml:"metrics"`
}

type ConfiguredMetric struct {
	Name           string            `yaml:"name"`
	Description    string            `yaml:"description"`
	Type           string            `yaml:"data_source_type"`
	Query          string            `yaml:"query"`
	ExpectedLabels []string          `yaml:"expected_labels"`
	Aggregation    AggregationConfig `yaml:"aggregation"`
}

type AggregationConfig struct {
	Labels []string `yaml:"labels"`
	Method string   `yaml:"method"`
}

func (c *Configure) Verify() bool {
	if len(c.Metrics) == 0 {
		return false
	}
	for _, metric := range c.Metrics {
		if !metric.Verify() {
			return false
		}
	}
	return true
}

func (m *ConfiguredMetric) Verify() bool {
	if (m.Name == "") || (m.Type == "") || (m.Query == "") {
		return false
	}
	// Convert ExpectedLabels to a map for faster lookups
	expectedLabels := make(map[string]struct{}, len(m.ExpectedLabels))
	for _, label := range m.ExpectedLabels {
		expectedLabels[label] = struct{}{}
	}

	// Check if Type is supported
	if _, ok := supportedTypes[DataSourceType(m.Type)]; !ok {
		return false
	}

	// Check if Aggregation.Method is supported
	if _, ok := supportedMethods[AggregateCollector(m.Aggregation.Method)]; !ok {
		return false
	}

	// Check if Aggregation.Labels are all in ExpectedLabels
	for _, label := range m.Aggregation.Labels {
		if _, ok := expectedLabels[label]; !ok {
			return false
		}
	}

	return true
}

func (m *ConfiguredMetric) GetPrometheusDesc() (*prometheus.Desc, error) {
	var labels []string
	// labels equals to ExpectedLabels - Aggregation.Label
	mapAggregationLabels := make(map[string]struct{}, len(m.Aggregation.Labels))
	for _, label := range m.Aggregation.Labels {
		mapAggregationLabels[label] = struct{}{}
	}
	for _, label := range m.ExpectedLabels {
		if _, ok := mapAggregationLabels[label]; !ok {
			labels = append(labels, label)
		}
	}

	return prometheus.NewDesc(
		m.Name,
		m.Description,
		labels,
		nil,
	), nil
}
