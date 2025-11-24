# Bifrost SDK

Go SDK for Hyperfluid data access with **two fluent APIs**:
1. **Progressive API** (Type-safe, resource management) - **NEW!** ✨
2. **Catalog-First API** (Simple queries)

## Quick Start

```bash
# Install
go get github.com/nudibranches-tech/bifrost-hyperfluid-sdk-dev/sdk
```

## Two APIs, Choose Your Style

### 🚀 Progressive API (Type-Safe Navigation)

Navigate the resource hierarchy with **type safety** and **contextual operations**:

```go
// Full path: Org → Harbor → DataDock → Catalog → Schema → Table
resp, err := client.
    Org(orgID).                    // OrgBuilder - can ListHarbors(), CreateHarbor()
    Harbor(harborID).              // HarborBuilder - can ListDataDocks()
    DataDock(dataDockID).          // DataDockBuilder - can WakeUp(), Sleep(), RefreshCatalog()
    Catalog(catalogName).          // CatalogBuilder - can ListSchemas()
    Schema(schemaName).            // SchemaBuilder - can ListTables()
    Table(tableName).              // TableQueryBuilder - can query
    Limit(10).
    Get(ctx)
```

**Benefits:**
- ✅ Type-safe: Each level is a different type
- ✅ IDE autocomplete shows only valid methods
- ✅ Resource management (WakeUp, Sleep, RefreshCatalog)
- ✅ Listing at each level (ListHarbors, ListSchemas, etc.)

**See:** [PROGRESSIVE_API.md](PROGRESSIVE_API.md) for complete documentation

---

### 📦 Catalog-First API (Simple Queries)

Jump directly to tables when you just need data:

```go
import (
    "context"
    "fmt"
    "github.com/nudibranches-tech/bifrost-hyperfluid-sdk-dev/sdk"
    "github.com/nudibranches-tech/bifrost-hyperfluid-sdk-dev/sdk/utils"
)

func main() {
    // Configure the client
    config := utils.Configuration{
        BaseURL: "https://bifrost.hyperfluid.cloud",
        OrgID:   "your-org-id",
        Token:   "your-token",
        // or use Keycloak for token management
    }

    // Create a new client
    client := sdk.NewClient(config)

    // Simple query with fluent API
    resp, err := client.
        Catalog("sales").
        Schema("public").
        Table("orders").
        Limit(10).
        Get(context.Background())

    if err != nil {
        // Handle error
    }

    fmt.Println(resp.Data)
}
```

**Benefits:**
- ✅ Simple and concise
- ✅ Perfect for data queries
- ✅ No need for intermediate IDs

---

## Progressive API Examples

### Resource Management

```go
// List resources at each level
harbors, err := client.Org(orgID).ListHarbors(ctx)
datadocks, err := client.Org(orgID).Harbor(harborID).ListDataDocks(ctx)
schemas, err := datadock.Catalog("postgres").ListSchemas(ctx)
tables, err := schema.ListTables(ctx)

// Create resources
client.Org(orgID).CreateHarbor(ctx, "my-harbor")
harbor.CreateDataDock(ctx, datadockConfig)

// DataDock lifecycle
datadock := client.Org(orgID).Harbor(harborID).DataDock(dataDockID)
datadock.RefreshCatalog(ctx)  // Update metadata
datadock.WakeUp(ctx)          // Bring online
datadock.Sleep(ctx)           // Save costs
datadock.Update(ctx, config)  // Update config
```

### Queries with Full Path

```go
// Simple query
resp, err := client.
    Org(orgID).
    Harbor(harborID).
    DataDock(dataDockID).
    Catalog("postgres").
    Schema("public").
    Table("users").
    Limit(10).
    Get(ctx)

// Complex query with filters
resp, err := client.
    Org(orgID).
    Harbor(harborID).
    DataDock(dataDockID).
    Catalog("sales").
    Schema("public").
    Table("orders").
    Select("id", "customer", "total").
    Where("status", "=", "completed").
    Where("total", ">", 1000).
    OrderBy("created_at", "DESC").
    Limit(100).
    Get(ctx)
```

---

## Catalog-First API Examples

### Advanced Queries

