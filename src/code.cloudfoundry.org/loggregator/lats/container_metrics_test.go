package lats_test

import (
	"code.cloudfoundry.org/loggregator/lats/helpers"

	"code.cloudfoundry.org/loggregator/plumbing/conversion"
	v2 "code.cloudfoundry.org/loggregator/plumbing/v2"

	uuid "github.com/nu7hatch/gouuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry/sonde-go/events"
)

var _ = Describe("Container Metrics Endpoint", func() {
	var (
		appID string
	)

	BeforeEach(func() {
		guid, err := uuid.NewV4()
		Expect(err).ToNot(HaveOccurred())

		appID = guid.String()
	})

	It("can receive container metrics", func() {
		envelope := createContainerMetric(appID)
		helpers.EmitToMetronV1(envelope)

		f := func() []*events.ContainerMetric {
			return helpers.RequestContainerMetrics(appID)
		}
		Eventually(f).Should(ContainElement(envelope.ContainerMetric))
	})

	Describe("emit v2 and consume via reverse log proxy", func() {
		It("can receive container metrics", func() {
			envelope := createContainerMetric(appID)
			v2Env := conversion.ToV2(envelope, false)
			helpers.EmitToMetronV2(v2Env)

			f := func() []*v2.Envelope {
				return helpers.ReadContainerFromRLP(appID, false)
			}
			Eventually(f).Should(ContainElement(v2Env))
		})

		It("can receive container metrics with preferred tags", func() {
			envelope := createContainerMetric(appID)
			v2Env := conversion.ToV2(envelope, true)
			helpers.EmitToMetronV2(v2Env)

			f := func() []*v2.Envelope {
				return helpers.ReadContainerFromRLP(appID, true)
			}
			Eventually(f).Should(ContainElement(v2Env))
		})
	})
})
