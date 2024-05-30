# Supported Port for eks or ks3

|  Port          |   Protocol    | Description 
| -------------- | ------------- | -------------
| 6443           | TCP           | Kubernetes API server. Required for communication between the Kubernetes
                                   control plane and the nodes in the cluster.
| 22             | TCP           | SSH access to the instances. Necessary for administrative access and management.
| 8080           | TCP           | SPDK Proxy for the storage node. Facilitates communication between the storage nodes and the
                                   management node.
| 2375           | TCP           | Docker Engine API. Allows the management node to communicate with Docker engines running on other 
                                   nodes.
|  8             | ICMP          | Allows ICMP Echo requests. Used for ping operations to check the availability and responsiveness of 
                                   management nodes.
|  5000          | TCP           | Caching node. Enables communication with caching services running on the node.
