package validations

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/openshift/rosa/pkg/aws/mocks"
)

var _ = Describe("GetRole", func() {
	var (
		mockCtrl   *gomock.Controller
		mockIamAPI *mocks.MockIAMAPI
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockIamAPI = mocks.NewMockIAMAPI(mockCtrl)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("When IAM client returns a role", func() {
		BeforeEach(func() {
			mockIamAPI.EXPECT().GetRole(gomock.Any()).Return(&iam.GetRoleOutput{
				Role: &iam.Role{
					RoleName: aws.String("role-name"),
				},
			}, nil)
		})

		It("should return the role", func() {
			role, err := GetRole(mockIamAPI, "role-name")
			Expect(err).ToNot(HaveOccurred())
			Expect(role).To(Equal(&iam.GetRoleOutput{
				Role: &iam.Role{
					RoleName: aws.String("role-name"),
				},
			}))
		})
	})

	Context("When IAM client returns an error", func() {
		BeforeEach(func() {
			mockIamAPI.EXPECT().GetRole(gomock.Any()).Return(nil, errors.New("some-error"))
		})

		It("should return the error", func() {
			role, err := GetRole(mockIamAPI, "role-name")
			Expect(err).To(HaveOccurred())
			Expect(role).To(BeNil())
		})
	})
})
