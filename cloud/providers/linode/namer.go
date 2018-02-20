package linode

import (
	"strings"

	"github.com/kubernetes-incubator/external-storage/lib/controller"
	core "k8s.io/api/core/v1"
)

type LinodeNamer struct{}

var _ controller.Namer = &LinodeNamer{}

func NewLinodeNamer() controller.Namer {
	return &LinodeNamer{}
}
func (n *LinodeNamer) GetProvisionedVolumeNameForClaim(claim *core.PersistentVolumeClaim) string {
	claimId := string(claim.UID)
	claimId = "p" + strings.Replace(claimId, "-", "", -1)
	if len(claimId) > 32 {
		claimId = claimId[:32]
	}
	return claimId
}
