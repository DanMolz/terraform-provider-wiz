package wiz

import (
	"encoding/json"
)

// Outpost struct
type Outpost struct {
	ID                     string               `json:"id"`
	Name                   string               `json:"name"`
	Enabled                bool                 `json:"enabled"`
	ServiceType            string               `json:"serviceType"`
	AllowedRegions         []string             `json:"allowedRegions"`
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
	Config       OutpostConfig       `json:"config"`
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

// OutpostConfig struct
type OutpostConfig struct {
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

// CreateOutpostInput struct
type CreateOutpostInput struct {
	Name           string                    `json:"name,omitempty"`
	ServiceType    string                    `json:"serviceType,omitempty"`
	Enabled        *bool                     `json:"enabled,omitempty"`
	SelfManaged    *bool                     `json:"selfManaged,omitempty"`
	Config         OutpostConfigInput        `json:"config"`
	ManagedConfig  OutpostManagedConfigInput `json:"managedConfig"`
	AllowedRegions []string                  `json:"allowedRegions,omitempty"`
}

type OutpostConfigInput struct {
	AwsConfig OutpostAWSConfigInput `json:"awsConfig"`
}

// OutpostAWSConfig struct
type OutpostAWSConfigInput struct {
	RoleARN           string `json:"roleARN,omitempty"`
	ExternalID        string `json:"externalID,omitempty"`
	StateBucketName   string `json:"stateBucketName,omitempty"`
	SettingsRegion    string `json:"settingsRegion,omitempty"`
	AccessKey         string `json:"accessKey,omitempty"`
	SecretKey         string `json:"secretKey,omitempty"`
	DisableNatGateway *bool  `json:"disableNatGateway,omitempty"`
	ResultsBucketName string `json:"resultsBucketName,omitempty"`
	SubscriptionID    string `json:"subscriptionID,omitempty"`
}

// OutpostManagedConfig struct
type OutpostManagedConfigInput struct {
	KubernetesLoggingEnabled         *bool `json:"kubernetesLoggingEnabled"`
	ManualNetwork                    *bool `json:"manualNetwork"`
	KubernetesCloudMonitoringEnabled *bool `json:"kubernetesCloudMonitoringEnabled"`
}

// CreateOutpostPayload struct
type CreateOutpostPayload struct {
	Outpost Outpost `json:"outpost,omitempty"`
}

// DeleteOutpostInput struct
type DeleteOutpostInput struct {
	ID string `json:"id"`
}

// DeleteOutpostPayload struct
type DeleteOutpostPayload struct {
	Stub string `json:"_stub,omitempty"`
}

// UpdateOutpostInput struct
type UpdateOutpostInput struct {
	ID    string             `json:"id"`
	Patch UpdateOutpostPatch `json:"patch"`
}

// UpdateOutpostPatch struct
type UpdateOutpostPatch struct {
	Query   json.RawMessage `json:"query,omitempty"`
	Enabled *bool           `json:"enabled,omitempty"`
	Name    string          `json:"name,omitempty"`
}

// UpdateOutpostPayload struct
type UpdateOutpostPayload struct {
	Outpost Outpost `json:"outpost,omitempty"`
}
