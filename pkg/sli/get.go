package sli

import (
	"context"
	"github.com/bradfordwagner/go-util/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"strconv"
	"strings"
)

type GetInterface interface {
	// Gets the metrics from the configmap
	Get(ctx context.Context, configmapName string, metrics map[string]*PushMetric) (err error)
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
func (g *get) Get(ctx context.Context, configmapName string, metrics map[string]*PushMetric) (err error) {
	l := log.Log()
	configMap, err := g.kc.CoreV1().ConfigMaps("").Get(ctx, configmapName, metav1.GetOptions{})

	if strings.Contains(err.Error(), "the server could not find the requested resource") {
		// configmap does not exist, do not throw an error
		err = nil
		l.With("configmap", configmapName).Warn("Configmap does not exist")
	} else if err != nil {
		return
	}

	for k, v := range configMap.Data {
		// convert to float64
		var f float64
		f, err = strconv.ParseFloat(v, 64)
		if err != nil {
			return
		}

		// check configmap field matches metric name
		metric, ok := metrics[k]
		if ok {
			metric.Value = f
		}
	}

	return
}
