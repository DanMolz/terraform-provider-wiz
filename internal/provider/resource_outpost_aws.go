package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"wiz.io/hashicorp/terraform-provider-wiz/internal"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/client"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/utils"
	"wiz.io/hashicorp/terraform-provider-wiz/internal/wiz"
)

func resourceWizOutpostAWS() *schema.Resource {
	return &schema.Resource{
		Description: "This resource allows you to create, read, update, and delete Wiz Outpost Configuration.",
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "Wiz identifier for the Output.",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "Name of the Output.",
				Required:    true,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "Whether to enable the Outpost.",
				Optional:    true,
				Default:     true,
			},
			"self_managed": {
				Type:        schema.TypeBool,
				Description: "Whether to enable self managed Outpost.",
				Optional:    true,
				Default:     false,
			},
			"orchestrator_role_arn": {
				Type:        schema.TypeString,
				Description: "The role is used to setup and monitor the Outpost deployment in-account",
				Required:    true,
			},
			"configuration_bucket_name": {
				Type:        schema.TypeString,
				Description: "The Configuration Bucket is used to configure the EKS Clusters running as part of the Outpost Connector",
				Required:    true,
			},
			"configuration_bucket_region": {
				Type:        schema.TypeString,
				Description: "The region where the Configuration Bucket has been created",
				Required:    true,
			},
			"results_bucket_name": {
				Type:        schema.TypeString,
				Description: "The region where the Configuration Bucket has been created",
				Optional:    true,
			},
			"disable_nat_gateway": {
				Type:        schema.TypeBool,
				Description: "Whether to disable NAT Gateway.",
				Optional:    true,
				Default:     false,
			},
			"manual_network_management": {
				Type:        schema.TypeBool,
				Description: "Whether to enable manual network configuration.",
				Optional:    true,
				Default:     false,
			},
			"kubernetes_logging_enabled": {
				Type:        schema.TypeBool,
				Description: "Whether to enable Kubernetes Logging.",
				Computed:    true,
			},
			"kubernetes_cloud_monitoring_enabled": {
				Type:        schema.TypeBool,
				Description: "Whether to enable Kubernetes Cloud Monitoring.",
				Computed:    true,
			},
			"allowed_regions": {
				Type:        schema.TypeList,
				Description: "List of allowed regions for the Outpost.",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
		CreateContext: resourceWizOutpostAWSCreate,
		ReadContext:   resourceWizOutpostAWSRead,
		UpdateContext: resourceWizOutpostAWSUpdate,
		DeleteContext: resourceWizOutpostAWSDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// CreateOutpost struct
type CreateOutpost struct {
	CreateOutpost wiz.CreateOutpostPayload `json:"createOutpost"`
}

func resourceWizOutpostAWSCreate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizOutpostCreate called...")

	// define the graphql query
	query := `mutation CreateOutpost($input: CreateOutpostInput!) {
    createOutpost(input: $input) {
      outpost {
        id
      }
    }
  }`

	// populate the graphql variables
	vars := &wiz.CreateOutpostInput{}
	vars.Name = d.Get("name").(string)
	vars.ServiceType = "AWS"
	enabled := d.Get("enabled").(bool)
	vars.Enabled = &enabled
	selfManaged := d.Get("self_managed").(bool)
	vars.SelfManaged = &selfManaged
	vars.Config.AwsConfig.RoleARN = d.Get("orchestrator_role_arn").(string)
	vars.Config.AwsConfig.StateBucketName = d.Get("configuration_bucket_name").(string)
	vars.Config.AwsConfig.SettingsRegion = d.Get("configuration_bucket_region").(string)
	vars.Config.AwsConfig.ResultsBucketName = d.Get("results_bucket_name").(string)
	disableNatGateway := d.Get("disable_nat_gateway").(bool)
	vars.Config.AwsConfig.DisableNatGateway = &disableNatGateway
	manualNetwork := d.Get("manual_network_management").(bool)
	vars.ManagedConfig.ManualNetwork = &manualNetwork
	kubernetesLoggingEnabled := d.Get("kubernetes_logging_enabled").(bool)
	vars.ManagedConfig.KubernetesLoggingEnabled = &kubernetesLoggingEnabled
	kubernetesCloudMonitoringEnabled := d.Get("kubernetes_cloud_monitoring_enabled").(bool)
	vars.ManagedConfig.KubernetesCloudMonitoringEnabled = &kubernetesCloudMonitoringEnabled
	vars.AllowedRegions = utils.ConvertListToString(d.Get("allowed_regions").([]interface{}))

	// process the request
	data := &CreateOutpost{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "outpost", "create")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	// set the id
	d.SetId(data.CreateOutpost.Outpost.ID)

	return resourceWizOutpostAWSRead(ctx, d, m)
}

// ReadOutpostPayload struct
type ReadOutpostPayload struct {
	Outpost wiz.Outpost `json:"outpost"`
}

func resourceWizOutpostAWSRead(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizOutpostRead called...")

	// check the id
	if d.Id() == "" {
		tflog.Debug(ctx, "outpost ID Missing")
		return nil
	}

	// define the graphql query
	query := `query OutpostDetails($id: ID!) {
      outpost(id: $id) {
        id
        name
        enabled
        serviceType
        allowedRegions
        selfManaged
        externalInternetAccess
        selfManagedConfig {
          disableAutomaticConfigurationBucketSync
          imageRepository
          version {
            id
            images {
              url
            }
          }
        }
        managedConfig {
          kubernetesLoggingEnabled
          manualNetwork
          kubernetesCloudMonitoringEnabled
        }
        externalInternetAccess
        status
        errorCode
        createdAt
        addedBy {
          id
          name
          email
        }
        clusters {
          ...OutpostClusterDetails
        }
        customConfig {
          podAnnotations
          resourceTags
          namespacePrefix
        }
        config {
          ... on OutpostAWSConfig {
            roleARN
            externalID
            stateBucketName
            settingsRegion
            accessKey
            secretKey
            disableNatGateway
            resultsBucketName
            subscriptionID
          }
          ... on OutpostGCPConfig {
            orchestratorKey
            workerAccountEmail
            stateBucketName
            disableNatGateway
          }
          ... on OutpostAzureConfig {
            tenantID
            subscriptionID
            keyVaultName
            applicationKeyVaultName
            orchestratorClientID
            orchestratorClientSecret
            workerClientID
            workerClientSecret
            scannerAppID
            deployPremiumServiceBus
            enablePrivateCluster
            environment
            stateStorageAccountName
            globalResourceGroupName
          }
          ... on OutpostOCIConfig {
            compartmentOCID
            vaultOCID
            keyOCID
            stateBucketName
            orchestrator {
              fingerprint
              privateKey
            }
          }
          ... on OutpostAlibabaConfig {
            workerResourcesGroupID
            stateBucketName
            outpostCredentials {
              orchestratorAccessKeyID
              orchestratorAccessKeySecret
            }
          }
        }
      }
    }

        fragment OutpostClusterDetails on OutpostCluster {
      id
      region
      createdAt
      httpProxyConfig {
        httpProxyURL
        httpsProxyURL
        vpcCIDRs
      }
      nodeGroups {
        nodeGroupId
        type
        maxNodeCount
        minNodeCount
      }
      config {
        ... on OutpostClusterAWSConfig {
          clusterName
          sqsURL
          kubernetesServiceAccountName
        }
        ... on OutpostClusterAzureConfig {
          clusterName
          servicebusQueueName
          servicebusNamespace
          resourceGroupName
          storageAccountNames
          subscriptionId
          serviceAuthorizedIPRanges
        }
        ... on OutpostClusterGCPConfig {
          clusterName
          projectId
          clusterZone
          topicName
          pubSubSubscription
        }
        ... on OutpostClusterOCIConfig {
          clusterName
          streamOCID
        }
        ... on OutpostClusterAlibabaConfig {
          clusterName
          queueName
        }
      }
      addedBy {
        id
        name
        email
      }
    }`

	// populate the graphql variables
	vars := &internal.QueryVariables{}
	vars.ID = d.Id()

	tflog.Debug(ctx, "resourceWizOutpostRead debug 1...")

	// process the request
	data := &ReadOutpostPayload{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "outpost", "read")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	tflog.Info(ctx, "resourceWizOutpostRead debug 3...")

	// set the resource parameters
	err := d.Set("name", data.Outpost.Name)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	err = d.Set("enabled", data.Outpost.Enabled)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	err = d.Set("self_managed", data.Outpost.SelfManaged)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	err = d.Set("orchestrator_role_arn", data.Outpost.Config.RoleARN)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	err = d.Set("configuration_bucket_name", data.Outpost.Config.StateBucketName)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	err = d.Set("configuration_bucket_region", data.Outpost.Config.SettingsRegion)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	err = d.Set("results_bucket_name", data.Outpost.Config.ResultsBucketName)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	err = d.Set("disable_nat_gateway", data.Outpost.Config.DisableNatGateway)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	err = d.Set("manual_network_management", data.Outpost.ManagedConfig.ManualNetwork)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	err = d.Set("kubernetes_logging_enabled", data.Outpost.ManagedConfig.KubernetesLoggingEnabled)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	err = d.Set("kubernetes_cloud_monitoring_enabled", data.Outpost.ManagedConfig.KubernetesCloudMonitoringEnabled)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	return diags
}

// UpdateOutpost struct
type UpdateOutpost struct {
	UpdateOutpost wiz.UpdateOutpostPayload `json:"updateOutpost"`
}

func resourceWizOutpostAWSUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizOutpostUpdate called...")

	// check the id
	if d.Id() == "" {
		return nil
	}

	// define the graphql query
	query := `mutation UpdateOutpost($input: UpdateOutpostInput!) {
      updateOutpost(input: $input) {
        outpost {
          id
          name
          enabled
          status
          errorCode
          managedConfig {
            kubernetesLoggingEnabled
            manualNetwork
          }
          selfManagedConfig {
            disableAutomaticConfigurationBucketSync
            version {
              id
              images {
                url
              }
            }
          }
          customConfig {
            podAnnotations
            resourceTags
            namespacePrefix
          }
          config {
            ... on OutpostAWSConfig {
              roleARN
              externalID
              stateBucketName
              settingsRegion
              accessKey
              secretKey
              disableNatGateway
              subscriptionID
            }
            ... on OutpostGCPConfig {
              orchestratorKey
              workerAccountEmail
              stateBucketName
              disableNatGateway
            }
            ... on OutpostAzureConfig {
              tenantID
              subscriptionID
              keyVaultName
              applicationKeyVaultName
              orchestratorClientID
              orchestratorClientSecret
              workerClientID
              workerClientSecret
              scannerAppID
              environment
              stateStorageAccountName
              globalResourceGroupName
              deployPremiumServiceBus
              enablePrivateCluster
            }
            ... on OutpostAlibabaConfig {
              stateBucketName
            }
          }
        }
      }
    }`

	// populate the graphql variables
	vars := &wiz.UpdateOutpostInput{}
	vars.ID = d.Id()

	// these can optionally be included in the patch
	if d.HasChange("enabled") {
		vars.Patch.Enabled = utils.ConvertBoolToPointer(d.Get("enabled").(bool))
	}
	if d.HasChange("name") {
		vars.Patch.Name = d.Get("name").(string)
	}

	// process the request
	data := &UpdateOutpost{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "outpost", "update")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	return resourceWizOutpostAWSRead(ctx, d, m)
}

// DeleteOutpost struct
type DeleteOutpost struct {
	DeleteOutpost wiz.DeleteOutpostPayload `json:"deleteOutpost"`
}

func resourceWizOutpostAWSDelete(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizOutpostDelete called...")

	// check the id
	if d.Id() == "" {
		return nil
	}

	// define the graphql query
	query := `mutation UninstallOutpost($input: UninstallOutpostInput!) {
      uninstallOutpost(input: $input) {
        outpost {
          id
          status
        }
      }
    }`

	// populate the graphql variables
	vars := &wiz.DeleteOutpostInput{}
	vars.ID = d.Id()

	// process the request
	data := &DeleteOutpost{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "outpost", "delete")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		tflog.Debug(ctx, fmt.Sprintf("Diags Count: %d", len(diags)))
		return diags
	}

	return diags
}
