package cloudlycke

import (
	"context"
	"net/http"

	"k8s.io/apimachinery/pkg/types"
	cloudprovider "k8s.io/cloud-provider"
	"k8s.io/klog"
)

type zones struct {
	client *http.Client
}

func newZones(c *http.Client) cloudprovider.Zones {
	return &zones{
		c,
	}
}

// GetZone returns the Zone containing the current failure zone and locality region that the program is running in
// In most cases, this method is called from the kubelet querying a local metadata service to acquire its zone.
// For the case of external cloud providers, use GetZoneByProviderID or GetZoneByNodeName since GetZone
// can no longer be called from the kubelets.
func (z *zones) GetZone(ctx context.Context) (cloudprovider.Zone, error) {
	klog.V(5).Info("GetZone()")
	return cloudprovider.Zone{
		FailureDomain: "laptop",
		Region:        "virtualbox",
	}, nil
}

// GetZoneByProviderID returns the Zone containing the current zone and locality region of the node specified by providerID
// This method is particularly used in the context of external cloud providers where node initialization must be done
// outside the kubelets.
func (z *zones) GetZoneByProviderID(ctx context.Context, providerID string) (cloudprovider.Zone, error) {
	klog.V(5).Infof("GetZoneByProviderID(%v)", providerID)
	return cloudprovider.Zone{
		FailureDomain: "virtualbox",
		Region:        "virtualbox",
	}, nil
}

// GetZoneByNodeName returns the Zone containing the current zone and locality region of the node specified by node name
// This method is particularly used in the context of external cloud providers where node initialization must be done
// outside the kubelets.
func (z *zones) GetZoneByNodeName(ctx context.Context, nodeName types.NodeName) (cloudprovider.Zone, error) {
	klog.V(5).Infof("GetZoneByNodeName(%v)", nodeName)
	return cloudprovider.Zone{
		FailureDomain: "virtualbox",
		Region:        "virtualbox",
	}, nil
}