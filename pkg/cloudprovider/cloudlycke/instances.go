package cloudlycke

import (
	"context"
	"net/http"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	cloudprovider "k8s.io/cloud-provider"
	"k8s.io/klog"
)

type instances struct {
	client *http.Client
}

func newInstances(c *http.Client) cloudprovider.Instances {
	return &instances{
		client: c,
	}
}

// NodeAddresses returns the addresses of the specified instance.
func (i *instances) NodeAddresses(ctx context.Context, name types.NodeName) ([]v1.NodeAddress, error) {
	klog.V(5).Infof("NodeAddresses(%v)", name)

	var addrs []v1.NodeAddress

	if string(name) == "master-c2-1" {
		nodeAddr := v1.NodeAddress{
			Type:    v1.NodeInternalIP,
			Address: "192.168.20.10",
		}
		nodeExternalAddr := v1.NodeAddress{
			Type:    v1.NodeExternalIP,
			Address: "192.168.20.10",
		}
		nodeHostName := v1.NodeAddress{
			Type:    v1.NodeHostName,
			Address: "master-c2-1",
		}
		addrs = append(addrs, nodeAddr)
		addrs = append(addrs, nodeExternalAddr)
		addrs = append(addrs, nodeHostName)
	} else if string(name) == "node-c2-1" {
		nodeAddr := v1.NodeAddress{
			Type:    v1.NodeInternalIP,
			Address: "192.168.20.11",
		}
		nodeExternalAddr := v1.NodeAddress{
			Type:    v1.NodeExternalIP,
			Address: "192.168.20.11",
		}
		nodeHostName := v1.NodeAddress{
			Type:    v1.NodeHostName,
			Address: "node-c2-1",
		}

		addrs = append(addrs, nodeAddr)
		addrs = append(addrs, nodeExternalAddr)
		addrs = append(addrs, nodeHostName)
	}

	return addrs, nil
}

// NodeAddressesByProviderID returns the addresses of the specified instance.
// The instance is specified using the providerID of the node. The
// ProviderID is a unique identifier of the node. This will not be called
// from the node whose nodeaddresses are being queried. i.e. local metadata
// services cannot be used in this method to obtain nodeaddresses
func (i *instances) NodeAddressesByProviderID(ctx context.Context, providerID string) ([]v1.NodeAddress, error) {
	klog.V(5).Infof("NodeAddressesByProviderID(%v)", providerID)

	// TODO: Add split function to get the instance ID from provider ID and do a "lookup"

	var addrs []v1.NodeAddress

	if providerID == "cloudlycke://m-c2-1" {
		nodeAddr := v1.NodeAddress{
			Type:    v1.NodeInternalIP,
			Address: "192.168.20.10",
		}
		nodeExternalAddr := v1.NodeAddress{
			Type:    v1.NodeExternalIP,
			Address: "192.168.20.10",
		}
		nodeHostName := v1.NodeAddress{
			Type:    v1.NodeHostName,
			Address: "master-c2-1",
		}

		addrs = append(addrs, nodeAddr)
		addrs = append(addrs, nodeExternalAddr)
		addrs = append(addrs, nodeHostName)
	} else if providerID == "cloudlycke://n-c2-1" {
		nodeAddr := v1.NodeAddress{
			Type:    v1.NodeInternalIP,
			Address: "192.168.20.11",
		}
		nodeExternalAddr := v1.NodeAddress{
			Type:    v1.NodeExternalIP,
			Address: "192.168.20.11",
		}
		nodeHostName := v1.NodeAddress{
			Type:    v1.NodeHostName,
			Address: "node-c2-1",
		}
		addrs = append(addrs, nodeAddr)
		addrs = append(addrs, nodeExternalAddr)
		addrs = append(addrs, nodeHostName)
	}

	return addrs, nil
}

// InstanceID returns the cloud provider ID of the node with the specified NodeName.
// Note that if the instance does not exist, we must return ("", cloudprovider.InstanceNotFound)
// cloudprovider.InstanceNotFound should NOT be returned for instances that exist but are stopped/sleeping
func (i *instances) InstanceID(ctx context.Context, nodeName types.NodeName) (string, error) {
	klog.V(5).Infof("InstanceID(%v)", nodeName)

	var instanceID string

	if string(nodeName) == "master-c2-1" {
		instanceID = "cloudlycke://m-c2-1"
	} else if string(nodeName) == "node-c2-1" {
		instanceID = "cloudlycke://n-c2-1"
	}

	return instanceID, nil
}

// InstanceType returns the type of the specified instance.
func (i *instances) InstanceType(ctx context.Context, name types.NodeName) (string, error) {
	klog.V(5).Infof("InstanceType(%v)", name)

	var instanceType string

	if string(name) == "master-c2-1" {
		instanceType = "vbox.vm.512mb.1cpu"
	} else if string(name) == "node-c2-1" {
		instanceType = "vbox.vm.1g.2cpu"
	}

	return instanceType, nil
}

// InstanceTypeByProviderID returns the type of the specified instance.
func (i *instances) InstanceTypeByProviderID(ctx context.Context, providerID string) (string, error) {
	klog.V(5).Infof("InstanceTypeByProviderID(%v)", providerID)

	var instanceType string

	if providerID == "cloudlycke://m-c2-1" {
		instanceType = "vbox.vm.512mb.1cpu"
	} else if providerID == "cloudlycke://n-c2-1" {
		instanceType = "vbox.vm.1g.2cpu"
	}

	return instanceType, nil
}

// AddSSHKeyToAllInstances adds an SSH public key as a legal identity for all instances
// expected format for the key is standard ssh-keygen format: <protocol> <blob>
func (i *instances) AddSSHKeyToAllInstances(ctx context.Context, user string, keyData []byte) error {
	klog.V(5).Info("AddSSHKeyToAllInstances(%v, %v)", user, keyData)
	return cloudprovider.NotImplemented
}

// CurrentNodeName returns the name of the node we are currently running on
// On most clouds (e.g. GCE) this is the hostname, so we provide the hostname
func (i *instances) CurrentNodeName(ctx context.Context, hostname string) (types.NodeName, error) {
	klog.V(5).Infof("CurrentNodeName(%v)", hostname)

	var nodeName types.NodeName

	if hostname == "master-c2-1" {
		nodeName = "master-c2-1"
	} else if hostname == "node-c2-1" {
		nodeName = "node-c2-1"
	}

	return nodeName, nil
}

// InstanceExistsByProviderID returns true if the instance for the given provider exists.
// If false is returned with no error, the instance will be immediately deleted by the cloud controller manager.
// This method should still return true for instances that exist but are stopped/sleeping.
func (i *instances) InstanceExistsByProviderID(ctx context.Context, providerID string) (bool, error) {
	klog.V(5).Infof("InstanceExistsByProviderID(%v)", providerID)

	var exists bool

	if providerID == "cloudlycke://m-c2-1" {
		exists = true
	} else if providerID == "cloudlycke://n-c2-1" {
		exists = true
	}

	return exists, nil
}

// InstanceShutdownByProviderID returns true if the instance is shutdown in cloudprovider
func (i *instances) InstanceShutdownByProviderID(ctx context.Context, providerID string) (bool, error) {
	klog.V(5).Infof("InstanceShutdownByProviderID(%v)", providerID)

	var shutdown bool

	if providerID == "cloudlycke://m-c2-1" {
		shutdown = false
	} else if providerID == "cloudlycke://n-c2-1" {
		shutdown = true
	}

	return shutdown, nil
}