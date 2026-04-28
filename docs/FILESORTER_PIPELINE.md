# Filesorter Pipeline

## Overview

The Filesorter pipeline is a specialized pipeline type that combines three processing steps to automatically classify and route documents:

1. **HephaistosPdfeed** - OCR and text extraction from PDFs
2. **Labelize** - LLM-based document classification using Mistral models
3. **CategoryRouterMetadataAware** - Intelligent file routing based on classification results

This pipeline is ideal for document management workflows where files need to be automatically categorized and organized based on their content.

## Architecture

```
Source S3 Bucket → OCR/Text Extraction → AI Classification → Category-based Routing → Destination S3
                   (HephaistosPdfeed)     (Labelize)         (CategoryRouter)
```

## Prerequisites

Before creating a filesorter pipeline, you need:

- **Source MinIO Data Dock** - Where input files are stored
- **Source Bucket Data Container** - The bucket containing files to process
- **Iceberg Catalog Data Container** - For metadata storage
- **Trino Data Dock** - For querying and managing metadata
- **Destination S3 Data Dock** (optional) - For sorted output files
- **Destination Bucket Data Container** (optional) - The bucket for sorted files
- **Mistral API Key** - For OCR and document classification

## Data Structures

### FileSorterOutputParameters

The main configuration structure for filesorter pipelines:

```go
type FileSorterOutputParameters struct {
    // Required: UUID of the Iceberg catalog Data Container
    DcIcebergId openapi_types.UUID `json:"dc_iceberg_id"`

    // Required: UUID of the Trino Data Dock for metadata storage
    DdTrinoInt openapi_types.UUID `json:"dd_trino_int"`

    // Optional: UUID of the destination bucket Data Container for sorted files
    DestinationBucketId *openapi_types.UUID `json:"destination_bucket_id"`

    // Optional: Destination prefix for sorted files
    DestinationPrefix *string `json:"destination_prefix"`

    // Optional: UUID of the destination S3 Data Dock for sorted files
    DestinationS3Dd *openapi_types.UUID `json:"destination_s3_dd"`

    // Required: Labels for document classification
    Labels []LabelDefinition `json:"labels"`

    // Required: Mistral API key for OCR and classification
    ModelApiKey string `json:"model_api_key"`

    // Required: Model name (e.g., "mistral-medium", "mistral-large-latest")
    ModelName string `json:"model_name"`

    // Required: Category router configuration
    Router CategoryRouterConfig `json:"router"`

    // Optional: Source S3 connection ID for router (references input.dd_minio)
    SourceFileConnectionId *string `json:"source_file_connection_id"`

    // Required: Trino schema name
    TrinoSchema string `json:"trino_schema"`

    // Required: Trino table name for metadata
    TrinoTable string `json:"trino_table"`
}
```

### LabelDefinition

Defines a classification category for the AI model:

```go
type LabelDefinition struct {
    // Label category key (e.g., "Factures/Eau", "Contrats/Entretien")
    Category string `json:"category"`

    // Human-readable description for LLM classification
    Description string `json:"description"`
}
```

### CategoryRouterConfig

Controls how files are routed based on classification:

```go
type CategoryRouterConfig struct {
    // Optional: Default category if no label matches (e.g., "inconnu")
    DefaultCategory *string `json:"default_category"`

    // Optional: Regex to extract prefix from source path (single capture group)
    // Example: "(condominium-\\d+)/"
    PrefixRegex *string `json:"prefix_regex"`
}
```

### PipelineInputParameters

Defines the source data location:

```go
type PipelineInputParameters struct {
    // Required: UUID of the bucket Data Container
    DcBucketId openapi_types.UUID `json:"dc_bucket_id"`

    // Required: UUID of the source MinIO Data Dock
    DdMinio openapi_types.UUID `json:"dd_minio"`

    // Optional: Source folder path within the bucket
    Folder *string `json:"folder"`

    // Optional: Default routing destination
    DefaultRoute *DestinationConfig `json:"default_route,omitempty"`

    // Optional: FileRouter routing rules: label -> list of destinations
    RoutingRules *map[string][]DestinationConfig `json:"routing_rules"`
}
```

### PipelineParameters

General pipeline configuration:

