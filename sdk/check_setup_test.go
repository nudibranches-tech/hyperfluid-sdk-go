package bifrost

import (
	"bifrost-for-developers/sdk/core"
	"context"
	"testing"
	"time"
)

var globalConfiguration *Configuration

func TestSetup_SanityCheck(t *testing.T) {
	t.Log("🔧 Verifying globalConfigurationuration...")

	Init()
	globalConfiguration = core.GetConfiguration()
	t.Logf("OrgID: %s", globalConfiguration.OrgID)
	t.Logf("BaseURL: %s", globalConfiguration.BaseURL)
	t.Logf("KeycloakUsername: %s", globalConfiguration.KeycloakUsername)
	t.Logf("KeycloakPassword: %s", globalConfiguration.KeycloakPassword)
	t.Logf("KeycloakBaseURL: %s", globalConfiguration.KeycloakBaseURL)
	t.Logf("KeycloakClientID: %s", globalConfiguration.KeycloakClientID)
	t.Logf("KeycloakRealm: %s", globalConfiguration.KeycloakRealm)
	t.Logf("TestCatalog: %s", globalConfiguration.TestCatalog)
	t.Logf("TestSchema: %s", globalConfiguration.TestSchema)
	t.Logf("TestTable: %s", globalConfiguration.TestTable)
	t.Logf("TestColumns: %s", globalConfiguration.TestColumns)

	if globalConfiguration == nil {
		t.Fatal("❌ globalConfigurationuration not initialized")
	}

	if globalConfiguration.BaseURL == "" {
		t.Fatal("❌ Missing HYPERFLUID_BASE_URL in .env\n   Add your base URL to continue.")
	}
	if globalConfiguration.OrgID == "" {
		t.Fatal("❌ Missing HYPERFLUID_ORG_ID in .env\n   Add your organization ID to continue.")
	}
	t.Log("✅ Base URL & Organization ID globalConfigurationured")

	if globalConfiguration.KeycloakUsername != "" && globalConfiguration.KeycloakPassword != "" {
		t.Log("🔑 Fetching token from Keycloak...")
		context, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if _, err := core.FetchKeycloakToken(context); err != nil {
			t.Logf("⚠️  Keycloak fetch failed: %v", err)
			// Fall back to token from .env if available
			globalConfiguration = core.GetConfiguration()
			if globalConfiguration.Token != "" {
				t.Log("✅ Using token from .env as fallback")
				t.Log("🎉 Setup complete! Ready to use.")
				return
			}
			t.Fatal("❌ Authentication failed: Keycloak unavailable and no HYPERFLUID_TOKEN in .env")
		}

		t.Log("✅ Token obtained (auto-fetched from Keycloak)")
		globalConfiguration = core.GetConfiguration()
		if globalConfiguration.Token == "" {
			t.Fatal("❌ Token fetch succeeded but token is empty")
		}
		t.Log("🎉 Setup complete! Ready to use.")
		return
	}

	if globalConfiguration.Token != "" {
		t.Log("✅ Authentication token found (from .env)")

		return
	}

	t.Log("🔧 Checking additional .env globalConfigurationuration...")
	if globalConfiguration.PostgresHost == "" || globalConfiguration.PostgresPort == 0 || globalConfiguration.PostgresUser == "" || globalConfiguration.PostgresDatabase == "" {
		t.Log("❌ Missing HYPERFLUID_POSTGRES_* in .env\n   Add your Postgres information to use PostgreSQL.")
	}
	if globalConfiguration.KeycloakUsername == "" || globalConfiguration.KeycloakPassword == "" || globalConfiguration.KeycloakBaseURL == "" || globalConfiguration.KeycloakClientID == "" || globalConfiguration.KeycloakRealm == "" {
		t.Log("❌ Missing KEYCLOAK_* in .env\n   Add your Keycloak information to refresh your token automatically.")
	}
	if globalConfiguration.TestCatalog == "" || globalConfiguration.TestSchema == "" || globalConfiguration.TestTable == "" {
		t.Log("❌ Missing BIFROST_TEST_* in .env\n   Add your Test data before running integration tests.")
	}
	t.Log("✅ Configuration checks complete")
	t.Log("🎉 Ready to use.")
}
