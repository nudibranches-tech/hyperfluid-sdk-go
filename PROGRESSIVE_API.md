# Progressive Fluent API - Type-Safe Navigation

## Concept

The Progressive Fluent API uses **typed builders** to provide:

1. **Type Safety**: Each level returns a different type with specific methods
2. **IDE Autocomplete**: Your IDE knows exactly what methods are available at each level
3. **Forced Order**: You MUST navigate in the correct order: Org → Harbor → DataDock → Catalog → Schema → Table
4. **Contextual Methods**: Each level has operations specific to that resource

## Architecture Hierarchy

```
Organization (Org)
    └── Harbor
        └── DataDock
            └── Catalog
                └── Schema
                    └── Table
```

## API Design

### Level 1: Organization (`OrgBuilder`)

**Navigation:**
- `Harbor(id)` → HarborBuilder

**Operations:**
- `ListHarbors(ctx)` → List all harbors in org
- `CreateHarbor(ctx, name)` → Create new harbor
- `ListDataDocks(ctx)` → List all datadocks across all harbors
- `RefreshAllDataDocks(ctx)` → Trigger refresh on all datadocks

**Example:**
```go
// List harbors
harbors, err := client.Org(orgID).ListHarbors(ctx)

// Create harbor
resp, err := client.Org(orgID).CreateHarbor(ctx, "my-harbor")
```

---

### Level 2: Harbor (`HarborBuilder`)

**Navigation:**
- `DataDock(id)` → DataDockBuilder

**Operations:**
- `ListDataDocks(ctx)` → List datadocks in this harbor
- `CreateDataDock(ctx, config)` → Create new datadock
- `Delete(ctx)` → Delete this harbor

**Example:**
```go
// List datadocks in harbor
datadocks, err := client.
    Org(orgID).
    Harbor(harborID).
    ListDataDocks(ctx)

// Create datadock
config := map[string]interface{}{
    "name": "postgres-prod",
    "connection_kind": map[string]interface{}{
        "Trino": map[string]interface{}{
            "host": "postgres.example.com",
            "port": 5432,
        },
    },
}
resp, err := client.
    Org(orgID).
    Harbor(harborID).
    CreateDataDock(ctx, config)
```

---

### Level 3: DataDock (`DataDockBuilder`)

**Navigation:**
- `Catalog(name)` → CatalogBuilder

**Operations:**
- `GetCatalog(ctx)` → Get full catalog metadata
- `RefreshCatalog(ctx)` → Trigger catalog introspection
- `WakeUp(ctx)` → Bring datadock online
- `Sleep(ctx)` → Put datadock to sleep
- `Get(ctx)` → Get datadock details
- `Update(ctx, config)` → Update configuration
- `Delete(ctx)` → Delete datadock

**Example:**
```go
datadock := client.
    Org(orgID).
    Harbor(harborID).
    DataDock(dataDockID)

// Get catalog metadata
catalog, err := datadock.GetCatalog(ctx)

// Refresh catalog
resp, err := datadock.RefreshCatalog(ctx)

// Lifecycle management
resp, err := datadock.WakeUp(ctx)
resp, err := datadock.Sleep(ctx)

// Get details
details, err := datadock.Get(ctx)

// Update
config := map[string]interface{}{"description": "Updated description"}
resp, err := datadock.Update(ctx, config)
```

---

### Level 4: Catalog (`CatalogBuilder`)

**Navigation:**
- `Schema(name)` → SchemaBuilder

**Operations:**
- `ListSchemas(ctx)` → List all schemas in catalog

**Example:**
```go
// List schemas
schemas, err := client.
    Org(orgID).
    Harbor(harborID).
    DataDock(dataDockID).
    Catalog("my_catalog").
    ListSchemas(ctx)

// Navigate to schema
schemaBuilder := client.
    Org(orgID).
    Harbor(harborID).
    DataDock(dataDockID).
    Catalog("my_catalog").
    Schema("public")
```

---

### Level 5: Schema (`SchemaBuilder`)

**Navigation:**
- `Table(name)` → TableQueryBuilder

**Operations:**
- `ListTables(ctx)` → List all tables in schema

**Example:**
```go
// List tables
tables, err := client.
    Org(orgID).
    Harbor(harborID).
    DataDock(dataDockID).
    Catalog("my_catalog").
    Schema("public").
    ListTables(ctx)

// Navigate to table
tableBuilder := client.
    Org(orgID).
    Harbor(harborID).
    DataDock(dataDockID).
    Catalog("my_catalog").
    Schema("public").
    Table("users")
```

