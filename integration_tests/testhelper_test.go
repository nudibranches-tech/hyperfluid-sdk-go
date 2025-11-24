package developpementtests

import (
	"github.com/nudibranches-tech/bifrost-hyperfluid-sdk-dev/sdk/utils"
	"os"
	"strconv"
	"time"
)

// getTestConfig loads configuration from environment variables,
// optionally overriding with values provided in overrideConfig.
func getTestConfig(overrideConfig *utils.Configuration) (utils.Configuration, error) {
	// Start with config loaded from environment variables
	config := loadTestConfigFromEnv()

	// Apply overrides
	if overrideConfig != nil {
		if overrideConfig.BaseURL != "" {
			config.BaseURL = overrideConfig.BaseURL
		}
		if overrideConfig.OrgID != "" {
			config.OrgID = overrideConfig.OrgID
		}
		if overrideConfig.Token != "" {
			config.Token = overrideConfig.Token
		}
		if overrideConfig.RequestTimeout != 0 {
			config.RequestTimeout = overrideConfig.RequestTimeout
		}
		if overrideConfig.KeycloakUsername != "" {
			config.KeycloakUsername = overrideConfig.KeycloakUsername
		}
		if overrideConfig.KeycloakPassword != "" {
			config.KeycloakPassword = overrideConfig.KeycloakPassword
		}
		if overrideConfig.KeycloakBaseURL != "" {
			config.KeycloakBaseURL = overrideConfig.KeycloakBaseURL
		}
		if overrideConfig.KeycloakClientID != "" {
			config.KeycloakClientID = overrideConfig.KeycloakClientID
		}
		if overrideConfig.KeycloakRealm != "" {
			config.KeycloakRealm = overrideConfig.KeycloakRealm
		}
		if overrideConfig.KeycloakClientSecret != "" {
			config.KeycloakClientSecret = overrideConfig.KeycloakClientSecret
		}
		config.SkipTLSVerify = overrideConfig.SkipTLSVerify // bools must be handled explicitly
		if overrideConfig.MaxRetries != 0 {
			config.MaxRetries = overrideConfig.MaxRetries
		}
	}

	return config, nil
}

// loadTestConfigFromEnv loads configuration solely from environment variables.
func loadTestConfigFromEnv() utils.Configuration {
	config := utils.Configuration{
		BaseURL:        getEnv("HYPERFLUID_BASE_URL", ""),
		OrgID:          getEnv("HYPERFLUID_ORG_ID", ""),
		Token:          getEnv("HYPERFLUID_TOKEN", ""),
		RequestTimeout: time.Duration(getEnvInt("HYPERFLUID_REQUEST_TIMEOUT", 30)) * time.Second,
		SkipTLSVerify:  getEnv("HYPERFLUID_SKIP_TLS_VERIFY", "false") == "true",
		MaxRetries:     getEnvInt("HYPERFLUID_MAX_RETRIES", 3),
	}

	// Keycloak config
	config.KeycloakUsername = getEnv("KEYCLOAK_USERNAME", "")
	config.KeycloakPassword = getEnv("KEYCLOAK_PASSWORD", "")
	config.KeycloakBaseURL = getEnv("KEYCLOAK_BASE_URL", "")
	config.KeycloakClientID = getEnv("KEYCLOAK_CLIENT_ID", "")
	config.KeycloakRealm = getEnv("KEYCLOAK_REALM", "")
	config.KeycloakClientSecret = getEnv("KEYCLOAK_CLIENT_SECRET", "")

	return config
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return fallback
}
