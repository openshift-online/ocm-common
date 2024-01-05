package accountroles_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/openshift-online/ocm-common/pkg/rosa/accountroles"
	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
)

var _ = Describe("Account Role functions", func() {
	var _ = Describe("Validates GetPathFromAccountRole function", func() {
		It("Generates the version ID on stable channel", func() {
			fakeCluster, err := cmv1.NewCluster().
				AWS(
					cmv1.NewAWS().
						STS(
							cmv1.NewSTS().
								RoleARN("arn:partition:service:region:account-id:resource-type/test-path/resource-id").
								SupportRoleARN("support").
								InstanceIAMRoles(
									cmv1.NewInstanceIAMRoles().
										MasterRoleARN("controlplane").
										WorkerRoleARN("worker"),
								),
						),
				).Build()
			Expect(err).ToNot(HaveOccurred())
			path, err := GetPathFromAccountRole(fakeCluster, AccountRoles[InstallerAccountRole].Name)
			Expect(err).ToNot(HaveOccurred())
			Expect(path).To(Equal("/test-path/"))
		})
	})

	var _ = Describe("Validates GetAccountRolesArnsMap function", func() {
		It("Checks the account roles are retrieved from a cluster", func() {
			fakeCluster, err := cmv1.NewCluster().
				AWS(
					cmv1.NewAWS().
						STS(
							cmv1.NewSTS().
								RoleARN("installer").
								SupportRoleARN("support").
								InstanceIAMRoles(
									cmv1.NewInstanceIAMRoles().
										MasterRoleARN("controlplane").
										WorkerRoleARN("worker"),
								),
						),
				).Build()
			Expect(err).ToNot(HaveOccurred())
			accRolesMap := GetAccountRolesArnsMap(fakeCluster)
			Expect(accRolesMap[AccountRoles[InstallerAccountRole].Name]).To(Equal("installer"))
			Expect(accRolesMap[AccountRoles[SupportAccountRole].Name]).To(Equal("support"))
			Expect(accRolesMap[AccountRoles[ControlPlaneAccountRole].Name]).To(Equal("controlplane"))
			Expect(accRolesMap[AccountRoles[WorkerAccountRole].Name]).To(Equal("worker"))
		})
	})
})
