package utils

import (
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws/arn"
)

func GetPathFromArn(arnStr string) (string, error) {
	parse, err := arn.Parse(arnStr)
	if err != nil {
		return "", err
	}
	resource := parse.Resource
	firstIndex := strings.Index(resource, "/")
	lastIndex := strings.LastIndex(resource, "/")
	if firstIndex == lastIndex {
		return "", nil
	}
	path := resource[firstIndex : lastIndex+1]
	return path, nil
}
