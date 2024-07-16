package e2e

import (
	"time"

	ginkgo "github.com/onsi/ginkgo/v2"
	"k8s.io/kubernetes/test/e2e/framework"
)

var _ = ginkgo.Describe("SPDKCSI-SNAPSHOT", func() {
	f := framework.NewDefaultFramework("spdkcsi")

	ginkgo.Context("Test SPDK CSI Snapshot", func() {
		ginkgo.It("Test SPDK CSI Snapshot", func() {
			ginkgo.By("create source pvc and write data", func() {
				deployPVC()
				deployTestPod()
				defer deleteTestPod()

				err := waitForTestPodReady(f.ClientSet, 3*time.Minute)
				if err != nil {
					ginkgo.Fail(err.Error())
				}
				writeDataToPod(f)
			})

			ginkgo.By("create snapshot and check data persistency", func() {
				deploySnapshot()
				defer deleteSnapshot()
				defer deletePVC()

				err := waitForTestPodReady(f.ClientSet, 3*time.Minute)
				if err != nil {
					ginkgo.Fail(err.Error())
				}
				err = compareDataInPod(f)
				if err != nil {
					ginkgo.Fail(err.Error())
				}
			})
		})
	})
})