```go
type PipelineParameters struct {
    // Optional: Whether the pipeline is enabled
    // Only applies to CronJob pipelines. When false, the CronJob will be suspended.
    // Defaults to true (enabled) if not specified.
    Enabled *bool `json:"enabled,omitempty"`

    // Pipeline name
    Name string `json:"name"`

    // Whether this is a one-off execution (true) or continuous (false)
    OneOff bool `json:"one_off"`

    // Organization UUID
    OrgId openapi_types.UUID `json:"org_id"`

    // Organization slug
    OrgSlug string `json:"org_slug"`

    // Pipeline type (use "filesorter" for filesorter pipelines)
    Type string `json:"type"`
}
```

## Complete Example

Here's a full example of creating a filesorter pipeline:

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/google/uuid"
    "github.com/nudibranches-tech/hyperfluid-sdk-go/sdk"
    "github.com/nudibranches-tech/hyperfluid-sdk-go/sdk/controlplaneapiclient"
    "github.com/oapi-codegen/runtime/types"
)

func createFileSorterPipeline(ctx context.Context) error {
    // Initialize SDK client
    client, err := sdk.NewClientFromServiceAccountFile(
        "/path/to/service_account.json",
        sdk.ServiceAccountOptions{
            BaseURL:         "https://api.hyperfluid.cloud",
            ControlPlaneURL: "https://console.hyperfluid.cloud",
        },
    )
    if err != nil {
        return fmt.Errorf("failed to create client: %w", err)
    }

    // Get Control Plane client
    cp, err := client.ControlPlane()
    if err != nil {
        return fmt.Errorf("failed to get control plane client: %w", err)
    }

    // Parse UUIDs for required resources
    orgId := uuid.MustParse("your-org-uuid")
    sourceBucketId := uuid.MustParse("source-bucket-uuid")
    sourceMinioId := uuid.MustParse("source-minio-uuid")
    icebergCatalogId := uuid.MustParse("iceberg-catalog-uuid")
    trinoDataDockId := uuid.MustParse("trino-datadock-uuid")
    destBucketId := uuid.MustParse("dest-bucket-uuid")
    destS3DdId := uuid.MustParse("dest-s3-datadock-uuid")

    // Optional parameters
    destPrefix := "sorted/"
    sourceFolder := "inbox/"
    defaultCategory := "inconnu"
    prefixRegex := "(condominium-\\d+)/"

    // Define classification labels
    labels := []controlplaneapiclient.LabelDefinition{
        {
            Category:    "Factures/Eau",
            Description: "Water utility bills and invoices",
        },
        {
            Category:    "Factures/Electricite",
            Description: "Electricity utility bills and invoices",
        },
        {
            Category:    "Contrats/Entretien",
            Description: "Maintenance contracts and service agreements",
        },
        {
            Category:    "Contrats/Assurance",
            Description: "Insurance policies and contracts",
        },
    }

    // Configure the router
    router := controlplaneapiclient.CategoryRouterConfig{
        DefaultCategory: &defaultCategory,
        PrefixRegex:     &prefixRegex,
    }

    // Create output parameters with filesorter type
    outputParams := controlplaneapiclient.PipelineOutputParameters2{
        Type:                   controlplaneapiclient.Filesorter,
        DcIcebergId:           types.UUID(icebergCatalogId),
        DdTrinoInt:            types.UUID(trinoDataDockId),
        DestinationBucketId:   (*types.UUID)(&destBucketId),
        DestinationPrefix:     &destPrefix,
        DestinationS3Dd:       (*types.UUID)(&destS3DdId),
        Labels:                labels,
        ModelApiKey:           "your-mistral-api-key",
        ModelName:             "mistral-large-latest",
        Router:                router,
        TrinoSchema:           "document_metadata",
        TrinoTable:            "classified_documents",
    }

    // Marshal output parameters as union type
    output := controlplaneapiclient.PipelineOutputParameters{}
    if err := output.FromPipelineOutputParameters2(outputParams); err != nil {
        return fmt.Errorf("failed to create output parameters: %w", err)
    }

    // Create input parameters
    input := controlplaneapiclient.PipelineInputParameters{
        DcBucketId: types.UUID(sourceBucketId),
        DdMinio:    types.UUID(sourceMinioId),
        Folder:     &sourceFolder,
    }

    // Create pipeline parameters
    enabled := true // Enable the pipeline (default behavior)
    pipeline := controlplaneapiclient.PipelineParameters{
        Name:    "invoice-classifier-pipeline",
        OneOff:  false, // Continuous processing
        OrgId:   types.UUID(orgId),
        OrgSlug: "your-org-slug",
        Type:    "filesorter",
        Enabled: &enabled, // Optional: control whether the CronJob is active
    }

    // Create the pipeline request
    request := controlplaneapiclient.CreatePipelineJSONRequestBody{
        Input:    input,
        Output:   output,
        Pipeline: pipeline,
    }

    // Submit the pipeline creation request
    resp, err := cp.CreatePipelineWithResponse(ctx, request)
    if err != nil {
        return fmt.Errorf("failed to create pipeline: %w", err)
    }

    if resp.StatusCode() != 200 && resp.StatusCode() != 201 {
        return fmt.Errorf("pipeline creation failed with status %d", resp.StatusCode())
    }

    fmt.Println("Filesorter pipeline created successfully!")
    return nil
}

