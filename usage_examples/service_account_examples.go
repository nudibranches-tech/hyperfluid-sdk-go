package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/nudibranches-tech/hyperfluid-sdk-go/sdk"
	"github.com/nudibranches-tech/hyperfluid-sdk-go/sdk/utils"
)

// This file demonstrates how to use Hyperfluid service accounts with the Bifrost SDK.
// Service accounts are the recommended way for service-to-service authentication.
//
// Service Account JSON Format:
//
//	{
//	  "client_id": "hf-org-sa-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
//	  "client_secret": "your-client-secret",
//	  "issuer": "https://auth.hyperfluid.cloud/realms/your-org",
//	  "auth_uri": "https://auth.hyperfluid.cloud/realms/your-org/protocol/openid-connect/auth",
//	  "token_uri": "https://auth.hyperfluid.cloud/realms/your-org/protocol/openid-connect/token"
//	}

// runServiceAccountFromFileExample demonstrates loading a service account from a file.
// This is the recommended approach for Kubernetes deployments using mounted secrets.
//
// Kubernetes Secret Example:
//
//	apiVersion: v1
//	kind: Secret
//	metadata:
//	  name: hyperfluid-service-account
//	type: Opaque
//	stringData:
//	  service_account.json: |
//	    {
//	      "client_id": "hf-org-sa-...",
//	      "client_secret": "...",
//	      "issuer": "https://auth.hyperfluid.cloud/realms/my-org",
//	      "auth_uri": "...",
//	      "token_uri": "..."
//	    }
//
// Pod Volume Mount:
//
//	volumes:
//	  - name: hyperfluid-credentials
//	    secret:
//	      secretName: hyperfluid-service-account
//	containers:
//	  - name: app
//	    volumeMounts:
//	      - name: hyperfluid-credentials
//	        mountPath: /var/run/secrets/hyperfluid
//	        readOnly: true
func runServiceAccountFromFileExample() {
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("🔐 Service Account Example 1: Load from File")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// Get the service account file path from environment or use default
	saPath := getEnv("HYPERFLUID_SERVICE_ACCOUNT_FILE", "/var/run/secrets/hyperfluid/service_account.json")
	apiURL := getEnv("HYPERFLUID_API_URL", "")

	if apiURL == "" {
		fmt.Println("⚠️  Skipping: HYPERFLUID_API_URL not set")
		fmt.Println()
		return
	}

	fmt.Printf("📁 Service account file: %s\n", saPath)
	fmt.Printf("🌐 API URL: %s\n", apiURL)
	fmt.Println()

	// Method 1: One-liner using NewClientFromServiceAccountFile
	client, err := sdk.NewClientFromServiceAccountFile(saPath, sdk.ServiceAccountOptions{
		BaseURL:        apiURL,
		OrgID:          getEnv("HYPERFLUID_ORG_ID", ""),
		DataDockID:     getEnv("HYPERFLUID_DATADOCK_ID", ""),
		RequestTimeout: 30,
		MaxRetries:     3,
		SkipTLSVerify:  globalSkipTLSVerify,
	})
	if err != nil {
		fmt.Printf("❌ Failed to create client: %v\n", err)
		fmt.Println()
		return
	}

	fmt.Println("✅ Client created successfully from service account file!")
	runSampleQuery(client)
	fmt.Println()
}

// runServiceAccountFromJSONExample demonstrates loading a service account from a JSON string.
// This is useful when the service account is provided via environment variables.
//
// Kubernetes ConfigMap/Secret as Environment Variable:
//
//	env:
//	  - name: HYPERFLUID_SERVICE_ACCOUNT
//	    valueFrom:
//	      secretKeyRef:
//	        name: hyperfluid-service-account
//	        key: service_account.json
func runServiceAccountFromJSONExample() {
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("🔐 Service Account Example 2: Load from JSON String")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// Get the service account JSON from environment variable
	saJSON := os.Getenv("HYPERFLUID_SERVICE_ACCOUNT")
	apiURL := getEnv("HYPERFLUID_API_URL", "")

	if saJSON == "" {
		fmt.Println("⚠️  Skipping: HYPERFLUID_SERVICE_ACCOUNT env var not set")
		fmt.Println()
		return
	}

	if apiURL == "" {
		fmt.Println("⚠️  Skipping: HYPERFLUID_API_URL not set")
		fmt.Println()
		return
	}

	fmt.Println("📋 Service account loaded from HYPERFLUID_SERVICE_ACCOUNT env var")
	fmt.Printf("🌐 API URL: %s\n", apiURL)
	fmt.Println()

	// Create client from JSON string
	client, err := sdk.NewClientFromServiceAccountJSON(saJSON, sdk.ServiceAccountOptions{
		BaseURL:        apiURL,
		OrgID:          getEnv("HYPERFLUID_ORG_ID", ""),
		DataDockID:     getEnv("HYPERFLUID_DATADOCK_ID", ""),
		RequestTimeout: 30,
		MaxRetries:     3,
		SkipTLSVerify:  globalSkipTLSVerify,
	})
	if err != nil {
		fmt.Printf("❌ Failed to create client: %v\n", err)
		fmt.Println()
		return
	}

	fmt.Println("✅ Client created successfully from JSON string!")
	runSampleQuery(client)
	fmt.Println()
}

