package developpementtests

import (
	"testing"

	bifrost "bifrost-for-developers/sdk"
	"bifrost-for-developers/sdk/core/utils"
)

func TestUnit_InitFromEnv(t *testing.T) {
	if testing.Short() {
		t.Skip("⏭️  Skipping in short mode")
	}
	bifrost.Init()
	var globalConfiguration = bifrost.GetGlobalConfiguration()

	t.Run("RequiredFieldsSet", func(t *testing.T) {
		if globalConfiguration == nil {
			t.Fatal("❌ globalConfigurationuration not initialized")
		}
		if globalConfiguration.OrgID == "" {
			t.Error("❌ Organization ID not set")
		}
		if globalConfiguration.BaseURL == "" {
			t.Error("❌ Base URL not set")
		}
		t.Logf("⚙️  OrgID: %s | URL: %s", globalConfiguration.OrgID, globalConfiguration.BaseURL)
	})
}

func TestUnit_InvalidRequest(t *testing.T) {
	if testing.Short() {
		t.Skip("⏭️  Skipping in short mode")
	}

	bifrost.Init()

	t.Run("UnsupportedRequestType", func(t *testing.T) {
		result := <-bifrost.Request(utils.BifrostRequest{Type: utils.RequestType("invalid"), GraphQLPayload: nil, OpenAPIPayload: nil, PostgresPayload: nil})
		if result.Error == nil {
			t.Error("❌ Expected error for invalid request type")
		}
	})
}
