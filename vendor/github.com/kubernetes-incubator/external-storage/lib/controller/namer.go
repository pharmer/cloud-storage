package controller

import "k8s.io/api/core/v1"

type Namer interface {
	GetProvisionedVolumeNameForClaim(claim *v1.PersistentVolumeClaim) string
}

type DefaultNamer struct {}

var _ Namer = &DefaultNamer{}


func NewDefaultNamer() Namer  {
	return &DefaultNamer{}
}

func (n *DefaultNamer) GetProvisionedVolumeNameForClaim(claim *v1.PersistentVolumeClaim) string  {
	return "pvc-" + string(claim.UID)
}