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

var _ = Describe("ResourceExisting", func() {
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

	Context("NAT Gateway", func() {
		It("should return true when NAT gateway is available", func() {
			mockEC2Client.EXPECT().
				DescribeNatGateways(gomock.Any(), gomock.Any()).
				Return(&ec2.DescribeNatGatewaysOutput{
					NatGateways: []types.NatGateway{
						{
							NatGatewayId: aws.String("nat-0cdb4695c3b5fa603"),
							State:        types.NatGatewayStateAvailable,
						},
					},
				}, nil)

			exists := client.ResourceExisting("nat-0cdb4695c3b5fa603")
			Expect(exists).To(BeTrue())
		})

		It("should return false when NAT gateway is pending", func() {
			mockEC2Client.EXPECT().
				DescribeNatGateways(gomock.Any(), gomock.Any()).
				Return(&ec2.DescribeNatGatewaysOutput{
					NatGateways: []types.NatGateway{
						{
							NatGatewayId: aws.String("nat-0cdb4695c3b5fa603"),
							State:        types.NatGatewayStatePending,
						},
					},
				}, nil)

			exists := client.ResourceExisting("nat-0cdb4695c3b5fa603")
			Expect(exists).To(BeFalse())
		})

		It("should return false when NAT gateway is failed", func() {
			mockEC2Client.EXPECT().
				DescribeNatGateways(gomock.Any(), gomock.Any()).
				Return(&ec2.DescribeNatGatewaysOutput{
					NatGateways: []types.NatGateway{
						{
							NatGatewayId: aws.String("nat-0cdb4695c3b5fa603"),
							State:        types.NatGatewayStateFailed,
						},
					},
				}, nil)

			exists := client.ResourceExisting("nat-0cdb4695c3b5fa603")
			Expect(exists).To(BeFalse())
		})

		It("should return false when NAT gateway is not found", func() {
			mockEC2Client.EXPECT().
				DescribeNatGateways(gomock.Any(), gomock.Any()).
				Return(&ec2.DescribeNatGatewaysOutput{
					NatGateways: []types.NatGateway{},
				}, nil)

			exists := client.ResourceExisting("nat-0cdb4695c3b5fa603")
			Expect(exists).To(BeFalse())
		})

		It("should return false when DescribeNatGateways returns NotFound error", func() {
			mockEC2Client.EXPECT().
				DescribeNatGateways(gomock.Any(), gomock.Any()).
				Return(nil, &smithy.GenericAPIError{
					Code:    awserrors.InvalidNatGatewayID,
					Message: "The NAT gateway 'nat-0cdb4695c3b5fa603' does not exist",
				})

			exists := client.ResourceExisting("nat-0cdb4695c3b5fa603")
			Expect(exists).To(BeFalse())
		})
	})
})

var _ = Describe("ResourceDeleted", func() {
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

	Context("NAT Gateway", func() {
		It("should return true when NAT gateway is deleted", func() {
			mockEC2Client.EXPECT().
				DescribeNatGateways(gomock.Any(), gomock.Any()).
				Return(&ec2.DescribeNatGatewaysOutput{
					NatGateways: []types.NatGateway{
						{
							NatGatewayId: aws.String("nat-0cdb4695c3b5fa603"),
							State:        types.NatGatewayStateDeleted,
						},
					},
				}, nil)

			deleted := client.ResourceDeleted("nat-0cdb4695c3b5fa603")
			Expect(deleted).To(BeTrue())
		})

		It("should return false when NAT gateway is available", func() {
			mockEC2Client.EXPECT().
				DescribeNatGateways(gomock.Any(), gomock.Any()).
				Return(&ec2.DescribeNatGatewaysOutput{
					NatGateways: []types.NatGateway{
						{
							NatGatewayId: aws.String("nat-0cdb4695c3b5fa603"),
							State:        types.NatGatewayStateAvailable,
						},
					},
				}, nil)

			deleted := client.ResourceDeleted("nat-0cdb4695c3b5fa603")
			Expect(deleted).To(BeFalse())
		})

		It("should return true when DescribeNatGateways returns NotFound error", func() {
			mockEC2Client.EXPECT().
				DescribeNatGateways(gomock.Any(), gomock.Any()).
				Return(nil, &smithy.GenericAPIError{
					Code:    awserrors.InvalidNatGatewayID,
					Message: "The NAT gateway 'nat-0cdb4695c3b5fa603' does not exist",
				})

			deleted := client.ResourceDeleted("nat-0cdb4695c3b5fa603")
			Expect(deleted).To(BeTrue())
		})
	})
})

