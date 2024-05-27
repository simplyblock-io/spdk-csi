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

import "testing"

func TestNewNVMf(t *testing.T) {
	clusterID := "clusterID"
	clusterIP := "clusterIP"
	clusterSecret := "clusterSecret"

	node := NewNVMf(clusterID, clusterIP, clusterSecret)

	if node == nil {
		t.Fatal("NewNVMf returned nil")
	}

	if node.clusterID != clusterID {
		t.Errorf("Expected clusterID %s, but got %s", clusterID, node.clusterID)
	}

	if node.clusterIP != clusterIP {
		t.Errorf("Expected clusterIP %s, but got %s", clusterIP, node.clusterIP)
	}

	if node.clusterSecret != clusterSecret {
		t.Errorf("Expected clusterSecret %s, but got %s", clusterSecret, node.clusterSecret)
	}
}
