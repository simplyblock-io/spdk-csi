package e2e

import (
	"time"

	ginkgo "github.com/onsi/ginkgo/v2"
	"k8s.io/kubernetes/test/e2e/framework"
)

const (
	iscsiConfigMapData = `{
	"simplybk": {
		"uuid": "79276661-5f8a-405d-ab6d-651b88326206",
		"ip": "44.204.254.210"
	 }
}`
)

var _ = ginkgo.Describe("SPDKCSI-ISCSI", func() {
	f := framework.NewDefaultFramework("spdkcsi")
	ginkgo.BeforeEach(func() {
		deployConfigs(iscsiConfigMapData)
		deployCsi()
	})

	ginkgo.AfterEach(func() {
		deleteCsi()
		deleteConfigs()
	})

	ginkgo.Context("Test SPDK CSI ISCSI", func() {
		ginkgo.It("Test SPDK CSI ISCSI", func() {
			ginkgo.By("checking controller statefulset is running", func() {
				err := waitForControllerReady(f.ClientSet, 4*time.Minute)
				if err != nil {
					ginkgo.Fail(err.Error())
				}
			})

			ginkgo.By("checking node daemonset is running", func() {
				err := waitForNodeServerReady(f.ClientSet, 2*time.Minute)
				if err != nil {
					ginkgo.Fail(err.Error())
				}
			})

			ginkgo.By("create a PVC and verify dynamic PV", func() {
				deployPVC()
				defer deletePVC()
				err := verifyDynamicPVCreation(f.ClientSet, "spdkcsi-pvc", 5*time.Minute)
				if err != nil {
					ginkgo.Fail(err.Error())
				}
			})

			ginkgo.By("create a PVC and bind it to a pod", func() {
				deployPVC()
				deployTestPod()
				defer deletePVCAndTestPod()
				err := waitForTestPodReady(f.ClientSet, 5*time.Minute)
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
	})
})
