package cluster

// AWSVolume represents AWS-specific volume configuration for cluster nodes
type AWSVolume struct {
	Type      string  `json:"type"`
	KMSKeyARN *string `json:"kmsKeyARN"`
}

// AWS represents AWS-specific configuration for cluster nodes
type AWS struct {
	RootVolume AWSVolume `json:"rootVolume"`
}

// OsDisk represents GCP-specific OS disk configuration
type OsDisk struct {
	DiskSizeGB int64 `json:"diskSizeGB"`
}

// GCP represents GCP-specific configuration for cluster nodes
type GCP struct {
	OsDisk OsDisk `json:"osDisk"`
	Type   string `json:"type"`
}

// Platform represents cloud platform-specific configuration that can be applied
// to either control plane or compute nodes
type Platform struct {
	AWS AWS `json:"aws"`
	GCP GCP `json:"gcp"`
}

// Spec represents the specification for a node pool in the install config,
// which can be either control plane or compute nodes
type Spec struct {
	Name     string   `json:"name"`
	Platform Platform `json:"platform"`
}

// InstallConfig represents the cluster installation configuration including
// control plane and compute node specifications
type InstallConfig struct {
	ControlPlane Spec   `json:"controlPlane"`
	Compute      []Spec `json:"compute"`
}
