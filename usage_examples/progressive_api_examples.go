package main

import (
	"github.com/nudibranches-tech/bifrost-hyperfluid-sdk-dev/sdk"
	"context"
	"fmt"
)

// This file demonstrates the NEW progressive fluent API
// Each level has its own type with contextual methods!

func runProgressiveAPIExample1() {
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("🎯 Progressive API Example 1: List Harbors")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	config := getConfig()
	client := sdk.NewClient(config)

	orgID := getEnv("HYPERFLUID_ORG_ID", "")
	if orgID == "" {
		fmt.Println("⚠️  Skipping: HYPERFLUID_ORG_ID not set")
		fmt.Println()
		return
	}

	fmt.Println("📝 Type-safe navigation:")
	fmt.Println("   client.Org(orgID).ListHarbors(ctx)")
	fmt.Println()

	// NEW! Each level has its own type with specific methods
	resp, err := client.
		Org(orgID).                       // Returns OrgBuilder with org-specific methods
		ListHarbors(context.Background()) // Only available on OrgBuilder!

	handleResponse(resp, err)
	fmt.Println()
}

func runProgressiveAPIExample2() {
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("🎯 Progressive API Example 2: List DataDocks")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	config := getConfig()
	client := sdk.NewClient(config)

	orgID := getEnv("HYPERFLUID_ORG_ID", "")
	harborID := getEnv("HYPERFLUID_HARBOR_ID", "")

	if orgID == "" || harborID == "" {
		fmt.Println("⚠️  Skipping: HYPERFLUID_ORG_ID or HYPERFLUID_HARBOR_ID not set")
		fmt.Println()
		return
	}

	fmt.Println("📝 Type-safe navigation:")
	fmt.Println("   client.Org(orgID).Harbor(harborID).ListDataDocks(ctx)")
	fmt.Println()

	resp, err := client.
		Org(orgID).
		Harbor(harborID).                   // Returns HarborBuilder with harbor-specific methods
		ListDataDocks(context.Background()) // Only available on HarborBuilder!

	handleResponse(resp, err)
	fmt.Println()
}

func runProgressiveAPIExample3() {
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("🎯 Progressive API Example 3: Get DataDock Catalog")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	config := getConfig()
	client := sdk.NewClient(config)

	orgID := getEnv("HYPERFLUID_ORG_ID", "")
	harborID := getEnv("HYPERFLUID_HARBOR_ID", "")
	dataDockID := getEnv("HYPERFLUID_DATADOCK_ID", "")

	if orgID == "" || harborID == "" || dataDockID == "" {
		fmt.Println("⚠️  Skipping: Required env vars not set")
		fmt.Println()
		return
	}

	fmt.Println("📝 Type-safe navigation:")
	fmt.Println("   client.Org(orgID).Harbor(harborID).DataDock(dataDockID).GetCatalog(ctx)")
	fmt.Println()

	resp, err := client.
		Org(orgID).
		Harbor(harborID).
		DataDock(dataDockID).            // Returns DataDockBuilder with datadock-specific methods
		GetCatalog(context.Background()) // GetCatalog, RefreshCatalog, WakeUp, Sleep only on DataDockBuilder!

	handleResponse(resp, err)
	fmt.Println()
}

func runProgressiveAPIExample4() {
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("🎯 Progressive API Example 4: DataDock Operations")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	config := getConfig()
	client := sdk.NewClient(config)

	orgID := getEnv("HYPERFLUID_ORG_ID", "")
	harborID := getEnv("HYPERFLUID_HARBOR_ID", "")
	dataDockID := getEnv("HYPERFLUID_DATADOCK_ID", "")

	if orgID == "" || harborID == "" || dataDockID == "" {
		fmt.Println("⚠️  Skipping: Required env vars not set")
		fmt.Println()
		return
	}

	datadock := client.Org(orgID).Harbor(harborID).DataDock(dataDockID)

	fmt.Println("📝 Available operations on DataDockBuilder:")
	fmt.Println("   - GetCatalog()")
	fmt.Println("   - RefreshCatalog()")
	fmt.Println("   - WakeUp()")
	fmt.Println("   - Sleep()")
	fmt.Println("   - Get()")
	fmt.Println("   - Update()")
	fmt.Println("   - Delete()")
	fmt.Println()

	// Example: Refresh catalog
	fmt.Println("🔄 Refreshing catalog...")
	resp, err := datadock.RefreshCatalog(context.Background())
	handleResponse(resp, err)
	fmt.Println()
}

