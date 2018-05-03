package linode

import (
	"context"
	"fmt"
	"strconv"

	"github.com/golang/glog"
	"github.com/kubernetes-incubator/external-storage/lib/controller"
	"github.com/kubernetes-incubator/external-storage/lib/util"
	"github.com/taoh/linodego"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	//"k8s.io/kubernetes/pkg/kubelet/apis"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	flexvolumeVendor = "pharmer"
	flexvolumeDriver = "flexvolumes"

	// are we allowed to set this? else make up our own
	annCreatedBy = "kubernetes.io/createdby"
	createdBy    = "linode-provisioner"

	annVolumeID   = "linode.external-storage.incubator.kubernetes.io/VolumeID"
	annVolumeName = "linode.external-storage.incubator.kubernetes.io/VolumeName"

	// A PV annotation for the identity of the s3fsProvisioner that provisioned it
	annProvisionerID = "Provisioner_Id"
)

func NewLinodeProvisioner(ctx context.Context, client kubernetes.Interface, linodeClient *linodego.Client) controller.Provisioner {
	var identity types.UID

	provisioner := &linodeProvisioner{
		client:       client,
		linodeClient: linodeClient,
		ctx:          ctx,
		identity:     identity,
	}
	return provisioner
}

type linodeProvisioner struct {
	client       kubernetes.Interface
	linodeClient *linodego.Client
	ctx          context.Context
	identity     types.UID
}

var _ controller.Provisioner = &linodeProvisioner{}

// https://github.com/kubernetes-incubator/external-storage/blob/e26435c2ccd9ed5d2a60c838a902d22a3ec6ef5c/iscsi/targetd/provisioner/iscsi-provisioner.go#L102
// getAccessModes returns access modes Linode Block Storage volume supports.
func (p *linodeProvisioner) getAccessModes() []core.PersistentVolumeAccessMode {
	return []core.PersistentVolumeAccessMode{
		core.ReadWriteOnce,
	}
}

func (p *linodeProvisioner) Provision(options controller.VolumeOptions) (*core.PersistentVolume, error) {
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
	annotations[annVolumeID] = strconv.Itoa(vol.VolumeId)
	annotations[annVolumeName] = vol.Label.String()

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
					Options:  map[string]string{},
					ReadOnly: false,
				},
			},
		},
	}
	return pv, nil
}

func (p *linodeProvisioner) createVolume(volumeOptions controller.VolumeOptions) (*linodego.Volume, error) {
	zone, ok := volumeOptions.Parameters["zone"]
	if !ok {
		return nil, fmt.Errorf("error zone parameter missing")
	}

	volSize := volumeOptions.PVC.Spec.Resources.Requests[core.ResourceName(core.ResourceStorage)]
	volSizeBytes := volSize.Value()
	volszInt := util.RoundUpSize(volSizeBytes, util.GiB)

	args := map[string]string{
		"DatacenterID": zone,
	}
	resp, err := p.linodeClient.Volume.Create(int(volszInt), volumeOptions.PVName, args)
	if err != nil {
		glog.Errorf("Failed to create volume %s, error: %s", volumeOptions.PVName, err.Error())
		return nil, err
	}

	vrsp, err := p.linodeClient.Volume.List(resp.VolumeId.VolumeId)
	if err != nil {
		glog.Errorf("Faild to get volume %s, error: %v", volumeOptions, err.Error())
		return nil, err
	}
	return &vrsp.Volume[0], nil
}
