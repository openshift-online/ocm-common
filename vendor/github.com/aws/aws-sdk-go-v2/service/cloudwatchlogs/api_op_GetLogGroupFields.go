// Code generated by smithy-go-codegen DO NOT EDIT.

package cloudwatchlogs

import (
	"context"
	"fmt"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// Returns a list of the fields that are included in log events in the specified
// log group. Includes the percentage of log events that contain each field. The
// search is limited to a time period that you specify. You can specify the log
// group to search by using either logGroupIdentifier or logGroupName . You must
// specify one of these parameters, but you can't specify both. In the results,
// fields that start with @ are fields generated by CloudWatch Logs. For example,
// @timestamp is the timestamp of each log event. For more information about the
// fields that are generated by CloudWatch logs, see Supported Logs and Discovered
// Fields (https://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/CWL_AnalyzeLogData-discoverable-fields.html)
// . The response results are sorted by the frequency percentage, starting with the
// highest percentage. If you are using CloudWatch cross-account observability, you
// can use this operation in a monitoring account and view data from the linked
// source accounts. For more information, see CloudWatch cross-account
// observability (https://docs.aws.amazon.com/AmazonCloudWatch/latest/monitoring/CloudWatch-Unified-Cross-Account.html)
// .
func (c *Client) GetLogGroupFields(ctx context.Context, params *GetLogGroupFieldsInput, optFns ...func(*Options)) (*GetLogGroupFieldsOutput, error) {
	if params == nil {
		params = &GetLogGroupFieldsInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "GetLogGroupFields", params, optFns, c.addOperationGetLogGroupFieldsMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*GetLogGroupFieldsOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type GetLogGroupFieldsInput struct {

	// Specify either the name or ARN of the log group to view. If the log group is in
	// a source account and you are using a monitoring account, you must specify the
	// ARN. You must include either logGroupIdentifier or logGroupName , but not both.
	LogGroupIdentifier *string

	// The name of the log group to search. You must include either logGroupIdentifier
	// or logGroupName , but not both.
	LogGroupName *string

	// The time to set as the center of the query. If you specify time , the 8 minutes
	// before and 8 minutes after this time are searched. If you omit time , the most
	// recent 15 minutes up to the current time are searched. The time value is
	// specified as epoch time, which is the number of seconds since January 1, 1970,
	// 00:00:00 UTC .
	Time *int64

	noSmithyDocumentSerde
}

type GetLogGroupFieldsOutput struct {

	// The array of fields found in the query. Each object in the array contains the
	// name of the field, along with the percentage of time it appeared in the log
	// events that were queried.
	LogGroupFields []types.LogGroupField

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationGetLogGroupFieldsMiddlewares(stack *middleware.Stack, options Options) (err error) {
	if err := stack.Serialize.Add(&setOperationInputMiddleware{}, middleware.After); err != nil {
		return err
	}
	err = stack.Serialize.Add(&awsAwsjson11_serializeOpGetLogGroupFields{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsAwsjson11_deserializeOpGetLogGroupFields{}, middleware.After)
	if err != nil {
		return err
	}
	if err := addProtocolFinalizerMiddlewares(stack, options, "GetLogGroupFields"); err != nil {
		return fmt.Errorf("add protocol finalizers: %v", err)
	}

	if err = addlegacyEndpointContextSetter(stack, options); err != nil {
		return err
	}
	if err = addSetLoggerMiddleware(stack, options); err != nil {
		return err
	}
	if err = addClientRequestID(stack); err != nil {
		return err
	}
	if err = addComputeContentLength(stack); err != nil {
		return err
	}
	if err = addResolveEndpointMiddleware(stack, options); err != nil {
		return err
	}
	if err = addComputePayloadSHA256(stack); err != nil {
		return err
	}
	if err = addRetry(stack, options); err != nil {
		return err
	}
	if err = addRawResponseToMetadata(stack); err != nil {
		return err
	}
	if err = addRecordResponseTiming(stack); err != nil {
		return err
	}
	if err = addClientUserAgent(stack, options); err != nil {
		return err
	}
	if err = smithyhttp.AddErrorCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = addSetLegacyContextSigningOptionsMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opGetLogGroupFields(options.Region), middleware.Before); err != nil {
		return err
	}
	if err = addRecursionDetection(stack); err != nil {
		return err
	}
	if err = addRequestIDRetrieverMiddleware(stack); err != nil {
		return err
	}
	if err = addResponseErrorMiddleware(stack); err != nil {
		return err
	}
	if err = addRequestResponseLogging(stack, options); err != nil {
		return err
	}
	if err = addDisableHTTPSMiddleware(stack, options); err != nil {
		return err
	}
	return nil
}

func newServiceMetadataMiddleware_opGetLogGroupFields(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		OperationName: "GetLogGroupFields",
	}
}
