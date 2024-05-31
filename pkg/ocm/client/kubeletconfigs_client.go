package client

import (
	"context"

	v1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
)

//go:generate mockgen -source=kubeletconfigs_client.go -package=test -destination=test/mock_kubeletconfigs_client.go
type KubeletConfigsClient interface {
	CollectionClusterSubResource[v1.KubeletConfig, string]
}

func NewKubeletConfigsClient(collection *v1.ClustersClient) KubeletConfigsClient {
	return &CollectionClusterSubResourceImpl[v1.KubeletConfig, string]{
		getFunc: func(ctx context.Context, clusterId string, instanceId string) (OcmInstanceResponse[v1.KubeletConfig], error) {
			return collection.Cluster(clusterId).KubeletConfigs().KubeletConfig(instanceId).Get().SendContext(ctx)
		},
		updateFunc: func(ctx context.Context, clusterId string, instance *v1.KubeletConfig) (OcmInstanceResponse[v1.KubeletConfig], error) {
			return collection.Cluster(clusterId).KubeletConfigs().KubeletConfig(instance.ID()).Update().Body(instance).SendContext(ctx)
		},
		createFunc: func(ctx context.Context, clusterId string, instance *v1.KubeletConfig) (OcmInstanceResponse[v1.KubeletConfig], error) {
			return collection.Cluster(clusterId).KubeletConfigs().Add().Body(instance).SendContext(ctx)
		},
		deleteFunc: func(ctx context.Context, clusterId string, instanceId string) (OcmResponse, error) {
			return collection.Cluster(clusterId).KubeletConfigs().KubeletConfig(instanceId).Delete().SendContext(ctx)
		},
		listFunc: func(ctx context.Context, clusterId string, paging Paging) (OcmListResponse[v1.KubeletConfig], error) {
			resp, err := collection.Cluster(clusterId).KubeletConfigs().List().Size(paging.size).Page(paging.page).SendContext(ctx)
			if err != nil {
				return nil, err
			}
			return NewListResponse(resp.Status(), resp.Items().Slice()), nil
		},
	}
}
