/*
Copyright (c) 2024 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package connection

import (
	"testing"

	. "github.com/onsi/ginkgo/v2" // nolint
	. "github.com/onsi/gomega"    // nolint
)

func TestConnection(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Connection")
}

var _ = Describe("ConnectionBuilder", func() {
	Describe("getLogger", func() {
		It("Returns a default logger when no explicit logger is set", func() {
			builder := NewConnection()
			logger, err := builder.getLogger()
			Expect(err).ToNot(HaveOccurred())
			Expect(logger).ToNot(BeNil())
		})

		It("Returns the configured logger when one is set via WithLogger", func() {
			builder := NewConnection()

			// Get a default logger to use as our custom logger
			customLogger, err := builder.getLogger()
			Expect(err).ToNot(HaveOccurred())

			builder2 := NewConnection().WithLogger(customLogger)
			logger, err := builder2.getLogger()
			Expect(err).ToNot(HaveOccurred())
			Expect(logger).To(Equal(customLogger))
		})
	})
})
