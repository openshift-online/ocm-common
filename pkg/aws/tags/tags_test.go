package tags_test

import (
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/openshift-online/ocm-common/pkg/aws/tags"
)

var _ = Describe("Validate", func() {
	It("accepts nil and empty maps", func() {
		Expect(tags.Validate(nil)).To(Succeed())
		Expect(tags.Validate(map[string]string{})).To(Succeed())
	})

	It("accepts valid custom tags", func() {
		Expect(tags.Validate(map[string]string{
			"env":         "prod",
			"owner":       "platform",
			"cost.center": "123",
		})).To(Succeed())
	})

	It("accepts empty values", func() {
		Expect(tags.Validate(map[string]string{"flag": ""})).To(Succeed())
	})

	It("accepts ordinary spaces in keys and values", func() {
		Expect(tags.Validate(map[string]string{
			"cost center": "team a",
		})).To(Succeed())
	})

	It("rejects empty keys", func() {
		err := tags.Validate(map[string]string{"": "value"})
		Expect(err).To(MatchError(tags.ErrEmptyKey))
	})

	It("rejects keys that are too long", func() {
		key := strings.Repeat("a", tags.MaxKeyLength+1)
		err := tags.Validate(map[string]string{key: "v"})
		Expect(err).To(MatchError(tags.ErrKeyTooLong))
	})

	It("rejects values that are too long", func() {
		value := strings.Repeat("a", tags.MaxValueLength+1)
		err := tags.Validate(map[string]string{"k": value})
		Expect(err).To(MatchError(tags.ErrValueTooLong))
	})

	It("rejects tabs and other non-space whitespace in keys and values", func() {
		Expect(tags.Validate(map[string]string{"bad\tkey": "v"})).To(MatchError(tags.ErrKeyFormat))
		Expect(tags.Validate(map[string]string{"k": "bad\tvalue"})).To(MatchError(tags.ErrValueFormat))
		Expect(tags.Validate(map[string]string{"bad\nkey": "v"})).To(MatchError(tags.ErrKeyFormat))
		Expect(tags.Validate(map[string]string{"k": "bad\nvalue"})).To(MatchError(tags.ErrValueFormat))
	})

	DescribeTable("reserved keys",
		func(key string) {
			err := tags.Validate(map[string]string{key: "v"})
			Expect(err).To(MatchError(tags.ErrReservedKey))
		},
		Entry("aws", tags.ReservedKeyAWS),
		Entry("red-hat-managed", tags.ReservedKeyRedHatManaged),
		Entry("red-hat-clustertype", tags.ReservedKeyRedHatClusterType),
		Entry("Name", tags.ReservedKeyName),
		Entry("aws: prefix", "aws:cloudformation:stack-name"),
		Entry("AWS: prefix case", "AWS:Something"),
		Entry("kubernetes.io/cluster/ prefix", "kubernetes.io/cluster/mycluster"),
	)
})

var _ = Describe("Equal and diff", func() {
	It("Equal treats nil and empty as equal", func() {
		Expect(tags.Equal(nil, nil)).To(BeTrue())
		Expect(tags.Equal(nil, map[string]string{})).To(BeTrue())
		Expect(tags.Equal(map[string]string{"a": "1"}, map[string]string{"a": "1"})).To(BeTrue())
		Expect(tags.Equal(map[string]string{"a": "1"}, map[string]string{"a": "2"})).To(BeFalse())
		Expect(tags.Equal(map[string]string{"a": "1"}, map[string]string{"b": "1"})).To(BeFalse())
	})

	It("computeDiff returns empty when maps match", func() {
		d := tags.ComputeDiff(
			map[string]string{"a": "1"},
			map[string]string{"a": "1"},
		)
		Expect(d.Empty()).To(BeTrue())
		Expect(d.Add).To(BeEmpty())
		Expect(d.Update).To(BeEmpty())
		Expect(d.Remove).To(BeEmpty())
	})

	It("computeDiff reports add update and remove", func() {
		current := map[string]string{
			"keep":   "same",
			"change": "old",
			"gone":   "x",
		}
		desired := map[string]string{
			"keep":   "same",
			"change": "new",
			"new":    "y",
		}
		d := tags.ComputeDiff(current, desired)
		Expect(d.Empty()).To(BeFalse())
		Expect(d.Add).To(Equal(map[string]string{"new": "y"}))
		Expect(d.Update).To(Equal(map[string]string{"change": "new"}))
		Expect(d.Remove).To(ConsistOf("gone"))
	})

	It("computeDiff handles nil maps", func() {
		d := tags.ComputeDiff(nil, map[string]string{"a": "1"})
		Expect(d.Add).To(Equal(map[string]string{"a": "1"}))
		Expect(d.Remove).To(BeEmpty())

		d = tags.ComputeDiff(map[string]string{"a": "1"}, nil)
		Expect(d.Add).To(BeEmpty())
		Expect(d.Remove).To(ConsistOf("a"))
	})
})
