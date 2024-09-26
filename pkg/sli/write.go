package sli

import (
	"context"
	"github.com/bradfordwagner/go-util/log"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type writeInterface interface {
	// Upsert the metrics to the configmap
	Upsert(ctx context.Context, namespace, configmapName string, metrics MetricsMap) (err error)
}

type write struct {
	kc kubernetes.Interface
}

// enforce interface
var _ writeInterface = &write{}

// newWrite creates a new WriteInterface
func newWrite(kc kubernetes.Interface) writeInterface {
	return &write{kc: kc}
}

func (w *write) Upsert(ctx context.Context, namespace, configmapName string, metrics MetricsMap) (err error) {
	l := log.Log().With("namespace", namespace, "configmap", configmapName)

	// get existing configmap
	existingConfigMap, err := w.kc.CoreV1().ConfigMaps(namespace).Get(ctx, configmapName, metav1.GetOptions{})
	create := err != nil

	if create {
		l.Debug("creating")
		_, err = w.kc.CoreV1().ConfigMaps(namespace).Create(ctx, &v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name: configmapName,
			},
			Data: metrics.ToConfigmapData(),
		}, metav1.CreateOptions{})
	} else {
		l.Debug("updating")
		existingConfigMap.Data = metrics.ToConfigmapData()
		_, err = w.kc.CoreV1().ConfigMaps(namespace).Update(ctx, existingConfigMap, metav1.UpdateOptions{})
	}

	return
}
