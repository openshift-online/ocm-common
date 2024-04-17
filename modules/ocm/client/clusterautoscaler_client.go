package client

import (
	"context"
	v1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
)

//go:generate mockgen -source=clusterautoscaler_client.go -package=testing -destination=testing/mock_clusterautoscaler_client.go
type ClusterAutoscalerClient interface {
	SingleClusterSubResource[v1.ClusterAutoscaler]
}

func NewClusterAutoscalerClient(collection *v1.ClustersClient) ClusterAutoscalerClient {
	return &SingleClusterSubResourceImpl[v1.ClusterAutoscaler]{
		getFunc: func(ctx context.Context, clusterId string) (OcmInstanceResponse[v1.ClusterAutoscaler], error) {
			return collection.Cluster(clusterId).Autoscaler().Get().SendContext(ctx)
		},
		updateFunc: func(ctx context.Context, clusterId string, instance *v1.ClusterAutoscaler) (OcmInstanceResponse[v1.ClusterAutoscaler], error) {
			return collection.Cluster(clusterId).Autoscaler().Update().Body(instance).SendContext(ctx)
		},
		createFunc: func(ctx context.Context, clusterId string, instance *v1.ClusterAutoscaler) (OcmInstanceResponse[v1.ClusterAutoscaler], error) {
			return collection.Cluster(clusterId).Autoscaler().Post().Request(instance).SendContext(ctx)
		},
		deleteFunc: func(ctx context.Context, clusterId string) (OcmResponse, error) {
			return collection.Cluster(clusterId).Autoscaler().Delete().SendContext(ctx)
		},
	}
}
