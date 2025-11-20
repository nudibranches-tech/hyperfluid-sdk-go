package developpementtests

import (
	"bifrost-for-developers/sdk"
	"bifrost-for-developers/sdk/utils"
	"context"
	"net/url"
	"os"
	"testing"
)

func TestIntegration_GetData(t *testing.T) {
	if testing.Short() {
		t.Skip("⏭️  Skipping integration test in short mode")
	}

	config, err := getTestConfig(nil) // Pass nil to load from environment variables
	if err != nil {
		t.Fatalf("Failed to get test config: %v", err)
	}

	testCatalog := os.Getenv("BIFROST_TEST_CATALOG")
	testSchema := os.Getenv("BIFROST_TEST_SCHEMA")
	testTable := os.Getenv("BIFROST_TEST_TABLE")

	if testCatalog == "" || testSchema == "" || testTable == "" {
		t.Skip("⏭️  Skipping integration test because BIFROST_TEST_CATALOG, BIFROST_TEST_SCHEMA or BIFROST_TEST_TABLE are not set")
	}

	client := sdk.NewClient(config)
	table := client.GetCatalog(testCatalog).Table(testSchema, testTable)

	params := url.Values{}
	params.Add("_limit", "1")

	resp, err := table.GetData(context.Background(), params)
	if err != nil {
		t.Fatalf("GetData failed: %v", err)
	}

	if resp.Status != "ok" {
		t.Fatalf("Expected status 'ok', got '%s'. Error: %s", resp.Status, resp.Error)
	}

	data, ok := resp.Data.([]interface{})
	if !ok {
		t.Fatalf("Expected response data to be a slice, got %T", resp.Data)
	}

	if len(data) != 1 {
		t.Errorf("Expected 1 row, got %d", len(data))
	}

	t.Log("✅ Successfully retrieved data from the API using environment variables")
}

func TestIntegration_GetDataWithParameters(t *testing.T) {
	if testing.Short() {
		t.Skip("⏭️  Skipping integration test in short mode")
	}

	// This test explicitly provides configuration parameters, overriding environment variables.
	// You need to fill these with valid values for the test to pass.
	overrideConfig := utils.Configuration{
		BaseURL: "https://bifrost.hyperfluid.cloud", // Replace with your actual base URL
		OrgID:   "your_org_id",                      // Replace with your actual Org ID
		Token:   "your_token",                       // Replace with your actual token OR Keycloak details
		// Example with Keycloak:
		// KeycloakBaseURL:    "https://keycloak.example.com",
		// KeycloakRealm:      "your_realm",
		// KeycloakClientID:   "your_client_id",
		// KeycloakClientSecret: "your_client_secret",
	}

	testCatalog := "your_test_catalog" // Replace with your actual test catalog
	testSchema := "your_test_schema"   // Replace with your actual test schema
	testTable := "your_test_table"     // Replace with your actual test table

	// Skip if placeholder values are still present
	if overrideConfig.OrgID == "your_org_id" || testCatalog == "your_test_catalog" {
		t.Skip("⏭️  Skipping TestIntegration_GetDataWithParameters, please provide actual config values in the test code.")
	}

	config, err := getTestConfig(&overrideConfig)
	if err != nil {
		t.Fatalf("Failed to get test config: %v", err)
	}

	client := sdk.NewClient(config)
	table := client.GetCatalog(testCatalog).Table(testSchema, testTable)

	params := url.Values{}
	params.Add("_limit", "1")

	resp, err := table.GetData(context.Background(), params)
	if err != nil {
		t.Fatalf("GetData with parameters failed: %v", err)
	}

	if resp.Status != "ok" {
		t.Fatalf("Expected status 'ok', got '%s'. Error: %s", resp.Status, resp.Error)
	}

	t.Log("✅ Successfully retrieved data from the API using explicit parameters")
}
