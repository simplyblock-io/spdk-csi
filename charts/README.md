# Installation with Helm 3

Follow this guide to install the SPDK-CSI Driver for Kubernetes.

## Prerequisites

### [Install Helm](https://helm.sh/docs/intro/quickstart/#install-helm)

### Build image

```console
make image
cd deploy/spdk
sudo docker build -t spdkdev .
```
 **_NOTE:_**
Kubernetes nodes must pre-allocate hugepages in order for the node to report its hugepage capacity.
A node can pre-allocate huge pages for multiple sizes.

## Install latest CSI Driver via `helm install`

```console

helm repo add spdk-csi https://raw.githubusercontent.com/simplyblock-io/spdk-csi/master/charts

helm repo update

hehelm install spdk-csi spdk-csi/spdk-csi \
  --set csiConfig.simplybk.uuid=ace14718-81eb-441f-9d4c-d71ce6904196 \
  --set csiConfig.simplybk.ip=https://96xdzb9ne7.execute-api.us-east-1.amazonaws.com \
  --set csiSecret.simplybk.secret=k6U5moyrY5vCVtSiCcKo \
  --set logicalVolume.pool_name=testing1
```

## After installation succeeds, you can get a status of Chart

```console
helm status "spdk-csi"
```

## Delete Chart

If you want to delete your Chart, use this command

```bash
helm uninstall "spdk-csi" --namespace "spdk-csi"
```

If you want to delete the namespace, use this command

```bash
kubectl delete namespace spdk-csi
```
