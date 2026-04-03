package vpc_client_test

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	smithy "github.com/aws/smithy-go"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"

	"github.com/openshift-online/ocm-common/pkg/aws/aws_client"
	awserrors "github.com/openshift-online/ocm-common/pkg/aws/errors"
	. "github.com/openshift-online/ocm-common/pkg/test/vpc_client"
)

var _ = Describe("CreateAdditionalSecurityGroups", func() {
	var (
		mockCtrl      *gomock.Controller
		mockEC2Client *aws_client.MockEC2ClientAPI
		vpc           *VPC
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockEC2Client = aws_client.NewMockEC2ClientAPI(mockCtrl)
		vpc = NewVPC().
			ID("vpc-0fa4cc34953703260").
			CIDR("10.0.0.0/16").
			AWSclient(&aws_client.AWSClient{Ec2Client: mockEC2Client})
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	It("should create a new security group successfully", func() {
		mockEC2Client.EXPECT().
			CreateSecurityGroup(gomock.Any(), gomock.Any()).
			Return(&ec2.CreateSecurityGroupOutput{
				GroupId: aws.String("sg-newgroup"),
			}, nil)
		mockEC2Client.EXPECT().
			DescribeSecurityGroups(gomock.Any(), gomock.Any()).
			Return(&ec2.DescribeSecurityGroupsOutput{
				SecurityGroups: []types.SecurityGroup{{GroupId: aws.String("sg-newgroup")}},
			}, nil).AnyTimes()
		mockEC2Client.EXPECT().
			CreateTags(gomock.Any(), gomock.Any()).
			Return(&ec2.CreateTagsOutput{}, nil).AnyTimes()
		mockEC2Client.EXPECT().
			AuthorizeSecurityGroupIngress(gomock.Any(), gomock.Any()).
			Return(&ec2.AuthorizeSecurityGroupIngressOutput{}, nil).AnyTimes()

		sgIDs, err := vpc.CreateAdditionalSecurityGroups(1, "test-sg", "test description")
		Expect(err).To(BeNil())
		Expect(sgIDs).To(HaveLen(1))
		Expect(sgIDs[0]).To(Equal("sg-newgroup"))
	})

	It("should reuse existing security group on InvalidGroup.Duplicate", func() {
		sgName := "bastion-sg-0"
		mockEC2Client.EXPECT().
			CreateSecurityGroup(gomock.Any(), gomock.Any()).
			Return(nil, &smithy.GenericAPIError{
				Code:    awserrors.InvalidGroupDuplicate,
				Message: fmt.Sprintf("The security group '%s' already exists for VPC 'vpc-0fa4cc34953703260'", sgName),
			})
		mockEC2Client.EXPECT().
			DescribeSecurityGroups(gomock.Any(), gomock.Any()).
			Return(&ec2.DescribeSecurityGroupsOutput{
				SecurityGroups: []types.SecurityGroup{
					{
						GroupId:   aws.String("sg-existing"),
						GroupName: aws.String(sgName),
					},
				},
			}, nil)
		mockEC2Client.EXPECT().
			AuthorizeSecurityGroupIngress(gomock.Any(), gomock.Any()).
			Return(&ec2.AuthorizeSecurityGroupIngressOutput{}, nil).AnyTimes()

		sgIDs, err := vpc.CreateAdditionalSecurityGroups(1, "bastion-sg", "bastion description")
		Expect(err).To(BeNil())
		Expect(sgIDs).To(HaveLen(1))
		Expect(sgIDs[0]).To(Equal("sg-existing"))
	})

	It("should return error on non-duplicate failures", func() {
		mockEC2Client.EXPECT().
			CreateSecurityGroup(gomock.Any(), gomock.Any()).
			Return(nil, &smithy.GenericAPIError{
				Code:    "UnauthorizedOperation",
				Message: "You are not authorized to perform this operation.",
			})

		sgIDs, err := vpc.CreateAdditionalSecurityGroups(1, "test-sg", "test description")
		Expect(err).NotTo(BeNil())
		Expect(sgIDs).To(BeEmpty())
	})

	It("should return error when duplicate SG is not found in VPC", func() {
		sgName := "bastion-sg-0"
		mockEC2Client.EXPECT().
			CreateSecurityGroup(gomock.Any(), gomock.Any()).
			Return(nil, &smithy.GenericAPIError{
				Code:    awserrors.InvalidGroupDuplicate,
				Message: fmt.Sprintf("The security group '%s' already exists for VPC 'vpc-0fa4cc34953703260'", sgName),
			})
		mockEC2Client.EXPECT().
			DescribeSecurityGroups(gomock.Any(), gomock.Any()).
			Return(&ec2.DescribeSecurityGroupsOutput{
				SecurityGroups: []types.SecurityGroup{},
			}, nil)

		sgIDs, err := vpc.CreateAdditionalSecurityGroups(1, "bastion-sg", "bastion description")
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(ContainSubstring("reported as duplicate but not found"))
		Expect(sgIDs).To(BeEmpty())
	})
})
