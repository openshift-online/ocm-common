package aws_v2

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/openshift-online/ocm-common/pkg/aws/log"
)

func (client *AwsV2Client) DescribeLogGroupsByName(logGroupName string) (cloudwatchlogs.DescribeLogGroupsOutput, error) {
	output, err := client.cloudWatchLogsClient.DescribeLogGroups(context.TODO(), &cloudwatchlogs.DescribeLogGroupsInput{
		LogGroupNamePrefix: &logGroupName,
	})
	if err != nil {
		log.LogError("Got error describe log group:%s ", err)
	}
	return *output, err
}

func (client *AwsV2Client) DescribeLogStreamByName(logGroupName string) (cloudwatchlogs.DescribeLogStreamsOutput, error) {
	output, err := client.cloudWatchLogsClient.DescribeLogStreams(context.TODO(), &cloudwatchlogs.DescribeLogStreamsInput{
		LogGroupName: &logGroupName,
	})
	if err != nil {
		log.LogError("Got error describe log stream: %s", err)
	}
	return *output, err
}

func (client *AwsV2Client) DeleteLogGroupByName(logGroupName string) (cloudwatchlogs.DeleteLogGroupOutput, error) {
	output, err := client.cloudWatchLogsClient.DeleteLogGroup(context.TODO(), &cloudwatchlogs.DeleteLogGroupInput{
		LogGroupName: &logGroupName,
	})
	if err != nil {
		log.LogError("Got error delete log group: %s", err)
	}
	return *output, err
}
