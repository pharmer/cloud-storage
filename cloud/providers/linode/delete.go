package linode

import (
	"fmt"
	"strconv"

	"github.com/golang/glog"
	"github.com/kubernetes-incubator/external-storage/lib/controller"
	core "k8s.io/api/core/v1"
)

func (p *linodeProvisioner) Delete(volume *core.PersistentVolume) error {
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
	vid, err := strconv.Atoi(vol)
	if err != nil {
		return fmt.Errorf("can't parse the volume id %s", vol)
	}
	_, err = p.linodeClient.Volume.Delete(vid)
	if err != nil {
		glog.Errorf("Failed to delete volume %s, error: %s", volume, err.Error())
		return err
	}
	return nil
}

func (p *linodeProvisioner) provisioned(volume *core.PersistentVolume) (bool, error) {
	provisionerID, ok := volume.Annotations[annProvisionerID]
	if !ok {
		return false, fmt.Errorf("pv doesn't have an annotation %s", annProvisionerID)
	}
	return provisionerID == string(p.identity), nil
}
