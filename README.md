# The Cloudlycke Cloud Controller Manager

This repository contains the Cloudlycke Cloud Controller Manager, an [out-of-tree](https://kubernetes.io/blog/2019/04/17/the-future-of-cloud-providers-in-kubernetes/) and [by-the-book](https://kubernetes.io/docs/tasks/administer-cluster/developing-cloud-controller-manager/#out-of-tree) built Kubernetes cloud controller that implements the `k8s.io/cloud-provider` [Interface](https://github.com/kubernetes/cloud-provider/blob/v0.18.2/cloud.go#L43-L62). 

This `cloud-controller-manager` is built using the `v1.18.x` release of Kubernetes. This means that `v1.18.x` is used everywhere we have dependencies on Kubernetes. 

Cloudlycke is my cloud provider, which is backed by Vagrant. Not to bad, huh? 

I wanted my Kubernetes clusters in this cloud provider to be able to integrate with the underlying cloud. Mainly to show you the ins and outs of the Kubernetes `cloud-controller-manager`.
 
 All of the API calls to the Cloudlycke cloud provider is hardcoded to respond with a particular response to fit the scenarios. It does *not* communicate with Vagrant in any way, but it looks like that anyways.

I've written an in-depth [write-up](https://medium.com/@m.json/the-kubernetes-cloud-controller-manager-d440af0d2be5) that explains and explores the Cloud Controller Manager, from more of a theoretical and source code level.

If this is of any kind of interest to you and if you've spotted something that just isn't correct, please feel free to contribute with issues and PRs!

## Todo

* Implement the `LoadBalancer()` interface methods to show how that would look like.

## Detailed overview

![cloudlycke-cloud-controller](img/cloudlycke-cloud-controller.svg)

The environment consists of the following components:
* `vagrant`
* `ansible`, used as the provisioner in `vagrant`
* `VirtualBox`, Hypervisor

Vagrant will be used to provision the virtual machines ontop of VirtualBox, on these VMs we'll deploy two Kubernetes clusters with one all-in-one master node and one worker node each.

Ansible will be used with `vagrant` during provisioning, included in this repository there's two Ansible playbooks (and other ansible specific resources) located [here](vagrant/ansible).

The first cluster will be deplyed as-is and the second one will be configured in such a way that we'll need a cloud controller to initialize the k8s node(s). Needed configuration of the k8s control plane components:

* The API server will be configured with the following flag(s): `--cloud-provider=external`. This is not needed, but since there's still code in the API server that does cloud provider specific method calls ([#1](https://github.com/kubernetes/kubernetes/blob/9e991415386e4cf155a24b1da15becaa390438d8/cmd/kube-apiserver/app/server.go#L235) [#2](https://github.com/kubernetes/kubernetes/blob/9e991415386e4cf155a24b1da15becaa390438d8/cmd/kube-apiserver/app/server.go#L241)) i'll leave it here as documentation.
* The Controller Manager will be configured wth the following flag(s): `--cloud-provider=external`
* The Kubelets will be configured with the following flag(s): `--node-ip <VM IP> --cloud-provider=external --provider-id=cloudlycke://<ID>`. I added the `provider-id` flag to force the `kubelet` to set that on node initialization since i don't have something like a instance metadata service to query. Although that can definately be built in or hard coded basically.
* The Cloudlycke Cloud Controller will be configured with the following flag(s): `--cloud-provider=cloudlycke`

  _Please note that the container image, used in the [all-in-one manifest](/manifests/cloudlycke-ccm.yaml), is one that i've built and pushed to my private Docker Hub repository. Please see the [Dockerfile](/Dockerfile) to see how the image was built._

## Starting the Vagrant (cloud) environment and deploy Kubernetes

1. Install Ansible in a `virtualenv` and activate the environment
2. Run `vagrant up`
3. When `ansible` and `vagrant` is done check the `artifacts/` directory, you should have two kubeconfigs there called `admin-master-c1-1.conf` and `admin-master-c2-1.conf`. Basically one for each Kubernetes cluster.

## Running the Cloud Controller

1. Export the kubeconfig(s) `export KUBECONFIG=<PATH TO admin-master-c2-1.conf>`
2. Check the current status of the cluster nodes
```
kubectl get nodes

NAME          STATUS   ROLES    AGE   VERSION
master-c2-1   Ready    master   24m   v1.18.2
node-c2-1     Ready    <none>   19m   v1.18.2
```
3. Deploy nginx pods (deployment with 3 replicas) for demo purposes
``` 
kubectl run --image nginx --replicas 3 nginx-demo
```
4. Check the status of all pods across all namespaces
```
kubectl get pods -A
NAMESPACE     NAME                                  READY   STATUS    RESTARTS   AGE   IP              NODE          NOMINATED NODE   READINESS GATES
default       nginx-demo-5756474c97-m4b9t           0/1     Pending   0          25m   <none>          <none>        <none>           <none>
default       nginx-demo-5756474c97-qqjg4           0/1     Pending   0          25m   <none>          <none>        <none>           <none>
default       nginx-demo-5756474c97-rbhvt           0/1     Pending   0          25m   <none>          <none>        <none>           <none>
kube-system   coredns-66bff467f8-cthqz              0/1     Pending   0          32m   <none>          <none>        <none>           <none>
kube-system   coredns-66bff467f8-j24c2              0/1     Pending   0          32m   <none>          <none>        <none>           <none>
kube-system   etcd-master-c2-1                      1/1     Running   0          32m   192.168.20.10   master-c2-1   <none>           <none>
kube-system   kube-apiserver-master-c2-1            1/1     Running   0          32m   192.168.20.10   master-c2-1   <none>           <none>
kube-system   kube-controller-manager-master-c2-1   1/1     Running   0          32m   192.168.20.10   master-c2-1   <none>           <none>
kube-system   kube-flannel-ds-amd64-qhbdp           1/1     Running   0          32m   192.168.20.10   master-c2-1   <none>           <none>
kube-system   kube-flannel-ds-amd64-r22f7           1/1     Running   1          27m   192.168.20.11   node-c2-1     <none>           <none>
kube-system   kube-proxy-bqn2b                      1/1     Running   0          32m   192.168.20.10   master-c2-1   <none>           <none>
kube-system   kube-proxy-rmwx7                      1/1     Running   0          27m   192.168.20.11   node-c2-1     <none>           <none>
kube-system   kube-scheduler-master-c2-1            1/1     Running   0          32m   192.168.20.10   master-c2-1   <none>           <none>
```
Note that some of the pods are reporting status `Pending`. The ones that are running are primarily the `DaemonSet` created ones and the ones with toleration configured that allows them to be scheduled e.g. `node-role.kubernetes.io/master: ""`.

5. `kubectl describe pods nginx-demo-5756474c97-m4b9t`
```
Events:
  Type     Reason            Age                From               Message
  ----     ------            ----               ----               -------
  Warning  FailedScheduling  20s (x2 over 20s)  default-scheduler  0/2 nodes are available: 1 node(s) had taint {node-role.kubernetes.io/master: }, that the pod didn't tolerate, 1 node(s) had taint {node.cloudprovider.kubernetes.io/uninitialized: true}, that the pod didn't tolerate.
``` 
Note that the master node `master-c2-1` will be tainted and only allow pods with the correct toleration. The worker node `node-c2-1` is still awaiting initialization of our external cloud provider controller.

6. Now install the Cloudlycke CCM, before you'll do that you can do the following:
  * Take a note of the Node `node-c2-1` labels.
  * Take a note of the Node `node-c2-1` taints.
  ```
  kubectl apply -f mainfests/cloudlycke-ccm.yaml 
  ```
   Immediately after you're done applying the manifest(s) please tail the log of the deployed Cloudlycke CCM Pod for more info:
  ```
  kubectl logs -n kube-system -l k8s-app=cloudlycke-cloud-controller-manager -f
  ```
  
  You now should've observed at least three things about the worker node `node-c2-1` at least:
  * The taint `node.cloudprovider.kubernetes.io/uninitialized` have been removed
  * The node now have a couple of more labels with information about the node given from the cloud provider, these should be:
    ```
    ...
    labels:
      ...
      beta.kubernetes.io/instance-type: vbox.vm.1g.2cpu
      failure-domain.beta.kubernetes.io/region: virtualbox
      failure-domain.beta.kubernetes.io/zone: virtualbox
      node.kubernetes.io/instance-type: vbox.vm.1g.2cpu
      topology.kubernetes.io/region: virtualbox
      topology.kubernetes.io/zone: virtualbox
      ...
    ```
    If you wonder why there's beta labels in there you can track the promotion of cloud provider labels to GA at [this](https://github.com/kubernetes/enhancements/issues/837) issue, here's a KEP defining [standard topology labels](https://github.com/kubernetes/enhancements/pull/1660) that also might be of interest.
  * The Nginx Pods that earlier were in `Pending` state now should be `Running`, this is a consequence of that taint being removed.
  
  Regarding the labels, note that e.g. the OpenStack external CCM uses the [Nova (OpenStack compute service) instance metadata](https://github.com/kubernetes/cloud-provider-openstack/blob/v1.18.0/pkg/util/metadata/metadata.go) to add the instance information to the `Node` labels. 
  
  The responsible controller for these operations are the (Cloud) Node Controller.
  
That's basically it for now! There's a bunch of things that i haven't implemented yet in the CCM (like the `LoadBalancer()` methods), but the very basics are in place and observable.

### References

* [Developing Cloud Controller Manager](https://kubernetes.io/docs/tasks/administer-cluster/developing-cloud-controller-manager/)
* [Cloud Controller Manager Administration](https://kubernetes.io/docs/tasks/administer-cluster/running-cloud-controller/)
* [DigitalOcean `cloud-controller-manager`](https://github.com/digitalocean/digitalocean-cloud-controller-manager)
* [OpenStack `cloud-controller-manager`](https://github.com/kubernetes/cloud-provider-openstack)
