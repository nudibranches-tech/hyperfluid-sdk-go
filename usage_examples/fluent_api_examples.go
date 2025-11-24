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

	// NEW FLUENT API - Simple and intuitive!
	resp, err := client.
		Catalog(testCatalog).
		Schema(testSchema).
		Table(testTable).
		Limit(5).
		Get(context.Background())

	handleResponse(resp, err)
	fmt.Println()
}

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

// Helper functions

func handleResponse(resp *utils.Response, err error) {
	if err != nil {
		fmt.Printf("❌ Error: %s\n", err.Error())
		return
	}
	if resp.Status != utils.StatusOK {
		fmt.Printf("❌ Error: %s\n", resp.Error)
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

func getConfig() utils.Configuration {
	return utils.Configuration{
		BaseURL:        getEnv("HYPERFLUID_BASE_URL", ""),
		OrgID:          getEnv("HYPERFLUID_ORG_ID", ""),
		Token:          getEnv("HYPERFLUID_TOKEN", ""),
		RequestTimeout: 30 * time.Second,
		MaxRetries:     3,

		KeycloakBaseURL:      getEnv("KEYCLOAK_BASE_URL", ""),
		KeycloakRealm:        getEnv("KEYCLOAK_REALM", ""),
		KeycloakClientID:     getEnv("KEYCLOAK_CLIENT_ID", ""),
		KeycloakClientSecret: getEnv("KEYCLOAK_CLIENT_SECRET", ""),
		KeycloakUsername:     getEnv("KEYCLOAK_USERNAME", ""),
		KeycloakPassword:     getEnv("KEYCLOAK_PASSWORD", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

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
