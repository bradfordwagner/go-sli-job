package sli

import (
	"bytes"
	"context"
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

func (p *PushMetric) String() string {
	return fmt.Sprintf("%s=%f", p.Name, p.Value)
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
	w writeInterface
}

type PusherInterface interface {
	Push(ctx context.Context, opts PushOpts) error
}

// enforce interface
var _ PusherInterface = &pusher{}

func NewPusher(write writeInterface) PusherInterface {
	return &pusher{
		w: write,
	}
}

type PushOpts struct {
	Url           string
	Namespace     string
	ConfigmapName string
	Metrics       MetricsMap
}

func (p *pusher) Push(ctx context.Context, opts PushOpts) error {
	if len(opts.Metrics) == 0 {
		return nil
	}

	client := http.Client{}

	buf := bytes.NewBuffer([]byte{})
	for _, metric := range opts.Metrics {
		if err := writeMetric(buf, metric); err != nil {
			return err
		}
	}

	resp, err := client.Post(opts.Url, "text/plain", buf)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf(
			"Got %v response when pushing metrics to %v",
			resp.StatusCode,
			opts.Url,
		)
	}

	// write to configmap
	return p.w.Upsert(ctx, opts.Namespace, opts.ConfigmapName, opts.Metrics)
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
