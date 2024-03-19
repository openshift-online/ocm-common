package main

import (
	"fmt"

	awsV2 "github.com/openshift-online/ocm-common/pkg/aws/client"
)

func main() {
	awsclient2, err := awsV2.CreateAWSV2Client("default", "us-east-1")
	if err != nil {
		fmt.Printf("Error while creating client: %v\n", err)
	}

	// output, err := awsclient2.DescribeLogGroupsByName("test")
	// if err != nil {
	// 	fmt.Printf("Error while Log Groups By Name: %v\n", err)
	// } else {
	// 	fmt.Printf("Log Groups: %v\n", output.LogGroups)
	// }

	// output4, err := awsclient2.GetInstancesByInfraID("test")
	// if err != nil {
	// 	fmt.Printf("Error while getting Instances by InfraID: %v\n", err)
	// } else {
	// 	fmt.Printf("Instances: %v\n", output4)
	// }
	// output5, err := awsclient2.ListAvaliableRegionsFromAWS()
	// if err != nil {
	// 	fmt.Printf("Error while ListAvaliableRegionsFromAWS: %v\n", err)
	// } else {
	// 	fmt.Printf("ListAvaliableRegionsFromAWS: %v\n", output5)
	// }
	// output6, err := awsclient2.PolicyAttachedToRole("test", "test2")
	// if err != nil {
	// 	fmt.Printf("Error while PolicyAttachedToRole: %v\n", err)
	// } else {
	// 	fmt.Printf("PolicyAttachedToRole: %v\n", output6)
	// }
	// output7, err := awsclient2.ListRoleAttachedPolicies("301721915996-admin")
	// if err != nil {
	// 	fmt.Printf("Error while ListRoleAttachedPolicies: %v\n", err)
	// } else {
	// 	fmt.Printf("ListRoleAttachedPolicies: %v\n", output7)
	// }
	// err2 := awsclient2.DetachRolePolicies("test")
	// if err != nil {
	// 	fmt.Printf("Error while DetachRolePolicies: %v\n", err2)
	// }

	// err3 := awsclient2.DeleteRoleInstanceProfiles("test")
	// if err != nil {
	// 	fmt.Printf("Error while Log Groups By Name: %v\n", err3)
	// }

	// output8, err := awsclient2.CreatePolicyForAuditLogForward("test")
	// if err != nil {
	// 	fmt.Printf("Error while CreatePolicyForAuditLogForward: %v\n", err)
	// } else {
	// 	fmt.Printf("CreatePolicyForAuditLogForward: %v\n", output8)
	// }

	output9, err := awsclient2.CreateRegularRole("ohada")
	if err != nil {
		fmt.Printf("Error while CreateRegularRole: %v\n", err)
	} else {
		fmt.Printf("CreateRegularRole: %v\n", output9)
	}

}
