package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/nudibranches-tech/bifrost-hyperfluid-sdk-dev/sdk"
	"github.com/nudibranches-tech/bifrost-hyperfluid-sdk-dev/sdk/utils"
)

// This file demonstrates the new fluent API for the Bifrost SDK.
// The fluent API provides a more intuitive and user-friendly way to interact with the SDK.

/*
	func runSearchExample() {
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Println("🎯 Search Example: Full-Text Search")
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

		config := getConfig()
		client := sdk.NewClient(config)

		projectID := getEnv("BIFROST_DATADOCK_ID", "")
		catalog := getEnv("BIFROST_TEST_CATALOG", "iceberg")
		schema := getEnv("BIFROST_TEST_SCHEMA", "public")
		table := getEnv("BIFROST_TEST_TABLE", "text_files")

		if projectID == "" {
			fmt.Println("⚠️  Skipping: BIFROST_DATADOCK_ID not set")
			fmt.Println()
			return
		}

		fmt.Println("📝 Full-text search query:")
		fmt.Printf("   DataDock: %s\n", projectID)
		fmt.Printf("   Searching in: %s.%s.%s\n", catalog, schema, table)
		fmt.Println()

		// Search for content
		resp, _ := client.Search().
			Query("rapport ventes").
			DataDock(projectID).           // ✅ Use the actual UUID variable
			Catalog(catalog).              // ✅ Use actual catalog
			Schema(schema).                // ✅ Use actual schema
			Table(table).                  // ✅ Use actual table
			Columns("content", "summary"). // Adjust columns based on your table
			Limit(10).
			Execute(context.Background())

		fmt.Println(resp.Results)
		fmt.Println()
	}
*/

/*
func runFluentAPISimpleExample() {
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("🎯 Fluent API Example 1: Simple Query")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	config := getConfig()
	client := sdk.NewClient(config)

	testCatalog := getEnv("BIFROST_TEST_CATALOG", "")
	testSchema := getEnv("BIFROST_TEST_SCHEMA", "")
	testTable := getEnv("BIFROST_TEST_TABLE", "")

	if testCatalog == "" || testSchema == "" || testTable == "" {
		fmt.Println("⚠️  Skipping: BIFROST_TEST_CATALOG, BIFROST_TEST_SCHEMA, or BIFROST_TEST_TABLE not set")
		fmt.Println()
		return
	}

	fmt.Printf("📝 Fluent query: client.Catalog(%q).Schema(%q).Table(%q).Limit(5).Get(ctx)\n",
		testCatalog, testSchema, testTable)

	dataDockID := getEnv("HYPERFLUID_DATADOCK_ID", "")

	resp, err := client.
		DataDock(dataDockID).
		Catalog(testCatalog).
		Schema(testSchema).
		Table(testTable).
		Limit(5).
		Get(context.Background())

	handleResponse(resp, err)
	fmt.Println()
}
*/

/*
func runFluentAPIWithSelectExample() {
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("🎯 Fluent API Example 2: Query with SELECT")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	config := getConfig()
	client := sdk.NewClient(config)

	testCatalog := getEnv("BIFROST_TEST_CATALOG", "")
	testSchema := getEnv("BIFROST_TEST_SCHEMA", "")
	testTable := getEnv("BIFROST_TEST_TABLE", "")
	testColumns := getEnv("BIFROST_TEST_COLUMNS", "")

	if testCatalog == "" || testSchema == "" || testTable == "" {
		fmt.Println("⚠️  Skipping: Test environment variables not set")
		fmt.Println()
		return
	}

	if testColumns == "" {
		fmt.Println("⚠️  Skipping: BIFROST_TEST_COLUMNS not set")
		fmt.Println()
		return
	}

	fmt.Printf("📝 Fluent query with SELECT: .Select(%q).Limit(10).Get(ctx)\n", testColumns)

	// Select specific columns (comma-separated string to variadic args)
	cols := splitColumns(testColumns)

	resp, err := client.
		Catalog(testCatalog).
		Schema(testSchema).
		Table(testTable).
		Select(cols...).
		Limit(10).
		Get(context.Background())

	handleResponse(resp, err)
	fmt.Println()
}
*/

