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

helm install -n spdk-csi --create-namespace spdk-csi spdk-csi/spdk-csi \
  --set csiConfig.simplybk.uuid=ace14718-81eb-441f-9d4c-d71ce6904196 \
  --set csiConfig.simplybk.ip=https://96xdzb9ne7.execute-api.us-east-1.amazonaws.com \
  --set csiSecret.simplybk.secret=k6U5moyrY5vCVtSiCcKo \
  --set logicalVolume.pool_name=testing1
```

## After installation succeeds, you can get a status of Chart

```console
helm status "spdk-csi" --namespace "spdk-csi"
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

## latest chart configuration

The following table lists the configurable parameters of the latest Simplyblock CSI Driver chart and default values.

| Parameter                              | Description                                                                                                              | Default                                                                 |
| -------------------------------------- | ------------------------------------------------------------------------------------------------------------------------ | ----------------------------------------------------------------------- |
| `driverName`                           | alternative driver name                                                                                                  | `csi.spdk.io`                                                           |
| `image.spdkcsi.repository`             | simplyblock-csi-driver image                                                                                             | `simplyblock/spdkcsi`                                                   |
| `image.spdkcsi.tag`                    | simplyblock-csi-driver image tag                                                                                         | `latest`                                                                |
| `image.spdkcsi.pullPolicy`             | simplyblock-csi-driver image pull policy                                                                                 | `Always`                                                                |
| `image.csiProvisioner.repository`      | csi-provisioner docker image                                                                                             | `registry.k8s.io/sig-storage/csi-provisioner`                           |
| `image.csiProvisioner.tag`             | csi-provisioner docker image tag                                                                                         | `v4.0.1`                                                                |
| `image.csiProvisioner.pullPolicy`      | csi-provisioner image pull policy                                                                                        | `Always`                                                                |
| `image.csiAttacher.repository`         | csiAttacher docker image                                                                                                 | `gcr.io/k8s-staging-sig-storage/csi-attacher`                           |
| `image.csiAttacher.tag`                | csiAttacher docker image tag                                                                                             | `v4.5.1`                                                                |
| `image.csiAttacher.pullPolicy`         | csiAttacher image pull policy                                                                                            | `Always`                                                                |
| `image.nodeDriverRegistrar.repository` | csi-node-driver-registrar docker image                                                                                   | `registry.k8s.io/sig-storage/csi-node-driver-registrar`                 |
| `image.nodeDriverRegistrar.tag`        | csi-node-driver-registrar docker image tag                                                                               | `v2.10.1`                                                               |
| `image.nodeDriverRegistrar.pullPolicy` | csi-node-driver-registrar image pull policy                                                                              | `Always`                                                                |
| `image.csiSnapshotter.repository`      | csi-snapshotter docker image                                                                                             | `registry.k8s.io/sig-storage/csi-snapshotter`                           |
| `image.csiSnapshotter.tag`             | csi-snapshotter docker image tag                                                                                         | `v7.0.2`                                                                |
| `image.csiSnapshotter.pullPolicy`      | csi-snapshotter image pull policy                                                                                        | `Always`                                                                |
| `image.csiResizer.repository`          | csi-resizer  docker image                                                                                                | `gcr.io/k8s-staging-sig-storage/csi-resizer`                            |
| `image.csiResizer.tag`                 | csi-resizer docker image tag                                                                                             | `v1.10.1`                                                               |
| `image.csiResizer.pullPolicy`          | csi-resizer image pull policy                                                                                            | `Always`                                                                |
| `image.csiHealthMonitor.repository`    | csi-external-health-monitor-controller docker image                                                                      | `gcr.io/k8s-staging-sig-storage/csi-external-health-monitor-controller` |
| `image.csiHealthMonitor.tag`           | csi-external-health-monitor-controller docker image tag                                                                  | `v0.11.0`                                                               |
| `image.csiHealthMonitor.pullPolicy`    | csi-external-health-monitor-controller image pull policy                                                                 | `Always`                                                                |
| `image.simplyblock.repository`         | simplyblock spdk docker image                                                                                            | `simplyblock/simplyblock`                                               |
| `image.simplyblock.tag`                | simplyblock spdk docker image tag                                                                                        | `release_v1`                                                            |
| `image.simplyblock.pullPolicy`         | csi-snapshotter image pull policy                                                                                        | `Always`                                                                |
| `serviceAccount.create`                | whether to create service account of spdkcsi-controller                                                                  | `true`                                                                  |
| `rbac.create`                          | whether create rbac of spdkcsi-controller                                                                                | `true`                                                                  |
| `controller.replicas`                  | replica number of spdkcsi-controller                                                                                     | `1`                                                                     |
| `storageClass.create`                  | create storageclass                                                                                                      | `true`                                                                  |  |
| `externallyManagedConfigmap.create`    | Specifies whether a externallyManagedConfigmap should be created                                                         | `true`                                                                  |  |
| `externallyManagedSecret.create`       | Specifies whether a externallyManagedSecret should be created                                                            | `true`                                                                  |  |
| `csiConfig.simplybk.uuid`              | the UUID of the cluster on which the volumes are created                                                                 | ``                                                                      |  |
| `csiConfig.simplybk.ip`                | the HTTPS API Gateway endpoint connected to the management node                                                          | `https://o5ls1ykzbb.execute-api.eu-central-1.amazonaws.com`             |  |
| `csiSecret.simplybk.secret`            | the cluster secret associated with the cluster                                                                           | ``                                                                      |  |
| `csiSecret.simplybkPvc.crypto_key1`    | if an encrypted PVC is to be created, value of `crypto_key1`                                                             | ``                                                                      |  |
| `csiSecret.simplybkPvc.crypto_key2`    | if an encrypted PVC is to be created, value of `crypto_key2`                                                             | ``                                                                      |  |
| `logicalVolume.pool_name`              | the name of the pool against which the lvols needs to be provisioned. This Pool needs to be created before its passed here. | `testing1`                                                              |  |
| `logicalVolume.qos_rw_iops`            | the value of lvol parameter qos_rw_iops                                                                                  | `0`                                                                     |  |
| `logicalVolume.qos_rw_mbytes`          | the value of lvol parameter qos_rw_mbytes                                                                                | `0`                                                                     |  |
| `logicalVolume.qos_r_mbytes`           | the value of lvol parameter qos_r_mbytes                                                                                 | `0`                                                                     |  |
| `logicalVolume.qos_w_mbytes`           | the value of lvol parameter qos_w_mbytes                                                                                 | `0`                                                                     |  |
| `logicalVolume.compression`            | set to `True` if compression needs be enabled on lvols                                                                   | `False`                                                                 |  |
| `logicalVolume.encryption`             | set to `True` if encryption needs be enabled on lvols.                                                                   | `False`                                                                 |  |
| `logicalVolume.distr_ndcs`             | the value of distr_ndcs                                                                                                  | `1`                                                                     |  |
| `logicalVolume.distr_npcs`             | the value of distr_npcs                                                                                                  | `1`                                                                     |  |

## troubleshooting
 - Add `--wait -v=5 --debug` in `helm install` command to get detailed error
 - Use `kubectl describe` to acquire more info
