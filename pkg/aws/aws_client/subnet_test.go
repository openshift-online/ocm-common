package aws_client_test

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	smithy "github.com/aws/smithy-go"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"

	. "github.com/openshift-online/ocm-common/pkg/aws/aws_client"
	awserrors "github.com/openshift-online/ocm-common/pkg/aws/errors"
)

var _ = Describe("WaitForSubnetAccessibility", func() {
	var (
		mockCtrl      *gomock.Controller
		mockEC2Client *MockEC2ClientAPI
		client        *AWSClient
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockEC2Client = NewMockEC2ClientAPI(mockCtrl)
		client = &AWSClient{
			Ec2Client: mockEC2Client,
		}
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("Empty subnet list", func() {
		It("should return nil immediately when subnet list is empty", func() {
			// No mock expectations needed - should return before calling API

			err := client.WaitForSubnetAccessibility([]string{}, 10)
			Expect(err).To(BeNil())
		})

		It("should return nil immediately when subnet list is nil", func() {
			// No mock expectations needed - should return before calling API

			err := client.WaitForSubnetAccessibility(nil, 10)
			Expect(err).To(BeNil())
		})
	})

	Context("Single subnet", func() {
		It("should return nil when single subnet is immediately accessible", func() {
			mockEC2Client.EXPECT().
				DescribeSubnets(gomock.Any(), gomock.Any()).
				Return(&ec2.DescribeSubnetsOutput{
					Subnets: []types.Subnet{
						{
							SubnetId: aws.String("subnet-12345"),
							State:    types.SubnetStateAvailable,
						},
					},
				}, nil)

			err := client.WaitForSubnetAccessibility([]string{"subnet-12345"}, 10)
			Expect(err).To(BeNil())
		})

		It("should return nil when subnet becomes accessible after initial unavailability", func() {
			// First call returns no subnets (not accessible yet), second call returns the subnet
			gomock.InOrder(
				mockEC2Client.EXPECT().
					DescribeSubnets(gomock.Any(), gomock.Any()).
					Return(&ec2.DescribeSubnetsOutput{
						Subnets: []types.Subnet{},
					}, nil),
				mockEC2Client.EXPECT().
					DescribeSubnets(gomock.Any(), gomock.Any()).
					Return(&ec2.DescribeSubnetsOutput{
						Subnets: []types.Subnet{
							{
								SubnetId: aws.String("subnet-12345"),
								State:    types.SubnetStateAvailable,
							},
						},
					}, nil),
			)

			err := client.WaitForSubnetAccessibility([]string{"subnet-12345"}, 10)
			Expect(err).To(BeNil())
		})

		It("should timeout when subnet never becomes accessible", func() {
			mockEC2Client.EXPECT().
				DescribeSubnets(gomock.Any(), gomock.Any()).
				Return(&ec2.DescribeSubnetsOutput{
					Subnets: []types.Subnet{},
				}, nil).
				AnyTimes()

			startTime := time.Now()
			err := client.WaitForSubnetAccessibility([]string{"subnet-12345"}, 3)
			duration := time.Since(startTime)

			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("timeout after 3 seconds"))
			Expect(err.Error()).To(ContainSubstring("subnet-12345"))
			Expect(duration).To(BeNumerically(">=", 3*time.Second))
		})

		It("should continue polling when DescribeSubnets returns error", func() {
			// First call returns error, second call succeeds
			gomock.InOrder(
				mockEC2Client.EXPECT().
					DescribeSubnets(gomock.Any(), gomock.Any()).
					Return(nil, &smithy.GenericAPIError{
						Code:    awserrors.SubnetNotFound,
						Message: "Subnet not found",
					}),
				mockEC2Client.EXPECT().
					DescribeSubnets(gomock.Any(), gomock.Any()).
					Return(&ec2.DescribeSubnetsOutput{
						Subnets: []types.Subnet{
							{
								SubnetId: aws.String("subnet-12345"),
								State:    types.SubnetStateAvailable,
							},
						},
					}, nil),
			)

			err := client.WaitForSubnetAccessibility([]string{"subnet-12345"}, 10)
			Expect(err).To(BeNil())
		})
	})

	Context("Multiple subnets", func() {
		It("should return nil when all subnets are immediately accessible", func() {
			// Tests: len(returned) == len(requested) on first call
			mockEC2Client.EXPECT().
				DescribeSubnets(gomock.Any(), gomock.Any()).
				Return(&ec2.DescribeSubnetsOutput{
					Subnets: []types.Subnet{
						{
							SubnetId: aws.String("subnet-11111"),
							State:    types.SubnetStateAvailable,
						},
						{
							SubnetId: aws.String("subnet-22222"),
							State:    types.SubnetStateAvailable,
						},
						{
							SubnetId: aws.String("subnet-33333"),
							State:    types.SubnetStateAvailable,
						},
					},
				}, nil)

			err := client.WaitForSubnetAccessibility(
				[]string{"subnet-11111", "subnet-22222", "subnet-33333"}, 10)
			Expect(err).To(BeNil())
		})

		It("should return nil when all subnets become accessible after delay", func() {
			// Tests progressive appearance: 1 → 3 subnets
			// Function waits until len(returned) == len(requested)
			gomock.InOrder(
				mockEC2Client.EXPECT().
					DescribeSubnets(gomock.Any(), gomock.Any()).
					Return(&ec2.DescribeSubnetsOutput{
						Subnets: []types.Subnet{
							{
								SubnetId: aws.String("subnet-11111"),
								State:    types.SubnetStateAvailable,
							},
						},
					}, nil),
				mockEC2Client.EXPECT().
					DescribeSubnets(gomock.Any(), gomock.Any()).
					Return(&ec2.DescribeSubnetsOutput{
						Subnets: []types.Subnet{
							{
								SubnetId: aws.String("subnet-11111"),
								State:    types.SubnetStateAvailable,
							},
							{
								SubnetId: aws.String("subnet-22222"),
								State:    types.SubnetStateAvailable,
							},
							{
								SubnetId: aws.String("subnet-33333"),
								State:    types.SubnetStateAvailable,
							},
						},
					}, nil),
			)

			err := client.WaitForSubnetAccessibility(
				[]string{"subnet-11111", "subnet-22222", "subnet-33333"}, 10)
			Expect(err).To(BeNil())
		})

		It("should timeout when only some subnets become accessible", func() {
			// Always returns only 2 out of 3 subnets
			mockEC2Client.EXPECT().
				DescribeSubnets(gomock.Any(), gomock.Any()).
				Return(&ec2.DescribeSubnetsOutput{
					Subnets: []types.Subnet{
						{
							SubnetId: aws.String("subnet-11111"),
							State:    types.SubnetStateAvailable,
						},
						{
							SubnetId: aws.String("subnet-22222"),
							State:    types.SubnetStateAvailable,
						},
					},
				}, nil).
				AnyTimes()

			startTime := time.Now()
			err := client.WaitForSubnetAccessibility(
				[]string{"subnet-11111", "subnet-22222", "subnet-33333"}, 3)
			duration := time.Since(startTime)

			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("timeout after 3 seconds"))
			Expect(err.Error()).To(ContainSubstring("subnet-11111"))
			Expect(err.Error()).To(ContainSubstring("subnet-22222"))
			Expect(err.Error()).To(ContainSubstring("subnet-33333"))
			Expect(duration).To(BeNumerically(">=", 3*time.Second))
		})
	})

	Context("Duplicate subnet IDs", func() {
		It("should deduplicate subnet IDs before checking", func() {
			// Should only check for 2 unique subnets even though 5 IDs provided
			mockEC2Client.EXPECT().
				DescribeSubnets(gomock.Any(), gomock.Any()).
				Return(&ec2.DescribeSubnetsOutput{
					Subnets: []types.Subnet{
						{
							SubnetId: aws.String("subnet-11111"),
							State:    types.SubnetStateAvailable,
						},
						{
							SubnetId: aws.String("subnet-22222"),
							State:    types.SubnetStateAvailable,
						},
					},
				}, nil)

			// Pass duplicates - should still work
			err := client.WaitForSubnetAccessibility(
				[]string{"subnet-11111", "subnet-22222", "subnet-11111", "subnet-22222", "subnet-11111"}, 10)
			Expect(err).To(BeNil())
		})
	})
})
