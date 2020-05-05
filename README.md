# The Cloudlycke Kubernetes Cloud Controller

This repository contains the Cloudlycke `cloud-controller`, an [out-of-tree](https://kubernetes.io/blog/2019/04/17/the-future-of-cloud-providers-in-kubernetes/) and [by-the-book](https://kubernetes.io/docs/tasks/administer-cluster/developing-cloud-controller-manager/#out-of-tree) built Kubernetes cloud controller that implements the `cloud-provider` [Interface](https://github.com/kubernetes/cloud-provider/blob/v0.18.2/cloud.go#L43-L62).

This `cloud-controller` is built using the `v1.18.x` release of Kubernetes. This means that `v1.18.x` is used everywhere we have dependencies on Kubernetes. 

Cloudlycke is my fake cloud provider, which at the moment is Vagrant. I wanted my Kubernetes clusters in this cloud provider to be able to integrate with the underlying cloud. For the purpose of showing the ins and outs of the Kubernetes `cloud-controller` most of the API "calls" to Cloudlycke is hardcoded to respond with a particular response.

I've written an in-depth [write-up](**ADD LINK HERE**) that explains and explores the Cloud Controller.

Inspired by the DigitalOcean and OpenStack Cloud Controllers!

## Detailed overview

![cloudlycke-cloud-controller](img/cloudlycke-cloud-controller.svg)

The environment consists of the following components:
* `Vagrant`
* `Ansible`
* `VirtualBox`

Vagrant will be used to provision the virtual machines ontop of VirtualBox, on these VMs we'll deploy two Kubernetes clusters with one all-in-one master node and one worker node each.

Ansible will be used with `vagrant` during provisioning, included in this repository there's two Ansible playbooks (and other ansible specific resources) located [here](vagrant/ansible).

The first cluster will be deplyed as-is and the second one will be configured in such a way that we'll need a cloud controller to initialize the worker node(s). The magic sauce is the following configuration:

* The API server will be configured with the following flag(s): `--cloud-provider=external`
* The Controller Manager will be configured wth the following flag(s): `--cloud-provider=external`
* The Kubelets will be configured with the following flag(s): `--node-ip <VM IP> --cloud-provider=external --provider-id=cloudlycke://<ID>`
* The Cloudlycke Cloud Controller will be configured with the following flag(s): `--cloud-provider=cloudlycke`

  _Please note that the container image, used in the [all-in-one manifest](/manifests/cloudlycke-ccm.yaml), is one that i've built and pushed to my private Docker Hub repository. Please see the [Dockerfile](/Dockerfile) to see how the image was built._

## Starting the Vagrant (cloud) environment and deploy Kubernetes

1. Install Ansible in a `virtualenv` and activate the environment
2. Run `vagrant up`
3. When `ansible` and `vagrant` is done check the `artifacts/` directory, you should have two kubeconfigs there called `admin-master-c1-1.conf` and `admin-master-c2-1.conf`

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
  
  The responsible controller for these operations are the (Cloud) Node Controller.

7. 

### References