/*
func runFluentAPIComplexExample() {
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("🎯 Fluent API Example 3: Complex Query")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	config := getConfig()
	client := sdk.NewClient(config)

	testCatalog := getEnv("BIFROST_TEST_CATALOG", "")
	testSchema := getEnv("BIFROST_TEST_SCHEMA", "")
	testTable := getEnv("BIFROST_TEST_TABLE", "")

	if testCatalog == "" || testSchema == "" || testTable == "" {
		fmt.Println("⚠️  Skipping: Test environment variables not set")
		fmt.Println()
		return
	}

	fmt.Println("📝 Complex fluent query with:")
	fmt.Println("   - Multiple SELECT columns")
	fmt.Println("   - WHERE filters")
	fmt.Println("   - ORDER BY")
	fmt.Println("   - Pagination (LIMIT + OFFSET)")

	// Complex query with all features
	resp, err := client.
		Catalog(testCatalog).
		Schema(testSchema).
		Table(testTable).
		Select("id", "name", "created_at").
		Where("status", "=", "active").
		Where("amount", ">", 100).
		OrderBy("created_at", "DESC").
		Limit(20).
		Offset(0).
		Get(context.Background())

	handleResponse(resp, err)
	fmt.Println()
}
*/
// Helper functions
/*
func handleResponse(resp *utils.Response, err error) {
	if err != nil {
		fmt.Printf("❌ Error: %s\n", err.Error())
		// Also show the server response if available
		if resp != nil && resp.Error != "" {
			fmt.Printf("   Server said: %s\n", resp.Error)
			fmt.Printf("   HTTP Status: %d\n", resp.HTTPCode)
		}
		return
	}
	if resp.Status != utils.StatusOK {
		fmt.Printf("❌ Error: %s\n", resp.Error)
		fmt.Printf("   HTTP Status: %d\n", resp.HTTPCode)
		return
	}
	fmt.Println("✅ Success!")
	if dataSlice, isSlice := resp.Data.([]interface{}); isSlice {
		fmt.Printf("📦 %d records", len(dataSlice))
		if len(dataSlice) > 0 {
			fmt.Printf(" | First: %v", dataSlice[0])
		}
		fmt.Println()
	} else if dataMap, isMap := resp.Data.(map[string]interface{}); isMap {
		fmt.Printf("📦 Data: %v\n", dataMap)
	}
}
*/

func runS3Example() {
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("🎯 S3 File Retrieval (SSO Auth)")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	config := getConfig()
	client := sdk.NewClient(config)

	bucket := getEnv("S3_BUCKET", "")
	fileKey := getEnv("S3_FILE_KEY", "example/file.txt")

	if bucket == "" {
		fmt.Println("⚠️  Skipping: S3_BUCKET not set")
		fmt.Println()
		return
	}

	if config.KeycloakClientID == "" || config.KeycloakClientSecret == "" {
		fmt.Println("⚠️  Skipping: Keycloak SSO credentials not configured")
		fmt.Println()
		return
	}

	fmt.Printf("🔐 Using SSO authentication (Keycloak)\n")
	fmt.Printf("📥 Retrieving: s3://%s/%s\n", bucket, fileKey)

	// Get file from S3 using SSO
	s3Builder, err := client.S3()
	if err != nil {
		fmt.Printf("❌ Failed to create S3 builder: %v\n", err)
		return
	}

	obj, err := s3Builder.
		Bucket(bucket).
		Key(fileKey).
		Get(context.Background())
	if err != nil {
		fmt.Printf("❌ Failed to get object: %v\n", err)
		return
	}

	// Ensure Body is closed and check error
	defer func() {
		if cerr := obj.Body.Close(); cerr != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to close S3 object body: %v\n", cerr)
		}
	}()

	fmt.Printf("📄 File size: %d bytes\n", obj.Size)
	fmt.Printf("📄 Content type: %s\n", obj.ContentType)
	if obj.LastModified != nil {
		fmt.Printf("📄 Last modified: %s\n", obj.LastModified.Format(time.RFC3339))
	}

	// Optionally read first 100 bytes of content
	preview := make([]byte, 100)
	n, _ := obj.Body.Read(preview)
	if n > 0 {
		fmt.Printf("📄 Preview: %s...\n", string(preview[:n]))
	}

	fmt.Println()
}

func getConfig() utils.Configuration {
	return utils.Configuration{
		BaseURL:        getEnv("HYPERFLUID_BASE_URL", ""),
		OrgID:          getEnv("HYPERFLUID_ORG_ID", ""),
		Token:          getEnv("HYPERFLUID_TOKEN", ""),
		DataDockID:     getEnv("HYPERFLUID_DATADOCK_ID", ""),
		RequestTimeout: 30 * time.Second,
		MaxRetries:     3,

		KeycloakBaseURL:      getEnv("KEYCLOAK_BASE_URL", ""),
		KeycloakRealm:        getEnv("KEYCLOAK_REALM", ""),
		KeycloakClientID:     getEnv("KEYCLOAK_CLIENT_ID", ""),
		KeycloakClientSecret: getEnv("KEYCLOAK_CLIENT_SECRET", ""),
		KeycloakUsername:     getEnv("KEYCLOAK_USERNAME", ""),
		KeycloakPassword:     getEnv("KEYCLOAK_PASSWORD", ""),
		MinIOEndpoint:        getEnv("MINIO_ENDPOINT", ""),
		MinIOAccessKey:       getEnv("MINIO_ACCESS_KEY", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

/*
func splitColumns(cols string) []string {
	if cols == "" {
		return []string{}
	}
	// Simple split by comma
	var result []string
	current := ""
	for _, c := range cols {
		if c == ',' {
			if current != "" {
				result = append(result, current)
				current = ""
			}
		} else if c != ' ' {
			current += string(c)
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}
*/
