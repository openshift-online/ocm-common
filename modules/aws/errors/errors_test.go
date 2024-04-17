package aws

import (
	"errors"
	"github.com/aws/smithy-go"

	iamtypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Error Checks", func() {
	It("should identify NoSuchEntityException", func() {
		err := &iamtypes.NoSuchEntityException{}
		Expect(IsNoSuchEntityException(err)).To(BeTrue())
		Expect(IsNoSuchEntityException(errors.New("random error"))).To(BeFalse())
	})

	It("should identify EntityAlreadyExistsException", func() {
		err := &iamtypes.EntityAlreadyExistsException{}
		Expect(IsEntityAlreadyExistsException(err)).To(BeTrue())
		Expect(IsEntityAlreadyExistsException(errors.New("random error"))).To(BeFalse())
	})

	It("should identify IsThrottle", func() {
		err := &smithy.GenericAPIError{Code: "Throttling"}
		Expect(IsThrottle(err)).To(BeTrue())
		Expect(IsThrottle(errors.New("random error"))).To(BeFalse())
	})

	It("should identify IsAccessDeniedException", func() {
		err := &smithy.GenericAPIError{Code: "AccessDenied"}
		Expect(IsAccessDeniedException(err)).To(BeTrue())
		Expect(IsAccessDeniedException(errors.New("random error"))).To(BeFalse())
	})

	It("should identify IsForbiddenException", func() {
		err := &smithy.GenericAPIError{Code: "Forbidden"}
		Expect(IsForbiddenException(err)).To(BeTrue())
		Expect(IsForbiddenException(errors.New("random error"))).To(BeFalse())
	})

	It("should identify IsLimitExceededException", func() {
		err := &smithy.GenericAPIError{Code: "LimitExceeded"}
		Expect(IsLimitExceededException(err)).To(BeTrue())
		Expect(IsLimitExceededException(errors.New("random error"))).To(BeFalse())
	})

	It("should identify IsInvalidTokenException", func() {
		err := &smithy.GenericAPIError{Code: "InvalidClientTokenId"}
		Expect(IsInvalidTokenException(err)).To(BeTrue())
		Expect(IsInvalidTokenException(errors.New("random error"))).To(BeFalse())
	})

	It("should identify IsSubnetNotFoundError", func() {
		err := &smithy.GenericAPIError{Code: "InvalidSubnetID.NotFound"}
		Expect(IsSubnetNotFoundError(err)).To(BeTrue())
		Expect(IsSubnetNotFoundError(errors.New("random error"))).To(BeFalse())
	})

	It("should identify ErrorCode", func() {
		err := &smithy.GenericAPIError{Code: "AccessDenied"}
		Expect(IsErrorCode(err, "AccessDenied")).To(BeTrue())
		Expect(IsErrorCode(err, "DifferentCode")).To(BeFalse())
	})

	It("should identify IsDeleteConfictException", func() {
		err := &iamtypes.DeleteConflictException{}
		Expect(IsDeleteConfictException(err)).To(BeTrue())
		Expect(IsDeleteConfictException(errors.New("random error"))).To(BeFalse())
	})

})
