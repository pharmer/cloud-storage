package linode

import (
	"context"

	. "github.com/pharmer/cloud-storage/cloud"
	"github.com/taoh/linodego"
	"k8s.io/client-go/kubernetes"
)

type ExternalStorage struct {
	ctx    context.Context
	kc     kubernetes.Interface
	client *linodego.Client
}

var _ Interface = &ExternalStorage{}

const (
	UID = "linode"
)

func init() {
	RegisterCloudManager(UID, func(ctx context.Context) (Interface, error) {
		return New(ctx), nil
	})
}

func New(ctx context.Context) Interface {
	return &ExternalStorage{ctx: ctx}
}