```go
// Complex query with all features
resp, err := client.
    Catalog("sales").
    Schema("public").
    Table("orders").
    Select("id", "customer_name", "total_amount").
    Where("status", "=", "completed").
    Where("total_amount", ">", 1000).
    OrderBy("created_at", "DESC").
    Limit(100).
    Offset(0).
    Get(ctx)

// Override organization ID for specific query
resp, err := client.
    Org("different-org-id").
    Catalog("catalog").
    Schema("schema").
    Table("table").
    Get(ctx)

// Use raw parameters for advanced cases
resp, err := client.
    Catalog("catalog").
    Schema("schema").
    Table("table").
    RawParams(url.Values{"custom_param": {"value"}}).
    Get(ctx)

// Building queries step by step
query := client.
    Catalog("sales").
    Schema("public").
    Table("orders")

// Add filters dynamically
if status != "" {
    query = query.Where("status", "=", status)
}

// Add pagination
query = query.Limit(pageSize).Offset(page * pageSize)

// Execute
resp, err := query.Get(ctx)
```

## Configuration

### Required
- `HYPERFLUID_ORG_ID` - Your organization ID
- `HYPERFLUID_TOKEN` - API token (or use Keycloak)

### Optional
- `HYPERFLUID_BASE_URL` - API endpoint (default: `https://bifrost.hyperfluid.cloud`)

### Keycloak (alternative to token)
- `KEYCLOAK_BASE_URL` - Keycloak server
- `KEYCLOAK_REALM` - Realm name
- `KEYCLOAK_CLIENT_ID` - Client ID (required for both grant types)
- `KEYCLOAK_CLIENT_SECRET` - Client Secret (for Client Credentials Grant - preferred for services)
- `KEYCLOAK_USERNAME` - Your username (for Password Grant - fallback if Client Secret not provided)
- `KEYCLOAK_PASSWORD` - Your password (for Password Grant - fallback if Client Secret not provided)

**Note:** If `KEYCLOAK_CLIENT_SECRET` is provided, the SDK will prioritize the more secure Client Credentials Grant. Otherwise, it will fall back to the Password Grant if `KEYCLOAK_USERNAME` and `KEYCLOAK_PASSWORD` are configured.

## Project Structure

```
sdk/
  client.go        # Client object and entry points
  query_builder.go # Fluent API implementation
  request.go       # HTTP request handling
  auth.go          # Authentication (Keycloak support)
  utils/           # Utility functions and types
```

## Fluent API Methods

### Query Building Methods

- **`Catalog(name string)`** - Set the catalog name
- **`Schema(name string)`** - Set the schema name
- **`Table(name string)`** - Set the table name
- **`Org(orgID string)`** - Override the organization ID from config

### Query Parameter Methods

- **`Select(columns ...string)`** - Specify columns to retrieve
- **`Where(column, operator, value)`** - Add filter conditions
  - Supported operators: `=`, `>`, `<`, `>=`, `<=`, `!=`, `LIKE`, `IN`
- **`OrderBy(column, direction)`** - Add ordering (ASC/DESC)
- **`Limit(n int)`** - Set maximum rows to return
- **`Offset(n int)`** - Set number of rows to skip
- **`RawParams(url.Values)`** - Add custom query parameters

### Execution Methods

- **`Get(ctx)`** - Execute SELECT query and return results
- **`Count(ctx)`** - Get count of matching rows
- **`Post(ctx, data)`** - Insert new data
- **`Put(ctx, data)`** - Update existing data
- **`Delete(ctx)`** - Delete matching rows

## Error Handling

```go
resp, err := client.
    Catalog("catalog").
    Schema("schema").
    Table("table").
    Get(ctx)

if err != nil {
    // Check for specific error types
    if errors.Is(err, utils.ErrNotFound) {
        log.Println("Resource not found")
    } else if errors.Is(err, utils.ErrPermissionDenied) {
        log.Println("Permission denied")
    } else if errors.Is(err, utils.ErrAuthenticationFailed) {
        log.Println("Authentication failed")
    } else {
        log.Fatalf("Request failed: %v", err)
    }
}

if resp.Status != utils.StatusOK {
    log.Printf("API error: %s", resp.Error)
}
```

## License

Private SDK for internal use.
