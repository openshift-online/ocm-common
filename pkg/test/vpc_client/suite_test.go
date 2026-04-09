package vpc_client_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestVPCClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "VPC Client Suite")
}
