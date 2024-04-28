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
				Computed:    true,
				Description: "Wiz identifier for the Output.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the Output.",
			},
			"service_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Service type of the Outpost.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether to enable the Outpost.",
			},
			"self_managed": {
				Type:        schema.TypeBool,
				Description: "Whether to enable self managed Outpost.",
				Optional:    true,
				Default:     false,
			},
			"role_arn": {
				Type:        schema.TypeString,
				Description: "The AWS role arn for Outpost Bucket.",
				Computed:    true,
			},
			"state_blucket_name": {
				Type:        schema.TypeString,
				Description: "The Bucket name for the Outpost state.",
				Computed:    true,
			},
			"blucket_region": {
				Type:        schema.TypeString,
				Description: "The Bucket region for the Outpost state.",
				Computed:    true,
			},
			"manual_network": {
				Type:        schema.TypeBool,
				Description: "Whether to enable manual network configuration.",
				Optional:    true,
			},
			"kubernetes_logging_enabled": {
				Type:        schema.TypeBool,
				Description: "Whether to enable Kubernetes Logging.",
				Optional:    true,
			},
			"kubernetes_cloud_monitoring_enabled": {
				Type:        schema.TypeBool,
				Description: "Whether to enable Kubernetes Cloud Monitoring.",
				Optional:    true,
			},
			"allowed_regions": {
				Type:        schema.TypeList,
				Description: "List of allowed regions for the Outpost.",
				Computed:    true,
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

// CreateOutpostAWS struct
type CreateOutpostAWS struct {
	CreateOutpostAWS wiz.CreateOutpostAWSPayload `json:"createOutpostAWS"`
}

func resourceWizOutpostAWSCreate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizOutpostAWSCreate called...")

	// define the graphql query
	query := `mutation CreateOutpost($input: CreateOutpostInput!) {
    createOutpost(input: $input) {
      outpost {
        id
      }
    }
  }`

	// populate the graphql variables
	vars := &wiz.CreateOutpostAWSInput{}
	vars.Name = d.Get("name").(string)

	// process the request
	data := &CreateOutpostAWS{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "control", "create")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	// set the id
	d.SetId(data.CreateOutpostAWS.OutpostAWS.ID)

	return resourceWizOutpostAWSRead(ctx, d, m)
}

// ReadOutpostAWSPayload struct
type ReadOutpostAWSPayload struct {
	OutpostAWS wiz.OutpostAWS `json:"outpostAWS"`
}

func resourceWizOutpostAWSRead(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizOutpostAWSRead called...")

	// check the id
	if d.Id() == "" {
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

	// process the request
	// this query returns http 200 with a payload that contains errors and a null data body
	// error message: oops! an internal error has occurred. for reference purposes, this is your request id
	data := &ReadOutpostAWSPayload{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "control", "read")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		tflog.Info(ctx, "Error from API call, checking if resource was deleted outside Terraform.")
		if data.OutpostAWS.ID == "" {
			tflog.Debug(ctx, fmt.Sprintf("Response: (%T) %s", data, utils.PrettyPrint(data)))
			tflog.Info(ctx, "Resource not found, marking as new.")
			d.SetId("")
			d.MarkNewResource()
			return nil
		}
		return diags
	}

	// set the resource parameters
	err := d.Set("name", data.OutpostAWS.Name)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	return diags
}

// UpdateOutpostAWS struct
type UpdateOutpostAWS struct {
	UpdateOutpostAWS wiz.UpdateOutpostAWSPayload `json:"updateOutpostAWS"`
}

func resourceWizOutpostAWSUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizOutpostAWSUpdate called...")

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
	vars := &wiz.UpdateOutpostAWSInput{}
	vars.ID = d.Id()

	// these can optionally be included in the patch
	if d.HasChange("enabled") {
		vars.Patch.Enabled = utils.ConvertBoolToPointer(d.Get("enabled").(bool))
	}
	if d.HasChange("name") {
		vars.Patch.Name = d.Get("name").(string)
	}

	// process the request
	data := &UpdateOutpostAWS{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "control", "update")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	return resourceWizOutpostAWSRead(ctx, d, m)
}

// DeleteOutpostAWS struct
type DeleteOutpostAWS struct {
	DeleteOutpostAWS wiz.DeleteOutpostAWSPayload `json:"deleteOutpostAWS"`
}

func resourceWizOutpostAWSDelete(ctx context.Context, d *schema.ResourceData, m interface{}) (diags diag.Diagnostics) {
	tflog.Info(ctx, "resourceWizOutpostAWSDelete called...")

	// check the id
	if d.Id() == "" {
		return nil
	}

	// define the graphql query
	query := `mutation DeleteControl (
	    $input: DeleteControlInput!
	) {
	    deleteControl(
		input: $input
	    ) {
		_stub
	    }
	}`

	// populate the graphql variables
	vars := &wiz.DeleteOutpostAWSInput{}
	vars.ID = d.Id()

	// process the request
	data := &UpdateOutpostAWS{}
	requestDiags := client.ProcessRequest(ctx, m, vars, data, query, "control", "delete")
	diags = append(diags, requestDiags...)
	if len(diags) > 0 {
		return diags
	}

	return diags
}
