package aws_client

import (
	"context"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	elbtypes "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing/types"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	iamtypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/sts"

	"github.com/openshift-online/ocm-common/pkg/log"

	elb "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing"
	"github.com/aws/aws-sdk-go-v2/service/route53"

	CON "github.com/openshift-online/ocm-common/pkg/aws/consts"
)

type AWSClient interface {
	GetAWSAccountID() string
	GetRegion() string

	EC2() *ec2.Client
	Route53() *route53.Client
	CloudFormation() *cloudformation.Client
	ELB() *elb.Client

	RetrieveAccessKey() (*AccessKeyMod, error)

	DescribeLogGroupsByName(logGroupName string) (*cloudwatchlogs.DescribeLogGroupsOutput, error)
	DescribeLogStreamByName(logGroupName string) (*cloudwatchlogs.DescribeLogStreamsOutput, error)
	DeleteLogGroupByName(logGroupName string) (*cloudwatchlogs.DeleteLogGroupOutput, error)

	AllocateEIPAddress() (*ec2.AllocateAddressOutput, error)
	DisassociateAddress(associateID string) (*ec2.DisassociateAddressOutput, error)
	AllocateEIPAndAssociateInstance(instanceID string) (string, error)
	ReleaseAddress(allocationID string) (*ec2.ReleaseAddressOutput, error)

	DescribeLoadBalancers(vpcID string) ([]elbtypes.LoadBalancerDescription, error)
	DeleteELB(ELB elbtypes.LoadBalancerDescription) error

	CopyImage(sourceImageID string, sourceRegion string, name string) (string, error)
	DescribeImage(imageIDs []string, filters ...map[string][]string) (*ec2.DescribeImagesOutput, error)

	LaunchInstance(subnetID string, imageID string, count int, instanceType string, keyName string, securityGroupIds []string, wait bool) (*ec2.RunInstancesOutput, error)
	ListInstances(instanceIDs []string, filters ...map[string][]string) ([]types.Instance, error)
	WaitForInstanceReady(instanceID string, timeout time.Duration) error
	CheckInstanceState(instanceIDs ...string) (*ec2.DescribeInstanceStatusOutput, error)
	WaitForInstancesRunning(instanceIDs []string, timeout time.Duration) (allRunning bool, err error)
	WaitForInstancesTerminated(instanceIDs []string, timeout time.Duration) (allTerminated bool, err error)
	ListAvaliableInstanceTypesForRegion(region string, availabilityZones ...string) ([]string, error)
	ListAvaliableZonesForRegion(region string, zoneType string) ([]string, error)
	TerminateInstances(instanceIDs []string, wait bool, timeout time.Duration) error
	WaitForInstanceTerminated(instanceIDs []string, timeout time.Duration) error
	GetTagsOfInstanceProfile(instanceProfileName string) ([]iamtypes.Tag, error)
	GetInstancesByInfraID(infraID string) ([]types.Instance, error)
	ListAvaliableRegionsFromAWS() ([]types.Region, error)

	CreateInternetGateway() (*ec2.CreateInternetGatewayOutput, error)
	AttachInternetGateway(internetGatewayID string, vpcID string) (*ec2.AttachInternetGatewayOutput, error)
	DetachInternetGateway(internetGatewayID string, vpcID string) (*ec2.DetachInternetGatewayOutput, error)
	ListInternetGateWay(vpcID string) ([]types.InternetGateway, error)
	DeleteInternetGateway(internetGatewayID string) (*ec2.DeleteInternetGatewayOutput, error)

	CreateKeyPair(keyName string) (*ec2.CreateKeyPairOutput, error)
	DeleteKeyPair(keyName string) (*ec2.DeleteKeyPairOutput, error)

	CreateKMSKeys(tagKey string, tagValue string, description string, policy string, multiRegion bool) (keyID string, keyArn string, err error)
	DescribeKMSKeys(keyID string) (kms.DescribeKeyOutput, error)
	ScheduleKeyDeletion(kmsKeyId string, pendingWindowInDays int32) (*kms.ScheduleKeyDeletionOutput, error)
	GetKMSPolicy(keyID string, policyName string) (kms.GetKeyPolicyOutput, error)
	PutKMSPolicy(keyID string, policyName string, policy string) (kms.PutKeyPolicyOutput, error)
	TagKeys(kmsKeyId string, tagKey string, tagValue string) (*kms.TagResourceOutput, error)

	CreateNatGateway(subnetID string, allocationID string, vpcID string) (*ec2.CreateNatGatewayOutput, error)
	DeleteNatGateway(natGatewayID string, timeout ...int) (*ec2.DeleteNatGatewayOutput, error)
	ListNatGateWays(vpcID string) ([]types.NatGateway, error)

	ListNetWorkAcls(vpcID string) ([]types.NetworkAcl, error)
	AddNetworkAclEntry(networkAclId string, egress bool, protocol string, ruleAction string, ruleNumber int32, fromPort int32, toPort int32, cidrBlock string) (*ec2.CreateNetworkAclEntryOutput, error)
	DeleteNetworkAclEntry(networkAclId string, egress bool, ruleNumber int32) (*ec2.DeleteNetworkAclEntryOutput, error)

	DescribeNetWorkInterface(vpcID string) ([]types.NetworkInterface, error)
	DeleteNetworkInterface(networkinterface types.NetworkInterface) error

	DeleteOIDCProvider(providerArn string) error

	CreateIAMPolicy(policyName string, policyDocument string, tags map[string]string) (*iamtypes.Policy, error)
	GetIAMPolicy(policyArn string) (*iamtypes.Policy, error)
	DeleteIAMPolicy(arn string) error
	AttachIAMPolicy(roleName string, policyArn string) error
	DetachIAMPolicy(roleAName string, policyArn string) error
	GetCustomerIAMPolicies() ([]iamtypes.Policy, error)
	FilterNeedCleanPolicies(cleanRule func(iamtypes.Policy) bool) ([]iamtypes.Policy, error)
	DeletePolicy(arn string) error
	DeletePolicyVersions(policyArn string) error
	CleanPolicies(cleanRule func(iamtypes.Policy) bool) error

	ResourceExisting(resourceID string) bool
	ResourceDeleted(resourceID string) bool
	WaitForResourceExisting(resourceID string, timeout int) error
	WaitForResourceDeleted(resourceID string, timeout int) error

	CreateRole(roleName string, assumeRolePolicyDocument string, permissionBoundry string, tags map[string]string, path string) (iamtypes.Role, error)
	GetRole(roleName string) (*iamtypes.Role, error)
	DeleteRole(roleName string) error
	DeleteRoleAndPolicy(roleName string, managedPolicy bool) error
	ListRoles() ([]iamtypes.Role, error)
	IsPolicyAttachedToRole(roleName string, policyArn string) (bool, error)
	ListAttachedRolePolicies(roleName string) ([]iamtypes.AttachedPolicy, error)
	DetachRolePolicies(roleName string) error
	DeleteRoleInstanceProfiles(roleName string) error
	CreateIAMRole(roleName string, ProdENVTrustedRole string, StageENVTrustedRole string, StageIssuerTrustedRole string, externalID ...string) (iamtypes.Role, error)
	CreateRegularRole(roleName string) (iamtypes.Role, error)
	CreateRoleForAuditLogForward(roleName, awsAccountID string, oidcEndpointURL string) (iamtypes.Role, error)
	CreatePolicy(policyName string, statements ...map[string]interface{}) (string, error)
	CreatePolicyForAuditLogForward(policyName string) (string, error)

	CreateRouteTable(vpcID string) (*ec2.CreateRouteTableOutput, error)
	AssociateRouteTable(routeTableID string, subnetID string, vpcID string) (*ec2.AssociateRouteTableOutput, error)
	ListCustomerRouteTables(vpcID string) ([]types.RouteTable, error)
	ListRTAssociations(routeTableID string) ([]string, error)
	DisassociateRouteTableAssociation(associationID string) (*ec2.DisassociateRouteTableOutput, error)
	DisassociateRouteTableAssociations(routeTableID string) error
	CreateRoute(routeTableID string, targetID string) (*types.Route, error)
	DeleteRouteTable(routeTableID string) error
	DeleteRouteTableChain(routeTableID string) error

	CreateHostedZone(hostedZoneName string, vpcID string, private bool) (*route53.CreateHostedZoneOutput, error)
	GetHostedZone(hostedZoneID string) (*route53.GetHostedZoneOutput, error)
	ListHostedZoneByDNSName(hostedZoneName string) (*route53.ListHostedZonesByNameOutput, error)

	ListSecurityGroups(vpcID string) ([]types.SecurityGroup, error)
	ReleaseInboundOutboundRules(sgID string) error
	DeleteSecurityGroup(groupID string) (*ec2.DeleteSecurityGroupOutput, error)
	AuthorizeSecurityGroupIngress(groupID string, cidr string, protocol string, fromPort int32, toPort int32) (*ec2.AuthorizeSecurityGroupIngressOutput, error)
	CreateSecurityGroup(vpcID string, groupName string, sgDescription string) (*ec2.CreateSecurityGroupOutput, error)
	GetSecurityGroupWithID(sgID string) (*ec2.DescribeSecurityGroupsOutput, error)

	CreateSubnet(vpcID string, zone string, subnetCidr string) (*types.Subnet, error)
	ListSubnetByVpcID(vpcID string) ([]types.Subnet, error)
	DeleteSubnet(subnetID string) (*ec2.DeleteSubnetOutput, error)
	ListSubnetDetail(subnetIDs ...string) ([]types.Subnet, error)
	ListSubnetsByFilter(filter []types.Filter) ([]types.Subnet, error)

	TagResource(resourceID string, tags map[string]string) (*ec2.CreateTagsOutput, error)
	RemoveResourceTag(resourceID string, tagKey string, tagValue string) (*ec2.DeleteTagsOutput, error)

	DescribeVolumeByID(volumeID string) (*ec2.DescribeVolumesOutput, error)

	ListVPCByName(vpcName string) ([]types.Vpc, error)
	CreateVpc(cidr string, name ...string) (*ec2.CreateVpcOutput, error)
	ModifyVpcDnsAttribute(vpcID string, dnsAttribute string, status bool) (*ec2.ModifyVpcAttributeOutput, error)
	DeleteVpc(vpcID string) (*ec2.DeleteVpcOutput, error)
	DescribeVPC(vpcID string) (types.Vpc, error)
	ListEndpointAssociation(vpcID string) ([]types.VpcEndpoint, error)
	DeleteVPCEndpoints(vpcID string) error
}

