package digitalocean

import (
	"context"

	"github.com/digitalocean/godo"
	. "github.com/pharmer/cloud-storage/cloud"
	"k8s.io/client-go/kubernetes"
)

type ExternalStorage struct {
	ctx    context.Context
	kc     kubernetes.Interface
	client *godo.Client
}

var _ Interface = &ExternalStorage{}

const (
	UID = "digitalocean"
)

func init() {
	RegisterCloudManager(UID, func(ctx context.Context) (Interface, error) {
		return New(ctx), nil
	})
}

func New(ctx context.Context) Interface {
	return &ExternalStorage{ctx: ctx}
}
