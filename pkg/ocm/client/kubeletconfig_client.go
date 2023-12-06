package client

import (
	"context"
	"net/http"

	v1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
)

// KubeletConfigClient wraps the OCM SDK to provide a more unit testable friendly way of
// interacting with the OCM API
//
//go:generate mockgen -source=kubeletconfig_client.go -package=testing -destination=testing/mock_kubeletconfig_client.go
type KubeletConfigClient interface {

	// Exists returns true if the KubeletConfig for the cluster exists
	Exists(ctx context.Context, clusterId string) (bool, *v1.KubeletConfig, error)

	// Get returns the existing KubeletConfig for the cluster (if any)
	Get(ctx context.Context, clusterId string) (*v1.KubeletConfig, error)

	// Create attempts to create a KubeletConfig for the cluster
	Create(ctx context.Context, clusterId string, config *v1.KubeletConfig) (*v1.KubeletConfig, error)

	// Update attempts to update the existing KubeletConfig for the cluster
	Update(ctx context.Context, clusterId string, config *v1.KubeletConfig) (*v1.KubeletConfig, error)

	// Delete removes the existing KubeletConfig for the cluster
	Delete(ctx context.Context, clusterId string) error
}

type KubeletConfigClientImpl struct {
	collection *v1.ClustersClient
}

func (k *KubeletConfigClientImpl) Exists(ctx context.Context, clusterId string) (bool, *v1.KubeletConfig, error) {
	response, err := k.collection.Cluster(clusterId).KubeletConfig().Get().SendContext(ctx)
	if err != nil {
		// A 404 indicates that the resource does not exist
		if response.Status() == http.StatusNotFound {
			return false, nil, nil
		}
		return false, nil, err
	}
	return response.Status() == http.StatusOK, response.Body(), nil
}

func (k *KubeletConfigClientImpl) Get(ctx context.Context, clusterId string) (*v1.KubeletConfig, error) {
	response, err := k.collection.Cluster(clusterId).KubeletConfig().Get().SendContext(ctx)
	if err != nil {
		return nil, err
	}
	return response.Body(), nil
}

func (k *KubeletConfigClientImpl) Create(ctx context.Context, clusterId string, config *v1.KubeletConfig) (*v1.KubeletConfig, error) {
	response, err := k.collection.Cluster(clusterId).KubeletConfig().Post().Body(config).SendContext(ctx)
	if err != nil {
		return nil, err
	}
	return response.Body(), nil
}

func (k *KubeletConfigClientImpl) Update(ctx context.Context, clusterId string, config *v1.KubeletConfig) (*v1.KubeletConfig, error) {
	response, err := k.collection.Cluster(clusterId).KubeletConfig().Update().Body(config).SendContext(ctx)
	if err != nil {
		return nil, err
	}
	return response.Body(), nil
}

func (k *KubeletConfigClientImpl) Delete(ctx context.Context, clusterId string) error {
	_, err := k.collection.Cluster(clusterId).KubeletConfig().Delete().SendContext(ctx)
	return err
}

func NewKubeletConfigClient(collection *v1.ClustersClient) KubeletConfigClient {
	return &KubeletConfigClientImpl{collection: collection}
}

var _ KubeletConfigClient = &KubeletConfigClientImpl{}
