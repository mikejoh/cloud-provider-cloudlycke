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

	switch name {
	case "master-c2-1":
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
	case "node-c2-n1":
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
			Address: "node-c2-n1",
		}

		addrs = append(addrs, nodeAddr)
		addrs = append(addrs, nodeExternalAddr)
		addrs = append(addrs, nodeHostName)
	default:
		klog.V(5).Info("NodeAddresses switch failed!")
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

	var addrs []v1.NodeAddress

	switch providerID {
	case "cloudlycke://1":
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
	case "cloudlycke://2":
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
			Address: "node-c2-n1",
		}

		addrs = append(addrs, nodeAddr)
		addrs = append(addrs, nodeExternalAddr)
		addrs = append(addrs, nodeHostName)
	default:
		klog.V(5).Infof("ProviderID is %v", providerID)
	}

	return addrs, nil
}

// InstanceID returns the cloud provider ID of the node with the specified NodeName.
// Note that if the instance does not exist, we must return ("", cloudprovider.InstanceNotFound)
// cloudprovider.InstanceNotFound should NOT be returned for instances that exist but are stopped/sleeping
func (i *instances) InstanceID(ctx context.Context, nodeName types.NodeName) (string, error) {
	klog.V(5).Infof("InstanceID(%v)", nodeName)

	var instanceID string

	switch nodeName {
	case "master-c2-1":
		instanceID = "1"
	case "node-c2-n1":
		instanceID = "2"
	default:
		klog.V(5).Info("InstanceID switch failed!")
	}

	return instanceID, nil
}

// InstanceType returns the type of the specified instance.
func (i *instances) InstanceType(ctx context.Context, name types.NodeName) (string, error) {
	klog.V(5).Infof("InstanceType(%v)", name)

	var instanceType string

	switch name {
	case "master-c2-1":
		instanceType = "vbox.vm.512mb.1cpu"
	case "node-c2-n1":
		instanceType = "vbox.vm.1g.2cpu"
	default:
		klog.V(5).Info("InstanceType switch failed!")
	}

	return instanceType, nil
}

// InstanceTypeByProviderID returns the type of the specified instance.
func (i *instances) InstanceTypeByProviderID(ctx context.Context, providerID string) (string, error) {
	klog.V(5).Infof("InstanceTypeByProviderID(%v)", providerID)

	var instanceType string

	switch providerID {
	case "cloudlycke://1":
		instanceType = "vbox.vm.512mb.1cpu"
	case "cloudlycke://2":
		instanceType = "vbox.vm.1g.2cpu"
	default:
		klog.V(5).Info("InstanceTypeByProviderID switch failed!")
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

	switch hostname {
	case "master-c2-1":
		nodeName = "master-c2-1"
	case "node-c2-n1":
		nodeName = "node-c2-n1"
	}

	return nodeName, nil
}

// InstanceExistsByProviderID returns true if the instance for the given provider exists.
// If false is returned with no error, the instance will be immediately deleted by the cloud controller manager.
// This method should still return true for instances that exist but are stopped/sleeping.
func (i *instances) InstanceExistsByProviderID(ctx context.Context, providerID string) (bool, error) {
	klog.V(5).Infof("InstanceExistsByProviderID(%v)", providerID)

	var exists bool

	switch providerID {
	case "cloudlycke://1":
		exists = true
	case "cloudlycke://2":
		exists = true
	}

	return exists, nil
}

// InstanceShutdownByProviderID returns true if the instance is shutdown in cloudprovider
func (i *instances) InstanceShutdownByProviderID(ctx context.Context, providerID string) (bool, error) {
	klog.V(5).Infof("InstanceShutdownByProviderID(%v)", providerID)

	var shutdown bool

	switch providerID {
	case "cloudlycke://1":
		shutdown = false
	case "cloudlycke://2":
		shutdown = false
	}

	return shutdown, nil
}