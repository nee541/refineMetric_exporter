package main

import (
	"os"
	"testing"

	"github.com/nee541/refineMetric_exporter/data/configure"
	"gopkg.in/yaml.v2"
)

func TestReadConfiguration(t *testing.T) {
	// Parse the configuration file
	data, err := os.ReadFile("./metric_configuration.yaml")
	if err != nil {
		t.Fatal(err)
	}

	var config configure.Configure
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(config)
}
