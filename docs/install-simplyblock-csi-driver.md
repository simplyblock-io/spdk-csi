### Prepare SPDK storage node

Follow [deploy/spdk/README](deploy/spdk/README.md) to deploy SPDK storage service on localhost.

### Deploy SPDKCSI services

1. Launch Minikube test cluster.
  ```bash
    $ cd scripts
    $ sudo ./minikube.sh up

    # Create kubectl shortcut (assume kubectl version 1.25.0)
    $ sudo ln -s /var/lib/minikube/binaries/v1.25.0/kubectl /usr/local/bin/kubectl

    # Wait for Kubernetes ready
    $ kubectl get pods --all-namespaces
    NAMESPACE     NAME                          READY   STATUS    RESTARTS   AGE
    kube-system   coredns-6955765f44-dlb88      1/1     Running   0          81s
    ......                                              ......
    kube-system   kube-apiserver-spdkcsi-dev    1/1     Running   0          67s
    ......                                              ......
  ```

2. Install snapshot controller and CRD

  ```bash
    SNAPSHOT_VERSION="v3.0.3" ./scripts/install-snapshot.sh install

    # Check status
    $ kubectl get pod snapshot-controller-0
    NAME                    READY   STATUS    RESTARTS   AGE
    snapshot-controller-0   1/1     Running   0          6m14s
  ```

3. Deploy SPDK-CSI services
  ```bash
    $ cd deploy/kubernetes
    $ ./deploy.sh

    # Check status
    $ kubectl get pods
    NAME                   READY   STATUS    RESTARTS   AGE
    spdkcsi-controller-0   3/3     Running   0          3m16s
    spdkcsi-node-lzvg5     2/2     Running   0          3m16s
  ```

4. Deploy test pod
  ```bash
    $ cd deploy/kubernetes
    $ kubectl apply -f testpod.yaml

    # Check status
    $ kubectl get pv
    NAME                       CAPACITY   ...    STORAGECLASS   REASON   AGE
    persistentvolume/pvc-...   256Mi      ...    spdkcsi-sc              43s

    $ kubectl get pvc
    NAME                                ...   CAPACITY   ACCESS MODES   STORAGECLASS   AGE
    persistentvolumeclaim/spdkcsi-pvc   ...   256Mi      RWO            spdkcsi-sc     44s

    $ kubectl get pods
    NAME                   READY   STATUS    RESTARTS   AGE
    spdkcsi-test           1/1     Running   0          1m31s

    # Check attached spdk volume in test pod
    $ kubectl exec spdkcsi-test mount | grep spdkcsi
    /dev/disk/by-id/nvme-..._spdkcsi-sn on /spdkvol type ext4 (rw,relatime)
  ```

5. Deploy PVC snapshot
  ```bash
    # Create snapshot of the bound PVC
    $ cd deploy/kubernetes
    $ kubectl apply -f testsnapshot.yaml

    # Get details about the snapshot
    $ kubectl get volumesnapshot spdk-snapshot
    NAME            READYTOUSE   SOURCEPVC   ... SNAPSHOTCLASS         AGE
    spdk-snapshot   false        spdkcsi-pvc ... csi-spdk-snapclass    29s

    # Get details about the volumesnapshotcontent
    kubectl get volumesnapshotcontent
    $ kubectl get volumesnapshotcontent
    NAME        ...   READYTOUSE   RESTORESIZE   DELETIONPOLICY   DRIVER        VOLUMESNAPSHOTCLASS   VOLUMESNAPSHOT   AGE
    snapcontent-...   true         268435456     Delete           csi.simplyblock.io   csi-spdk-snapclass    spdk-snapshot    29s
  ```

### Teardown

1. Delete PVC snapshot
  ```bash
    cd deploy/kubernetes
    kubectl delete -f testsnapshot.yaml
  ```

2. Delete test pod
  ```bash
    $ cd deploy/kubernetes
    $ kubectl delete -f testpod.yaml
  ```

3. Delete SPDK-CSI services
  ```bash
    $ cd deploy/kubernetes
    $ ./deploy.sh teardown
  ```

4. Delete snapshot controller and CRD
  ```bash
  SNAPSHOT_VERSION="v3.0.3" ./scripts/install-snapshot.sh cleanup
  ```

5. Teardown Kubernetes test cluster
  ```bash
    $ cd scripts
    $ sudo ./minikube.sh clean
  ```
