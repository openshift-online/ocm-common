package client

import (
	"context"

	v1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
)

//go:generate mockgen -source=cluster_client.go -package=test -destination=test/mock_cluster_client.go
type ClusterClient interface {
	CollectionClusterResource[v1.Cluster, string]
}

func NewClusterClient(collection *v1.ClustersClient) ClusterClient {
	return &CollectionClusterResourceImpl[v1.Cluster, string]{
		getFunc: func(ctx context.Context, clusterId string) (OcmInstanceResponse[v1.Cluster], error) {
			return collection.Cluster(clusterId).Get().SendContext(ctx)
		},
		updateFunc: func(ctx context.Context, clusterId string, instance *v1.Cluster) (OcmInstanceResponse[v1.Cluster], error) {
			return collection.Cluster(clusterId).Update().Body(instance).SendContext(ctx)
		},
		createFunc: func(ctx context.Context, instance *v1.Cluster) (OcmInstanceResponse[v1.Cluster], error) {
			return collection.Add().Body(instance).SendContext(ctx)
		},
		deleteFunc: func(ctx context.Context, clusterId string) (OcmResponse, error) {
			return collection.Cluster(clusterId).Delete().SendContext(ctx)
		},
		listFunc: func(ctx context.Context, paging Paging) (OcmListResponse[v1.Cluster], error) {
			resp, err := collection.List().Size(paging.size).Page(paging.page).SendContext(ctx)
			if err != nil {
				return nil, err
			}
			return NewListResponse(resp.Status(), resp.Items().Slice()), nil
		},
	}
}
