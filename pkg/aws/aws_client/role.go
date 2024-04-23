package aws_client

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
)

func (client *AWSClient) CreateRole(roleName string,
	assumeRolePolicyDocument string,
	permissionBoundry string,
	tags map[string]string,
	path string,
) (types.Role, error) {
	var roleTags []types.Tag
	for tagKey, tagValue := range tags {
		roleTags = append(roleTags, types.Tag{
			Key:   &tagKey,
			Value: &tagValue,
		})
	}
	description := "This is created role for ocm-qe automation testing"
	input := &iam.CreateRoleInput{
		RoleName:                 &roleName,
		AssumeRolePolicyDocument: &assumeRolePolicyDocument,
		Description:              &description,
	}
	if path != "" {
		input.Path = &path
	}
	if permissionBoundry != "" {
		input.PermissionsBoundary = &permissionBoundry
	}
	if len(tags) != 0 {
		input.Tags = roleTags
	}

	resp, err := client.IamClient.CreateRole(context.TODO(), input)
	if err != nil {
		return *resp.Role, err
	}
	err = client.WaitForResourceExisting("role-"+*resp.Role.RoleName, 10) // add a prefix to meet the resourceExisting split rule
	return *resp.Role, err
}

func (client *AWSClient) GetRole(roleName string) (*types.Role, error) {
	input := &iam.GetRoleInput{
		RoleName: &roleName,
	}
	out, err := client.IamClient.GetRole(context.TODO(), input)
	return out.Role, err
}
func (client *AWSClient) DeleteRole(roleName string) error {

	input := &iam.DeleteRoleInput{
		RoleName: &roleName,
	}
	_, err := client.IamClient.DeleteRole(context.TODO(), input)
	return err
}

func (client *AWSClient) DeleteRoleAndPolicy(roleName string, managedPolicy bool) error {
	input := &iam.ListAttachedRolePoliciesInput{
		RoleName: &roleName,
	}
	output, err := client.IamClient.ListAttachedRolePolicies(client.ClientContext, input)
	if err != nil {
		return err
	}

	fmt.Println(output.AttachedPolicies)
	for _, policy := range output.AttachedPolicies {
		err = client.DetachIAMPolicy(roleName, *policy.PolicyArn)
		if err != nil {
			return err
		}
		if !managedPolicy {
			err = client.DeletePolicy(*policy.PolicyArn)
			if err != nil {
				return err
			}
		}

	}
	err = client.DeleteRole(roleName)
	return err
}

func (client *AWSClient) ListRoles() ([]types.Role, error) {
	input := &iam.ListRolesInput{}
	out, err := client.IamClient.ListRoles(context.TODO(), input)
	return out.Roles, err
}

func (client *AWSClient) IsPolicyAttachedToRole(roleName string, policyArn string) (bool, error) {
	policies, err := client.ListAttachedRolePolicies(roleName)
	if err != nil {
		return false, err
	}
	for _, policy := range policies {
		if aws.ToString(policy.PolicyArn) == policyArn {
			return true, nil
		}
	}
	return false, nil
}

func (client *AWSClient) ListAttachedRolePolicies(roleName string) ([]types.AttachedPolicy, error) {
	policies := []types.AttachedPolicy{}
	policyLister := iam.ListAttachedRolePoliciesInput{
		RoleName: &roleName,
	}
	policyOut, err := client.IamClient.ListAttachedRolePolicies(context.TODO(), &policyLister)
	if err != nil {
		return policies, err
	}
	return policyOut.AttachedPolicies, nil
}

func (client *AWSClient) DetachRolePolicies(roleName string) error {
	policies, err := client.ListAttachedRolePolicies(roleName)
	if err != nil {
		return err
	}
	for _, policy := range policies {
		policyDetacher := iam.DetachRolePolicyInput{
			PolicyArn: policy.PolicyArn,
			RoleName:  &roleName,
		}
		_, err := client.IamClient.DetachRolePolicy(context.TODO(), &policyDetacher)
		if err != nil {
			return err
		}
	}
	return nil
}

func (client *AWSClient) DeleteRoleInstanceProfiles(roleName string) error {
	inProfileLister := iam.ListInstanceProfilesForRoleInput{
		RoleName: &roleName,
	}
	out, err := client.IamClient.ListInstanceProfilesForRole(context.TODO(), &inProfileLister)
	if err != nil {
		return err
	}
	for _, inProfile := range out.InstanceProfiles {
		profileDeleter := iam.RemoveRoleFromInstanceProfileInput{
			InstanceProfileName: inProfile.InstanceProfileName,
			RoleName:            &roleName,
		}
		_, err = client.IamClient.RemoveRoleFromInstanceProfile(context.TODO(), &profileDeleter)
		if err != nil {
			return err
		}
	}

	return nil
}

