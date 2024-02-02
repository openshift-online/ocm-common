package models

type DefaultIngressSpec struct {
	RouteSelectors           map[string]string
	ExcludedNamespaces       []string
	WildcardPolicy           string
	NamespaceOwnershipPolicy string
}

func NewDefaultIngressSpec() DefaultIngressSpec {
	defaultIngressSpec := DefaultIngressSpec{}
	defaultIngressSpec.RouteSelectors = map[string]string{}
	defaultIngressSpec.ExcludedNamespaces = []string{}
	return defaultIngressSpec
}
