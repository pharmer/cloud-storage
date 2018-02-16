package linode

import (
	"context"

	"github.com/kubernetes-incubator/external-storage/lib/controller"
)

func (e *ExternalStorage) Init() (controller.Provisioner, error) {
	token, err := getCredential()
	if err != nil {
		return nil, err
	}
	e.client = token.getClient()
	linodeProvisionar := NewLinodeProvisioner(context.Background(), e.kc, e.client)
	return linodeProvisionar, nil
}
