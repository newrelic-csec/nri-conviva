package main

import (
	"io"
	"os"

	"gopkg.in/yaml.v3"

	sdk_log "github.com/newrelic/infra-integrations-sdk/v4/log"
)
const (
	DEFAULT_API_V3_URL = "https://api.conviva.com/insights/3.0/metrics"
	DEFAULT_START_OFFSET = "20m"
	DEFAULT_END_OFFSET = "10m"
	DEFAULT_GRANULARITY = "PT1M"
)

type ConfigMetric struct {
	Metric          string 				`yaml:"metric"`
	MetricGroup     string 				`yaml:"metricGroup"`
	Names			[]string			`yaml:"names"`
	Dimensions      []string			`yaml:"dimensions"`
	Filters			map[string][]string `yaml:"filters"`
	StartOffset     string              `yaml:"startOffset"`
	EndOffset       string              `yaml:"endOffset"`
	Granularity     string              `yaml:"granularity"`
}

type Config struct {
	ApiV3URL          string
	StartOffset       string			`yaml:"startOffset"`
	EndOffset         string			`yaml:"endOffset"`
	Granularity       string			`yaml:"granularity"`
	Metrics           []ConfigMetric    `yaml:"metrics"`
}

func applyDefaults(config *Config) {
	if config.ApiV3URL == "" {
		config.ApiV3URL = DEFAULT_API_V3_URL
	}

	if config.StartOffset == "" {
		config.StartOffset = DEFAULT_START_OFFSET
	}

	if config.EndOffset == "" {
		config.EndOffset = DEFAULT_END_OFFSET
	}

	if config.Granularity == "" {
		config.Granularity = DEFAULT_GRANULARITY
	}
}

func loadConfig(configPath string, log sdk_log.Logger) (*Config, error) {
	log.Debugf("loading Conviva configuration from config path %s...", configPath)

	if configPath == "" {
		cfg := &Config{}
		applyDefaults(cfg)
		return cfg, nil
	}
	
	fd, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	rawConfig, err := io.ReadAll(fd)
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
    
	err = yaml.Unmarshal(rawConfig, cfg)
	if err != nil {
		return nil, err
	}

	applyDefaults(cfg)

	log.Debugf("conviva config loaded")
	log.Debugf("configuration: %v", *cfg)

	return cfg, nil
}