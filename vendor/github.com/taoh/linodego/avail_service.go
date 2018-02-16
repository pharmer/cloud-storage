package linodego

import (
	"encoding/json"
	"net/url"
	"strconv"
)

// AvailService is an interface generated for "github.com/appscode/linodego".AvailService.
type AvailInterface interface {
	DataCenters() (*AvailDataCentersResponse, error)
	Distributions() (*AvailDistributionsResponse, error)
	FilterKernels(int, int) (*KernelsResponse, error)
	Kernels(map[string]string) (*KernelsResponse, error)
	LinodePlans() (*LinodePlansResponse, error)
	NodeBalancers() (*NodeBalancersResponse, error)
	StackScripts() (*StackScriptsResponse, error)
}

// Avail Service
type AvailService struct {
	client *Client
}

var _ AvailInterface = &AvailService{}

// Response for avail.datacenters API
type AvailDataCentersResponse struct {
	Response
	DataCenters []DataCenter
}

// Response for avail.distributions API
type AvailDistributionsResponse struct {
	Response
	Distributions []Distribution
}

// Response for avail.kernels API
type KernelsResponse struct {
	Response
	Kernels []Kernel
}

// Response for avail.linodeplans API
type LinodePlansResponse struct {
	Response
	LinodePlans []LinodePlan
}

// Response for avail.nodebalancers API
type NodeBalancersResponse struct {
	Response
	NodeBalancers []NodeBalancer
}

// Response for avail.stackscripts API
type StackScriptsResponse struct {
	Response
	StackScripts []StackScript
}

// Get DataCenters
func (t *AvailService) DataCenters() (*AvailDataCentersResponse, error) {
	u := &url.Values{}
	v := AvailDataCentersResponse{}
	if err := t.client.do("avail.datacenters", u, &v.Response); err != nil {
		return nil, err
	}
	v.DataCenters = make([]DataCenter, 0)
	if err := json.Unmarshal(v.RawData, &v.DataCenters); err != nil {
		return nil, err
	}
	return &v, nil
}

// Get Distributions
func (t *AvailService) Distributions() (*AvailDistributionsResponse, error) {
	u := &url.Values{}
	v := AvailDistributionsResponse{}
	if err := t.client.do("avail.distributions", u, &v.Response); err != nil {
		return nil, err
	}
	v.Distributions = make([]Distribution, 0)
	if err := json.Unmarshal(v.RawData, &v.Distributions); err != nil {
		return nil, err
	}
	return &v, nil
}

// Get Kernels
func (t *AvailService) Kernels(args map[string]string) (*KernelsResponse, error) {
	u := &url.Values{}
	// add optional parameters
	processOptionalArgs(args, u)
	v := KernelsResponse{}
	if err := t.client.do("avail.kernels", u, &v.Response); err != nil {
		return nil, err
	}
	v.Kernels = make([]Kernel, 0)
	if err := json.Unmarshal(v.RawData, &v.Kernels); err != nil {
		return nil, err
	}
	return &v, nil
}

// Get filtered Kernels
func (t *AvailService) FilterKernels(isxen int, iskvm int) (*KernelsResponse, error) {
	params := &url.Values{}
	v := KernelsResponse{}
	xen_s := strconv.Itoa(isxen)
	kvm_s := strconv.Itoa(iskvm)
	if isxen < 2 && iskvm < 2 {
		params.Add("isxen", xen_s)
		params.Add("iskvm", kvm_s)
	} else if isxen < 2 {
		params.Add("isxen", xen_s)
	} else if iskvm < 2 {
		params.Add("iskvm", kvm_s)
	}

	if err := t.client.do("avail.kernels", params, &v.Response); err != nil {
		return nil, err
	}
	v.Kernels = make([]Kernel, 5)
	if err := json.Unmarshal(v.RawData, &v.Kernels); err != nil {
		return nil, err
	}
	return &v, nil
}

// Get Linode Plans
func (t *AvailService) LinodePlans() (*LinodePlansResponse, error) {
	u := &url.Values{}
	v := LinodePlansResponse{}
	if err := t.client.do("avail.linodeplans", u, &v.Response); err != nil {
		return nil, err
	}
	v.LinodePlans = make([]LinodePlan, 0)
	if err := json.Unmarshal(v.RawData, &v.LinodePlans); err != nil {
		return nil, err
	}
	return &v, nil
}

// Get Node Balancers
func (t *AvailService) NodeBalancers() (*NodeBalancersResponse, error) {
	u := &url.Values{}
	v := NodeBalancersResponse{}
	if err := t.client.do("avail.nodebalancers", u, &v.Response); err != nil {
		return nil, err
	}
	v.NodeBalancers = make([]NodeBalancer, 0)
	if err := json.Unmarshal(v.RawData, &v.NodeBalancers); err != nil {
		return nil, err
	}
	return &v, nil
}

// Get All Stackscripts
func (t *AvailService) StackScripts() (*StackScriptsResponse, error) {
	u := &url.Values{}
	v := StackScriptsResponse{}
	if err := t.client.do("avail.stackscripts", u, &v.Response); err != nil {
		return nil, err
	}
	v.StackScripts = make([]StackScript, 0)
	if err := json.Unmarshal(v.RawData, &v.StackScripts); err != nil {
		return nil, err
	}
	return &v, nil
}