var _ = Describe("WaitForResourceExisting", func() {
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

	Context("NAT Gateway", func() {
		It("should return nil when resource becomes available", func() {
			// First call returns pending, second call returns available
			gomock.InOrder(
				mockEC2Client.EXPECT().
					DescribeNatGateways(gomock.Any(), gomock.Any()).
					Return(&ec2.DescribeNatGatewaysOutput{
						NatGateways: []types.NatGateway{
							{
								NatGatewayId: aws.String("nat-0cdb4695c3b5fa603"),
								State:        types.NatGatewayStatePending,
							},
						},
					}, nil),
				mockEC2Client.EXPECT().
					DescribeNatGateways(gomock.Any(), gomock.Any()).
					Return(&ec2.DescribeNatGatewaysOutput{
						NatGateways: []types.NatGateway{
							{
								NatGatewayId: aws.String("nat-0cdb4695c3b5fa603"),
								State:        types.NatGatewayStateAvailable,
							},
						},
					}, nil),
			)

			err := client.WaitForResourceExisting("nat-0cdb4695c3b5fa603", 10)
			Expect(err).To(BeNil())
		})

		It("should return timeout error when resource does not become available within timeout", func() {
			mockEC2Client.EXPECT().
				DescribeNatGateways(gomock.Any(), gomock.Any()).
				Return(&ec2.DescribeNatGatewaysOutput{
					NatGateways: []types.NatGateway{
						{
							NatGatewayId: aws.String("nat-0cdb4695c3b5fa603"),
							State:        types.NatGatewayStatePending,
						},
					},
				}, nil).
				AnyTimes()

			startTime := time.Now()
			err := client.WaitForResourceExisting("nat-0cdb4695c3b5fa603", 3)
			duration := time.Since(startTime)

			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("timeout after 3 seconds"))
			Expect(duration).To(BeNumerically(">=", 3*time.Second))
		})

		It("should timeout when NAT gateway enters failed state", func() {
			// First call returns pending, second call returns failed
			gomock.InOrder(
				mockEC2Client.EXPECT().
					DescribeNatGateways(gomock.Any(), gomock.Any()).
					Return(&ec2.DescribeNatGatewaysOutput{
						NatGateways: []types.NatGateway{
							{
								NatGatewayId: aws.String("nat-0cdb4695c3b5fa603"),
								State:        types.NatGatewayStatePending,
							},
						},
					}, nil),
				mockEC2Client.EXPECT().
					DescribeNatGateways(gomock.Any(), gomock.Any()).
					Return(&ec2.DescribeNatGatewaysOutput{
						NatGateways: []types.NatGateway{
							{
								NatGatewayId: aws.String("nat-0cdb4695c3b5fa603"),
								State:        types.NatGatewayStateFailed,
							},
						},
					}, nil).
					AnyTimes(),
			)

			startTime := time.Now()
			err := client.WaitForResourceExisting("nat-0cdb4695c3b5fa603", 60)
			duration := time.Since(startTime)

			// Should timeout because failed state returns false, not an error
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("timeout after 60 seconds"))
			// Should keep polling until timeout even in failed state
			Expect(duration).To(BeNumerically(">=", 60*time.Second))
		})
	})
})
