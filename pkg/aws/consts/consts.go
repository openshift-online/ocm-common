package consts

const (
	MaxAwsRoleLength = 64
)

const (
	StageENVTrustedRole    = "arn:aws:iam::644306948063:role/RH-Managed-OpenShift-Installer"
	StageIssuerTrustedRole = "arn:aws:sts::644306948063:assumed-role/RH-Managed-OpenShift-Installer/OCM"
	ProdENVTrustedRole     = "arn:aws:iam::710019948333:role/RH-Managed-OpenShift-Installer"
)
