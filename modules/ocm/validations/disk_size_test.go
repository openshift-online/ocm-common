package validations

import (
	. "github.com/onsi/ginkgo/v2/dsl/table"
	. "github.com/onsi/gomega"
)

var _ = DescribeTable("getAWSVolumeMaxSize",
	func(testSpec struct {
		name string
		args string
		want int
		err  error
	}) {
		got, err := getAWSVolumeMaxSize(testSpec.args)
		Expect(err).ToNot(HaveOccurred())
		Expect(got).To(Equal(testSpec.want), "got %v, want %v", got, testSpec.want)
	},
	Entry("valid version for 4.11", struct {
		name string
		args string
		want int
		err  error
	}{"valid version for 4.11", "4.11", 1024, nil}),
	Entry("valid version for 4.13", struct {
		name string
		args string
		want int
		err  error
	}{"valid version for 4.13", "4.13", 1024, nil}),
	Entry("invalid version for 4.14", struct {
		name string
		args string
		want int
		err  error
	}{"invalid version for 4.14", "4.14", 16384, nil}),
	Entry("invalid version for 4.15", struct {
		name string
		args string
		want int
		err  error
	}{"invalid version for 4.15", "4.15", 16384, nil}),
)
