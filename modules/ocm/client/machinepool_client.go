package client

import (
	"context"
	v1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
)

//go:generate mockgen -source=machinepool_client.go -package=testing -destination=testing/mock_machinepool_client.go
type MachinePoolClient interface {
	CollectionClusterSubResource[v1.MachinePool, string]
}

func NewMachinePoolClient(collection *v1.ClustersClient) MachinePoolClient {
	return &CollectionClusterSubResourceImpl[v1.MachinePool, string]{
		getFunc: func(ctx context.Context, clusterId string, instanceId string) (OcmInstanceResponse[v1.MachinePool], error) {
			return collection.Cluster(clusterId).MachinePools().MachinePool(instanceId).Get().SendContext(ctx)
		},
		updateFunc: func(ctx context.Context, clusterId string, instance *v1.MachinePool) (OcmInstanceResponse[v1.MachinePool], error) {
			return collection.Cluster(clusterId).MachinePools().MachinePool(instance.ID()).Update().Body(instance).SendContext(ctx)
		},
		createFunc: func(ctx context.Context, clusterId string, instance *v1.MachinePool) (OcmInstanceResponse[v1.MachinePool], error) {
			return collection.Cluster(clusterId).MachinePools().Add().Body(instance).SendContext(ctx)
		},
		deleteFunc: func(ctx context.Context, clusterId string, instanceId string) (OcmResponse, error) {
			return collection.Cluster(clusterId).MachinePools().MachinePool(instanceId).Delete().SendContext(ctx)
		},
		listFunc: func(ctx context.Context, clusterId string, paging Paging) (OcmListResponse[v1.MachinePool], error) {
			resp, err := collection.Cluster(clusterId).MachinePools().List().Size(paging.size).Page(paging.page).SendContext(ctx)
			if err != nil {
				return nil, err
			}
			return NewListResponse(resp.Status(), resp.Items().Slice()), nil
		},
	}
}
