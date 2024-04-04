### Caching nodes volume provisioning

We now also have the concept of caching-nodes. These will be Kubernetes nodes which reside within the kubernetes compute cluster as normal worker nodes and can have any PODs deployed. However, they have an NVMe disk locally attached. If there are multiple NVMe disks attached, it will use the first one.


### Preparing nodes
Caching nodes are a special kind of node that works as a cache with a local NVMe disk.

#### Step 0: Networking & tools

Make sure that the Kubernetes worker nodes to be used for cache has access to the simplyblock storage cluster. If you are using terraform to deploy the cluster. Please attach `container-instance-sg` security group to all the instances.

#### Step1: Install nvme cli tools

To attach NVMe device to the host machine, the CSI driver uses [nvme-cli]([url](https://github.com/linux-nvme/nvme-cli)). So lets install that
```
sudo yum install -y nvme-cli
sudo modprobe nvme-tcp
```

#### Step1: Setup hugepages

Before you prepare the caching nodes, please decide the amount of huge pages that you would like to allocate for simplyblock and set those hugepages accordingly. We suggest allocating at least 8GB of huge pages.

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
kubectl label nodes ip-10-0-4-118.us-east-2.compute.internal ip-10-0-4-176.us-east-2.compute.internal type=cache
```
Now the nodes are ready to deploy caching nodes.


### Driver deployment

During driver deployment, we will be deploying the caching nodes on all the nodes tagged with `type=cache`
```
kubectl apply -f caching-node.yaml
```

Once the caching nodes agents are deployed, we add the caching node simplyblock cluster

```
MGMT_IP=3.16.54.133
CLUSTER_ID=8d4aee39-3c4f-4c4b-8cef-08502327f2a3
CLUSTER_SECRET=GhbS1fwU8WSiABazTZUb

for node in $(kubectl get pods -l app=caching-node -owide | awk 'NR>1 {print $6}'); do
	echo "adding caching node: $node"

	curl --location "http://${MGMT_IP}/cachingnode/" \
		--header "Content-Type: application/json" \
		--header "Authorization: ${CLUSTER_ID} ${CLUSTER_SECRET}" \
		--data '{
		"cluster_id": "'"${CLUSTER_ID}"'",
		"node_ip": "'"${node}:5000"'",
		"iface_name": "eth0",
		"spdk_mem": "4g",
	}
	'
done
```

These steps are already added to `./deploy.sh` script.


### StorageClass

If the user wants to create a PVC that uses NVMe cache, a new storage class can be used with additional volume parameter as `type: cache`.


### Usage and Implementation

During dynamic volume provisioning, nodeSelector should be provided on pod, deployment, daemonset, statefulset. So that such pods are scheduled only on the nodes that has the `cache` label on it.

As shown below
```
    spec:
      nodeSelector:
        type: cache
```

On the controller server, when a new volume is requested, we create a `lvol` . This steps is exactly same as the current implementation.

On the node driver, during the volume mount, the following steps happens.
1. Get the caching node ID of the current node
2. Connect caching node with lvol. This will create a new NVMe device on the host machine. This device will be used to mount into pod.