type awsClient struct {
	Ec2Client            *ec2.Client
	Route53Client        *route53.Client
	StackFormationClient *cloudformation.Client
	ElbClient            *elb.Client
	StsClient            *sts.Client
	Region               string
	IamClient            *iam.Client
	ClientContext        context.Context
	KmsClient            *kms.Client
	CloudWatchLogsClient *cloudwatchlogs.Client
	AWSConfig            *aws.Config
}

type AccessKeyMod struct {
	AccessKeyId     string `ini:"aws_access_key_id,omitempty"`
	SecretAccessKey string `ini:"aws_secret_access_key,omitempty"`
}

func CreateAWSClient(profileName string, region string) (AWSClient, error) {
	var cfg aws.Config
	var err error

	if envCredential() {
		log.LogInfo("Got AWS_ACCESS_KEY_ID env settings, going to build the config with the env")
		cfg, err = config.LoadDefaultConfig(context.TODO(),
			config.WithRegion(region),
			config.WithCredentialsProvider(
				credentials.NewStaticCredentialsProvider(
					os.Getenv("AWS_ACCESS_KEY_ID"),
					os.Getenv("AWS_SECRET_ACCESS_KEY"),
					"")),
		)
	} else {
		if envAwsProfile() {
			file := os.Getenv("AWS_SHARED_CREDENTIALS_FILE")
			log.LogInfo("Got file path: %s from env variable AWS_SHARED_CREDENTIALS_FILE\n", file)
			cfg, err = config.LoadDefaultConfig(context.TODO(),
				config.WithRegion(region),
				config.WithSharedCredentialsFiles([]string{file}),
			)
		} else {
			cfg, err = config.LoadDefaultConfig(context.TODO(),
				config.WithRegion(region),
				config.WithSharedConfigProfile(profileName),
			)
		}

	}

	if err != nil {
		return nil, err
	}

	awsClient := &awsClient{
		Ec2Client:            ec2.NewFromConfig(cfg),
		Route53Client:        route53.NewFromConfig(cfg),
		StackFormationClient: cloudformation.NewFromConfig(cfg),
		ElbClient:            elb.NewFromConfig(cfg),
		Region:               region,
		StsClient:            sts.NewFromConfig(cfg),
		IamClient:            iam.NewFromConfig(cfg),
		ClientContext:        context.TODO(),
		KmsClient:            kms.NewFromConfig(cfg),
		AWSConfig:            &cfg,
	}
	return awsClient, nil
}
func (client *awsClient) GetRegion() string {
	return client.Region
}