func (client *AWSClient) CreateIAMRole(roleName string, ProdENVTrustedRole string, StageENVTrustedRole string, StageIssuerTrustedRole string,
	externalID ...string) (types.Role, error) {
	statement := map[string]interface{}{
		"Effect": "Allow",
		"Principal": map[string]interface{}{
			"Service": "ec2.amazonaws.com",
			"AWS": []string{
				ProdENVTrustedRole,
				StageENVTrustedRole,
				StageIssuerTrustedRole,
			},
		},
		"Action": "sts:AssumeRole",
	}

	if len(externalID) == 1 {
		statement["Condition"] = map[string]map[string]string{
			"StringEquals": {
				"sts:ExternalId": "aaaa",
			},
		}
	}

	assumeRolePolicyDocument, err := completeRolePolicyDocument(statement)
	if err != nil {
		fmt.Println("Failed to convert Role Policy Document into JSON: ", err)
		return types.Role{}, err
	}

	return client.CreateRole(roleName, string(assumeRolePolicyDocument), "", make(map[string]string), "/")
}

func (client *AWSClient) CreateRegularRole(roleName string) (types.Role, error) {

	statement := map[string]interface{}{
		"Effect": "Allow",
		"Principal": map[string]interface{}{
			"Service": "ec2.amazonaws.com",
		},
		"Action": "sts:AssumeRole",
	}

	assumeRolePolicyDocument, err := completeRolePolicyDocument(statement)
	if err != nil {
		fmt.Println("Failed to convert Role Policy Document into JSON: ", err)
		return types.Role{}, err
	}
	return client.CreateRole(roleName, assumeRolePolicyDocument, "", make(map[string]string), "/")
}

func (client *AWSClient) CreateRoleForAuditLogForward(roleName, awsAccountID string, oidcEndpointURL string) (types.Role, error) {
	statement := map[string]interface{}{
		"Effect": "Allow",
		"Principal": map[string]interface{}{
			"Federated": fmt.Sprintf("arn:aws:iam::%s:oidc-provider/%s", awsAccountID, oidcEndpointURL),
		},
		"Action": "sts:AssumeRoleWithWebIdentity",
		"Condition": map[string]interface{}{
			"StringEquals": map[string]interface{}{
				fmt.Sprintf("%s:sub", oidcEndpointURL): "system:serviceaccount:openshift-config-managed:cloudwatch-audit-exporter",
			},
		},
	}

	assumeRolePolicyDocument, err := completeRolePolicyDocument(statement)
	if err != nil {
		fmt.Println("Failed to convert Role Policy Document into JSON: ", err)
		return types.Role{}, err
	}

	return client.CreateRole(roleName, string(assumeRolePolicyDocument), "", make(map[string]string), "/")
}

func (client *AWSClient) CreatePolicy(policyName string, statements ...map[string]interface{}) (string, error) {
	timeCreation := time.Now().Local().String()
	description := fmt.Sprintf("Created by OCM QE at %s", timeCreation)
	document := map[string]interface{}{
		"Version":   "2012-10-17",
		"Statement": []map[string]interface{}{},
	}
	if len(statements) != 0 {
		for _, statement := range statements {
			document["Statement"] = append(document["Statement"].([]map[string]interface{}), statement)
		}
	}
	documentBytes, err := json.Marshal(document)
	if err != nil {
		err = fmt.Errorf("error to unmarshal the statement to string: %v", err)
		return "", err
	}
	documentStr := string(documentBytes)
	policyCreator := iam.CreatePolicyInput{
		PolicyDocument: &documentStr,
		PolicyName:     &policyName,
		Description:    &description,
	}
	outRes, err := client.IamClient.CreatePolicy(context.TODO(), &policyCreator)
	if err != nil {
		return "", err
	}
	policyArn := *outRes.Policy.Arn
	return policyArn, err
}

func (client *AWSClient) CreatePolicyForAuditLogForward(policyName string) (string, error) {

	statement := map[string]interface{}{
		"Effect":   "Allow",
		"Resource": "arn:aws:logs:*:*:*",
		"Action": []string{
			"logs:PutLogEvents",
			"logs:CreateLogGroup",
			"logs:PutRetentionPolicy",
			"logs:CreateLogStream",
			"logs:DescribeLogGroups",
			"logs:DescribeLogStreams",
		},
	}
	return client.CreatePolicy(policyName, statement)
}

func completeRolePolicyDocument(statement map[string]interface{}) (string, error) {
	rolePolicyDocument := map[string]interface{}{
		"Version":   "2012-10-17",
		"Statement": statement,
	}

	assumeRolePolicyDocument, err := json.Marshal(rolePolicyDocument)
	return string(assumeRolePolicyDocument), err
}
