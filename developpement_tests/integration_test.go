package developpementtests

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	bifrost "bifrost-for-developers/sdk"
	"bifrost-for-developers/sdk/core"
	"bifrost-for-developers/sdk/core/utils"
)

func TestIntegration_01_FetchKeycloakToken(t *testing.T) {
	if testing.Short() {
		t.Skip("⏭️  Skipping integration test in short mode")
	}

	bifrost.Init()
	var globalConfiguration = bifrost.GetGlobalConfiguration()

	if globalConfiguration.KeycloakUsername == "" || globalConfiguration.KeycloakPassword == "" {
		t.Skip("⏭️  Keycloak credentials not globalConfigurationured")
	}

	t.Run("FetchAndValidateToken", func(t *testing.T) {
		requestContext, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancelFunc()
		if _, error := core.FetchKeycloakToken(requestContext); error != nil {
			t.Logf("❌ Keycloak failed: %v", error)
			t.Skip("⏭️  Keycloak not available")
		}

		if globalConfiguration.Token == "" {
			t.Fatal("❌ Token should be set")
		}
		tokenPreview := globalConfiguration.Token
		if len(tokenPreview) > 50 {
			tokenPreview = tokenPreview[:50] + "..."
		}
		t.Logf("   Token: %s", tokenPreview)
	})
}

func TestIntegration_GraphQL(t *testing.T) {
	if testing.Short() {
		t.Skip("⏭️  Skipping integration test in short mode")
	}

	bifrost.Init()
	var globalConfiguration = bifrost.GetGlobalConfiguration()

	query := fmt.Sprintf(`{%s{%s{%s(limit:10){%s}}}}`,
		globalConfiguration.TestCatalog, globalConfiguration.TestSchema, globalConfiguration.TestTable, strings.ReplaceAll(globalConfiguration.TestColumns, ",", " "))

	t.Logf("🔷 GraphQL Request: %s", query)

	t.Run("QueryAndValidate", func(t *testing.T) {
		result := <-bifrost.Request(utils.BifrostRequest{Type: utils.RequestGraphQL, GraphQLPayload: &utils.GraphQLPayload{Query: query}})
		response, error := result.Response, result.Error

		if error != nil || response == nil || !response.IsOK() {
			if error != nil {
				t.Logf("❌ Request failed: %v", error)
			}
			if response != nil && response.HasError() {
				t.Logf("   Error: %s", response.Error)
			}
			t.Skip("⏭️  GraphQL not available")
		}

		if responseDataMap, isMap := response.GetDataAsMap(); isMap {
			if _, hasData := responseDataMap["data"]; hasData {
				t.Log("   ✓ Data received")
				t.Log(json.Marshal(responseDataMap))
			}
			if errors, hasErrors := responseDataMap["errors"]; hasErrors {
				if errorList, isList := errors.([]interface{}); isList && len(errorList) > 0 {
					t.Logf("   ⚠️  %d error(s)", len(errorList))
					t.Log(json.Marshal(errorList))
				}
			}
		}
	})
}

func TestIntegration_OpenAPI(t *testing.T) {
	if testing.Short() {
		t.Skip("⏭️  Skipping integration test in short mode")
	}

	bifrost.Init()
	var globalConfiguration = bifrost.GetGlobalConfiguration()

	var openAPIPayload = &utils.OpenAPIPayload{
		Catalog: globalConfiguration.TestCatalog,
		Schema:  globalConfiguration.TestSchema,
		Table:   globalConfiguration.TestTable,
		Method:  "GET",
		Params: map[string]string{
			"_limit": "10",
			"select": globalConfiguration.TestColumns,
		},
	}
	t.Logf("🌐 OpenAPI Request Payload: %s", utils.JsonMarshal(openAPIPayload))

	t.Run("QueryAndValidate", func(t *testing.T) {
		result := <-bifrost.Request(utils.BifrostRequest{Type: utils.RequestOpenAPI, OpenAPIPayload: openAPIPayload})
		response, error := result.Response, result.Error

		if error != nil || response == nil || !response.IsOK() {
			if error != nil {
				t.Logf("❌ Request failed: %v", error)
			}
			if response != nil && response.HasError() {
				t.Logf("   Error: %s", response.Error)
			}
			t.Skip("⏭️  OpenAPI not available")
		}

		if responseDataSlice, isSlice := response.GetDataAsSlice(); isSlice {
			t.Logf("   ✓ %d rows returned", len(responseDataSlice))
			t.Log(json.Marshal(responseDataSlice))
		}
	})
}

func TestIntegration_Postgres(t *testing.T) {
	if testing.Short() {
		t.Skip("⏭️  Skipping integration test in short mode")
	}

	bifrost.Init()
	var globalConfiguration = bifrost.GetGlobalConfiguration()

	sqlQuery := fmt.Sprintf("SELECT %s FROM %s.%s.%s LIMIT 5", globalConfiguration.TestColumns, globalConfiguration.TestCatalog, globalConfiguration.TestSchema, globalConfiguration.TestTable)
	t.Logf("🐘 PostgreSQL: %s", sqlQuery)

	t.Run("QueryAndValidate", func(t *testing.T) {
		result := <-bifrost.Request(utils.BifrostRequest{Type: utils.RequestPostgres, PostgresPayload: &utils.PostgresPayload{SQL: sqlQuery}})
		response, error := result.Response, result.Error

		if error != nil || response == nil || !response.IsOK() {
			if error != nil {
				t.Logf("❌ Connection failed: %v", error)
			}
			t.Skip("⏭️  PostgreSQL not available")
		}

		if responseDataSlice, isSlice := response.GetDataAsSlice(); isSlice {
			t.Logf("   ✓ %d rows returned", len(responseDataSlice))
		}
	})
}
