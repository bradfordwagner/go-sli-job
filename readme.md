# go-sli-job 

## Example
```go
package main

import (
	"context"
	"fmt"
	"github.com/bradfordwagner/go-util/log"
	"github.com/sethvargo/go-envconfig"
	"vmwh-sli/internal/config"
)

func main() {
	// setup logging
	l := log.Log()

	// init vars
	metricsPushUrl := ""
	configmapName := "my-configmap"
	namespace := "default"
	extraLabels := fmt.Sprintf(`probe_type="%s"`, "test")

	// setup context
	sliContext, err := sli.NewContext()
	if err != nil {
		l.With("error", err).Fatal("Error creating sli context")
	}

	// establish metrics
	errorsMetric := &sli.PushMetric{
		Name:        "sli_errors",
		Description: "Whether or not the SLI probe's request was a success",
		MetricType:  sli.CounterType,
		ExtraLabels: extraLabels,
	}
	requestsMetric := &sli.PushMetric{
		Name:        "sli_requests",
		Description: "how many sli probe were made",
		MetricType:  sli.CounterType,
		ExtraLabels: extraLabels,
	}
	metrics := sli.Map(errorsMetric, requestsMetric)

	// get existing metrics
	sliContext.Get.Get(ctx, namespace, configmapName, metrics)

	// "update" metrics
	errorsMetric.Value, requestsMetric.Value = 2, requestsMetric.Value+1

	// push metrics to telegraf and to configmap
	err = sliContext.Pusher.Push(ctx, sli.PushOpts{
		Url:           metricsPushUrl,
		Namespace:     namespace,
		ConfigmapName: configmapName,
		Metrics:       metrics,
	})
	if err != nil {
		l.With("error", err).Fatal("Failed to push metrics")
	}
}
```
