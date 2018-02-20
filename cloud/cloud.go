package cloud

import (
	"time"

	"github.com/golang/glog"
	"github.com/kubernetes-incubator/external-storage/lib/controller"
	"github.com/pharmer/cloud-storage/cmds/options"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	CredentialFileEnv         = "CRED_FILE_PATH"
	CredentialDefaultLocation = "/etc/kubernetes/cloud.json"
	SecretDefaultLocation     = "/var/run/secrets/pharmer/flexvolmues"
	RetryInterval             = 5 * time.Second
	RetryTimeout              = 10 * time.Minute
)

//var ErrNotSupported = errors.New("Not Supported")
//var ErrIncorrectArgNumber = errors.New("Incorrect number of args")

type Interface interface {
	Init() (controller.Provisioner, error)
	Namer() controller.Namer
}

type Client struct {
	Kube          kubernetes.Interface
	ServerVersion string
}

func InitializeClient(opt *options.Config) *Client {
	// Create the client according to whether we are running in or out-of-cluster
	var config *rest.Config
	var err error
	if opt.MasterUrl != "" || opt.Kubeconfig != "" {
		glog.Infof("Either master or kubeconfig specified. building kube config from that..")
		config, err = clientcmd.BuildConfigFromFlags(opt.MasterUrl, opt.Kubeconfig)
	} else {
		glog.Infof("Building kube configs for running in cluster...")
		config, err = rest.InClusterConfig()
	}
	if err != nil {
		glog.Fatalf("Failed to create config: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		glog.Fatalf("Failed to create client: %v", err)
	}

	// The controller needs to know what the server version is because out-of-tree
	// provisioners aren't officially supported until 1.5
	serverVersion, err := clientset.Discovery().ServerVersion()
	if err != nil {
		glog.Fatalf("Error getting server version: %v", err)
	}
	return &Client{
		Kube:          clientset,
		ServerVersion: serverVersion.GitVersion,
	}

}
