package sli

// Map takes a slice of PushMetric and returns a map of PushMetric
func Map(metrics ...*PushMetric) (m map[string]*PushMetric) {
	m = make(map[string]*PushMetric)
	for _, metric := range metrics {
		m[metric.Name] = metric
	}
	return
}
