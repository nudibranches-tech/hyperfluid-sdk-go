package core

import (
	"bifrost-for-developers/sdk/core/utils"
	"context"
	"time"

	"github.com/joho/godotenv"
)

var globalConfiguration *utils.Configuration

func InitFromEnv() {

	godotenv.Overload("../.env")

	globalConfiguration = &utils.Configuration{
		BaseURL: utils.GetEnvironmentVariable("HYPERFLUID_BASE_URL", "https://bifrost.hyperfluid.cloud"),
		OrgID:   utils.GetEnvironmentVariable("HYPERFLUID_ORG_ID", ""),
		Token:   GetToken(context.Background()),

		SkipTLSVerify:  utils.GetEnvironmentVariable("HYPERFLUID_SKIP_TLS_VERIFY", "true") == "true",
		RequestTimeout: time.Duration(utils.GetEnvironmentVariableInt("HYPERFLUID_REQUEST_TIMEOUT", 30)) * time.Second,
		MaxRetries:     utils.GetEnvironmentVariableInt("HYPERFLUID_MAX_RETRIES", 3),

		PostgresHost:     utils.GetEnvironmentVariable("HYPERFLUID_POSTGRES_HOST", "bifrost.hyperfluid.cloud"),
		PostgresPort:     utils.GetEnvironmentVariableInt("HYPERFLUID_POSTGRES_PORT", 5432),
		PostgresUser:     utils.GetEnvironmentVariable("HYPERFLUID_POSTGRES_USER", "hyperfluid"),
		PostgresDatabase: utils.GetEnvironmentVariable("HYPERFLUID_POSTGRES_DATABASE", ""),

		KeycloakUsername: utils.GetEnvironmentVariable("KEYCLOAK_USERNAME", "demo"),
		KeycloakPassword: utils.GetEnvironmentVariable("KEYCLOAK_PASSWORD", "demo"),
		KeycloakBaseURL:  utils.GetEnvironmentVariable("KEYCLOAK_BASE_URL", "https://keycloak.localhost:8443"),
		KeycloakClientID: utils.GetEnvironmentVariable("KEYCLOAK_CLIENT_ID", "authentication-client"),
		KeycloakRealm:    utils.GetEnvironmentVariable("KEYCLOAK_REALM", "nudibranches-tech"),

		TestCatalog: utils.GetEnvironmentVariable("BIFROST_TEST_CATALOG", ""),
		TestSchema:  utils.GetEnvironmentVariable("BIFROST_TEST_SCHEMA", ""),
		TestTable:   utils.GetEnvironmentVariable("BIFROST_TEST_TABLE", ""),
		TestColumns: utils.GetEnvironmentVariable("BIFROST_TEST_COLUMNS", ""),
	}
}

func GetConfiguration() *utils.Configuration {
	return globalConfiguration
}
