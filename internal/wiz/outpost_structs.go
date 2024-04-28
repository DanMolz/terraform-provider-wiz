package wiz

import (
	"encoding/json"
)

// OutpostAWS struct
type OutpostAWS struct {
	ID                     string               `json:"id"`
	Name                   string               `json:"name"`
	Enabled                bool                 `json:"enabled"`
	ServiceType            string               `json:"serviceType"`
	AllowedRegions         any                  `json:"allowedRegions"`
	SelfManaged            bool                 `json:"selfManaged"`
	ExternalInternetAccess string               `json:"externalInternetAccess"`
	SelfManagedConfig      any                  `json:"selfManagedConfig"`
	ManagedConfig          OutpostManagedConfig `json:"managedConfig"`
	Status                 string               `json:"status"`
	ErrorCode              any                  `json:"errorCode"`
	CreatedAt              string               `json:"createdAt"`
	AddedBy                struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"addedBy"`
	Clusters     []OutpostCluster    `json:"clusters"`
	CustomConfig OutpostCustomConfig `json:"customConfig"`
	Config       OutpostAWSConfig    `json:"config"`
}

// OutpostCluster struct
type OutpostCluster struct {
	ID              string `json:"id"`
	Region          string `json:"region"`
	CreatedAt       string `json:"createdAt"`
	HTTPProxyConfig any    `json:"httpProxyConfig"`
	NodeGroups      []struct {
		NodeGroupID  string `json:"nodeGroupId"`
		Type         string `json:"type"`
		MaxNodeCount int    `json:"maxNodeCount"`
		MinNodeCount int    `json:"minNodeCount"`
		Typename     string `json:"__typename"`
	} `json:"nodeGroups"`
	Config struct {
		ClusterName                  string `json:"clusterName"`
		SqsURL                       string `json:"sqsURL"`
		KubernetesServiceAccountName any    `json:"kubernetesServiceAccountName"`
		Typename                     string `json:"__typename"`
	} `json:"config"`
	AddedBy struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		Email    string `json:"email"`
		Typename string `json:"__typename"`
	} `json:"addedBy"`
}

// OutpostCustomConfig struct
type OutpostCustomConfig struct {
	PodAnnotations  any `json:"podAnnotations"`
	ResourceTags    any `json:"resourceTags"`
	NamespacePrefix any `json:"namespacePrefix"`
}

// OutpostManagedConfig struct
type OutpostManagedConfig struct {
	KubernetesLoggingEnabled         bool `json:"kubernetesLoggingEnabled"`
	ManualNetwork                    bool `json:"manualNetwork"`
	KubernetesCloudMonitoringEnabled bool `json:"kubernetesCloudMonitoringEnabled"`
}

// OutpostAWSConfig struct
type OutpostAWSConfig struct {
	RoleARN           string `json:"roleARN"`
	ExternalID        string `json:"externalID"`
	StateBucketName   string `json:"stateBucketName,omitempty"`
	SettingsRegion    string `json:"settingsRegion,omitempty"`
	AccessKey         string `json:"accessKey"`
	SecretKey         string `json:"secretKey"`
	DisableNatGateway bool   `json:"disableNatGateway,omitempty"`
	ResultsBucketName string `json:"resultsBucketName,omitempty"`
	SubscriptionID    string `json:"subscriptionID"`
}

// CreateOutpostAWSInput struct
type CreateOutpostAWSInput struct {
	Name        string `json:"name"`
	ServiceType string `json:"serviceType"`
	Enabled     bool   `json:"enabled"`
	SelfManaged string `json:"selfManaged"`
	Config      struct {
		AwsConfig struct {
			SettingsRegion  string `json:"settingsRegion"`
			StateBucketName string `json:"stateBucketName"`
			RoleARN         string `json:"roleARN"`
		} `json:"awsConfig"`
	} `json:"config"`
	ManagedConfig struct {
		ManualNetwork                    bool `json:"manualNetwork"`
		KubernetesLoggingEnabled         bool `json:"kubernetesLoggingEnabled"`
		KubernetesCloudMonitoringEnabled bool `json:"kubernetesCloudMonitoringEnabled"`
	} `json:"managedConfig"`
	AllowedRegions string `json:"allowedRegions"`
}

// CreateOutpostAWSPayload struct
type CreateOutpostAWSPayload struct {
	OutpostAWS OutpostAWS `json:"outpostAWS,omitempty"`
}

// DeleteControlInput struct
type DeleteOutpostAWSInput struct {
	ID string `json:"id"`
}

// DeleteOutpostAWSPayload struct
type DeleteOutpostAWSPayload struct {
	Stub string `json:"_stub,omitempty"`
}

// UpdateOutpostAWSInput struct
type UpdateOutpostAWSInput struct {
	ID    string                `json:"id"`
	Patch UpdateOutpostAWSPatch `json:"patch"`
}

// UpdateOutpostAWSPatch struct
type UpdateOutpostAWSPatch struct {
	Query   json.RawMessage `json:"query,omitempty"`
	Enabled *bool           `json:"enabled,omitempty"`
	Name    string          `json:"name,omitempty"`
}

// UpdateOutpostAWSPayload struct
type UpdateOutpostAWSPayload struct {
	OutpostAWS OutpostAWS `json:"outpostAWS,omitempty"`
}
