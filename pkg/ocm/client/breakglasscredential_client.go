package client

import (
	"context"
	"errors"

	v1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
)

// BreakGlassCredentialClient wraps the OCM SDK to provide a more unit testable friendly way of
// interacting with the OCM API
//
//go:generate mockgen -source=breakglasscredential_client.go -package=test -destination=test/mock_breakglasscredential_client.go
type BreakGlassCredentialClient interface {
	CollectionClusterSubResource[v1.BreakGlassCredential, string]
}

func NewBreakGlassCredentialClient(collection *v1.ClustersClient) BreakGlassCredentialClient {
	return &CollectionClusterSubResourceImpl[v1.BreakGlassCredential, string]{
		getFunc: func(ctx context.Context, clusterId, breakGlassCredentialId string) (OcmInstanceResponse[v1.BreakGlassCredential], error) {
			return collection.Cluster(clusterId).BreakGlassCredentials().BreakGlassCredential(breakGlassCredentialId).Get().SendContext(ctx)
		},
		createFunc: func(ctx context.Context, clusterId string, instance *v1.BreakGlassCredential) (OcmInstanceResponse[v1.BreakGlassCredential], error) {
			return collection.Cluster(clusterId).BreakGlassCredentials().Add().Body(instance).SendContext(ctx)
		},
		listFunc: func(ctx context.Context, clusterId string, paging Paging) (OcmListResponse[v1.BreakGlassCredential], error) {
			resp, err := collection.Cluster(clusterId).BreakGlassCredentials().List().Size(paging.size).Page(paging.page).SendContext(ctx)
			if err != nil {
				return nil, err
			}
			return NewListResponse(resp.Status(), resp.Items().Slice()), nil
		},
		deleteFunc: func(ctx context.Context, clusterId, breakGlassCredentialId string) (OcmResponse, error) {
			return nil, errors.New("Delete Break Glass Credential is not supported")
		},
		updateFunc: func(ctx context.Context, clusterId string, instance *v1.BreakGlassCredential) (OcmInstanceResponse[v1.BreakGlassCredential], error) {
			return nil, errors.New("Update Break Glass Credential is not supported")
		},
	}
}
