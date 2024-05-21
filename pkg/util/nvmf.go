/*
Copyright (c) Arm Limited and Contributors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package util

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"k8s.io/klog"
)

type NodeNVMf struct {
	client *rpcClient

	clusterID     string
	clusterIP     string
	clusterSecret string
}

// func newNVMf(client *rpcClient, targetType, targetAddr string) *nodeNVMf {
// config.Simplybk.Uuid, config.Simplybk.Ip, secret.Simplybk.Secret
func NewNVMf(clusterID, clusterIP, clusterSecret string) *NodeNVMf {
	client := rpcClient{
		ClusterID:     clusterID,
		ClusterIP:     clusterIP,
		ClusterSecret: clusterSecret,
		httpClient:    &http.Client{Timeout: cfgRPCTimeoutSeconds * time.Second},
	}
	return &NodeNVMf{
		client:        &client,
		clusterID:     clusterID,
		clusterIP:     clusterIP,
		clusterSecret: clusterSecret,
	}
}

func (node *NodeNVMf) Info() string {
	return node.client.info()
}

func (node *NodeNVMf) LvStores() ([]LvStore, error) {
	return node.client.lvStores()
}

// VolumeInfo returns a string:string map containing information necessary
// for CSI node(initiator) to connect to this target and identify the disk.
func (node *NodeNVMf) VolumeInfo(lvolID string) (map[string]string, error) {
	lvol, err := node.client.getVolumeInfo(lvolID)
	if err != nil {
		return nil, err
	}

	return lvol, nil
}

type CreateLVolData struct {
	LvolName    string `json:"name"`
	Size        string `json:"size"`
	LvsName     string `json:"pool"`
	Compression bool   `json:"comp"`
	Encryption  bool   `json:"crypto"`
	Snapshot    bool   `json:"snapshot"`
	MaxRWIOPS   string `json:"max_rw_iops"`
	MaxRWmBytes string `json:"max_rw_mbytes"`
	MaxRmBytes  string `json:"max_r_mbytes"`
	MaxWmBytes  string `json:"max_w_mbytes"`
	DistNdcs    int    `json:"distr_ndcs"`
	DistNpcs    int    `json:"distr_npcs"`
	CryptoKey1  string `json:"crypto_key1"`
	CryptoKey2  string `json:"crypto_key2"`
}

// CreateVolume creates a logical volume and returns volume ID
func (node *NodeNVMf) CreateVolume(params *CreateLVolData) (string, error) {
	lvolID, err := node.client.createVolume(params)
	if err != nil {
		return "", err
	}
	klog.V(5).Infof("volume created: %s", lvolID)
	return lvolID, nil
}

// GetVolume returns the volume id of the given volume name and lvstore name. return error if not found.
func (node *NodeNVMf) GetVolume(lvolName, poolName string) (string, error) {
	lvol, err := node.client.getVolume(fmt.Sprintf("%s/%s", poolName, lvolName))
	if err != nil {
		return "", err
	}
	return lvol.UUID, err
}

func (node *NodeNVMf) GetVolumeSize(lvolID string) (string, error) {
	lvol, err := node.client.getVolume(lvolID)
	if err != nil {
		return "", err
	}

	size := strconv.FormatInt(lvol.LvolSize, 10)
	return size, err
}

func (node *NodeNVMf) ListVolumes() ([]*BDev, error) {
	return node.client.listVolumes()
}

// func (node *NodeNVMf) isVolumeCreated(lvolID string) (bool, error) {
// 	return node.client.isVolumeCreated(lvolID)
// }

func (node *NodeNVMf) ResizeVolume(lvolID string, newSize int64) (bool, error) {
	return node.client.resizeVolume(lvolID, newSize)
}

func (node *NodeNVMf) ListSnapshots() ([]*SnapshotResp, error) {
	return node.client.listSnapshots()
}

func (node *NodeNVMf) CreateSnapshot(lvolID, snapshotName, poolName string) (string, error) {
	snapshotID, err := node.client.snapshot(lvolID, snapshotName, poolName)
	if err != nil {
		return "", err
	}
	klog.V(5).Infof("snapshot created: %s", snapshotID)
	return snapshotID, nil
}

func (node *NodeNVMf) DeleteVolume(lvolID string) error {
	err := node.client.deleteVolume(lvolID)
	if err != nil {
		return err
	}
	klog.V(5).Infof("volume deleted: %s", lvolID)
	return nil
}

func (node *NodeNVMf) DeleteSnapshot(snapshotID string) error {
	err := node.client.deleteSnapshot(snapshotID)
	if err != nil {
		return err
	}
	klog.V(5).Infof("snapshot deleted: %s", snapshotID)
	return nil
}

// PublishVolume exports a volume through NVMf target
func (node *NodeNVMf) PublishVolume(lvolID string) error {
	_, err := node.client.callSBCLI("GET", "/lvol/publish_volume/"+lvolID, nil)
	if err != nil {
		return err
	}
	// exists, err := node.isVolumeCreated(lvolID)
	// if err != nil {
	// 	return err
	// }
	// if !exists {
	// 	return ErrVolumeDeleted
	// }
	// published, err := node.isVolumePublished(lvolID)
	// if err != nil {
	// 	return err
	// }
	// if published {
	// 	return nil
	// }

	// err = node.createTransport()
	// if err != nil {
	// 	return err
	// }

	// err = node.createSubsystem(lvolID)
	// if err != nil {
	// 	return err
	// }

	// _, err = node.subsystemAddNs(lvolID)
	// if err != nil {
	// 	node.deleteSubsystem(lvolID) //nolint:errcheck // we can do few
	// 	return err
	// }

	// err = node.subsystemAddListener(lvolID)
	// if err != nil {
	// 	node.subsystemRemoveNs(lvolID) //nolint:errcheck // ditto
	// 	node.deleteSubsystem(lvolID)   //nolint:errcheck // ditto
	// 	return err
	// }

	klog.V(5).Infof("volume published: %s", lvolID)
	return nil
}

// func (node *NodeNVMf) isVolumePublished(lvolID string) (bool, error) {
// 	var isPublished bool
// 	out, err := node.client.callSBCLI("GET", "/lvol/is_volume_published/"+lvolID, nil)
// 	if err != nil {
// 		// querying nqn that does not exist, an invalid parameters error will be thrown
// 		if errorMatches(err, ErrInvalidParameters) {
// 			return false, nil
// 		}
// 		return false, err
// 	}
// 	isPublished, ok := out.(bool)
// 	if !ok {
// 		return false, fmt.Errorf("failed to convert the response to bool type. Interface: %v", out)
// 	}
// 	return isPublished, nil
// }

func (node *NodeNVMf) UnpublishVolume(lvolID string) error {
	_, err := node.client.callSBCLI("GET", "/lvol/unpublish_volume/"+lvolID, nil)
	if err != nil {
		return err
	}
	// exists, err := node.isVolumeCreated(lvolID)
	// if err != nil {
	// 	return err
	// }
	// if !exists {
	// 	return ErrVolumeDeleted
	// }
	// published, err := node.isVolumePublished(lvolID)
	// if err != nil {
	// 	return err
	// }
	// if !published {
	// 	// already unpublished
	// 	return nil
	// }
	// err = node.subsystemRemoveNs(lvolID)
	// if err != nil {
	// 	// we should try deleting subsystem even if we fail here
	// 	klog.Errorf("failed to remove namespace(nqn=%s): %s", node.getVolumeNqn(lvolID), err)
	// }
	// err = node.deleteSubsystem(lvolID)
	// if err != nil {
	// 	return err
	// }
	klog.V(5).Infof("volume unpublished: %s", lvolID)
	return nil
}
