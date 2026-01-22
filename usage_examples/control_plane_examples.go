package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nudibranches-tech/bifrost-hyperfluid-sdk-dev/sdk"
	"github.com/nudibranches-tech/bifrost-hyperfluid-sdk-dev/sdk/controlplaneapiclient"
)

// runControlPlaneListArchiveOperationsExample demonstrates listing archive operations for an organization.
func runControlPlaneListArchiveOperationsExample() {
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("🎛️  Control Plane Example 1: List Archive Operations")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	saPath := getEnv("HYPERFLUID_SERVICE_ACCOUNT_FILE", "")
	apiURL := getEnv("HYPERFLUID_API_URL", "")
	dataDockID := getEnv("HYPERFLUID_DATA_DOCK_ID", "")
	dataContainerID := getEnv("HYPERFLUID_DATA_CONTAINER_ID", "")

	if saPath == "" || apiURL == "" || dataDockID == "" || dataContainerID == "" {
		fmt.Println("⚠️  Skipping: Required environment variables not set")
		fmt.Println("   Set HYPERFLUID_SERVICE_ACCOUNT_FILE, HYPERFLUID_API_URL, HYPERFLUID_DATA_DOCK_ID, and HYPERFLUID_DATA_CONTAINER_ID")
		fmt.Println()
		return
	}

	client, err := sdk.NewClientFromServiceAccountFile(saPath, sdk.ServiceAccountOptions{
		BaseURL:        apiURL,
		RequestTimeout: 30,
		SkipTLSVerify:  globalSkipTLSVerify,
	})
	if err != nil {
		fmt.Printf("❌ Failed to create client: %v\n", err)
		fmt.Println()
		return
	}

	cp, err := client.ControlPlane()
	if err != nil {
		fmt.Printf("❌ Failed to get Control Plane client: %v\n", err)
		fmt.Println()
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fmt.Println("✅ Control Plane client created successfully!")
	fmt.Println()

	// Parse data dock ID to UUID
	dataDockUUID, err := uuid.Parse(dataDockID)
	if err != nil {
		fmt.Printf("❌ Invalid data dock ID: %v\n", err)
		fmt.Println()
		return
	}

	// Parse data container ID to UUID
	dataContainerUUID, err := uuid.Parse(dataContainerID)
	if err != nil {
		fmt.Printf("❌ Invalid data container ID: %v\n", err)
		fmt.Println()
		return
	}

	// List Archive Operations
	fmt.Printf("📋 Listing Archive Operations for DataDock %s, DataContainer %s...\n", dataDockID, dataContainerID)
	fmt.Println()

	resp, err := cp.ListArchiveOperationsWithResponse(ctx, dataDockUUID, dataContainerUUID, &controlplaneapiclient.ListArchiveOperationsParams{})
	if err != nil {
		fmt.Printf("❌ Failed to list archive operations: %v\n", err)
		fmt.Println()
		return
	}

	if resp.StatusCode() != 200 {
		fmt.Printf("❌ API returned status %d: %s\n", resp.StatusCode(), string(resp.Body))
		fmt.Println()
		return
	}

	if resp.JSON200 != nil {
		operations := *resp.JSON200
		fmt.Printf("✅ Found %d Archive Operation(s)\n", len(operations))
		fmt.Println()

		if len(operations) == 0 {
			fmt.Println("   No archive operations found.")
		}

		for i, op := range operations {
			fmt.Printf("   %d. Operation ID: %s\n", i+1, op.Id)
			fmt.Printf("      Type: %s\n", op.OperationType)
			fmt.Printf("      Status: %s\n", op.Status)
			fmt.Printf("      File: %s (Type: %s)\n", op.FileName, op.FileType)

			if op.FileSize != nil {
				fmt.Printf("      Size: %d bytes\n", *op.FileSize)
			}

			if op.FilePath != nil {
				fmt.Printf("      Path: %s\n", *op.FilePath)
			}

			if op.Prefix != nil {
				fmt.Printf("      Prefix: %s\n", *op.Prefix)
			}

			if op.ErrorMessage != nil && *op.ErrorMessage != "" {
				fmt.Printf("      Error: %s\n", *op.ErrorMessage)
			}

			fmt.Printf("      Created: %s\n", op.CreatedAt.Format(time.RFC3339))
			fmt.Printf("      Updated: %s\n", op.UpdatedAt.Format(time.RFC3339))

			if i < len(operations)-1 {
				fmt.Println()
			}
		}
	} else {
		fmt.Println("⚠️  No operations data in response")
	}
	fmt.Println()
}

// RunControlPlaneExamples runs all Control Plane examples.
// Call this from main.go to include these examples in the demo.
func RunControlPlaneExamples() {
	fmt.Println()
	fmt.Println("══════════════════════════════════════════════════════════")
	fmt.Println("🎛️  CONTROL PLANE API EXAMPLES")
	fmt.Println("══════════════════════════════════════════════════════════")
	fmt.Println()

	runControlPlaneListArchiveOperationsExample()
}