func (client *awsClient) GetAWSAccountID() string {
	input := &sts.GetCallerIdentityInput{}
	out, err := client.StsClient.GetCallerIdentity(client.ClientContext, input)
	if err != nil {
		return ""
	}
	return *out.Account
}

func (client *awsClient) EC2() *ec2.Client {
	return client.Ec2Client
}

func (client *awsClient) Route53() *route53.Client {
	return client.Route53Client
}
func (client *awsClient) CloudFormation() *cloudformation.Client {
	return client.StackFormationClient
}
func (client *awsClient) ELB() *elb.Client {
	return client.ElbClient
}

func GrantValidAccessKeys(userName string) (*AccessKeyMod, error) {
	var keysMod *AccessKeyMod
	var err error
	retryTimes := 3
	for retryTimes > 0 {
		if keysMod.AccessKeyId != "" {
			break
		}
		client, err := CreateAWSClient(userName, CON.DefaultAWSRegion)
		if err != nil {
			return nil, err
		}

		keysMod, err = client.RetrieveAccessKey()
		if err != nil {
			return nil, err
		}

		retryTimes--
	}
	return keysMod, err
}

func (client *awsClient) RetrieveAccessKey() (*AccessKeyMod, error) {
	cre, err := client.AWSConfig.Credentials.Retrieve(client.ClientContext)
	if err != nil {
		return nil, err
	}
	log.LogInfo(">>> Access key grant successfully")

	keysMod := &AccessKeyMod{
		AccessKeyId:     cre.AccessKeyID,
		SecretAccessKey: cre.SecretAccessKey,
	}
	return keysMod, nil
}