func main() {
    ctx := context.Background()
    if err := createFileSorterPipeline(ctx); err != nil {
        log.Fatal(err)
    }
}
```

## Common Label Categories

Here are some common document categories you might use:

### Financial Documents
```go
{Category: "Factures/Eau", Description: "Water utility bills"}
{Category: "Factures/Electricite", Description: "Electricity bills"}
{Category: "Factures/Gaz", Description: "Gas utility bills"}
{Category: "Factures/Telecom", Description: "Telephone and internet bills"}
```

### Legal Documents
```go
{Category: "Contrats/Bail", Description: "Rental lease agreements"}
{Category: "Contrats/Assurance", Description: "Insurance contracts"}
{Category: "Contrats/Entretien", Description: "Maintenance service agreements"}
```

### Administrative
```go
{Category: "Courrier/Administration", Description: "Administrative correspondence"}
{Category: "Documents/Identite", Description: "Identity documents"}
{Category: "Documents/Fiscaux", Description: "Tax documents"}
```

## Best Practices

1. **Label Design**: Create clear, specific label descriptions that help the AI model distinguish between categories
2. **Model Selection**: Use `mistral-large-latest` for best accuracy, or `mistral-medium` for cost optimization
3. **Prefix Regex**: Use capture groups to preserve organizational structure (e.g., apartment numbers, project codes)
4. **Default Category**: Always provide a default category for unclassifiable documents
5. **Metadata Storage**: Use a dedicated Trino schema for each pipeline to avoid table conflicts
6. **One-Off vs Continuous**: Set `OneOff: false` for continuous monitoring of new files
7. **Pipeline Control**: Use the `Enabled` field to suspend/resume CronJob pipelines without deleting them. Set to `false` to temporarily disable processing.

## Troubleshooting

### Pipeline Creation Fails
- Verify all UUID references exist and are accessible
- Check that the service account has permissions for all referenced resources
- Ensure the Mistral API key is valid

### Classification Not Working
- Review label descriptions - they should be detailed enough for the AI to understand
- Consider using a more powerful model (mistral-large-latest)
- Check that input files are valid PDFs

### Files Not Being Routed
- Verify the CategoryRouterConfig is correctly configured
- Check that destination S3 Data Dock and bucket exist
- Review the prefix_regex pattern matches your source path structure

### Pipeline Not Running
- Check if the pipeline is enabled using the `Enabled` field in PipelineParameters
- For CronJob pipelines, verify the `Enabled` field is set to `true` (or omitted for default enabled state)
- One-off pipelines ignore the `Enabled` field

## API Reference

See [client.gen.go:458](/home/mmoreiradj/git/hyperfluid-sdk-go/sdk/controlplaneapiclient/client.gen.go:458) for the complete `FileSorterOutputParameters` structure.

Pipeline type constant: [client.gen.go:111](/home/mmoreiradj/git/hyperfluid-sdk-go/sdk/controlplaneapiclient/client.gen.go:111)

```go
const (
    Filesorter PipelineOutputParameters2Type = "filesorter"
)
```
