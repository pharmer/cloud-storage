package cmds

import (
	"context"

	"github.com/golang/glog"
	"github.com/kubernetes-incubator/external-storage/lib/controller"
	"github.com/pharmer/cloud-storage/cloud"
	"github.com/pharmer/cloud-storage/cmds/options"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/util/wait"
)

func NewCmdInit() *cobra.Command {
	cfg := options.NewConfig()
	cmd := &cobra.Command{
		Use:               "init",
		Short:             "Initializes the driver.",
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			cfg.ValidateFlags()
			provider, err := cloud.GetCloudManager(cfg.Provider, context.Background())
			if err != nil {
				glog.Fatalf("Invalid provider specified: %v", err)
			}
			provisioner, err := provider.Init()
			if err != nil {

			}
			client := cloud.InitializeClient(cfg)

			pc := controller.NewProvisionController(
				client.Kube,
				cfg.Provisioner,
				provisioner,
				client.ServerVersion,
				provider.Namer(),
			)
			pc.Run(wait.NeverStop)

		},
	}
	cfg.AddFlags(cmd.Flags())
	return cmd
}
