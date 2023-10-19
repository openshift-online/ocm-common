package validations

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/openshift-online/ocm-common/pkg"
)

func GetRoleName(prefix string, role string) string {
	name := fmt.Sprintf("%s-%s-Role", prefix, role)
	if len(name) > pkg.MaxByteSize {
		name = name[0:pkg.MaxByteSize]
	}
	return name
}

func IsManagedRole(roleTags []*iam.Tag) bool {
    for _, tag := range roleTags {
        if aws.StringValue(tag.Key) == ManagedPolicies && aws.StringValue(tag.Value) == "true" {
            return true
        }
    }

    return false
}

func IamResourceHasTag(iamTags []*iam.Tag, tagKey string, tagValue string) bool {
	for _, tag := range iamTags {
		if aws.StringValue(tag.Key) == tagKey && aws.StringValue(tag.Value) == tagValue {
			return true
		}
	}

	return false
}
