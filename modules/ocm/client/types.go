package client

import (
	"context"
	"net/http"
)

// SingleClusterSubResource defines clients that operate on a resource for which
// there is exactly one per-cluster. Examples of this type of resource include
// KubeletConfig and ClusterAutoscaler
type SingleClusterSubResource[T any] interface {
	Get(ctx context.Context, clusterId string) (*T, error)
	Exists(ctx context.Context, clusterId string) (bool, *T, error)
	Create(ctx context.Context, clusterId string, instance *T) (*T, error)
	Update(ctx context.Context, clusterId string, instance *T) (*T, error)
	Delete(ctx context.Context, clusterId string) error
}

// CollectionClusterSubResource defines clients that operate on resources that have many-to-one
// relationship with a cluster. For example, MachinePools - there are many MachinePools for a single cluster
type CollectionClusterSubResource[T any, S any] interface {
	Get(ctx context.Context, clusterId string, instanceId S) (*T, error)
	Exists(ctx context.Context, clusterId string, instanceId S) (bool, *T, error)
	Create(ctx context.Context, clusterId string, instance *T) (*T, error)
	Update(ctx context.Context, clusterId string, instance *T) (*T, error)
	Delete(ctx context.Context, clusterId string, instanceId S) error
	List(ctx context.Context, clusterId string, paging Paging) ([]*T, bool, error)
}

// Paging encapsulates paging requests for list methods
type Paging struct {
	size int
	page int
}

func NewPaging(page int, size int) Paging {
	return Paging{
		size: size,
		page: page,
	}
}

func NewListResponse[T any](status int, items []*T) OcmListResponse[T] {
	return &DefaultListResponse[T]{
		status: status,
		items:  items,
	}
}

type DefaultListResponse[T any] struct {
	items  []*T
	status int
}

func (d DefaultListResponse[T]) HasItems() bool {
	return len(d.items) != 0
}

func (d DefaultListResponse[T]) Status() int {
	return d.status
}

func (d DefaultListResponse[T]) Items() []*T {
	return d.items
}

type OcmResponse interface {
	Status() int
}

type OcmListResponse[T any] interface {
	OcmResponse
	Items() []*T
	HasItems() bool
}

type OcmInstanceResponse[T any] interface {
	OcmResponse
	Body() *T
}

// SingleClusterSubResourceImpl provides the basic struct for SingleClusterSubResource clients
type SingleClusterSubResourceImpl[T any] struct {
	getFunc    func(ctx context.Context, clusterId string) (OcmInstanceResponse[T], error)
	updateFunc func(ctx context.Context, clusterId string, instance *T) (OcmInstanceResponse[T], error)
	createFunc func(ctx context.Context, clusterId string, instance *T) (OcmInstanceResponse[T], error)
	deleteFunc func(ctx context.Context, clusterId string) (OcmResponse, error)
}

// CollectionClusterSubResourceImpl provides the basic struct for CollectionClusterSubResource clients
type CollectionClusterSubResourceImpl[T any, S any] struct {
	getFunc    func(ctx context.Context, clusterId string, instanceId S) (OcmInstanceResponse[T], error)
	updateFunc func(ctx context.Context, clusterId string, instance *T) (OcmInstanceResponse[T], error)
	createFunc func(ctx context.Context, clusterId string, instance *T) (OcmInstanceResponse[T], error)
	deleteFunc func(ctx context.Context, clusterId string, instanceId S) (OcmResponse, error)
	listFunc   func(ctx context.Context, clusterId string, paging Paging) (OcmListResponse[T], error)
}

func (c *CollectionClusterSubResourceImpl[T, S]) Get(ctx context.Context, clusterId string, instanceId S) (*T, error) {
	response, err := c.getFunc(ctx, clusterId, instanceId)
	if err != nil {
		return nil, err
	}
	return response.Body(), nil
}

func (c *CollectionClusterSubResourceImpl[T, S]) Exists(ctx context.Context, clusterId string, instanceId S) (bool, *T, error) {
	return exists(c.getFunc(ctx, clusterId, instanceId))
}

func (c *CollectionClusterSubResourceImpl[T, S]) Create(ctx context.Context, clusterId string, instance *T) (*T, error) {
	response, err := c.createFunc(ctx, clusterId, instance)
	if err != nil {
		return nil, err
	}
	return response.Body(), nil
}

func (c *CollectionClusterSubResourceImpl[T, S]) Update(ctx context.Context, clusterId string, instance *T) (*T, error) {
	response, err := c.updateFunc(ctx, clusterId, instance)
	if err != nil {
		return nil, err
	}
	return response.Body(), nil
}

func (c *CollectionClusterSubResourceImpl[T, S]) Delete(ctx context.Context, clusterId string, instanceId S) error {
	_, err := c.deleteFunc(ctx, clusterId, instanceId)
	return err
}

func (c *CollectionClusterSubResourceImpl[T, S]) List(ctx context.Context, clusterId string, paging Paging) ([]*T, bool, error) {
	response, err := c.listFunc(ctx, clusterId, paging)
	if err != nil {
		return make([]*T, 0), false, err
	}
	return response.Items(), response.HasItems(), nil
}

var _ SingleClusterSubResource[interface{}] = &SingleClusterSubResourceImpl[interface{}]{}
var _ CollectionClusterSubResource[interface{}, string] = &CollectionClusterSubResourceImpl[interface{}, string]{}

func (s *SingleClusterSubResourceImpl[T]) Get(ctx context.Context, clusterId string) (*T, error) {
	response, err := s.getFunc(ctx, clusterId)
	if err != nil {
		return nil, err
	}

	return response.Body(), nil
}

func exists[T any](response OcmInstanceResponse[T], err error) (bool, *T, error) {
	if err != nil {
		// A 404 indicates that the resource does not exist
		if response.Status() == http.StatusNotFound {
			return false, nil, nil
		}
		return false, nil, err
	}
	return response.Status() == http.StatusOK, response.Body(), nil
}

func (s *SingleClusterSubResourceImpl[T]) Exists(ctx context.Context, clusterId string) (bool, *T, error) {
	return exists(s.getFunc(ctx, clusterId))
}

func (s *SingleClusterSubResourceImpl[T]) Create(ctx context.Context, clusterId string, instance *T) (*T, error) {
	response, err := s.createFunc(ctx, clusterId, instance)
	if err != nil {
		return nil, err
	}
	return response.Body(), nil
}

func (s *SingleClusterSubResourceImpl[T]) Update(ctx context.Context, clusterId string, instance *T) (*T, error) {
	response, err := s.updateFunc(ctx, clusterId, instance)
	if err != nil {
		return nil, err
	}

	return response.Body(), nil
}

func (s *SingleClusterSubResourceImpl[T]) Delete(ctx context.Context, clusterId string) error {
	_, err := s.deleteFunc(ctx, clusterId)
	return err
}
