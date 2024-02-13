### Caching nodes volume provisioning

We now also have the concept of caching-nodes. These will be Kubernetes nodes which reside within the kubernetes compute cluster as normal worker nodes and can have any PODs deployed. However, they have an NVMe disk locally attached. If there are multiple NVMe disks attached, it will use the first one.


### Preparing nodes
Caching nodes are a special kind of node that provides storage by mounting a local NVMe disk. To prepare the nodes run the below steps. Before you prepare the caching nodes, please decide the amount of huge pages that you would like to allocate for simplyblock and set those values in `/etc/sysctl.conf`. We suggest allocating atleast 8GB of huge pages.

```
echo "======= setting huge pages ======="
echo "vm.nr_hugepages=2048" >>/etc/sysctl.conf # 
sysctl -p

# confirm it by running
cat /proc/meminfo | grep -i hug

# reboot
sudo reboot
```

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
for node in $(kubectl get pods -l app=caching-node -owide | awk 'NR>1 {print $6}'); do
	echo "adding caching node: $node"

	curl --location "http://${MGMT_IP}/cachingnode/" \
		--header "Content-Type: application/json" \
		--header "Authorization: ${CLUSTER_ID} ${CLUSTER_SECRET}" \
		--data '{
		"cluster_id": "'"${CLUSTER_ID}"'",
		"node_ip": "'"${node}:5000"'",
		"iface_name": "eth0",
		"spdk_mem": 8589934592,
	}
	'
done
```

These steps are already added to `./deploy.sh` script.


### StorageClass

If the user wants to create a PVC that uses NVMe cache, a new storage class can be used with addtional volume parameter as `type: cache`.


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
