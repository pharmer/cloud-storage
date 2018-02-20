package packet

import (
	"fmt"
	"testing"

	"github.com/packethost/packngo"
)

func TestVol(t *testing.T) {
	client := getClient()

	vol, _, err := client.Volumes.Get("")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(vol)

	createRequest := &packngo.VolumeCreateRequest{
		Size:         int(10),
		BillingCycle: "hourly",
		ProjectID:    "",
		FacilityID:   "ewr1",
		Description:  "testvol",
		PlanID:       "storage_1",
	}
	v, _, err := client.Volumes.Create(createRequest, "")
	fmt.Println(v, err)
}

func TestVolDel(t *testing.T) {
	client := getClient()
	_, err := client.Volumes.Delete("")
	fmt.Println(err)
}

func getClient() *packngo.Client {
	return packngo.NewClient("", "", nil)
}
