package sli

import (
	"context"
	"github.com/bradfordwagner/go-util/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"strconv"
)

type GetInterface interface {
	// Gets the metrics from the configmap
	Get(ctx context.Context, namespace, configmapName string, metrics MetricsMap)
}

type get struct {
	kc kubernetes.Interface
}

// enforce interface
var _ GetInterface = &get{}

// NewGet creates a new GetInterface
func NewGet(kc kubernetes.Interface) GetInterface {
	return &get{kc: kc}
}

// Get gets the metrics from the configmap
func (g *get) Get(ctx context.Context, namespace, configmapName string, metrics MetricsMap) {
	l := log.Log().With("namespace", namespace, "configmap", configmapName)
	configMap, err := g.kc.CoreV1().ConfigMaps(namespace).Get(ctx, configmapName, metav1.GetOptions{})

	// dne, use provided values
	if err != nil {
		l.With("err", err).Warn("Configmap does not exist")
		return
	}

	// note if we add a new metric to the map
	// it will reset all other metrics
	var failedParse bool
	data := configMap.Data
	updates := make(map[string]float64)
	for k := range metrics {
		// when missing reset all metrics to zero
		// ignore configmap values
		v, ok := data[k]
		if !ok {
			failedParse = true
			continue
		}

		// convert to float
		// if it fails ignore configmap values
		var f float64
		f, err = strconv.ParseFloat(v, 64)
		if err != nil {
			l.With("k", k, "v", v, "err", err).Warn("Failed to parse value")
			failedParse = true
			continue
		}
		updates[k] = f
	}

	// on successful parse update metrics
	if !failedParse {
		l.Info("Successfully parsed configmap")
		for k := range updates {
			metrics[k].Value = updates[k]
		}
	} else {
		l.Warn("Failed to parse configmap, this will likely reset your metrics to nil values")
	}

	return
}
