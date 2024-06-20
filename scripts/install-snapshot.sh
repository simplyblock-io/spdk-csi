#!/bin/bash -e

set -ex

# This script installs snapshot CRDs

SNAPSHOT_VERSION=${SNAPSHOT_VERSION:-"v7.0.1"}

SNAPSHOTTER_URL="https://raw.githubusercontent.com/kubernetes-csi/external-snapshotter/${SNAPSHOT_VERSION}"

# controller
SNAPSHOT_RBAC="${SNAPSHOTTER_URL}/deploy/kubernetes/snapshot-controller/rbac-snapshot-controller.yaml"
SNAPSHOT_CONTROLLER="${SNAPSHOTTER_URL}/deploy/kubernetes/snapshot-controller/setup-snapshot-controller.yaml"

# snapshot CRDs
SNAPSHOTCLASS="${SNAPSHOTTER_URL}/client/config/crd/snapshot.storage.k8s.io_volumesnapshotclasses.yaml"
VOLUME_SNAPSHOT_CONTENT="${SNAPSHOTTER_URL}/client/config/crd/snapshot.storage.k8s.io_volumesnapshotcontents.yaml"
VOLUME_SNAPSHOT="${SNAPSHOTTER_URL}/client/config/crd/snapshot.storage.k8s.io_volumesnapshots.yaml"


function install_snapshot_crds() {
    local namespace=$1
    if [ -z "${namespace}" ]; then
        namespace="default"
    fi

    kubectl apply -f "${SNAPSHOT_RBAC}" -n "${namespace}"
    kubectl apply -f "${SNAPSHOT_CONTROLLER}" -n "${namespace}"

    kubectl apply -f "${SNAPSHOTCLASS}" -n "${namespace}"
    kubectl apply -f "${VOLUME_SNAPSHOT_CONTENT}" -n "${namespace}"
    kubectl apply -f "${VOLUME_SNAPSHOT}" -n "${namespace}"
}

function delete_snapshot_crds() {
    local namespace=$1
    if [ -z "${namespace}" ]; then
        namespace="default"
    fi

    kubectl delete -f "${SNAPSHOTCLASS}" -n "${namespace}" --ignore-not-found
    kubectl delete -f "${VOLUME_SNAPSHOT_CONTENT}" -n "${namespace}" --ignore-not-found
    kubectl delete -f "${VOLUME_SNAPSHOT}" -n "${namespace}" --ignore-not-found
}

# parse the kubernetes version
# v1.17.2 -> kube_version 1 -> 1  (Major)
# v1.17.2 -> kube_version 2 -> 17 (Minor)
function kube_version() {
    echo "${KUBE_VERSION}" | sed 's/^v//' | cut -d'.' -f"${1}"
}

if ! get_kube_version=$(kubectl version) ||
   [[ -z "${get_kube_version}" ]]; then
    echo "could not get Kubernetes server version"
    echo "hint: check if you have specified the right host or port"
    exit 1
fi

KUBE_VERSION=$(echo "${get_kube_version}" | grep "^Server Version" | cut -d' ' -f3)
KUBE_MAJOR=$(kube_version 1)
KUBE_MINOR=$(kube_version 2)

# skip snapshot operation if kube version is less than 1.17.0
if [[ "${KUBE_MAJOR}" -lt 1 ]] || [[ "${KUBE_MAJOR}" -eq 1  &&  "${KUBE_MINOR}" -lt 17 ]]; then
    echo "skipping: Kubernetes server version is < 1.17.0"
    exit 1
fi

case "${1:-}" in
install)
    install_snapshot_crds "$2"
    ;;
delete)
    delete_snapshot_crds "$2"
    ;;
*)
    echo "usage:" >&2
    echo "  $0 install [namespace]" >&2
    echo "  $0 delete-crd [namespace]" >&2
    ;;
esac

