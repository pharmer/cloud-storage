package packet

import (
	"context"
	"fmt"

	//	"github.com/golang/glog"
	"github.com/kubernetes-incubator/external-storage/lib/controller"
	"github.com/kubernetes-incubator/external-storage/lib/util"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"

	//"k8s.io/kubernetes/pkg/kubelet/apis"
	"github.com/appscode/go/wait"
	"github.com/packethost/packngo"
	. "github.com/pharmer/cloud-storage/cloud"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	flexvolumeVendor = "pharmer"
	flexvolumeDriver = "flexvolumes"

	// are we allowed to set this? else make up our own
	annCreatedBy = "kubernetes.io/createdby"
	createdBy    = "packet-provisioner"

	annVolumeID   = "packet.external-storage.incubator.kubernetes.io/VolumeID"
	annVolumeName = "packet.external-storage.incubator.kubernetes.io/VolumeName"

	// A PV annotation for the identity of the s3fsProvisioner that provisioned it
	annProvisionerID = "Provisioner_Id"
)

func NewPacketProvisioner(ctx context.Context, client kubernetes.Interface, packClient *packngo.Client) controller.Provisioner {
	var identity types.UID
	provisioner := &packetProvisioner{
		client:     client,
		packClient: packClient,
		ctx:        ctx,
		identity:   identity,
	}
	return provisioner
}

type packetProvisioner struct {
	client     kubernetes.Interface
	packClient *packngo.Client
	ctx        context.Context
	identity   types.UID
}

var _ controller.Provisioner = &packetProvisioner{}

func (p *packetProvisioner) getAccessModes() []core.PersistentVolumeAccessMode {
	return []core.PersistentVolumeAccessMode{
		core.ReadWriteOnce,
	}
}

func (p *packetProvisioner) Provision(options controller.VolumeOptions) (*core.PersistentVolume, error) {
	if !util.AccessModesContainedInAll(p.getAccessModes(), options.PVC.Spec.AccessModes) {
		return nil, fmt.Errorf("Invalid Access Modes: %v, Supported Access Modes: %v", options.PVC.Spec.AccessModes, p.getAccessModes())
	}

	vol, err := p.createVolume(options)
	if err != nil {
		return nil, err
	}

	annotations := make(map[string]string)
	annotations[annCreatedBy] = createdBy

	annotations[annProvisionerID] = string(p.identity)
	annotations[annVolumeID] = vol.ID
	annotations[annVolumeName] = vol.Name

	labels := make(map[string]string)
	//labels[apis.LabelZoneFailureDomain] =

	pv := &core.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name:        options.PVName,
			Labels:      labels,
			Annotations: annotations,
		},
		Spec: core.PersistentVolumeSpec{
			PersistentVolumeReclaimPolicy: options.PersistentVolumeReclaimPolicy,
			AccessModes:                   options.PVC.Spec.AccessModes,
			Capacity: core.ResourceList{
				core.ResourceName(core.ResourceStorage): resource.MustParse(fmt.Sprintf("%dGi", vol.Size)),
			},
			PersistentVolumeSource: core.PersistentVolumeSource{
				FlexVolume: &core.FlexPersistentVolumeSource{
					Driver:   fmt.Sprintf("%s/%s", flexvolumeVendor, flexvolumeDriver),
					FSType:   "ext4",
					Options:  map[string]string{},
					ReadOnly: false,
				},
			},
		},
	}
	return pv, nil

}

func (p *packetProvisioner) createVolume(volumeOptions controller.VolumeOptions) (*packngo.Volume, error) {
	zone, ok := volumeOptions.Parameters["zone"]
	if !ok {
		return nil, fmt.Errorf("error zone parameter missing")
	}
	volSize := volumeOptions.PVC.Spec.Resources.Requests[core.ResourceName(core.ResourceStorage)]
	volSizeBytes := volSize.Value()
	volszInt := util.RoundUpSize(volSizeBytes, util.GiB)

	projectId, err := getProjectID()
	if err != nil {
		return nil, err
	}
	createRequest := &packngo.VolumeCreateRequest{
		Size:         int(volszInt),
		BillingCycle: "hourly",
		// ProjectID:    projectId,
		FacilityID:  zone,
		Description: volumeOptions.PVName,
		PlanID:      "storage_1",
	}
	v, _, err := p.packClient.Volumes.Create(createRequest, projectId)
	if err != nil {
		return nil, err
	}
	return waitVolumeActive(v.ID, p.packClient)
}

func waitVolumeActive(id string, c *packngo.Client) (*packngo.Volume, error) {
	// 15 minutes = 180 * 5sec-retry
	err := wait.PollImmediate(RetryInterval, RetryTimeout, func() (bool, error) {
		v, _, err := c.Volumes.Get(id)
		if err != nil {
			return false, nil
		}
		if v.State == "active" {
			return true, nil
		}
		return false, nil
	})

	if err != nil {
		return nil, fmt.Errorf("volume %s is still not active after timeout", id)
	}
	v, _, err := c.Volumes.Get(id)
	return v, err
}
