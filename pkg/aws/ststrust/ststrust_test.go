package ststrust_test

import (
	"encoding/json"
	"errors"
	"net/url"
	"os"
	"path/filepath"
	"runtime"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/openshift-online/ocm-common/pkg/aws/ststrust"
)

const (
	externalIDA = "223B9588-36A5-ECA4-BE8D-7C673B77CEC1"
	externalIDB = "333B9588-36A5-ECA4-BE8D-7C673B77CDCD"
)

func loadFixture(name string) string {
	_, filename, _, ok := runtime.Caller(0)
	Expect(ok).To(BeTrue())
	path := filepath.Join(filepath.Dir(filename), "testdata", name)
	data, err := os.ReadFile(path)
	Expect(err).NotTo(HaveOccurred())
	return string(data)
}

func policyWithExternalID(externalID string) string {
	policy := map[string]interface{}{
		"Version": "2012-10-17",
		"Statement": []map[string]interface{}{
			{
				"Effect": "Allow",
				"Principal": map[string]interface{}{
					"AWS": "arn:aws:iam::123456789012:role/test",
				},
				"Action": "sts:AssumeRole",
				"Condition": map[string]interface{}{
					"StringEquals": map[string]interface{}{
						"sts:ExternalId": externalID,
					},
				},
			},
		},
	}
	out, err := json.Marshal(policy)
	Expect(err).NotTo(HaveOccurred())
	return string(out)
}

func policyWithMultipleExternalIDs(ids []string) string {
	policy := map[string]interface{}{
		"Version": "2012-10-17",
		"Statement": []map[string]interface{}{
			{
				"Effect": "Allow",
				"Principal": map[string]interface{}{
					"AWS": "arn:aws:iam::123456789012:role/test",
				},
				"Action": "sts:AssumeRole",
				"Condition": map[string]interface{}{
					"StringEquals": map[string]interface{}{
						"sts:ExternalId": ids,
					},
				},
			},
		},
	}
	out, err := json.Marshal(policy)
	Expect(err).NotTo(HaveOccurred())
	return string(out)
}

func policyWithTwoStatements(idA, idB string) string {
	policy := map[string]interface{}{
		"Version": "2012-10-17",
		"Statement": []map[string]interface{}{
			{
				"Effect": "Allow",
				"Principal": map[string]interface{}{
					"AWS": "arn:aws:iam::123456789012:role/test",
				},
				"Action": "sts:AssumeRole",
				"Condition": map[string]interface{}{
					"StringEquals": map[string]interface{}{
						"sts:ExternalId": idA,
					},
				},
			},
			{
				"Effect": "Allow",
				"Principal": map[string]interface{}{
					"AWS": "arn:aws:iam::123456789012:role/other",
				},
				"Action": "sts:AssumeRole",
				"Condition": map[string]interface{}{
					"StringEquals": map[string]interface{}{
						"sts:ExternalId": idB,
					},
				},
			},
		},
	}
	out, err := json.Marshal(policy)
	Expect(err).NotTo(HaveOccurred())
	return string(out)
}

