package sli

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

type PushMetric struct {
	Name        string
	Description string
	MetricType  MetricType
	Value       float64

	//  in the format of key="value"
	ExtraLabels string
}

type MetricType string

const (
	GaugeType   MetricType = "gauge"
	CounterType MetricType = "counter"
)

func (m MetricType) String() string {
	return string(m)
}

type pusher struct {
	Url string
}

// enforce interface
var _ PusherInterface = &pusher{}

func NewPusher() PusherInterface {
	return &pusher{}
}

func (p *pusher) Push(url string, metrics map[string]*PushMetric) error {
	if len(metrics) == 0 {
		return nil
	}

	client := http.Client{}

	buf := bytes.NewBuffer([]byte{})
	for _, metric := range metrics {
		if err := writeMetric(buf, metric); err != nil {
			return err
		}
	}

	resp, err := client.Post(url, "text/plain", buf)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf(
			"Got %v response when pushing metrics to %v",
			resp.StatusCode,
			p.Url,
		)
	}

	return nil
}

func writeMetric(w io.Writer, metric *PushMetric) error {
	// https://github.com/prometheus/docs/blob/main/content/docs/instrumenting/exposition_formats.md#text-format-example
	if _, err := fmt.Fprintf(w, "# HELP %s %s\n", metric.Name, metric.Description); err != nil {
		return err
	}

	if _, err := fmt.Fprintf(w, "# TYPE %s %s\n", metric.Name, metric.MetricType.String()); err != nil {
		return err
	}

	if _, err := fmt.Fprintf(w, "%s{%s} %f\n", metric.Name, metric.ExtraLabels, metric.Value); err != nil {
		return err
	}

	return nil
}

type PusherInterface interface {
	Push(url string, metrics map[string]*PushMetric) error
}
