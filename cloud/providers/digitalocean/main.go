package digitalocean

import (
	"context"

	dv "github.com/kubernetes-incubator/external-storage/digitalocean/pkg/volume"
	"github.com/kubernetes-incubator/external-storage/lib/controller"
)

func (e *ExternalStorage) Init() (controller.Provisioner, error) {
	token, err := getCredential()
	if err != nil {
		return nil, err
	}
	e.client = token.getClient()
	doProvisionar := dv.NewDigitalOceanProvisioner(context.Background(), e.kc, e.client)
	return doProvisionar, nil
}

func (e *ExternalStorage) Namer() controller.Namer {
	return controller.NewDefaultNamer()
}
