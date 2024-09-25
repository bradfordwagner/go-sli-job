package sli

import (
	"github.com/bradfordwagner/go-kubeclient/kube"
	"k8s.io/client-go/kubernetes"
)

// Context is the context for the SLI package
type Context struct {
	KubeClient kubernetes.Interface
	Pusher     PusherInterface
}

// NewContext creates a new context for the SLI package
func NewContext() (c *Context, err error) {
	client, err := kube.Client()
	if err != nil {
		return
	}

	return &Context{
		KubeClient: client,
		Pusher:     NewPusher(),
	}, nil
}
