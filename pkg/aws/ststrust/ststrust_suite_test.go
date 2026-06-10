package ststrust_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestSTSTrust(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "STS trust policy external ID")
}
