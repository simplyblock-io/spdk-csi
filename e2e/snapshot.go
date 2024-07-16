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
				// do not delete pvc here, since we need it for snapshot

				err := waitForTestPodReady(f.ClientSet, 3*time.Minute)
				if err != nil {
					ginkgo.Fail(err.Error())
				}
				// write data to source pvc
				writeDataToPod(f)
			})

			ginkgo.By("create snapshot and check data persistency", func() {
				// deploy volume snapshot and new pvc which data source is the snapshot
				deploySnapshot()
				defer deleteSnapshot()
				defer deletePVC()

				err := waitForTestPodReady(f.ClientSet, 3*time.Minute)
				if err != nil {
					ginkgo.Fail(err.Error())
				}
				// check if data in snapshot pvc is the same as source pvc
				err = compareDataInPod(f)
				if err != nil {
					ginkgo.Fail(err.Error())
				}
			})
		})
	})
})
