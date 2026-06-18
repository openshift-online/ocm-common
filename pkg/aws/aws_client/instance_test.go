package aws_client_test

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"

	. "github.com/openshift-online/ocm-common/pkg/aws/aws_client"
)

var _ = Describe("WaitForInstancesRunning", func() {
	var (
		mockCtrl      *gomock.Controller
		mockEC2Client *MockEC2ClientAPI
		client        *AWSClient
		ctx           context.Context
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockEC2Client = NewMockEC2ClientAPI(mockCtrl)
		client = &AWSClient{Ec2Client: mockEC2Client}
		ctx = context.Background()
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	It("returns true when the SDK waiter reports all instances are status-ok", func() {
		instanceID := "i-0abc123"

		// The InstanceStatusOkWaiter calls DescribeInstanceStatus until
		// all instances have InstanceStatus.Status == ok.
		// The SDK waiter injects an extra option function (3rd variadic arg).
		mockEC2Client.EXPECT().
			DescribeInstanceStatus(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(&ec2.DescribeInstanceStatusOutput{
				InstanceStatuses: []types.InstanceStatus{
					{
						InstanceId:     aws.String(instanceID),
						InstanceState:  &types.InstanceState{Name: types.InstanceStateNameRunning},
						InstanceStatus: &types.InstanceStatusSummary{Status: types.SummaryStatusOk},
						SystemStatus:   &types.InstanceStatusSummary{Status: types.SummaryStatusOk},
					},
				},
			}, nil).
			AnyTimes()

		allRunning, err := client.WaitForInstancesRunning(ctx, []string{instanceID}, 10)
		Expect(err).NotTo(HaveOccurred())
		Expect(allRunning).To(BeTrue())
	})

	It("propagates errors returned by the SDK waiter", func() {
		// Use a short deadline so the waiter gives up after the first failed
		// call instead of retrying for the full timeout duration.
		shortCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
		defer cancel()

		mockEC2Client.EXPECT().
			DescribeInstanceStatus(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil, fmt.Errorf("api error: request throttled")).
			AnyTimes()

		allRunning, err := client.WaitForInstancesRunning(shortCtx, []string{"i-0abc123"}, 10)
		Expect(err).To(HaveOccurred())
		Expect(allRunning).To(BeFalse())
	})
})

var _ = Describe("WaitForInstancesTerminated", func() {
	var (
		mockCtrl      *gomock.Controller
		mockEC2Client *MockEC2ClientAPI
		client        *AWSClient
		ctx           context.Context
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockEC2Client = NewMockEC2ClientAPI(mockCtrl)
		client = &AWSClient{Ec2Client: mockEC2Client}
		ctx = context.Background()
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	It("returns true when the SDK waiter reports all instances are terminated", func() {
		instanceID := "i-dead"

		// The InstanceTerminatedWaiter calls DescribeInstances until
		// all instances have State.Name == terminated.
		// The SDK waiter injects an extra option function (3rd variadic arg).
		mockEC2Client.EXPECT().
			DescribeInstances(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(&ec2.DescribeInstancesOutput{
				Reservations: []types.Reservation{
					{
						Instances: []types.Instance{
							{
								InstanceId: aws.String(instanceID),
								State:      &types.InstanceState{Name: types.InstanceStateNameTerminated},
							},
						},
					},
				},
			}, nil).
			AnyTimes()

		allTerminated, err := client.WaitForInstancesTerminated(ctx, []string{instanceID}, 10)
		Expect(err).NotTo(HaveOccurred())
		Expect(allTerminated).To(BeTrue())
	})

	It("propagates errors returned by the SDK waiter", func() {
		// Use a short deadline so the waiter gives up quickly.
		shortCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
		defer cancel()

		mockEC2Client.EXPECT().
			DescribeInstances(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil, fmt.Errorf("api error: unauthorized")).
			AnyTimes()

		allTerminated, err := client.WaitForInstancesTerminated(shortCtx, []string{"i-dead"}, 10)
		Expect(err).To(HaveOccurred())
		Expect(allTerminated).To(BeFalse())
	})
})
