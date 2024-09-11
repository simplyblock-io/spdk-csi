### Caching nodes volume provisioning

We now also have the concept of caching-nodes. These will be Kubernetes nodes which reside within the kubernetes compute cluster as normal worker nodes and can have any PODs deployed. However, they have an NVMe disk locally attached. If there are multiple NVMe disks attached, it will use the first one.


### Preparing nodes
Caching nodes are a special kind of node that works as a cache with a local NVMe disk.

#### Step 0: Networking & tools

Make sure that the Kubernetes worker nodes to be used for cache has access to the simplyblock storage cluster. If you are using terraform to deploy the cluster. Please attach `container-instance-sg` security group to all the instances.

#### Step1: Install nvme cli tools and nbd

To attach NVMe device to the host machine, the CSI driver uses [nvme-cli]([url](https://github.com/linux-nvme/nvme-cli)). So lets install that
```
sudo yum install -y nvme-cli
sudo modprobe nvme-tcp
sudo modprobe nbd
```

#### Step1: Setup hugepages

Before you prepare the caching nodes, please decide the amount of huge pages that you would like to allocate for simplyblock and set those hugepages accordingly. 
It is recommended to use a minimum of 1 GiB + 0.5% of the size of the local SSD, which you want to use as a cache. For example, if your local SSD has a size of 1.9 TiB, and you want to use it entirely as a write-through cache, you need to assign 10.5 GiB of RAM. If you only want to utilize 1 TiB (52.9% of the SSD), you assign 6 GiB of RAM and the cache will be automatically resized to fit the available (assigned) memory. 

>[!IMPORTANT]
>One huge page contains 2 MiB of memory. A value of e.g. 4096 therefore is equal to 8 GiB of huge page memory.

```
sudo sysctl -w vm.nr_hugepages=4096
```

confirm the hugepage changes by running
cat /proc/meminfo | grep -i hug


and restart kubelet
```
sudo systemctl restart kubelet
```

conform if huge pages are added to the cluster or not.
```
kubectl describe node ip-10-0-2-184.us-east-2.compute.internal | grep hugepages-2Mi
```
this output should show 8GB. This worker node can allocate 8GB of hugepages to pods which is required in case of SPDK pods.

#### Step2: Mount the SSD to be used for caching
If the instance comes with a default NVMe disk, it can be used. Or an additional EBS or SSD volume can be mounted. the disks can be viewed by running:

```
sudo yum install pciutils
lspci
```


#### Step3: Tag the kubernetes nodes

After the nodes are prepared, label the kubernetes nodes
```
kubectl label nodes ip-10-0-4-118.us-east-2.compute.internal ip-10-0-4-176.us-east-2.compute.internal type=simplyblock-cache
```
Now the nodes are ready to deploy caching nodes.

### StorageClass

If the user wants to create a PVC that uses NVMe cache, a new storage class can be used with additional volume parameter as `type: cache`.


### Usage and Implementation

During dynamic volume provisioning, nodeSelector should be provided on pod, deployment, daemonset, statefulset. So that such pods are scheduled only on the nodes that has the `simplyblock-cache` label on it.

As shown below
```
    spec:
      nodeSelector:
        type: simplyblock-cache
```

On the controller server, when a new volume is requested, we create a `lvol` . This steps is exactly same as the current implementation.

On the node driver, during the volume mount, the following steps happens.
1. Get the caching node ID of the current node
2. Connect caching node with lvol. This will create a new NVMe device on the host machine. This device will be used to mount into pod.