// runServiceAccountManualExample demonstrates the two-step approach:
// 1. Load the service account
// 2. Create the client
//
// This is useful when you need to inspect or modify the service account before use.
func runServiceAccountManualExample() {
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("🔐 Service Account Example 3: Two-Step Approach")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	saPath := getEnv("HYPERFLUID_SERVICE_ACCOUNT_FILE", "")
	apiURL := getEnv("HYPERFLUID_API_URL", "")

	if saPath == "" || apiURL == "" {
		fmt.Println("⚠️  Skipping: Required environment variables not set")
		fmt.Println()
		return
	}

	// Step 1: Load the service account
	sa, err := sdk.LoadServiceAccount(saPath)
	if err != nil {
		fmt.Printf("❌ Failed to load service account: %v\n", err)
		fmt.Println()
		return
	}

	// Inspect the service account (e.g., for logging/debugging)
	fmt.Printf("📋 Service Account Details:\n")
	fmt.Printf("   Client ID: %s\n", sa.ClientID)
	fmt.Printf("   Issuer: %s\n", sa.Issuer)
	fmt.Println()

	// Parse the issuer to see what realm we're connecting to
	baseURL, realm, err := sa.ParseIssuer()
	if err != nil {
		fmt.Printf("❌ Failed to parse issuer: %v\n", err)
		fmt.Println()
		return
	}
	fmt.Printf("   Auth Server: %s\n", baseURL)
	fmt.Printf("   Realm: %s\n", realm)
	fmt.Println()

	// Step 2: Create the client
	client, err := sdk.NewClientFromServiceAccount(sa, sdk.ServiceAccountOptions{
		BaseURL:        apiURL,
		OrgID:          getEnv("HYPERFLUID_ORG_ID", ""),
		DataDockID:     getEnv("HYPERFLUID_DATADOCK_ID", ""),
		RequestTimeout: 30,
		MaxRetries:     3,
		SkipTLSVerify:  globalSkipTLSVerify,
	})
	if err != nil {
		fmt.Printf("❌ Failed to create client: %v\n", err)
		fmt.Println()
		return
	}

	fmt.Println("✅ Client created successfully!")
	runSampleQuery(client)
	fmt.Println()
}

// runServiceAccountWithSkipTLSExample demonstrates using SkipTLSVerify for development.
// WARNING: Never use SkipTLSVerify in production!
func runServiceAccountWithSkipTLSExample() {
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("🔐 Service Account Example 4: Development (Skip TLS)")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	saPath := getEnv("HYPERFLUID_SERVICE_ACCOUNT_FILE", "")
	apiURL := getEnv("HYPERFLUID_API_URL", "")

	if saPath == "" || apiURL == "" {
		fmt.Println("⚠️  Skipping: Required environment variables not set")
		fmt.Println()
		return
	}

	fmt.Println("⚠️  WARNING: SkipTLSVerify is enabled - for development only!")
	fmt.Println()

	client, err := sdk.NewClientFromServiceAccountFile(saPath, sdk.ServiceAccountOptions{
		BaseURL:       apiURL,
		SkipTLSVerify: true, // Only for development!
	})
	if err != nil {
		fmt.Printf("❌ Failed to create client: %v\n", err)
		fmt.Println()
		return
	}

	fmt.Println("✅ Client created with TLS verification disabled")
	runSampleQuery(client)
	fmt.Println()
}

// runSampleQuery executes a simple query to verify the client is working.
func runSampleQuery(client *sdk.Client) {
	testCatalog := getEnv("BIFROST_TEST_CATALOG", "")
	testSchema := getEnv("BIFROST_TEST_SCHEMA", "")
	testTable := getEnv("BIFROST_TEST_TABLE", "")
	dataDockID := getEnv("HYPERFLUID_DATADOCK_ID", "")

	if testCatalog == "" || testSchema == "" || testTable == "" {
		fmt.Println("ℹ️  Skipping sample query: test table not configured")
		return
	}

	fmt.Printf("📊 Running sample query on %s.%s.%s...\n", testCatalog, testSchema, testTable)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := client.
		DataDock(dataDockID).
		Catalog(testCatalog).
		Schema(testSchema).
		Table(testTable).
		Limit(5).
		Get(ctx)

	handleServiceAccountResponse(resp, err)
}

func handleServiceAccountResponse(resp *utils.Response, err error) {
	if err != nil {
		fmt.Printf("❌ Query error: %v\n", err)
		if resp != nil && resp.Error != "" {
			fmt.Printf("   Server response: %s (HTTP %d)\n", resp.Error, resp.HTTPCode)
		}
		return
	}
	if resp.Status != utils.StatusOK {
		fmt.Printf("❌ Query failed: %s (HTTP %d)\n", resp.Error, resp.HTTPCode)
		return
	}
	fmt.Println("✅ Query successful!")
	if dataSlice, ok := resp.Data.([]interface{}); ok {
		fmt.Printf("📦 Retrieved %d records\n", len(dataSlice))
	}
}

// RunServiceAccountExamples runs all service account examples.
// Call this from main.go to include these examples in the demo.
func RunServiceAccountExamples() {
	fmt.Println()
	fmt.Println("══════════════════════════════════════════════════════════")
	fmt.Println("🔐 SERVICE ACCOUNT AUTHENTICATION EXAMPLES")
	fmt.Println("══════════════════════════════════════════════════════════")
	fmt.Println()

	runServiceAccountFromFileExample()
	runServiceAccountFromJSONExample()
	runServiceAccountManualExample()
	runServiceAccountWithSkipTLSExample()
}
