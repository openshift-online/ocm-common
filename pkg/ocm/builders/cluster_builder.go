package builders

import (
	"fmt"
	"reflect"

	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"

	idpUtils "github.com/openshift-online/ocm-common/pkg/idp/utils"
	"github.com/openshift-online/ocm-common/pkg/ocm/consts"
	"github.com/openshift-online/ocm-common/pkg/ocm/models"
	"github.com/openshift-online/ocm-common/pkg/ocm/utils"
)

func CreateClusterSpec(config models.ClusterSpec) (*cmv1.Cluster, error) {
	clusterProperties := map[string]string{}

	if config.CustomProperties != nil {
		for key, value := range config.CustomProperties {
			clusterProperties[key] = value
		}
	}

	// Make sure we don't have a custom properties collision
	if _, present := clusterProperties[consts.CreatorArn]; present {
		return nil, fmt.Errorf(
			"Custom properties key %s collides with a property needed by rosa",
			consts.CreatorArn,
		)
	}

	if config.Creator.ARN != "" {
		clusterProperties[consts.CreatorArn] = config.Creator.ARN
	}

	// Create the cluster:
	clusterBuilder := cmv1.NewCluster().
		Name(config.Name).
		MultiAZ(config.MultiAZ).
		Product(
			cmv1.NewProduct().
				ID("rosa"),
		).
		Region(
			cmv1.NewCloudRegion().
				ID(config.Region),
		).
		FIPS(config.FIPS).
		EtcdEncryption(config.EtcdEncryption).
		Properties(clusterProperties)

	if config.DisableWorkloadMonitoring != nil {
		clusterBuilder = clusterBuilder.DisableUserWorkloadMonitoring(*config.DisableWorkloadMonitoring)
	}

	if config.Flavour != "" {
		clusterBuilder = clusterBuilder.Flavour(
			cmv1.NewFlavour().
				ID(config.Flavour),
		)
	}

	if config.Version != "" {
		clusterBuilder = clusterBuilder.Version(
			cmv1.NewVersion().
				ID(config.Version).
				ChannelGroup(config.ChannelGroup),
		)
	}

	if !config.Expiration.IsZero() {
		clusterBuilder = clusterBuilder.ExpirationTimestamp(config.Expiration)
	}

	if config.Hypershift.Enabled {
		hyperShiftBuilder := cmv1.NewHypershift().Enabled(true)
		clusterBuilder.Hypershift(hyperShiftBuilder)
	}

	if config.ComputeMachineType != "" || config.ComputeNodes != 0 || len(config.AvailabilityZones) > 0 ||
		config.Autoscaling || len(config.ComputeLabels) > 0 {
		clusterNodesBuilder := cmv1.NewClusterNodes()
		if config.ComputeMachineType != "" {
			clusterNodesBuilder = clusterNodesBuilder.ComputeMachineType(
				cmv1.NewMachineType().ID(config.ComputeMachineType),
			)
		}
		if machinePoolRootDisk := config.MachinePoolRootDisk; machinePoolRootDisk != nil &&
			machinePoolRootDisk.Size != 0 {
			machineTypeRootVolumeBuilder := cmv1.NewRootVolume().
				AWS(cmv1.NewAWSVolume().
					Size(machinePoolRootDisk.Size))
			clusterNodesBuilder = clusterNodesBuilder.ComputeRootVolume(
				(machineTypeRootVolumeBuilder),
			)
		}
		if config.Autoscaling {
			clusterNodesBuilder = clusterNodesBuilder.AutoscaleCompute(
				cmv1.NewMachinePoolAutoscaling().
					MinReplicas(config.MinReplicas).
					MaxReplicas(config.MaxReplicas))
		} else if config.ComputeNodes != 0 {
			clusterNodesBuilder = clusterNodesBuilder.Compute(config.ComputeNodes)
		}
		if len(config.AvailabilityZones) > 0 {
			clusterNodesBuilder = clusterNodesBuilder.AvailabilityZones(config.AvailabilityZones...)
		}
		if len(config.ComputeLabels) > 0 {
			clusterNodesBuilder = clusterNodesBuilder.ComputeLabels(config.ComputeLabels)
		}
		clusterBuilder = clusterBuilder.Nodes(clusterNodesBuilder)
	}

	if config.NetworkType != "" || config.NoCni ||
		!utils.IsEmptyCIDR(config.MachineCIDR) ||
		!utils.IsEmptyCIDR(config.ServiceCIDR) ||
		!utils.IsEmptyCIDR(config.PodCIDR) ||
		config.HostPrefix != 0 {
		networkBuilder := cmv1.NewNetwork()
		if config.NetworkType != "" {
			networkBuilder = networkBuilder.Type(config.NetworkType)
		}
		if config.NoCni {
			networkBuilder = networkBuilder.Type("Other")
		}
		if !utils.IsEmptyCIDR(config.MachineCIDR) {
			networkBuilder = networkBuilder.MachineCIDR(config.MachineCIDR.String())
		}
		if !utils.IsEmptyCIDR(config.ServiceCIDR) {
			networkBuilder = networkBuilder.ServiceCIDR(config.ServiceCIDR.String())
		}
		if !utils.IsEmptyCIDR(config.PodCIDR) {
			networkBuilder = networkBuilder.PodCIDR(config.PodCIDR.String())
		}
		if config.HostPrefix != 0 {
			networkBuilder = networkBuilder.HostPrefix(config.HostPrefix)
		}
		clusterBuilder = clusterBuilder.Network(networkBuilder)
	}

	if config.Creator.AccountID == "" {
		return nil, fmt.Errorf("account ID can not be nil when building cluster spec")
	}
	awsBuilder := cmv1.NewAWS().
		AccountID(config.Creator.AccountID)

	if len(config.AdditionalComputeSecurityGroupIds) > 0 {
		awsBuilder = awsBuilder.AdditionalComputeSecurityGroupIds(config.AdditionalComputeSecurityGroupIds...)
	}

	if len(config.AdditionalInfraSecurityGroupIds) > 0 {
		awsBuilder = awsBuilder.AdditionalInfraSecurityGroupIds(config.AdditionalInfraSecurityGroupIds...)
	}

	if len(config.AdditionalControlPlaneSecurityGroupIds) > 0 {
		awsBuilder = awsBuilder.AdditionalControlPlaneSecurityGroupIds(config.AdditionalControlPlaneSecurityGroupIds...)
	}

	if config.SubnetIds != nil {
		awsBuilder = awsBuilder.SubnetIDs(config.SubnetIds...)
	}

	if config.PrivateLink != nil {
		awsBuilder = awsBuilder.PrivateLink(*config.PrivateLink)
		if *config.PrivateLink {
			*config.Private = true
		}
	}

	if config.BillingAccount != "" {
		awsBuilder = awsBuilder.BillingAccountID(config.BillingAccount)
	}

	if config.Ec2MetadataHttpTokens != "" {
		awsBuilder = awsBuilder.Ec2MetadataHttpTokens(config.Ec2MetadataHttpTokens)
	}

	if config.RoleARN != "" {
		stsBuilder := cmv1.NewSTS().RoleARN(config.RoleARN)
		if config.ExternalID != "" {
			stsBuilder = stsBuilder.ExternalID(config.ExternalID)
		}
		if config.SupportRoleARN != "" {
			stsBuilder = stsBuilder.SupportRoleARN(config.SupportRoleARN)
		}
		if len(config.OperatorIAMRoles) > 0 {
			roles := []*cmv1.OperatorIAMRoleBuilder{}
			for _, role := range config.OperatorIAMRoles {
				roles = append(roles, cmv1.NewOperatorIAMRole().
					Name(role.Name).
					Namespace(role.Namespace).
					RoleARN(role.RoleARN),
				)
			}
			stsBuilder = stsBuilder.OperatorIAMRoles(roles...)
		}
		if config.OidcConfigId != "" {
			stsBuilder = stsBuilder.OidcConfig(cmv1.NewOidcConfig().ID(config.OidcConfigId))
		}
		instanceIAMRolesBuilder := cmv1.NewInstanceIAMRoles()
		if config.ControlPlaneRoleARN != "" {
			instanceIAMRolesBuilder.MasterRoleARN(config.ControlPlaneRoleARN)
		}
		if config.WorkerRoleARN != "" {
			instanceIAMRolesBuilder.WorkerRoleARN(config.WorkerRoleARN)
		}
		stsBuilder = stsBuilder.InstanceIAMRoles(instanceIAMRolesBuilder)

		mode := false
		if config.Mode == models.ModeAuto {
			mode = true
		}
		stsBuilder.AutoMode(mode)
		awsBuilder = awsBuilder.STS(stsBuilder)
	} else {
		if config.Creator.AccessKeyID == "" || config.Creator.SecretAccessKey == "" {
			return nil, fmt.Errorf("access keys can not be empty for non STS cluster spec")
		}
		awsBuilder = awsBuilder.
			AccessKeyID(config.Creator.AccessKeyID).
			SecretAccessKey(config.Creator.SecretAccessKey)
	}
	if config.KMSKeyArn != "" {
		awsBuilder = awsBuilder.KMSKeyArn(config.KMSKeyArn)
	}
	if len(config.Tags) > 0 {
		awsBuilder = awsBuilder.Tags(config.Tags)
	}

	if config.AuditLogRoleARN != nil {
		auditLogBuiler := cmv1.NewAuditLog().RoleArn(*config.AuditLogRoleARN)
		awsBuilder = awsBuilder.AuditLog(auditLogBuiler)
	}

	// etcd encryption kms key arn
	if config.EtcdEncryptionKMSArn != "" {
		awsBuilder = awsBuilder.EtcdEncryption(cmv1.NewAwsEtcdEncryption().KMSKeyARN(config.EtcdEncryptionKMSArn))
	}

	// shared vpc
	if config.PrivateHostedZoneID != "" {
		awsBuilder = awsBuilder.PrivateHostedZoneID(config.PrivateHostedZoneID)
		awsBuilder = awsBuilder.PrivateHostedZoneRoleARN(config.SharedVPCRoleArn)
	}
	if config.BaseDomain != "" {
		clusterBuilder = clusterBuilder.DNS(cmv1.NewDNS().BaseDomain(config.BaseDomain))
	}

	clusterBuilder = clusterBuilder.AWS(awsBuilder)

	if config.Private != nil {
		if *config.Private {
			clusterBuilder = clusterBuilder.API(
				cmv1.NewClusterAPI().
					Listening(cmv1.ListeningMethodInternal),
			)
		} else {
			clusterBuilder = clusterBuilder.API(
				cmv1.NewClusterAPI().
					Listening(cmv1.ListeningMethodExternal),
			)
		}
	}

	if config.DisableSCPChecks != nil && *config.DisableSCPChecks {
		clusterBuilder = clusterBuilder.CCS(cmv1.NewCCS().
			Enabled(true).
			DisableSCPChecks(true),
		)
	}

	if config.HTTPProxy != nil || config.HTTPSProxy != nil {
		proxyBuilder := cmv1.NewProxy()
		if config.HTTPProxy != nil {
			proxyBuilder.HTTPProxy(*config.HTTPProxy)
		}
		if config.HTTPSProxy != nil {
			proxyBuilder.HTTPSProxy(*config.HTTPSProxy)
		}
		if config.NoProxy != nil {
			proxyBuilder.NoProxy(*config.NoProxy)
		}
		clusterBuilder = clusterBuilder.Proxy(proxyBuilder)
	}

	if config.AdditionalTrustBundle != nil {
		clusterBuilder = clusterBuilder.AdditionalTrustBundle(*config.AdditionalTrustBundle)
	}

	if config.ClusterAdminUser != "" {
		hashedPwd, err := idpUtils.GenerateHTPasswdCompatibleHash(config.ClusterAdminPassword)
		if err != nil {
			return nil, fmt.Errorf("Failed to get access keys for user '%s': %v",
				models.AdminUserName, err)
		}
		var htpasswdUsers []*cmv1.HTPasswdUserBuilder
		htpasswdUsers = append(htpasswdUsers, cmv1.NewHTPasswdUser().
			Username(config.ClusterAdminUser).HashedPassword(hashedPwd))
		htpassUserList := cmv1.NewHTPasswdUserList().Items(htpasswdUsers...)
		htPasswdIDP := cmv1.NewHTPasswdIdentityProvider().Users(htpassUserList)
		clusterBuilder = clusterBuilder.Htpasswd(htPasswdIDP)
	}

	if !reflect.DeepEqual(config.DefaultIngress, models.NewDefaultIngressSpec()) {
		defaultIngress := cmv1.NewIngress().Default(true)
		if len(config.DefaultIngress.RouteSelectors) != 0 {
			defaultIngress.RouteSelectors(config.DefaultIngress.RouteSelectors)
		}
		if len(config.DefaultIngress.ExcludedNamespaces) != 0 {
			defaultIngress.ExcludedNamespaces(config.DefaultIngress.ExcludedNamespaces...)
		}
		if !utils.Contains([]string{"", consts.SkipSelectionOption}, config.DefaultIngress.WildcardPolicy) {
			defaultIngress.RouteWildcardPolicy(cmv1.WildcardPolicy(config.DefaultIngress.WildcardPolicy))
		}
		if !utils.Contains([]string{"", consts.SkipSelectionOption}, config.DefaultIngress.NamespaceOwnershipPolicy) {
			defaultIngress.RouteNamespaceOwnershipPolicy(
				cmv1.NamespaceOwnershipPolicy(config.DefaultIngress.NamespaceOwnershipPolicy))
		}
		clusterBuilder.Ingresses(cmv1.NewIngressList().Items(defaultIngress))
	}

	if config.AutoscalerConfig != nil {
		clusterBuilder.Autoscaler(BuildClusterAutoscaler(config.AutoscalerConfig))
	}

	clusterSpec, err := clusterBuilder.Build()
	if err != nil {
		return nil, fmt.Errorf("Failed to create description of cluster: %v", err)
	}

	return clusterSpec, nil
}
