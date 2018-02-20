package packet

import (
	"fmt"

	"github.com/golang/glog"
	"github.com/kubernetes-incubator/external-storage/lib/controller"
	core "k8s.io/api/core/v1"
)

func (p *packetProvisioner) Delete(volume *core.PersistentVolume) error {
	glog.Infof("Delete called for volume: %s", volume.Name)

	provisioned, err := p.provisioned(volume)
	if err != nil {
		return fmt.Errorf("error determining if this provisioner was the one to provision volume %q: %v", volume.Name, err)
	}

	if !provisioned {
		strerr := fmt.Sprintf("this provisioner id %s didn't provision volume %q and so can't delete it; id %s did & can", p.identity, volume.Name, volume.Annotations[annProvisionerID])
		return &controller.IgnoredError{Reason: strerr}
	}

	vol, ok := volume.Annotations[annVolumeID]
	if !ok {
		return fmt.Errorf("pv doesn't have an annotation %s", annVolumeID)
	}

	_, err = p.packClient.Volumes.Delete(vol)

	return err
}

func (p *packetProvisioner) provisioned(volume *core.PersistentVolume) (bool, error) {
	provisionerID, ok := volume.Annotations[annProvisionerID]
	if !ok {
		return false, fmt.Errorf("pv doesn't have an annotation %s", annProvisionerID)
	}
	return provisionerID == string(p.identity), nil
}