func runProgressiveAPIExample5() {
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("🎯 Progressive API Example 5: Full Navigation Path")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	config := getConfig()
	client := sdk.NewClient(config)

	orgID := getEnv("HYPERFLUID_ORG_ID", "")
	harborID := getEnv("HYPERFLUID_HARBOR_ID", "")
	dataDockID := getEnv("HYPERFLUID_DATADOCK_ID", "")
	catalogName := getEnv("BIFROST_TEST_CATALOG", "")
	schemaName := getEnv("BIFROST_TEST_SCHEMA", "")
	tableName := getEnv("BIFROST_TEST_TABLE", "")

	if orgID == "" || harborID == "" || dataDockID == "" ||
		catalogName == "" || schemaName == "" || tableName == "" {
		fmt.Println("⚠️  Skipping: Required env vars not set")
		fmt.Println()
		return
	}

	fmt.Println("📝 Complete type-safe path:")
	fmt.Println("   client")
	fmt.Println("     .Org(orgID)              → OrgBuilder")
	fmt.Println("     .Harbor(harborID)        → HarborBuilder")
	fmt.Println("     .DataDock(dataDockID)    → DataDockBuilder")
	fmt.Println("     .Catalog(catalogName)    → CatalogBuilder")
	fmt.Println("     .Schema(schemaName)      → SchemaBuilder")
	fmt.Println("     .Table(tableName)        → TableQueryBuilder")
	fmt.Println("     .Limit(10)")
	fmt.Println("     .Get(ctx)")
	fmt.Println()

	resp, err := client.
		Org(orgID).
		Harbor(harborID).
		DataDock(dataDockID).
		Catalog(catalogName).
		Schema(schemaName).
		Table(tableName).
		Limit(10).
		Get(context.Background())

	handleResponse(resp, err)
	fmt.Println()
}

func runProgressiveAPIExample6() {
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("🎯 Progressive API Example 6: Complex Query with Path")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	config := getConfig()
	client := sdk.NewClient(config)

	orgID := getEnv("HYPERFLUID_ORG_ID", "")
	harborID := getEnv("HYPERFLUID_HARBOR_ID", "")
	dataDockID := getEnv("HYPERFLUID_DATADOCK_ID", "")
	catalogName := getEnv("BIFROST_TEST_CATALOG", "")
	schemaName := getEnv("BIFROST_TEST_SCHEMA", "")
	tableName := getEnv("BIFROST_TEST_TABLE", "")

	if orgID == "" || harborID == "" || dataDockID == "" ||
		catalogName == "" || schemaName == "" || tableName == "" {
		fmt.Println("⚠️  Skipping: Required env vars not set")
		fmt.Println()
		return
	}

	fmt.Println("📝 Full path with query:")
	fmt.Println("   Navigate: Org → Harbor → DataDock → Catalog → Schema → Table")
	fmt.Println("   Query: Select, Where, OrderBy, Limit")
	fmt.Println()

	resp, err := client.
		Org(orgID).
		Harbor(harborID).
		DataDock(dataDockID).
		Catalog(catalogName).
		Schema(schemaName).
		Table(tableName).
		Select("id", "name", "created_at"). // Query methods
		Where("status", "=", "active").
		OrderBy("created_at", "DESC").
		Limit(20).
		Get(context.Background())

	handleResponse(resp, err)
	fmt.Println()
}

func runProgressiveAPIListingExample() {
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("🎯 Progressive API Example: Listing Resources")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	config := getConfig()
	client := sdk.NewClient(config)

	orgID := getEnv("HYPERFLUID_ORG_ID", "")
	harborID := getEnv("HYPERFLUID_HARBOR_ID", "")
	dataDockID := getEnv("HYPERFLUID_DATADOCK_ID", "")
	catalogName := getEnv("BIFROST_TEST_CATALOG", "")
	schemaName := getEnv("BIFROST_TEST_SCHEMA", "")

	if orgID == "" {
		fmt.Println("⚠️  Skipping: HYPERFLUID_ORG_ID not set")
		fmt.Println()
		return
	}

	fmt.Println("📋 Listing resources at each level:")
	fmt.Println()

	// List harbors in org
	fmt.Println("1. List Harbors in Organization:")
	if harbors, err := client.Org(orgID).ListHarbors(context.Background()); err == nil {
		fmt.Printf("   ✓ Found harbors\n")
		_ = harbors
	} else {
		fmt.Printf("   ✗ Error: %v\n", err)
	}

	if harborID == "" {
		fmt.Println("   (Set HYPERFLUID_HARBOR_ID to continue)")
		fmt.Println()
		return
	}

	// List datadocks in harbor
	fmt.Println("2. List DataDocks in Harbor:")
	if datadocks, err := client.Org(orgID).Harbor(harborID).ListDataDocks(context.Background()); err == nil {
		fmt.Printf("   ✓ Found datadocks\n")
		_ = datadocks
	} else {
		fmt.Printf("   ✗ Error: %v\n", err)
	}

	if dataDockID == "" || catalogName == "" {
		fmt.Println("   (Set HYPERFLUID_DATADOCK_ID and BIFROST_TEST_CATALOG to continue)")
		fmt.Println()
		return
	}

	// List schemas in catalog
	fmt.Println("3. List Schemas in Catalog:")
	if schemas, err := client.Org(orgID).Harbor(harborID).DataDock(dataDockID).
		Catalog(catalogName).ListSchemas(context.Background()); err == nil {
		fmt.Printf("   ✓ Found %d schemas: %v\n", len(schemas), schemas)
	} else {
		fmt.Printf("   ✗ Error: %v\n", err)
	}

	if schemaName == "" {
		fmt.Println("   (Set BIFROST_TEST_SCHEMA to continue)")
		fmt.Println()
		return
	}

	// List tables in schema
	fmt.Println("4. List Tables in Schema:")
	if tables, err := client.Org(orgID).Harbor(harborID).DataDock(dataDockID).
		Catalog(catalogName).Schema(schemaName).ListTables(context.Background()); err == nil {
		fmt.Printf("   ✓ Found %d tables: %v\n", len(tables), tables)
	} else {
		fmt.Printf("   ✗ Error: %v\n", err)
	}

	fmt.Println()
}
