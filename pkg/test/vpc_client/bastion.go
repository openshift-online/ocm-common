package vpc_client

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"net"

	CON "github.com/openshift-online/ocm-common/pkg/aws/consts"
	"github.com/openshift-online/ocm-common/pkg/log"
)

// LaunchBastion will launch a bastion instance on the indicated zone.
// If set imageID to empty, it will find the bastion image using filter with specific name.
func (vpc *VPC) LaunchBastion(imageID string, zone string, userData string) (*types.Instance, error) {
	var inst *types.Instance
	if imageID == "" {

		var err error
		imageID, err = vpc.FindProxyLaunchImage()
		if err != nil {
			log.LogError("Cannot find bastion image of region %s in map bastionImageMap, please indicate it as parameter", vpc.Region)
			return nil, err
		}
	}
	if userData == "" {
		log.LogError("Userdata can not be empty, pleas provide the correct userdata")
		return nil, errors.New("userData should not be empty")
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
		[]string{SGID}, true, userData)

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

func (vpc *VPC) PrepareBastionProxy(zone string, cidrBlock string) (*types.Instance, error) {
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
		if cidrBlock == "" {
			cidrBlock = CON.RouteDestinationCidrBlock
		}
		_, _, err = net.ParseCIDR(cidrBlock)
		if err != nil {
			log.LogError("CIDR IP address format is invalid")
			return nil, err
		}

		userData := fmt.Sprintf(`#!/bin/bash
		yum update -y
		yum install -y squid
		cd /etc/squid/
		sudo mv ./squid.conf ./squid.conf.bak
		sudo touch squid.conf
		echo http_port 3128 >> /etc/squid/squid.conf
		echo acl allowed_ips src %s >> /etc/squid/squid.conf
		echo http_access allow allowed_ips >> /etc/squid/squid.conf
		echo http_access deny all >> /etc/squid/squid.conf
		systemctl start squid
		systemctl enable squid`, cidrBlock)

		encodeUserData := base64.StdEncoding.EncodeToString([]byte(userData))
		return vpc.LaunchBastion("", zone, encodeUserData)

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
	return nil
}
