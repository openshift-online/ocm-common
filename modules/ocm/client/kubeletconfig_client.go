package client

import (
	"context"
	v1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
)

// KubeletConfigClient wraps the OCM SDK to provide a more unit testable friendly way of
// interacting with the OCM API
//
//go:generate mockgen -source=kubeletconfig_client.go -package=testing -destination=testing/mock_kubeletconfig_client.go
type KubeletConfigClient interface {
	SingleClusterSubResource[v1.KubeletConfig]
}

func NewKubeletConfigClient(collection *v1.ClustersClient) KubeletConfigClient {
	return &SingleClusterSubResourceImpl[v1.KubeletConfig]{
		getFunc: func(ctx context.Context, clusterId string) (OcmInstanceResponse[v1.KubeletConfig], error) {
			return collection.Cluster(clusterId).KubeletConfig().Get().SendContext(ctx)
		},
		updateFunc: func(ctx context.Context, clusterId string, instance *v1.KubeletConfig) (OcmInstanceResponse[v1.KubeletConfig], error) {
			return collection.Cluster(clusterId).KubeletConfig().Update().Body(instance).SendContext(ctx)
		},
		createFunc: func(ctx context.Context, clusterId string, instance *v1.KubeletConfig) (OcmInstanceResponse[v1.KubeletConfig], error) {
			return collection.Cluster(clusterId).KubeletConfig().Post().Body(instance).SendContext(ctx)
		},
		deleteFunc: func(ctx context.Context, clusterId string) (OcmResponse, error) {
			return collection.Cluster(clusterId).KubeletConfig().Delete().SendContext(ctx)
		},
	}
}
