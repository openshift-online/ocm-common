package aws_client

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/openshift-online/ocm-common/pkg/log"
)

func (client *AWSClient) CreateCapacityReservation(crInput *ec2.CreateCapacityReservationInput, endDate time.Time, creator string) (*ec2.CreateCapacityReservationOutput, error) {
	crInput.EndDateType = types.EndDateTypeLimited
	crInput.EndDate = aws.Time(endDate)
	crInput.TagSpecifications = []types.TagSpecification{
		{
			ResourceType: types.ResourceTypeCapacityReservation,
			Tags: []types.Tag{
				{
					Key:   aws.String("Creator"),
					Value: aws.String(creator),
				},
				{
					Key:   aws.String("Timestamp"),
					Value: aws.String(time.Now().Format(time.RFC3339)),
				},
			},
		},
	}
	output, err := client.Ec2Client.CreateCapacityReservation(context.Background(), crInput)
	if err != nil {
		log.LogError("Got error creating capacity reservation: %s", err)
	}
	return output, err
}

func (client *AWSClient) CancelCapacityReservation(crId string) (*ec2.CancelCapacityReservationOutput, error) {
	output, err := client.Ec2Client.CancelCapacityReservation(context.Background(), &ec2.CancelCapacityReservationInput{
		CapacityReservationId: aws.String(crId),
	})
	if err != nil {
		log.LogError("Delete capacity reservation failed with error: %s", err)
	}
	return output, err
}
