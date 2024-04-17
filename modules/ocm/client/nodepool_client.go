package client

import (
	"context"
	v1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
)

//go:generate mockgen -source=nodepool_client.go -package=testing -destination=testing/mock_nodepool_client.go
type NodePoolClient interface {
	CollectionClusterSubResource[v1.NodePool, string]
}

func NewNodePoolClient(collection *v1.ClustersClient) NodePoolClient {
	return &CollectionClusterSubResourceImpl[v1.NodePool, string]{
		getFunc: func(ctx context.Context, clusterId string, instanceId string) (OcmInstanceResponse[v1.NodePool], error) {
			return collection.Cluster(clusterId).NodePools().NodePool(instanceId).Get().SendContext(ctx)
		},
		updateFunc: func(ctx context.Context, clusterId string, instance *v1.NodePool) (OcmInstanceResponse[v1.NodePool], error) {
			return collection.Cluster(clusterId).NodePools().NodePool(instance.ID()).Update().Body(instance).SendContext(ctx)
		},
		createFunc: func(ctx context.Context, clusterId string, instance *v1.NodePool) (OcmInstanceResponse[v1.NodePool], error) {
			return collection.Cluster(clusterId).NodePools().Add().Body(instance).SendContext(ctx)
		},
		deleteFunc: func(ctx context.Context, clusterId string, instanceId string) (OcmResponse, error) {
			return collection.Cluster(clusterId).NodePools().NodePool(instanceId).Delete().SendContext(ctx)
		},
		listFunc: func(ctx context.Context, clusterId string, paging Paging) (OcmListResponse[v1.NodePool], error) {
			response, err := collection.Cluster(clusterId).NodePools().List().Size(paging.size).Page(paging.page).SendContext(ctx)
			if err != nil {
				return nil, err
			}
			return NewListResponse(response.Status(), response.Items().Slice()), nil
		},
	}
}
