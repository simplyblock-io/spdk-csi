### case#1: volume mount/unmount failed

#### simplyblock CSI driver requires nvme-cli to be installed.

If the kubernetes worker doesn't have nvme-cli installed, you might get error like this during the volume mount stage.
```
  Warning  FailedMount             17s (x7 over 49s)  kubelet                  MountVolume.MountDevice failed for volume "pvc-ac6f4696-18a8-4e98-a0d4-28237eef5ad9" : rpc error: code = Internal desc = exit status 254
```

This because the CSI driver internally nvme-cli to connect to the remote nvme volume. The CLI can be installed by running
```
sudo yum install -y nvme-cli;
sudo modprobe nvme-tcp
```

#### the worker nodes should be able to contact simplyblock storage node

Since SPDK uses remove nvme connection, make sure that there is network connectivity between the simplyblock storage node and the kubernetes worker nodes.

If not, we get errors like this during VolumeStage step.

```
  Warning  FailedMount             2s (x8 over 88s)  kubelet                  MountVolume.MountDevice failed for volume "pvc-d191cb36-c542-4d1a-806d-28c1999a50c6" : rpc error: code = Internal desc = exit status 146
```

#### check the events of pvc or pod

For any other types of failures in general, you could check the events of the Pod or PVC.

```console
kubectl describe pod test-pod
kubectl describe pvc test-pvc
```

For a more detailed analysis during the `CreateVolume` stage, you could refer to the controller.

```
kubectl -n spdk-csi logs -f spdkcsi-controller-0 -c spdkcsi-controller
```

Or you could refer to the node driver logs to get more info during the volume mount stage

```
kubectl -n spdk-csi logs -f spdkcsi-node-8vh5n -c spdkcsi-node
```