var _ = Describe("STS external ID trust policy", func() {
	Describe("ValidateSTSExternalIDFormat", func() {
		It("accepts a valid external ID", func() {
			Expect(ststrust.ValidateSTSExternalIDFormat(externalIDA)).To(Succeed())
		})

		It("rejects empty external ID", func() {
			err := ststrust.ValidateSTSExternalIDFormat("")
			Expect(errors.Is(err, ststrust.ErrExternalIDEmpty)).To(BeTrue())
		})

		It("rejects too-short external ID", func() {
			err := ststrust.ValidateSTSExternalIDFormat("a")
			Expect(errors.Is(err, ststrust.ErrExternalIDFormat)).To(BeTrue())
		})

		It("rejects too-long external ID", func() {
			longID := make([]byte, ststrust.MaxSTSExternalIDLength+1)
			for i := range longID {
				longID[i] = 'a'
			}
			err := ststrust.ValidateSTSExternalIDFormat(string(longID))
			Expect(errors.Is(err, ststrust.ErrExternalIDFormat)).To(BeTrue())
		})

		It("rejects external ID with invalid characters", func() {
			err := ststrust.ValidateSTSExternalIDFormat("invalid id with spaces")
			Expect(errors.Is(err, ststrust.ErrExternalIDFormat)).To(BeTrue())
		})
	})

	Describe("CollectSTSExternalIDsFromTrustPolicy", func() {
		It("returns nil for empty policy", func() {
			ids, err := ststrust.CollectSTSExternalIDsFromTrustPolicy("")
			Expect(err).NotTo(HaveOccurred())
			Expect(ids).To(BeNil())
		})

		It("returns nil for installer fixture without ExternalId", func() {
			ids, err := ststrust.CollectSTSExternalIDsFromTrustPolicy(loadFixture("sts_installer_trust_policy.json"))
			Expect(err).NotTo(HaveOccurred())
			Expect(ids).To(BeEmpty())
		})

		It("returns nil for support fixture without ExternalId", func() {
			ids, err := ststrust.CollectSTSExternalIDsFromTrustPolicy(loadFixture("sts_support_trust_policy.json"))
			Expect(err).NotTo(HaveOccurred())
			Expect(ids).To(BeEmpty())
		})

		It("collects a single ExternalId", func() {
			ids, err := ststrust.CollectSTSExternalIDsFromTrustPolicy(policyWithExternalID(externalIDA))
			Expect(err).NotTo(HaveOccurred())
			Expect(ids).To(Equal([]string{externalIDA}))
		})

		It("collects ExternalIds from percent-encoded policy JSON", func() {
			escaped := url.PathEscape(policyWithExternalID(externalIDA))
			ids, err := ststrust.CollectSTSExternalIDsFromTrustPolicy(escaped)
			Expect(err).NotTo(HaveOccurred())
			Expect(ids).To(Equal([]string{externalIDA}))
		})

		It("preserves plus signs in ExternalId values", func() {
			const idWithPlus = "ab+cde"
			ids, err := ststrust.CollectSTSExternalIDsFromTrustPolicy(policyWithExternalID(idWithPlus))
			Expect(err).NotTo(HaveOccurred())
			Expect(ids).To(Equal([]string{idWithPlus}))
		})

		It("collects multiple ExternalIds from an array condition", func() {
			ids, err := ststrust.CollectSTSExternalIDsFromTrustPolicy(policyWithMultipleExternalIDs([]string{externalIDA, externalIDB}))
			Expect(err).NotTo(HaveOccurred())
			Expect(ids).To(Equal([]string{externalIDA, externalIDB}))
		})

		It("collects ExternalIds across multiple statements", func() {
			ids, err := ststrust.CollectSTSExternalIDsFromTrustPolicy(policyWithTwoStatements(externalIDA, externalIDB))
			Expect(err).NotTo(HaveOccurred())
			Expect(ids).To(Equal([]string{externalIDA, externalIDB}))
		})

		It("collects ExternalId from StringEqualsIfExists", func() {
			policy := map[string]interface{}{
				"Version": "2012-10-17",
				"Statement": []map[string]interface{}{
					{
						"Effect":  "Allow",
						"Action":  []string{"sts:AssumeRole"},
						"Condition": map[string]interface{}{
							"StringEqualsIfExists": map[string]interface{}{
								"sts:ExternalId": externalIDA,
							},
						},
					},
				},
			}
			out, err := json.Marshal(policy)
			Expect(err).NotTo(HaveOccurred())
			ids, err := ststrust.CollectSTSExternalIDsFromTrustPolicy(string(out))
			Expect(err).NotTo(HaveOccurred())
			Expect(ids).To(Equal([]string{externalIDA}))
		})

		It("returns an error for invalid JSON", func() {
			_, err := ststrust.CollectSTSExternalIDsFromTrustPolicy("{not-json")
			Expect(err).To(HaveOccurred())
		})

		It("collects from policies with Action as a JSON array", func() {
			policy := `{
				"Version": "2012-10-17",
				"Statement": [{
					"Effect": "Allow",
					"Action": ["sts:AssumeRole"],
					"Condition": {
						"StringEquals": { "sts:ExternalId": "` + externalIDA + `" }
					}
				}]
			}`
			ids, err := ststrust.CollectSTSExternalIDsFromTrustPolicy(policy)
			Expect(err).NotTo(HaveOccurred())
			Expect(ids).To(Equal([]string{externalIDA}))
		})

		It("ignores non-AssumeRole statements", func() {
			policy := map[string]interface{}{
				"Version": "2012-10-17",
				"Statement": []map[string]interface{}{
					{
						"Effect": "Allow",
						"Action": "s3:GetObject",
						"Condition": map[string]interface{}{
							"StringEquals": map[string]interface{}{
								"sts:ExternalId": externalIDA,
							},
						},
					},
				},
			}
			out, err := json.Marshal(policy)
			Expect(err).NotTo(HaveOccurred())
			ids, err := ststrust.CollectSTSExternalIDsFromTrustPolicy(string(out))
			Expect(err).NotTo(HaveOccurred())
			Expect(ids).To(BeEmpty())
		})
	})

	Describe("ExternalIDMatchesTrustPolicy", func() {
		It("returns true when entered is in the policy", func() {
			match, err := ststrust.ExternalIDMatchesTrustPolicy(externalIDA, policyWithExternalID(externalIDA))
			Expect(err).NotTo(HaveOccurred())
			Expect(match).To(BeTrue())
		})

		It("returns false when entered is not in the policy", func() {
			match, err := ststrust.ExternalIDMatchesTrustPolicy(externalIDA, policyWithExternalID(externalIDB))
			Expect(err).NotTo(HaveOccurred())
			Expect(match).To(BeFalse())
		})

		It("returns error when entered is empty", func() {
			_, err := ststrust.ExternalIDMatchesTrustPolicy("", policyWithExternalID(externalIDA))
			Expect(errors.Is(err, ststrust.ErrExternalIDEmpty)).To(BeTrue())
		})
	})

	Describe("ApplySTSExternalIDToTrustPolicy", func() {
		It("injects ExternalId into installer fixture", func() {
			base := loadFixture("sts_installer_trust_policy.json")
			updated, err := ststrust.ApplySTSExternalIDToTrustPolicy(base, externalIDA)
			Expect(err).NotTo(HaveOccurred())
			ids, err := ststrust.CollectSTSExternalIDsFromTrustPolicy(updated)
			Expect(err).NotTo(HaveOccurred())
			Expect(ids).To(Equal([]string{externalIDA}))
		})

		It("injects ExternalId into support fixture", func() {
			base := loadFixture("sts_support_trust_policy.json")
			updated, err := ststrust.ApplySTSExternalIDToTrustPolicy(base, externalIDA)
			Expect(err).NotTo(HaveOccurred())
			ids, err := ststrust.CollectSTSExternalIDsFromTrustPolicy(updated)
			Expect(err).NotTo(HaveOccurred())
			Expect(ids).To(Equal([]string{externalIDA}))
		})

		It("is idempotent when ExternalId already present", func() {
			policy := policyWithExternalID(externalIDA)
			updated, err := ststrust.ApplySTSExternalIDToTrustPolicy(policy, externalIDA)
			Expect(err).NotTo(HaveOccurred())
			Expect(updated).To(Equal(policy))
		})

		It("fails when existing policy defines a different ExternalId set", func() {
			policy := policyWithExternalID(externalIDB)
			_, err := ststrust.ApplySTSExternalIDToTrustPolicy(policy, externalIDA)
			Expect(errors.Is(err, ststrust.ErrExternalIDConflictOnInject)).To(BeTrue())
		})

		It("allows injection when entered matches one of multiple IDs", func() {
			policy := policyWithMultipleExternalIDs([]string{externalIDA, externalIDB})
			updated, err := ststrust.ApplySTSExternalIDToTrustPolicy(policy, externalIDA)
			Expect(err).NotTo(HaveOccurred())
			Expect(updated).To(Equal(policy))
		})

		It("rejects empty external ID", func() {
			_, err := ststrust.ApplySTSExternalIDToTrustPolicy(loadFixture("sts_installer_trust_policy.json"), "")
			Expect(errors.Is(err, ststrust.ErrExternalIDEmpty)).To(BeTrue())
		})

		It("rejects empty trust policy document", func() {
			_, err := ststrust.ApplySTSExternalIDToTrustPolicy("", externalIDA)
			Expect(err).To(MatchError(ContainSubstring("trust policy document is empty")))
		})

		It("rejects invalid external ID format", func() {
			_, err := ststrust.ApplySTSExternalIDToTrustPolicy(loadFixture("sts_installer_trust_policy.json"), "x")
			Expect(errors.Is(err, ststrust.ErrExternalIDFormat)).To(BeTrue())
		})
	})

	Describe("ValidateEnteredForRoleTrustPolicies", func() {
		It("succeeds when entered matches both roles", func() {
			installer := policyWithMultipleExternalIDs([]string{externalIDA, externalIDB})
			support := policyWithExternalID(externalIDA)
			err := ststrust.ValidateEnteredForRoleTrustPolicies(externalIDA, installer, support)
			Expect(err).NotTo(HaveOccurred())
		})

		It("fails when support role does not contain entered ID", func() {
			installer := policyWithExternalID(externalIDA)
			support := policyWithExternalID(externalIDB)
			err := ststrust.ValidateEnteredForRoleTrustPolicies(externalIDA, installer, support)
			var mismatch *ststrust.ExternalIDMismatchError
			Expect(errors.As(err, &mismatch)).To(BeTrue())
			Expect(mismatch.RoleLabel).To(Equal("support role"))
		})

		It("fails when installer trust policy has no ExternalId", func() {
			err := ststrust.ValidateEnteredForRoleTrustPolicies(
				externalIDA,
				loadFixture("sts_installer_trust_policy.json"),
				"",
			)
			Expect(errors.Is(err, ststrust.ErrNoTrustPolicyExternalID)).To(BeTrue())
			Expect(err.Error()).To(ContainSubstring("installer role"))
		})

		It("fails when support trust policy has no ExternalId", func() {
			err := ststrust.ValidateEnteredForRoleTrustPolicies(
				externalIDA,
				"",
				loadFixture("sts_support_trust_policy.json"),
			)
			Expect(errors.Is(err, ststrust.ErrNoTrustPolicyExternalID)).To(BeTrue())
			Expect(err.Error()).To(ContainSubstring("support role"))
		})

		It("succeeds when only support policy is provided and matches", func() {
			support := policyWithExternalID(externalIDA)
			err := ststrust.ValidateEnteredForRoleTrustPolicies(externalIDA, "", support)
			Expect(err).NotTo(HaveOccurred())
		})

		It("fails when neither installer nor support policy is provided", func() {
			err := ststrust.ValidateEnteredForRoleTrustPolicies(externalIDA, "", "")
			Expect(errors.Is(err, ststrust.ErrNoTrustPolicyExternalID)).To(BeTrue())
		})

		It("fails when installer role does not contain entered ID", func() {
			installer := policyWithExternalID(externalIDB)
			support := policyWithExternalID(externalIDA)
			err := ststrust.ValidateEnteredForRoleTrustPolicies(externalIDA, installer, support)
			var mismatch *ststrust.ExternalIDMismatchError
			Expect(errors.As(err, &mismatch)).To(BeTrue())
			Expect(mismatch.RoleLabel).To(Equal("installer role"))
			Expect(errors.Is(err, ststrust.ErrExternalIDNotInTrustPolicy)).To(BeTrue())
			Expect(mismatch.Error()).To(ContainSubstring(externalIDA))
		})

		It("rejects invalid entered format before checking policies", func() {
			err := ststrust.ValidateEnteredForRoleTrustPolicies("x", policyWithExternalID(externalIDA), "")
			Expect(errors.Is(err, ststrust.ErrExternalIDFormat)).To(BeTrue())
		})
	})

	Describe("DiscoverSTSExternalID", func() {
		It("returns empty when no ExternalIds exist", func() {
			id, err := ststrust.DiscoverSTSExternalID(
				loadFixture("sts_installer_trust_policy.json"),
				loadFixture("sts_support_trust_policy.json"),
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(id).To(BeEmpty())
		})

		It("returns sole union member from installer policy only", func() {
			installer := policyWithExternalID(externalIDA)
			id, err := ststrust.DiscoverSTSExternalID(installer, "")
			Expect(err).NotTo(HaveOccurred())
			Expect(id).To(Equal(externalIDA))
		})

		It("returns sole union member from support policy only", func() {
			support := policyWithExternalID(externalIDA)
			id, err := ststrust.DiscoverSTSExternalID("", support)
			Expect(err).NotTo(HaveOccurred())
			Expect(id).To(Equal(externalIDA))
		})

		It("returns sole union member from support fixture with ExternalId", func() {
			support := loadFixture("sts_support_trust_policy.json")
			updated, err := ststrust.ApplySTSExternalIDToTrustPolicy(support, externalIDA)
			Expect(err).NotTo(HaveOccurred())
			id, err := ststrust.DiscoverSTSExternalID("", updated)
			Expect(err).NotTo(HaveOccurred())
			Expect(id).To(Equal(externalIDA))
		})

		It("returns intersection member when union is ambiguous", func() {
			installer := policyWithMultipleExternalIDs([]string{externalIDA, externalIDB})
			support := policyWithExternalID(externalIDA)
			id, err := ststrust.DiscoverSTSExternalID(installer, support)
			Expect(err).NotTo(HaveOccurred())
			Expect(id).To(Equal(externalIDA))
		})

		It("returns empty when union has multiple IDs with no single intersection", func() {
			installer := policyWithExternalID(externalIDA)
			support := policyWithExternalID(externalIDB)
			id, err := ststrust.DiscoverSTSExternalID(installer, support)
			Expect(err).NotTo(HaveOccurred())
			Expect(id).To(BeEmpty())
		})
	})

	Describe("CanInjectSTSExternalID", func() {
		It("allows injection into installer fixture without ExternalId", func() {
			err := ststrust.CanInjectSTSExternalID(loadFixture("sts_installer_trust_policy.json"), externalIDA)
			Expect(err).NotTo(HaveOccurred())
		})

		It("allows injection into support fixture without ExternalId", func() {
			err := ststrust.CanInjectSTSExternalID(loadFixture("sts_support_trust_policy.json"), externalIDA)
			Expect(err).NotTo(HaveOccurred())
		})

		It("allows when entered is in existing set", func() {
			err := ststrust.CanInjectSTSExternalID(policyWithMultipleExternalIDs([]string{externalIDA, externalIDB}), externalIDA)
			Expect(err).NotTo(HaveOccurred())
		})

		It("rejects empty external ID", func() {
			err := ststrust.CanInjectSTSExternalID(loadFixture("sts_installer_trust_policy.json"), "")
			Expect(errors.Is(err, ststrust.ErrExternalIDEmpty)).To(BeTrue())
		})

		It("fails when entered is not in existing ExternalId set", func() {
			err := ststrust.CanInjectSTSExternalID(policyWithExternalID(externalIDB), externalIDA)
			Expect(errors.Is(err, ststrust.ErrExternalIDConflictOnInject)).To(BeTrue())
		})
	})

	Describe("DiscoverSTSExternalID error paths", func() {
		It("returns error when installer policy JSON is invalid", func() {
			_, err := ststrust.DiscoverSTSExternalID("{bad", policyWithExternalID(externalIDA))
			Expect(err).To(HaveOccurred())
		})

		It("returns error when support policy JSON is invalid", func() {
			_, err := ststrust.DiscoverSTSExternalID(policyWithExternalID(externalIDA), "{bad")
			Expect(err).To(HaveOccurred())
		})
	})
})
