package vpc_client

import (
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/openshift-online/ocm-common/module/aws/consts"
	"github.com/openshift-online/ocm-common/module/aws/log"
)

func (vpc *VPC) DeleteVPCSecurityGroups(customizedOnly bool) error {
	needCleanGroups := []types.SecurityGroup{}
	securityGroups, err := vpc.AWSClient.ListSecurityGroups(vpc.VpcID)
	if customizedOnly {
		for _, sg := range securityGroups {
			for _, tag := range sg.Tags {
				if *tag.Key == "Name" && (*tag.Value == consts.ProxySecurityGroupName || *tag.Value == consts.AdditionalSecurityGroupName) {
					needCleanGroups = append(needCleanGroups, sg)
				}
			}
		}
	} else {
		needCleanGroups = securityGroups
	}
	if err != nil {
		return err
	}
	for _, sg := range needCleanGroups {
		_, err = vpc.AWSClient.DeleteSecurityGroup(*sg.GroupId)
		if err != nil {
			return err
		}
	}
	return nil
}

func (vpc *VPC) CreateAndAuthorizeDefaultSecurityGroupForProxy() (string, error) {
	var groupID string
	var err error
	protocol := consts.TCPProtocol
	resp, err := vpc.AWSClient.CreateSecurityGroup(vpc.VpcID, consts.ProxySecurityGroupName, consts.ProxySecurityGroupDescription)
	if err != nil {
		log.LogError("Create proxy security group failed for vpc %s: %s", vpc.VpcID, err)
		return "", err
	}
	groupID = *resp.GroupId
	log.LogInfo("SG %s created for vpc %s", groupID, vpc.VpcID)
	cidrPortsMap := map[string]int32{
		vpc.CIDRValue: 8080,
		"0.0.0.0/0":   22,
	}
	for cidr, port := range cidrPortsMap {
		_, err = vpc.AWSClient.AuthorizeSecurityGroupIngress(groupID, cidr, protocol, port, port)
		if err != nil {
			log.LogError("Authorize CIDR %s with port %v failed to SG %s of vpc %s: %s",
				cidr, port, groupID, vpc.VpcID, err)
			return groupID, err
		}
	}
	log.LogInfo("Authorize SG %s successfully for proxy.", groupID)

	return groupID, err
}
