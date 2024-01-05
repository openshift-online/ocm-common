package utils

import (
	"fmt"

	"github.com/openshift-online/ocm-common/pkg/ocm/consts"
)

func CreateVersionId(version string, channelGroup string) string {
	versionId := fmt.Sprintf("openshift-v%s", version)
	if channelGroup != consts.DefaultChannelGroup {
		versionId = fmt.Sprintf("%s-%s", versionId, channelGroup)
	}
	return versionId
}
