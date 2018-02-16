package options

import (
	//	"log"
	"os"
	//	"github.com/pharmer/cloud-storage/util"
	"strings"

	"github.com/golang/glog"
	"github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/util/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

type Config struct {
	Provider    string
	Provisioner string
	MasterUrl   string
	Kubeconfig  string
}

const (
	keyProvider = "provider"
	envProvider = "PROVIDER"
)

func NewConfig() *Config {
	//	var err error
	provider := os.Getenv(envProvider)
	/*if provider == "" {
		provider, err = util.ReadSecretKeyFromFile(cloud.SecretDefaultLocation, keyProvider)
		if err != nil {
			log.Fatalln(err)
		}
	}*/

	return &Config{
		Provider:    provider,
		Provisioner: "external/pharmer",
		MasterUrl:   "",
		Kubeconfig:  "",
	}
}

func (c *Config) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&c.Provisioner, "provisioner", c.Provisioner, "Name of the provisioner. The provisioner will only provision volumes for claims that request a StorageClass with a provisioner field set equal to this name.")
	fs.StringVar(&c.MasterUrl, "master", c.MasterUrl, "Master URL to build a client config from. Either this or kubeconfig needs to be set if the provisioner is being run out of cluster.")
	fs.StringVar(&c.Kubeconfig, "kubeconfig", c.Kubeconfig, "Absolute path to the kubeconfig file. Either this or master needs to be set if the provisioner is being run out of cluster.")
}

func (c *Config) ValidateFlags() error {
	allErrs := field.ErrorList{}
	fldPath := field.NewPath("provisioner")
	if len(c.Provisioner) == 0 {
		allErrs = append(allErrs, field.Required(fldPath, c.Provisioner))
	}
	if len(c.Provisioner) > 0 {
		for _, msg := range validation.IsQualifiedName(strings.ToLower(c.Provisioner)) {
			allErrs = append(allErrs, field.Invalid(fldPath, c.Provisioner, msg))
		}
	}
	if len(allErrs) != 0 {
		glog.Fatalf("Invalid provisioner specified: %v", allErrs)
	}
	glog.Infof("Provisioner %s specified", c.Provisioner)
	return nil
}
