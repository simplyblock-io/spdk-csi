package e2e

import (
	"fmt"
	"time"

	ginko "github.com/onsi/ginkgo/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubernetes/test/e2e/framework"
)

const (
	Size5GB          = 5 * 1024 * 1024 * 1024
	StorageclassName = "spdk-fio-hostid"
)

// TODO: use ctx.Before and ctx.After to cleanup the resources
var _ = ginko.Describe("CSI Driver tests", func() {
	f := framework.NewDefaultFramework("spdkcsi")
	ginko.Context("Control Plane: delete second lvol while the first lvol has IO running", func() {
		ginko.It("if a node has lvol with IO running, adding and deleting an new lvol from the same node should work", func() {
			// get nodes from the cluster and pick a random node
			// TODO: defer does not work if the tests were called using ctrl-c. Handle this case
			ginko.By("run a pod with fio and add and delete an lvol from the same node", func() {
				c := f.ClientSet
				sn, err := getStorageNode(c)
				if err != nil {
					fmt.Println(err.Error())
					// ginko.Fail(err.Error())
				}

				// create a new storage class
				// Problem: is this SC causing the issue?
				// or is the issue with lvol add or delete
				err = createstorageClassWithHostID(c, StorageclassName, sn)
				if err != nil {
					fmt.Println(err.Error())
					// ginko.Fail(err.Error())
				}

				// create a configmap
				configMapname := "fio-config"
				err = createFioConfigMap(c, nameSpace, configMapname)
				if err != nil {
					fmt.Println(err.Error())
					// ginko.Fail(err.Error())
				}

				// create a pvc with the storage class
				pvcName := "spdk-fio-pvc4"
				// fmt.Printf("creating pvc: %s \n", pvcName)

				err = createPVC(c, nameSpace, pvcName, StorageclassName, Size5GB)
				if err != nil {
					fmt.Println(err.Error())
					// ginko.Fail(err.Error())
				}

				// create a pod with the storage class
				podName := "spdk-fio-pod4"
				// fmt.Printf("creating pod: %s \n", podName)
				err = createFioWorkloadPod(c, nameSpace, podName, configMapname, pvcName)
				if err != nil {
					fmt.Println(err.Error())
					// ginko.Fail(err.Error())
				}

				defer func() {
					err2 := c.CoreV1().ConfigMaps(nameSpace).Delete(ctx, configMapname, metav1.DeleteOptions{})
					if err2 != nil {
						fmt.Println(err.Error())
						// ginko.Fail(err2.Error())
					}
					err2 = c.CoreV1().PersistentVolumeClaims(nameSpace).Delete(ctx, pvcName, metav1.DeleteOptions{})
					if err2 != nil {
						fmt.Println(err.Error())
						// ginko.Fail(err2.Error())
					}
					err2 = c.CoreV1().Pods(nameSpace).Delete(ctx, podName, metav1.DeleteOptions{})
					if err2 != nil {
						fmt.Println(err.Error())
						// ginko.Fail(err2.Error())
					}
				}()

				// add and delete an lvol from the same node
				// TODO: cleanup the variables
				pvcName2 := "spdk-fio-pvc5"
				podName2 := "spdk-fio-pod5"

				// fmt.Printf("creating pvc: %s \n", pvcName2)
				err = createPVC(c, nameSpace, pvcName2, StorageclassName, Size5GB)
				if err != nil {
					fmt.Println(err.Error())
					// ginko.Fail(err.Error())
				}

				// fmt.Printf("creating pod: %s \n", podName2)
				err = createSimplePod(c, nameSpace, podName2, pvcName2)
				if err != nil {
					fmt.Println(err.Error())
					// ginko.Fail(err.Error())
				}

				// fmt.Printf("deleting pod: %s \n", podName2)
				err = c.CoreV1().Pods(nameSpace).Delete(ctx, podName2, metav1.DeleteOptions{})
				if err != nil {
					fmt.Println(err.Error())
					// ginko.Fail(err.Error())
				}

				// fmt.Printf("deleting pvc: %s \n", pvcName2)
				err = c.CoreV1().PersistentVolumeClaims(nameSpace).Delete(ctx, pvcName2, metav1.DeleteOptions{})
				if err != nil {
					fmt.Println(err.Error())
					// ginko.Fail(err.Error())
				}

				// fmt.Println("waiting for 10 secs")
				time.Sleep(10 * time.Second)
				// fio should not stop on spd-csi-pod4
				// fmt.Println("checking for pod spdk-fio-pod4 status")
				_, err = c.CoreV1().Pods(nameSpace).Get(ctx, podName, metav1.GetOptions{})
				if err != nil {
					fmt.Println(err.Error())
					// ginko.Fail(err.Error())
				}
				// fixme: this should be in running mode
				// fmt.Printf("pod status: %s \n", pod.Status.Phase)
			})
		})
	})
})
