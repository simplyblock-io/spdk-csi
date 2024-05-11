package e2e

import (
	"time"

	ginkgo "github.com/onsi/ginkgo/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubernetes/test/e2e/framework"
)

var _ = ginkgo.Describe("SPDKCSI-NVMEOF", func() {
	f := framework.NewDefaultFramework("spdkcsi")

	ginkgo.Context("Test SPDK CSI Dynamic Volume Provisioning", func() {
		ginkgo.It("CSI driver components should function properly", func() {
			ginkgo.By("checking controller statefulset is running", func() {
				err := waitForControllerReady(f.ClientSet, 4*time.Minute)
				if err != nil {
					ginkgo.Fail(err.Error())
				}
			})

			ginkgo.By("checking node daemonset is running", func() {
				err := waitForNodeServerReady(f.ClientSet, 3*time.Minute)
				if err != nil {
					ginkgo.Fail(err.Error())
				}
			})
		})

		ginkgo.It("Test the flow for Dynamic volume provisioning", func() {
			ginkgo.By("creating a PVC and verify dynamic PV", func() {
				deployPVC()
				defer deletePVC()
				err := verifyDynamicPVCreation(f.ClientSet, "spdkcsi-pvc", 5*time.Minute)
				if err != nil {
					ginkgo.Fail(err.Error())
				}
			})

			ginkgo.By("creating a PVC and binding it to a pod", func() {
				deployPVC()
				deployTestPod()
				defer deletePVCAndTestPod()
				err := waitForTestPodReady(f.ClientSet, 3*time.Minute)
				if err != nil {
					ginkgo.Fail(err.Error())
				}
			})
		})

		ginkgo.It("Test the flow for Caching nodes", func() {
			ginkgo.By("creating a caching PVC and bind it to a pod", func() {
				deployCachePVC()
				deployCacheTestPod()
				defer deleteCachePVCAndCacheTestPod()
				err := waitForCacheTestPodReady(f.ClientSet, 3*time.Minute)
				if err != nil {
					ginkgo.Fail(err.Error())
				}
			})

			ginkgo.By("check data persistency after the pod is removed and recreated", func() {
				deployPVC()
				deployTestPod()
				defer deletePVCAndTestPod()

				err := waitForTestPodReady(f.ClientSet, 3*time.Minute)
				if err != nil {
					ginkgo.Fail(err.Error())
				}

				err = checkDataPersist(f)
				if err != nil {
					ginkgo.Fail(err.Error())
				}
			})
		})

		ginkgo.It("Test multiple PVCs", func() {
			ginkgo.By("create multiple pvcs and a pod with multiple pvcs attached, and check data persistence after the pod is removed and recreated", func() {
				deployMultiPvcs()
				deployTestPodWithMultiPvcs()
				defer func() {
					deleteMultiPvcsAndTestPodWithMultiPvcs()
					if err := waitForTestPodGone(f.ClientSet); err != nil {
						ginkgo.Fail(err.Error())
					}
					for _, pvcName := range []string{"spdkcsi-pvc1", "spdkcsi-pvc2", "spdkcsi-pvc3"} {
						if err := waitForPvcGone(f.ClientSet, pvcName); err != nil {
							ginkgo.Fail(err.Error())
						}
					}
				}()
				err := waitForTestPodReady(f.ClientSet, 3*time.Minute)
				if err != nil {
					ginkgo.Fail(err.Error())
				}

				err = checkDataPersistForMultiPvcs(f)
				if err != nil {
					ginkgo.Fail(err.Error())
				}
			})
		})

		ginkgo.It("if a node has lvol with IO running, adding and deleting an lvol from the same node should work", func() {
			// get nodes from the cluster and pick a random node
			ginkgo.By("run a pod with fio", func() {
				c := f.ClientSet
				sn, err := getStorageNode(c)
				if err != nil {
					ginkgo.Fail(err.Error())
				}
				// create a new storage class
				storageClassName := "spdk-csi-hostid"
				err = createstorageClassWithHostID(c, storageClassName, sn)
				if err != nil {
					ginkgo.Fail(err.Error())
				}
				defer c.StorageV1().StorageClasses().Delete(ctx, storageClassName, metav1.DeleteOptions{})

				// create a configmap
				configMapname := "fio-config"
				err = createFioConfigMap(c, nameSpace, configMapname)
				if err != nil {
					ginkgo.Fail(err.Error())
				}
				defer c.CoreV1().ConfigMaps(nameSpace).Delete(ctx, configMapname, metav1.DeleteOptions{})

				// create a pvc with the storage class
				pvcName := "spdk-csi-pvc4"
				err = createPVC(c, nameSpace, pvcName, storageClassName, 256*1024*1024)
				if err != nil {
					ginkgo.Fail(err.Error())
				}
				defer c.CoreV1().PersistentVolumeClaims(nameSpace).Delete(ctx, pvcName, metav1.DeleteOptions{})

				// create a pod with the storage class
				podName := "spdk-csi-pod4"
				err = createFioWorkloadPod(c, nameSpace, podName, configMapname, pvcName)
				if err != nil {
					ginkgo.Fail(err.Error())
				}
				defer c.CoreV1().Pods(nameSpace).Delete(ctx, podName, metav1.DeleteOptions{})
			})

			ginkgo.By("add and delete an lvol from the same node", func() {
				c := f.ClientSet
				pvcName2 := "spdk-csi-pvc5"
				podName2 := "spdk-csi-pod5"
				storageClassName := "spdk-csi-hostid"

				err := createPVC(c, nameSpace, pvcName2, storageClassName, 256*1024*1024)
				if err != nil {
					ginkgo.Fail(err.Error())
				}

				err = createSimplePod(c, nameSpace, podName2, pvcName2)
				if err != nil {
					ginkgo.Fail(err.Error())
				}

				err = c.CoreV1().Pods(nameSpace).Delete(ctx, podName2, metav1.DeleteOptions{})
				if err != nil {
					ginkgo.Fail(err.Error())
				}

				err = c.CoreV1().PersistentVolumeClaims(nameSpace).Delete(ctx, pvcName2, metav1.DeleteOptions{})
				if err != nil {
					ginkgo.Fail(err.Error())
				}

				// fio should not stop on spd-csi-pod4
				err = waitForPodRunning(ctx, c, nameSpace, "spdk-csi-pod4", 5*time.Minute)
				ginkgo.Fail(err.Error())
			})
		})
	})
})