---

### Level 6: Table (`TableQueryBuilder`)

**Query Building:**
- `Select(columns...)` → Add SELECT columns
- `Where(column, operator, value)` → Add filter
- `OrderBy(column, direction)` → Add sorting
- `Limit(n)` → Set limit
- `Offset(n)` → Set offset
- `RawParams(params)` → Add custom params

**Execution:**
- `Get(ctx)` → Execute query and get results

**Example:**
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

// Complex query
resp, err := client.
    Org(orgID).
    Harbor(harborID).
    DataDock(dataDockID).
    Catalog("postgres").
    Schema("public").
    Table("orders").
    Select("id", "customer_name", "total").
    Where("status", "=", "completed").
    Where("total", ">", 1000).
    OrderBy("created_at", "DESC").
    Limit(100).
    Offset(0).
    Get(ctx)
```

---

## Complete Example

```go
package main

import (
    "github.com/nudibranches-tech/bifrost-hyperfluid-sdk-dev/sdk"
    "github.com/nudibranches-tech/bifrost-hyperfluid-sdk-dev/sdk/utils"
    "context"
    "fmt"
)

func main() {
    // Configure client
    config := utils.Configuration{
        BaseURL: "https://bifrost.hyperfluid.cloud",
        OrgID:   "your-org-id",
        Token:   "your-token",
    }
    client := sdk.NewClient(config)
    ctx := context.Background()

    // Navigate the hierarchy
    org := client.Org("my-org-id")
    harbor := org.Harbor("my-harbor-id")
    datadock := harbor.DataDock("my-datadock-id")
    catalog := datadock.Catalog("postgres")
    schema := catalog.Schema("public")
    table := schema.Table("users")

    // Execute query
    resp, err := table.
        Select("id", "name", "email").
        Where("active", "=", true).
        Limit(50).
        Get(ctx)

    if err != nil {
        panic(err)
    }

    fmt.Printf("Status: %s\n", resp.Status)
    fmt.Printf("Data: %v\n", resp.Data)
}
```

---

## Comparison: Progressive vs Catalog-First APIs

### Progressive API (New - Type-Safe)
```go
// Forces hierarchy: Org → Harbor → DataDock → Catalog → Schema → Table
resp, err := client.
    Org(orgID).                    // OrgBuilder
    Harbor(harborID).              // HarborBuilder
    DataDock(dataDockID).          // DataDockBuilder
    Catalog(catalogName).          // CatalogBuilder
    Schema(schemaName).            // SchemaBuilder
    Table(tableName).              // TableQueryBuilder
    Limit(10).
    Get(ctx)
```

**Pros:**
- ✅ Type-safe: Each level is a different type
- ✅ IDE autocomplete shows only valid methods
- ✅ Contextual operations (WakeUp, RefreshCatalog, etc.)
- ✅ Listing at each level (ListHarbors, ListSchemas, etc.)
- ✅ Forces correct navigation order

**Cons:**
- ❌ More verbose (must specify full path)
- ❌ Requires knowing IDs for intermediate levels

---

### Catalog-First API (Legacy - Simple)
```go
// Direct to table via catalog/schema/table
resp, err := client.
    Catalog(catalogName).          // QueryBuilder
    Schema(schemaName).            // QueryBuilder
    Table(tableName).              // QueryBuilder
    Limit(10).
    Get(ctx)
```

**Pros:**
- ✅ Simple and concise
- ✅ Good for direct table queries
- ✅ No need to know org/harbor/datadock IDs

**Cons:**
- ❌ No access to org/harbor/datadock operations
- ❌ Can't list resources
- ❌ No lifecycle management (WakeUp, Sleep, etc.)

---

## When to Use Each

### Use Progressive API When:
- Managing resources (creating harbors, datadocks)
- Listing resources (harbors, schemas, tables)
- DataDock lifecycle operations (WakeUp, Sleep, RefreshCatalog)
- You have full path information (org, harbor, datadock IDs)
- Building admin/management tools

### Use Catalog-First API When:
- Simple data queries
- You only care about catalog/schema/table
- Quick scripts and one-liners
- You don't need resource management

---

## Benefits of Progressive API

### 1. Type Safety
```go
// ❌ This won't compile:
client.Org(orgID).Sleep(ctx)  // Sleep() not available on OrgBuilder!

