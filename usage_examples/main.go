package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

var (
	runAll            = flag.Bool("all", false, "Run all examples")
	runFluent         = flag.Bool("fluent", false, "Run fluent API examples")
	runS3             = flag.Bool("s3", false, "Run S3 examples")
	runSearch         = flag.Bool("search", false, "Run search examples")
	runControlPlane   = flag.Bool("control-plane", false, "Run Control Plane API examples")
	runServiceAccount = flag.Bool("service-account", false, "Run service account examples")
	listExamples      = flag.Bool("list", false, "List available examples")
	skipTLSVerify     = flag.Bool("skip-tls-verify", false, "Skip TLS certificate verification (WARNING: Use only for development)")
)

// Global variable to share skipTLSVerify setting across examples
var globalSkipTLSVerify bool

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "Bifrost SDK Usage Examples")
		fmt.Fprintln(os.Stderr, "\nOptions:")
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, "\nExamples:")
		fmt.Fprintln(os.Stderr, "  Run all examples:")
		fmt.Fprintf(os.Stderr, "    %s --all\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "\n  Run specific examples:")
		fmt.Fprintf(os.Stderr, "    %s --fluent --search\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "    %s --control-plane\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "\n  Development mode (skip TLS verification):")
		fmt.Fprintf(os.Stderr, "    %s --control-plane --skip-tls-verify\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "\n  List available examples:")
		fmt.Fprintf(os.Stderr, "    %s --list\n", os.Args[0])
	}
	flag.Parse()

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("🚀 Bifrost SDK - Usage Examples")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	if *listExamples {
		printAvailableExamples()
		return
	}

	if err := godotenv.Load("../.env"); err != nil {
		log.Printf("⚠️  Warning: .env not loaded: %v\n", err)
	} else {
		fmt.Println("✓ .env file loaded successfully")
	}
	fmt.Println()

	// Set global skip TLS verify flag
	globalSkipTLSVerify = *skipTLSVerify
	if globalSkipTLSVerify {
		fmt.Println("⚠️  WARNING: TLS certificate verification is DISABLED")
		fmt.Println("   This should only be used for development/testing!")
		fmt.Println()
	}

	// If no flags specified, show usage
	noFlagsSet := !*runAll && !*runFluent && !*runS3 && !*runSearch && !*runControlPlane && !*runServiceAccount
	if noFlagsSet {
		fmt.Println("ℹ️  No examples selected. Use --help to see available options.")
		fmt.Println()
		printAvailableExamples()
		return
	}

	// Track which examples to run
	runFluentExamples := *runAll || *runFluent
	runS3Examples := *runAll || *runS3
	runSearchExamples := *runAll || *runSearch
	runControlPlaneExamples := *runAll || *runControlPlane
	runServiceAccountExamples := *runAll || *runServiceAccount

	// Run selected examples
	if runFluentExamples {
		fmt.Println("══════════════════════════════════════════════")
		fmt.Println("📊 FLUENT API EXAMPLES")
		fmt.Println("══════════════════════════════════════════════")
		fmt.Println()
		runFluentAPIWithSelectExample()
	}

	if runS3Examples {
		fmt.Println("══════════════════════════════════════════════")
		fmt.Println("📦 S3 EXAMPLES")
		fmt.Println("══════════════════════════════════════════════")
		fmt.Println()
		runS3Example()
	}

	if runSearchExamples {
		fmt.Println("══════════════════════════════════════════════")
		fmt.Println("🔍 SEARCH EXAMPLES")
		fmt.Println("══════════════════════════════════════════════")
		fmt.Println()
		runSearchExample()
	}

	if runServiceAccountExamples {
		RunServiceAccountExamples()
	}

	if runControlPlaneExamples {
		RunControlPlaneExamples()
	}

	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("🎉 All examples completed!")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}

func printAvailableExamples() {
	examples := []struct {
		name        string
		flag        string
		description string
		envVars     []string
	}{
		{
			name:        "Fluent API",
			flag:        "--fluent",
			description: "Demonstrates fluent query building with method chaining",
			envVars:     []string{"HYPERFLUID_SERVICE_ACCOUNT_FILE", "HYPERFLUID_API_URL", "BIFROST_TEST_CATALOG"},
		},
		{
			name:        "S3 Operations",
			flag:        "--s3",
			description: "S3 bucket operations (list, upload, download)",
			envVars:     []string{"HYPERFLUID_SERVICE_ACCOUNT_FILE", "HYPERFLUID_API_URL"},
		},
		{
			name:        "Search API",
			flag:        "--search",
			description: "Full-text search across data containers",
			envVars:     []string{"HYPERFLUID_SERVICE_ACCOUNT_FILE", "HYPERFLUID_API_URL"},
		},
		{
			name:        "Service Accounts",
			flag:        "--service-account",
			description: "Service account authentication patterns",
			envVars:     []string{"HYPERFLUID_SERVICE_ACCOUNT_FILE", "HYPERFLUID_API_URL"},
		},
		{
			name:        "Control Plane API",
			flag:        "--control-plane",
			description: "Archive operations (import/export tracking and management)",
			envVars:     []string{"HYPERFLUID_SERVICE_ACCOUNT_FILE", "HYPERFLUID_API_URL", "HYPERFLUID_DATA_DOCK_ID", "HYPERFLUID_DATA_CONTAINER_ID"},
		},
	}

	fmt.Println("Available Examples:")
	fmt.Println(strings.Repeat("─", 80))
	for _, ex := range examples {
		fmt.Printf("\n📌 %s (%s)\n", ex.name, ex.flag)
		fmt.Printf("   %s\n", ex.description)
		fmt.Printf("   Required env vars: %s\n", strings.Join(ex.envVars, ", "))
	}
	fmt.Println()
	fmt.Println(strings.Repeat("─", 80))
	fmt.Println("\nUsage:")
	fmt.Println("  go run . --all                 # Run all examples")
	fmt.Println("  go run . --control-plane       # Run Control Plane examples only")
	fmt.Println("  go run . --fluent --search     # Run multiple specific examples")
	fmt.Println()
}
