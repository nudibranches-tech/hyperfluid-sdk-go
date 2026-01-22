# Bifrost SDK Usage Examples

This directory contains runnable examples demonstrating various features of the Bifrost SDK.

## Quick Start

### 1. Set up environment variables

Create a `.env` file in the project root with your credentials:

```bash
# Required for most examples
HYPERFLUID_SERVICE_ACCOUNT_FILE=/path/to/service_account.json
HYPERFLUID_API_URL=https://api.hyperfluid.cloud

# Optional - for Control Plane examples
HYPERFLUID_CONTROL_PLANE_URL=https://console.hyperfluid.cloud
HYPERFLUID_ORG_ID=your-org-uuid
HYPERFLUID_HARBOR_ID=your-harbor-uuid
HYPERFLUID_DATADOCK_ID=your-datadock-uuid

# Optional - for data query examples
BIFROST_TEST_CATALOG=your-catalog-name
BIFROST_TEST_SCHEMA=your-schema-name
BIFROST_TEST_TABLE=your-table-name
```

### 2. Run examples

```bash
# List available examples
go run ./usage_examples --list

# Run all examples
go run ./usage_examples --all

# Run specific examples
go run ./usage_examples --control-plane
go run ./usage_examples --fluent --search
go run ./usage_examples --service-account

# Or build and run
go build -o examples ./usage_examples
./examples --control-plane
```

## Available Examples

### 📊 Fluent API (`--fluent`)
Demonstrates fluent query building with method chaining.

**Required env vars:**
- `HYPERFLUID_SERVICE_ACCOUNT_FILE`
- `HYPERFLUID_API_URL`
- `BIFROST_TEST_CATALOG`

**Example:**
```bash
go run ./usage_examples --fluent
```

### 📦 S3 Operations (`--s3`)
S3 bucket operations including list, upload, and download.

**Required env vars:**
- `HYPERFLUID_SERVICE_ACCOUNT_FILE`
- `HYPERFLUID_API_URL`

**Example:**
```bash
go run ./usage_examples --s3
```

### 🔍 Search API (`--search`)
Full-text search across data containers.

**Required env vars:**
- `HYPERFLUID_SERVICE_ACCOUNT_FILE`
- `HYPERFLUID_API_URL`

**Example:**
```bash
go run ./usage_examples --search
```

### 🔐 Service Accounts (`--service-account`)
Different patterns for loading and using service accounts.

**Required env vars:**
- `HYPERFLUID_SERVICE_ACCOUNT_FILE`
- `HYPERFLUID_API_URL`

**Example:**
```bash
go run ./usage_examples --service-account
```

### 🎛️ Control Plane API (`--control-plane`)
Administrative operations for data docks and archive operations.

**Features:**
- List archive operations (imports/exports)
- View operation status, file details, and timestamps
- Track import/export progress

**Required env vars:**
- `HYPERFLUID_SERVICE_ACCOUNT_FILE`
- `HYPERFLUID_API_URL`
- `HYPERFLUID_DATA_DOCK_ID`
- `HYPERFLUID_DATA_CONTAINER_ID`

**Example:**
```bash
go run ./usage_examples --control-plane
```

**Code snippet:**
```go
// Create SDK client
client, err := sdk.NewClientFromServiceAccountFile(
    "/path/to/sa.json",
    sdk.ServiceAccountOptions{
        BaseURL:         "https://api.hyperfluid.cloud",
        ControlPlaneURL: "https://console.hyperfluid.cloud",
    },
)

// Get Control Plane client (lazy-init, cached, automatic OAuth2)
cp, err := client.ControlPlane()

// List archive operations
dataDockUUID, _ := uuid.Parse("your-datadock-id")
dataContainerUUID, _ := uuid.Parse("your-container-id")

resp, err := cp.ListArchiveOperationsWithResponse(
    ctx,
    dataDockUUID,
    dataContainerUUID,
    &controlplaneapiclient.ListArchiveOperationsParams{},
)

for _, op := range *resp.JSON200 {
    fmt.Printf("Operation: %s - Status: %s\n", op.OperationType, op.Status)
}
```

## Command-Line Options

```
-all                    Run all examples
-fluent                 Run fluent API examples
-s3                     Run S3 examples
-search                 Run search examples
-service-account        Run service account examples
-control-plane          Run Control Plane API examples
-list                   List available examples
-skip-tls-verify        Skip TLS certificate verification (WARNING: Use only for development)
-help                   Show help message
```

### Development Mode

For development environments with self-signed certificates, you can disable TLS verification:

```bash
go run ./usage_examples --control-plane --skip-tls-verify
```

**⚠️ WARNING:** Never use `--skip-tls-verify` in production environments. This flag disables certificate validation and makes your connection vulnerable to man-in-the-middle attacks.

## Service Account Format

Your service account JSON should look like this:

```json
{
  "client_id": "hf-org-sa-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
  "client_secret": "your-client-secret",
  "issuer": "https://auth.hyperfluid.cloud/realms/your-org",
  "auth_uri": "https://auth.hyperfluid.cloud/realms/your-org/protocol/openid-connect/auth",
  "token_uri": "https://auth.hyperfluid.cloud/realms/your-org/protocol/openid-connect/token"
}
```

## Files

- `main.go` - Entry point with command-line flag handling
- `fluent_api_examples.go` - Fluent query API examples
- `service_account_examples.go` - Service account authentication examples
- `control_plane_examples.go` - Control Plane API examples
- `progressive_api_examples.go` - Progressive builder examples

## Troubleshooting

### Examples are skipped
If you see "⚠️ Skipping" messages, it means the required environment variables are not set. Check your `.env` file and ensure all required variables are configured.

### Authentication errors
- Verify your service account file path is correct
- Ensure the service account has the necessary permissions
- Check that the API URLs are correct

### Control Plane examples fail
- Make sure `HYPERFLUID_ORG_ID` is set to a valid organization UUID
- Verify your service account has permissions to access the Control Plane API
- Check that `HYPERFLUID_CONTROL_PLANE_URL` is set (or defaults to `HYPERFLUID_API_URL`)
