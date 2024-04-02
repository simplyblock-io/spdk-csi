#!/bin/bash

set -ex

CLUSTER_ID='67f704a0-4359-422a-9417-ed79d1cf1edc'
MGMT_IP='34.200.226.58'
CLUSTER_SECRET='nRPkQEAm07Uw6SlpSNIb'

# list in creation order
files=(driver config-map nodeserver-config-map secret controller-rbac node-rbac controller node storageclass caching-node)

if [ "$1" = "teardown" ]; then
	# delete in reverse order
	for ((i = ${#files[@]} - 1; i >= 0; i--)); do
		echo "=== kubectl delete -f ${files[i]}.yaml"
		kubectl delete -f "${files[i]}.yaml"
	done
	exit 0
else
	for ((i = 0; i <= ${#files[@]} - 1; i++)); do
		echo "=== kubectl apply -f ${files[i]}.yaml"
		kubectl apply -f "${files[i]}.yaml"
	done
fi

echo ""
echo "Deploying Caching node..."

output=$(kubectl get nodes -l type=cache | wc -l)

if [ $output -lt 2 ]; then
    echo "No caching nodes found. Exiting..."
    exit 0
fi


## check if the caching nodes has required huge pages
echo "-- caching nodes --"
kubectl get nodes -l type=cache


kubectl apply -f caching-node.yaml
kubectl wait --timeout=3m --for=condition=ready pod -l app=caching-node

for node in $(kubectl get pods -l app=caching-node -owide | awk 'NR>1 {print $(NF-3)}'); do
	echo "adding caching node: $node"

	curl --location "http://${MGMT_IP}/cachingnode/" \
		--header "Content-Type: application/json" \
		--header "Authorization: ${CLUSTER_ID} ${CLUSTER_SECRET}" \
		--data '{
		"cluster_id": "'"${CLUSTER_ID}"'",
		"node_ip": "'"${node}:5000"'",
		"iface_name": "eth0",
		"spdk_mem": "2g"
	}
	'
done
