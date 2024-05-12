package e2e

import (
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubernetes/test/e2e/framework"
)

const (
	SIZE_5GB           = 5 * 1024 * 1024 * 1024
	STORAGE_CLASS_NAME = "spdk-fio-hostid"
)

// TODO: use ctx.Before and ctx.After to cleanup the resources
var _ = Describe("CSI Driver tests", func() {
	f := framework.NewDefaultFramework("spdkcsi")
	Context("Control Plane: delete second lvol while the first lvol has IO running", func() {
		It("if a node has lvol with IO running, adding and deleting an lvol from the same node should work", func() {
			// get nodes from the cluster and pick a random node
			// TODO: defer does not work if the tests were called using ctrl-c. Handle this case
			By("run a pod with fio and add and delete an lvol from the same node", func() {
				c := f.ClientSet
				sn, err := getStorageNode(c)
				if err != nil {
					Fail(err.Error())
				}

				// create a new storage class
				err = createstorageClassWithHostID(c, STORAGE_CLASS_NAME, sn)
				if err != nil {
					Fail(err.Error())
				}
				defer c.StorageV1().StorageClasses().Delete(ctx, STORAGE_CLASS_NAME, metav1.DeleteOptions{})

				// create a configmap
				configMapname := "fio-config"
				err = createFioConfigMap(c, nameSpace, configMapname)
				if err != nil {
					Fail(err.Error())
				}
				defer c.CoreV1().ConfigMaps(nameSpace).Delete(ctx, configMapname, metav1.DeleteOptions{})

				// create a pvc with the storage class
				pvcName := "spdk-fio-pvc4"
				fmt.Println("creating pvc: ", pvcName)

				err = createPVC(c, nameSpace, pvcName, STORAGE_CLASS_NAME, SIZE_5GB)
				if err != nil {
					Fail(err.Error())
				}
				defer c.CoreV1().PersistentVolumeClaims(nameSpace).Delete(ctx, pvcName, metav1.DeleteOptions{})

				// create a pod with the storage class
				podName := "spdk-fio-pod4"
				fmt.Println("creating pod: ", podName)
				err = createFioWorkloadPod(c, nameSpace, podName, configMapname, pvcName)
				if err != nil {
					Fail(err.Error())
				}
				defer c.CoreV1().Pods(nameSpace).Delete(ctx, podName, metav1.DeleteOptions{})

				// add and delete an lvol from the same node
				// TODO: cleanup the variables
				pvcName2 := "spdk-fio-pvc5"
				podName2 := "spdk-fio-pod5"

				fmt.Println("creating pvc: ", pvcName2)
				err = createPVC(c, nameSpace, pvcName2, STORAGE_CLASS_NAME, SIZE_5GB)
				if err != nil {
					Fail(err.Error())
				}

				fmt.Println("creating pod: ", podName2)
				err = createSimplePod(c, nameSpace, podName2, pvcName2)
				if err != nil {
					Fail(err.Error())
				}

				fmt.Println("deleting pod: ", podName2)
				err = c.CoreV1().Pods(nameSpace).Delete(ctx, podName2, metav1.DeleteOptions{})
				if err != nil {
					Fail(err.Error())
				}

				fmt.Println("deleting pvc: ", pvcName2)
				err = c.CoreV1().PersistentVolumeClaims(nameSpace).Delete(ctx, pvcName2, metav1.DeleteOptions{})
				if err != nil {
					Fail(err.Error())
				}

				fmt.Println("waiting for 10 secs")
				time.Sleep(10 * time.Second)
				// fio should not stop on spd-csi-pod4
				fmt.Println("checking for pod spdk-fio-pod4 status")
				pod, err := c.CoreV1().Pods(nameSpace).Get(ctx, podName, metav1.GetOptions{})
				if err != nil {
					Fail(err.Error())
				}
				// fixme: this should be in running mode
				fmt.Println("pod status: ", pod.Status.Phase)
			})
		})
	})
})
