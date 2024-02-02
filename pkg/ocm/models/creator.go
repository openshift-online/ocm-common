package models

type Creator struct {
	ARN             string
	AccountID       string
	AccessKeyID     string
	SecretAccessKey string
	IsSTS           bool
	IsGovcloud      bool
}
