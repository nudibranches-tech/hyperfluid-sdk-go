package main

import (
	"bifrost-for-developers/sdk"
	"bifrost-for-developers/sdk/utils"
	"context"
	"fmt"
	"net/url"
	"os"
	"time"
)

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

func runPostgresExample() {
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("🎯 Example 1: PostgreSQL Query")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

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

	table := client.GetCatalog(testCatalog).Table(testSchema, testTable)

	params := url.Values{}
	params.Add("_limit", "5")

	fmt.Printf("📝 GET /%s/%s/%s?_limit=5\n", testCatalog, testSchema, testTable)

	resp, err := table.GetData(context.Background(), params)
	handleResponse(resp, err)
	fmt.Println()
}

func runGraphQLExample() {
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("🎯 Example 2: GraphQL Query")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("⚠️  Note: GraphQL is not yet supported in the current SDK version")
	fmt.Println("         Use the REST API via runOpenAPIExample() instead")
	fmt.Println()
}

func runOpenAPIExample() {
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("🎯 Example 3: OpenAPI (REST) Query")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	config := getConfig()
	client := sdk.NewClient(config)

	testCatalog := getEnv("BIFROST_TEST_CATALOG", "")
	testSchema := getEnv("BIFROST_TEST_SCHEMA", "")
	testTable := getEnv("BIFROST_TEST_TABLE", "")
	testColumns := getEnv("BIFROST_TEST_COLUMNS", "*")

	if testCatalog == "" || testSchema == "" || testTable == "" {
		fmt.Println("⚠️  Skipping: BIFROST_TEST_CATALOG, BIFROST_TEST_SCHEMA, or BIFROST_TEST_TABLE not set")
		fmt.Println()
		return
	}

	table := client.GetCatalog(testCatalog).Table(testSchema, testTable)

	params := url.Values{}
	params.Add("_limit", "10")
	if testColumns != "" && testColumns != "*" {
		params.Add("select", testColumns)
	}

	fmt.Printf("📝 GET /%s/%s/%s?_limit=10&select=%s\n", testCatalog, testSchema, testTable, testColumns)

	resp, err := table.GetData(context.Background(), params)
	handleResponse(resp, err)
	fmt.Println()
}
