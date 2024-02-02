package models

import (
	"net"
	"time"

	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
)

var NetworkTypes = []string{"OpenShiftSDN", "OVNKubernetes"}

// ClusterSpec is the configuration for a cluster spec.
type ClusterSpec struct {
	// Basic configs
	Name                      string
	Region                    string
	MultiAZ                   bool
	Version                   string
	ChannelGroup              string
	Expiration                time.Time
	Flavour                   string
	DisableWorkloadMonitoring *bool

	//Encryption
	FIPS                 bool
	EtcdEncryption       bool
	KMSKeyArn            string
	EtcdEncryptionKMSArn string
	// Scaling config
	ComputeMachineType string
	ComputeNodes       int
	Autoscaling        bool
	AutoscalerConfig   *AutoscalerConfig
	MinReplicas        int
	MaxReplicas        int
	ComputeLabels      map[string]string

	// SubnetIDs
	SubnetIds []string

	// AvailabilityZones
	AvailabilityZones []string

	// Network config
	NetworkType string
	MachineCIDR net.IPNet
	ServiceCIDR net.IPNet
	PodCIDR     net.IPNet
	HostPrefix  int
	Private     *bool
	PrivateLink *bool

	// Properties
	CustomProperties map[string]string

	// User-defined tags for AWS resources
	Tags map[string]string

	// Simulate creating a cluster but don't actually create it
	DryRun *bool

	// Disable SCP checks in the installer by setting credentials mode as mint
	DisableSCPChecks *bool

	// STS
	IsSTS               bool
	RoleARN             string
	ExternalID          string
	SupportRoleARN      string
	OperatorIAMRoles    []OperatorIAMRole
	ControlPlaneRoleARN string
	WorkerRoleARN       string
	OidcConfigId        string
	Mode                string

	Creator Creator

	NodeDrainGracePeriodInMinutes float64

	EnableProxy               bool
	HTTPProxy                 *string
	HTTPSProxy                *string
	NoProxy                   *string
	AdditionalTrustBundleFile *string
	AdditionalTrustBundle     *string

	// HyperShift options:
	Hypershift     Hypershift
	BillingAccount string
	NoCni          bool

	// Audit Log Forwarding
	AuditLogRoleARN *string

	Ec2MetadataHttpTokens cmv1.Ec2MetadataHttpTokens

	// Cluster Admin
	ClusterAdminUser     string
	ClusterAdminPassword string

	// Default Ingress Attributes
	DefaultIngress DefaultIngressSpec

	// Machine pool's storage
	MachinePoolRootDisk *Volume

	// Shared VPC
	PrivateHostedZoneID string
	SharedVPCRoleArn    string
	BaseDomain          string

	// Worker Machine Pool attributes
	AdditionalComputeSecurityGroupIds []string

	// Infra Machine Pool attributes
	AdditionalInfraSecurityGroupIds []string

	// Control Plane Machine Pool attributes
	AdditionalControlPlaneSecurityGroupIds []string
}

// Volume represents a volume property for a disk
type Volume struct {
	Size int
}

type OperatorIAMRole struct {
	Name      string
	Namespace string
	RoleARN   string
	Path      string
}

type Hypershift struct {
	Enabled bool
}
