package utils

import (
	"net"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("IsEmptyCIDR", func() {
	When("the CIDR is empty", func() {
		It("should return true", func() {
			emptyCIDR := net.IPNet{}
			Expect(IsEmptyCIDR(emptyCIDR)).To(BeTrue())
		})
	})

	When("the CIDR is not empty", func() {
		It("should return false", func() {
			nonEmptyCIDR := net.IPNet{
				IP:   net.IPv4(192, 168, 1, 1),
				Mask: net.IPMask{255, 255, 255, 0},
			}
			Expect(IsEmptyCIDR(nonEmptyCIDR)).To(BeFalse())
		})
	})
})
