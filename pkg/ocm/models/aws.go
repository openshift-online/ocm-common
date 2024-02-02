package models

var Mode string
var Modes = []string{ModeAuto, ModeManual}

const (
	ModeAuto   = "auto"
	ModeManual = "manual"

	AdminUserName        = "osdCcsAdmin"
	OsdCcsAdminStackName = "osdCcsAdminIAMUser"

	// Since CloudFormation stacks are region-dependent, we hard-code OCM's default region and
	// then use it to ensure that the user always gets the stack from the same region.
	DefaultRegion = "us-east-1"
	Inline        = "inline"
	Attached      = "attached"

	GovPartition = "aws-us-gov"

	AwsMaxFilterLength = 200
)
