### testing


if the resources are not deleted properly using these commands to delete the resources.

```
kubectl -n spdk-csi delete pod spdk-fio-pod4
kubectl -n spdk-csi delete pvc spdk-fio-pvc4
kubectl spdk-csi delete sc spdk-fio-hostid
kubectl -n spdk-csi delete cm fio-config
```