// ✅ This works:
client.Org(orgID).Harbor(harborID).DataDock(dataDockID).Sleep(ctx)
```

### 2. IDE Autocomplete
When you type `client.Org(orgID).`, your IDE shows:
- `Harbor(id)`
- `ListHarbors(ctx)`
- `CreateHarbor(ctx, name)`
- `ListDataDocks(ctx)`
- `RefreshAllDataDocks(ctx)`

### 3. Discoverable APIs
No need to read docs - just follow the types!

```go
org := client.Org(orgID)
// IDE: What can I do with org?
//   - Harbor()
//   - ListHarbors()
//   - CreateHarbor()
//   - etc.

harbor := org.Harbor(harborID)
// IDE: What can I do with harbor?
//   - DataDock()
//   - ListDataDocks()
//   - CreateDataDock()
//   - Delete()

datadock := harbor.DataDock(dataDockID)
// IDE: What can I do with datadock?
//   - Catalog()
//   - GetCatalog()
//   - RefreshCatalog()
//   - WakeUp()
//   - Sleep()
//   - etc.
```

### 4. Resource Management
```go
// List all resources
harbors, _ := client.Org(orgID).ListHarbors(ctx)
datadocks, _ := client.Org(orgID).Harbor(harborID).ListDataDocks(ctx)
schemas, _ := datadock.Catalog("postgres").ListSchemas(ctx)
tables, _ := schema.ListTables(ctx)

// Lifecycle operations
datadock.RefreshCatalog(ctx)  // Update schema metadata
datadock.WakeUp(ctx)          // Bring online
datadock.Sleep(ctx)           // Save costs

// Create resources
client.Org(orgID).CreateHarbor(ctx, "new-harbor")
harbor.CreateDataDock(ctx, datadockConfig)
```

---

## Migration Guide

If you're using the old catalog-first API and want to migrate:

### Before (Catalog-First):
```go
resp, err := client.
    Catalog("postgres").
    Schema("public").
    Table("users").
    Limit(10).
    Get(ctx)
```

### After (Progressive):
```go
// Option 1: Full path (recommended)
resp, err := client.
    Org(orgID).
    Harbor(harborID).
    DataDock(dataDockID).
    Catalog("postgres").
    Schema("public").
    Table("users").
    Limit(10).
    Get(ctx)

// Option 2: Reuse builders
datadock := client.Org(orgID).Harbor(harborID).DataDock(dataDockID)

// Now use datadock for multiple queries
users, _ := datadock.Catalog("postgres").Schema("public").Table("users").Get(ctx)
orders, _ := datadock.Catalog("postgres").Schema("public").Table("orders").Get(ctx)
```

**Note:** The catalog-first API still works! No breaking changes.

---

## Advanced Patterns

### Reusing Builders
```go
// Create reusable builders
org := client.Org(orgID)
harbor := org.Harbor(harborID)
datadock := harbor.DataDock(dataDockID)
catalog := datadock.Catalog("postgres")
schema := catalog.Schema("public")

// Use them multiple times
users := schema.Table("users").Limit(10).Get(ctx)
orders := schema.Table("orders").Where("status", "=", "pending").Get(ctx)
products := schema.Table("products").OrderBy("name", "ASC").Get(ctx)
```

### Dynamic Navigation
```go
func queryTable(orgID, harborID, dataDockID, catalog, schema, table string) (*utils.Response, error) {
    return client.
        Org(orgID).
        Harbor(harborID).
        DataDock(dataDockID).
        Catalog(catalog).
        Schema(schema).
        Table(table).
        Get(context.Background())
}
```

### Resource Discovery
```go
// Discover all resources in an org
org := client.Org(orgID)

harbors, _ := org.ListHarbors(ctx)
for _, harbor := range harbors {
    harborID := harbor["id"].(string)
    datadocks, _ := org.Harbor(harborID).ListDataDocks(ctx)

    for _, datadock := range datadocks {
        datadockID := datadock["id"].(string)
        catalog, _ := org.Harbor(harborID).DataDock(datadockID).GetCatalog(ctx)
        fmt.Printf("Catalog: %v\n", catalog)
    }
}
```

---

## Summary

The Progressive Fluent API provides:
- ✅ **Type safety** with distinct types for each level
- ✅ **IDE support** with accurate autocomplete
- ✅ **Resource management** operations at each level
- ✅ **Forced navigation** order for correctness
- ✅ **Contextual methods** specific to each resource type
- ✅ **Backward compatible** with catalog-first API

Perfect for building robust, maintainable applications with excellent developer experience! 🚀
