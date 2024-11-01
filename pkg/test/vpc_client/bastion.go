package vpc_client

import (
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"net"

	CON "github.com/openshift-online/ocm-common/pkg/aws/consts"
	"github.com/openshift-online/ocm-common/pkg/log"
)

// LaunchBastion will launch a bastion instance on the indicated zone.
// If set imageID to empty, it will find the bastion image using filter with specific name.
func (vpc *VPC) LaunchBastion(imageID string, zone string, userData ...string) (*types.Instance, error) {
	var inst *types.Instance
	if imageID == "" {

		var err error
		imageID, err = vpc.FindProxyLaunchImage()
		if err != nil {
			log.LogError("Cannot find bastion image of region %s in map bastionImageMap, please indicate it as parameter", vpc.Region)
			return nil, err
		}
	}
	pubSubnet, err := vpc.PreparePublicSubnet(zone)
	if err != nil {
		log.LogError("Error preparing a subnet in current zone %s with image ID %s: %s", zone, imageID, err)
		return nil, err
	}
	SGID, err := vpc.CreateAndAuthorizeDefaultSecurityGroupForProxy(3128)
	if err != nil {
		log.LogError("Prepare SG failed for the bastion preparation %s", err)
		return inst, err
	}

	key, err := vpc.CreateKeyPair(fmt.Sprintf("%s-bastion", CON.InstanceKeyNamePrefix))
	if err != nil {
		log.LogError("Create key pair failed %s", err)
		return inst, err
	}
	instOut, err := vpc.AWSClient.LaunchInstance(pubSubnet.ID, imageID, 1, "t3.medium", *key.KeyName,
		[]string{SGID}, true, userData...)

	if err != nil {
		log.LogError("Launch bastion instance failed %s", err)
		return inst, err
	} else {
		log.LogInfo("Launch bastion instance %s succeed", *instOut.Instances[0].InstanceId)
	}
	tags := map[string]string{
		"Name": CON.BastionName,
	}
	instID := *instOut.Instances[0].InstanceId
	_, err = vpc.AWSClient.TagResource(instID, tags)
	if err != nil {
		return inst, fmt.Errorf("tag instance %s failed:%s", instID, err)
	}

	publicIP, err := vpc.AWSClient.AllocateEIPAndAssociateInstance(instID)
	if err != nil {
		log.LogError("Prepare EIP failed for the bastion preparation %s", err)
		return inst, err
	}
	log.LogInfo("Prepare EIP successfully for the bastion preparation. Launch with IP: %s", publicIP)
	inst = &instOut.Instances[0]
	inst.PublicIpAddress = &publicIP
	return inst, nil
}

func (vpc *VPC) PrepareBastionProxy(zone string, cidrIp string) (*types.Instance, error) {
	filters := []map[string][]string{
		{
			"vpc-id": {
				vpc.VpcID,
			},
		},
		{
			"tag:Name": {
				CON.BastionName,
			},
		},
	}

	insts, err := vpc.AWSClient.ListInstances([]string{}, filters...)
	if err != nil {
		return nil, err
	}
	if len(insts) == 0 {
		log.LogInfo("Didn't found an existing bastion, going to launch one")
		if cidrIp == "" {
			cidrIp = CON.RouteDestinationCidrBlock
		}
		_, _, err = net.ParseCIDR(cidrIp)
		if err != nil {
			log.LogError("CIDR IP address format is invalid")
			return nil, err
		}
		userData := fmt.Sprintf("#!/bin/bash\n\t\t"+
			"yum update -y\n\t\t"+
			"yum install -y squid\n\t\t"+
			"cd /etc/squid/\n\t\t"+
			"sudo mv ./squid.conf ./squid.conf.bak\n\t\t"+
			"sudo touch squid.conf\n\t\t"+
			"echo \"http_port 3128\" >> /etc/squid/squid.conf\n\t\t"+
			"echo \"acl allowed_ips src %s\" >> /etc/squid/squid.conf\n\t\t"+
			"echo \"http_access allow allowed_ips\" >> /etc/squid/squid.conf\n\t\t"+
			"echo \"http_access deny all\" >> /etc/squid/squid.conf\n\t\t"+
			"systemctl start squid\n\t\t"+
			"systemctl enable squid", cidrIp)

		encodeUserData := base64.StdEncoding.EncodeToString([]byte(userData))
		regionZone := fmt.Sprintf("%s%s", vpc.Region, zone)
		return vpc.LaunchBastion("", regionZone, encodeUserData)

	}
	log.LogInfo("Found existing bastion: %s", *insts[0].InstanceId)
	return &insts[0], nil
}

func (vpc *VPC) DestroyBastionProxy(instance types.Instance) error {
	var instanceIDs []string
	instanceIDs = append(instanceIDs, *instance.InstanceId)
	err := vpc.AWSClient.TerminateInstances(instanceIDs, true, 10)
	if err != nil {
		log.LogError("Terminate instance failed")
		return err
	}

	var keyNames []string
	keyNames = append(keyNames, *instance.KeyName)
	err = vpc.DeleteKeyPair(keyNames)
	if err != nil {
		log.LogError("Delete key pair failed")
		return err
	}

	err = vpc.DeleteVPCSecurityGroups(true)
	if err != nil {
		log.LogError("Delete VPC security group failed")
		return err
	}

	_, err = vpc.AWSClient.DeleteSubnet(*instance.SubnetId)
	if err != nil {
		log.LogError("Delete VPC public subnet failed")
		return err
	}

	return nil
}
