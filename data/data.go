package data

type Fetcher interface {
	// Fetch fetches the metrics from the configured endpoint.
	Fetch() (TSBD, error)
}

// type TSDB interface {
// 	Select(where ...WhereOption) ([]byte, error)
// }

// type WhereCondition struct {
// 	Label string
// 	Value string
// }

// type WhereOption func(*WhereCondition)

// func Where(label, value string) WhereOption {
// 	return func(condition *WhereCondition) {
// 		condition.Label = label
// 		condition.Value = value
// 	}
// }

type TSBD interface {
	GroupByLabels(labels []string) (map[string][]float64, error)
}