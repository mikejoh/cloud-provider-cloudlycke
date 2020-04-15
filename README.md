# The Cloudlycke Cloud Controller manager

This repository contains the `cloudlycke-cloud-controller-manager` an out-of-tree built, cloud controller that is encapsulates the business logic of my fake external cloud provider called Cloudlycke. The purpose of this repository is to show how the Cloud Controller works in detail and can be used as an example to create your own cloud controller to be used with your cloud provider.

See the following blog post for more in-depth information and discussion: **ADD LINK HERE**

## TODO:
* Successfully started the `cloudlycke-cloud-controller-manager` in `minikube`, gives us the following log at the moment:
```
I0414 08:16:33.845655       1 node_controller.go:110] Sending events to api server.                                                           
I0414 08:16:33.846444       1 cloud.go:55] Instances()                                                                                        
I0414 08:16:33.846708       1 controllermanager.go:247] Started "cloud-node"                                                                  
I0414 08:16:33.847005       1 controllermanager.go:237] Starting "cloud-node-lifecycle"                                                       
I0414 08:16:33.846964       1 cloud.go:55] Instances()                                                                                        
I0414 08:16:33.851153       1 instances.go:43] InstanceID()                                                                                   
W0414 08:16:33.851358       1 node_controller.go:527] Cannot find valid providerID for node name "minikube", assuming non existence           
I0414 08:16:33.851456       1 node_controller.go:239] The node minikube is no longer present according to the cloud provider, do not process. 
I0414 08:16:33.851610       1 cloud.go:55] Instances()                                                                                        
I0414 08:16:33.854183       1 instances.go:43] InstanceID()                                                                                   
W0414 08:16:33.854784       1 node_controller.go:527] Cannot find valid providerID for node name "minikube", assuming non existence           
I0414 08:16:33.855022       1 node_controller.go:239] The node minikube is no longer present according to the cloud provider, do not process.
``` 