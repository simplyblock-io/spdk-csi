### Storage nodes volume provisioning

We now also have the concept of storage-nodes on kubernetes. These will be Kubernetes nodes which reside within the kubernetes compute cluster as normal worker nodes and can have any PODs deployed. However, they have either NVMe disk locally attached with partitions or ebs volumes without partitions.


### Preparing nodes

#### Step 0: Networking & tools

Make sure that the Kubernetes worker nodes used as storage-node has access to the simplyblock cluster. If you are using terraform to deploy the cluster. Please attach `container-instance-sg` security group to all the instances.

#### Step1: Install nvme cli tools and nbd

To attach NVMe device to the host machine, the CSI driver uses [nvme-cli]([url](https://github.com/linux-nvme/nvme-cli)). So lets install that
```
sudo yum install -y nvme-cli
sudo modprobe nvme-tcp
sudo modprobe nbd
```

#### Step1: Setup hugepages

Before you prepare the storage nodes, please decide the amount of huge pages that you would like to allocate for simplyblock and set those hugepages accordingly. We suggest allocating at least 8GB of huge pages. 

>[!IMPORTANT]
>The storage node requires at least 2.2% of the size of the nvme cache + 50 MiB of RAM. This should be the minimum configured as hugepage
>memory.

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

#### Step2: Mount the SSD or EBS to be used by the storage node
If the instance comes with a default NVMe disk, it can be used with minimum of 2 partitions and 2 device where one is used for Journal manager and the other storage node. Or 2 additional EBS one for Journal Manager and the other for the Storage. the disks can be viewed by running:

```
sudo yum install pciutils
lspci
```


#### Step3: Tag the kubernetes nodes

After the nodes are prepared, label the kubernetes nodes
```
kubectl label nodes ip-10-0-4-118.us-east-2.compute.internal ip-10-0-4-176.us-east-2.compute.internal type=simplyblock-storage-plane
```
Now the nodes are ready to deploy storage nodes.
