package vpc_client

import (
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/openshift-online/ocm-common/pkg/log"
)

func (vpc *vpc) CreateKeyPair(keyName string) (*ec2.CreateKeyPairOutput, error) {
	output, err := vpc.AWSClient.CreateKeyPair(keyName)
	if err != nil {
		return nil, err
	}
	log.LogInfo("create key pair: %v successfully\n", *output.KeyPairId)

	return output, nil
}

func (vpc *vpc) DeleteKeyPair(keyNames []string) error {
	for _, key := range keyNames {
		_, err := vpc.AWSClient.DeleteKeyPair(key)
		if err != nil {
			return err
		}
	}
	log.LogInfo("delete key pair successfully\n")
	return nil
}
