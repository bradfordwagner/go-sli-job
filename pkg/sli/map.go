package sli

import (
	"fmt"
	"strconv"
)

// Map takes a slice of PushMetric and returns a map of PushMetric
func Map(metrics ...*PushMetric) (m MetricsMap) {
	m = make(map[string]*PushMetric)
	for _, metric := range metrics {
		m[metric.Name] = metric
	}
	return
}

type MetricsMap map[string]*PushMetric

func (m MetricsMap) ToConfigmapData() map[string]string {
	data := make(map[string]string)
	for k, v := range m {
		data[k] = fmt.Sprintf("%f", v.Value)
	}
	return data
}

func (m MetricsMap) ExtractFromConfigmapData(data map[string]string) (err error) {
	for k, v := range data {
		var f float64
		f, err = strconv.ParseFloat(v, 64)
		if err != nil {
			return
		}
		m[k].Value = f
	}
	return
}
