package from_prom

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type Response struct {
	Status string `json:"status"`
	Data   Data   `json:"data"`
}

type Data struct {
	ResultType string   `json:"resultType"`
	Result     []Result `json:"result"`
}

type Result struct {
	Metric map[string]string `json:"metric"`
	Value  []interface{}     `json:"value"`
}

func (r *Result) GetFloatValue() (float64, error) {
	if str, ok := r.Value[1].(string); !ok {
		return 0, fmt.Errorf("error converting value to float64: %v", r.Value[1])
	} else {
		val, err := strconv.ParseFloat(strings.TrimSpace(str), 64)
		if err != nil {
			return 0, err
		}
		return val, nil
	}
}

func generateKeyFromLabels(metric map[string]string, labels []string) (string, error) {
	keyData := []string{}
	for _, label := range labels {
		if value, exists := metric[label]; exists {
			keyData = append(keyData, value)
		} else {
			return "", fmt.Errorf("label %s not found in metric", label)
		}
	}

	bytes, err := json.Marshal(keyData)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func (d *Data) GroupByLabels(labels []string) (map[string][]float64, error) {
	grouped := make(map[string][]float64)
	for _, result := range d.Result {
		key, err := generateKeyFromLabels(result.Metric, labels)
		if err != nil {
			return nil, err
		}
		value, err := result.GetFloatValue()
		if err != nil {
			return nil, err
		}
		grouped[key] = append(grouped[key], value)
	}
	return grouped, nil
}
