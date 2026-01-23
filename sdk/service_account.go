package sdk

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"

	"github.com/nudibranches-tech/hyperfluid-sdk-go/sdk/utils"
)

// ServiceAccount represents the Hyperfluid service account credentials.
// This is the standard format distributed by Hyperfluid for service-to-service authentication.
//
// Example JSON file:
//
//	{
//	  "client_id": "hf-org-sa-9e4132be-9498-4e75-8d8b-9699c92c4673",
//	  "client_secret": "8ek2Muno5b5sHeLUV8yk6pUYoPHPk6oZ",
//	  "issuer": "https://auth.hyperfluid.cloud/realms/nudibranches-tech",
//	  "auth_uri": "https://auth.hyperfluid.cloud/realms/nudibranches-tech/protocol/openid-connect/auth",
//	  "token_uri": "https://auth.hyperfluid.cloud/realms/nudibranches-tech/protocol/openid-connect/token"
//	}
type ServiceAccount struct {
	// ClientID is the OAuth2 client identifier for the service account.
	ClientID string `json:"client_id"`

	// ClientSecret is the OAuth2 client secret for authentication.
	ClientSecret string `json:"client_secret"`

	// Issuer is the OIDC issuer URL (e.g., "https://auth.hyperfluid.cloud/realms/my-org").
	// Used to derive the Keycloak base URL and realm.
	Issuer string `json:"issuer"`

	// AuthURI is the OAuth2 authorization endpoint (typically not used for service accounts).
	AuthURI string `json:"auth_uri"`

	// TokenURI is the OAuth2 token endpoint used to obtain access tokens.
	TokenURI string `json:"token_uri"`
}

// LoadServiceAccount loads a ServiceAccount from a JSON file at the given path.
// This is the recommended way to load credentials in production environments,
// especially when using Kubernetes secrets mounted as files.
//
// Example:
//
//	// Load from a mounted Kubernetes secret
//	sa, err := sdk.LoadServiceAccount("/var/run/secrets/hyperfluid/service_account.json")
//	if err != nil {
//	    log.Fatalf("Failed to load service account: %v", err)
//	}
//
//	// Create client with additional options
//	client, err := sdk.NewClientFromServiceAccount(sa, sdk.ServiceAccountOptions{
//	    BaseURL: "https://api.hyperfluid.cloud",
//	})
func LoadServiceAccount(path string) (*ServiceAccount, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open service account file: %w", err)
	}
	defer func() { _ = file.Close() }()

	return LoadServiceAccountFromReader(file)
}

// LoadServiceAccountFromJSON loads a ServiceAccount from a JSON string.
// This is useful when the service account is provided via environment variables.
//
// Example:
//
//	// Load from environment variable
//	saJSON := os.Getenv("HYPERFLUID_SERVICE_ACCOUNT")
//	sa, err := sdk.LoadServiceAccountFromJSON(saJSON)
//	if err != nil {
//	    log.Fatalf("Failed to parse service account: %v", err)
//	}
func LoadServiceAccountFromJSON(jsonStr string) (*ServiceAccount, error) {
	return LoadServiceAccountFromReader(strings.NewReader(jsonStr))
}

// LoadServiceAccountFromReader loads a ServiceAccount from an io.Reader.
// This provides maximum flexibility for loading from various sources.
//
// Example:
//
//	// Load from an embedded file or any io.Reader
//	sa, err := sdk.LoadServiceAccountFromReader(myReader)
func LoadServiceAccountFromReader(r io.Reader) (*ServiceAccount, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read service account data: %w", err)
	}

	var sa ServiceAccount
	if err := json.Unmarshal(data, &sa); err != nil {
		return nil, fmt.Errorf("failed to parse service account JSON: %w", err)
	}

	if err := sa.Validate(); err != nil {
		return nil, fmt.Errorf("invalid service account: %w", err)
	}

	return &sa, nil
}

// Validate checks that the ServiceAccount has all required fields populated.
func (sa *ServiceAccount) Validate() error {
	if sa.ClientID == "" {
		return fmt.Errorf("client_id is required")
	}
	if sa.ClientSecret == "" {
		return fmt.Errorf("client_secret is required")
	}
	if sa.Issuer == "" && sa.TokenURI == "" {
		return fmt.Errorf("either issuer or token_uri is required")
	}
	return nil
}

// ParseIssuer extracts the Keycloak base URL and realm from the issuer URL.
// The issuer URL format is: https://<host>/realms/<realm>
//
// Returns:
//   - baseURL: The Keycloak server URL (e.g., "https://auth.hyperfluid.cloud")
//   - realm: The Keycloak realm name (e.g., "nudibranches-tech")
//   - error: If the issuer URL cannot be parsed
func (sa *ServiceAccount) ParseIssuer() (baseURL, realm string, err error) {
	if sa.Issuer == "" {
		// Try to extract from token_uri as fallback
		if sa.TokenURI != "" {
			return parseKeycloakURL(sa.TokenURI)
		}
		return "", "", fmt.Errorf("issuer is empty and no token_uri available")
	}
	return parseKeycloakURL(sa.Issuer)
}

// parseKeycloakURL extracts base URL and realm from a Keycloak URL.
// Supports both issuer format (https://host/realms/realm) and
// token URL format (https://host/realms/realm/protocol/openid-connect/token).
func parseKeycloakURL(rawURL string) (baseURL, realm string, err error) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse URL: %w", err)
	}

	// Validate scheme is present and valid
	if parsed.Scheme == "" {
		return "", "", fmt.Errorf("URL missing scheme (http/https): %s", rawURL)
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return "", "", fmt.Errorf("URL has invalid scheme %q, expected http or https: %s", parsed.Scheme, rawURL)
	}

	// Path format: /realms/<realm> or /realms/<realm>/protocol/...
	parts := strings.Split(strings.Trim(parsed.Path, "/"), "/")
	if len(parts) < 2 || parts[0] != "realms" {
		return "", "", fmt.Errorf("URL does not contain /realms/<realm> pattern: %s", rawURL)
	}

	realm = parts[1]
	baseURL = fmt.Sprintf("%s://%s", parsed.Scheme, parsed.Host)

	return baseURL, realm, nil
}

// ServiceAccountOptions provides additional configuration when creating a client
// from a service account. These options supplement the authentication credentials
// from the service account file.
type ServiceAccountOptions struct {
	// BaseURL is the Hyperfluid API base URL (required).
	// Example: "https://api.hyperfluid.cloud"
	BaseURL string

	// ControlPlaneURL is the Control Plane API base URL (optional).
	// If not set, defaults to BaseURL.
	// Example: "https://console.hyperfluid.cloud"
	ControlPlaneURL string

	// OrgID is the default organization ID for API requests (optional).
	// If set, this will be used as the default for operations requiring an org ID.
	OrgID string

	// DataDockID is the default DataDock ID for query operations (optional).
	// If set, this will be used as the default for query operations.
	DataDockID string

	// SkipTLSVerify disables TLS certificate verification (optional).
	// WARNING: Only use this for development/testing. Never in production.
	SkipTLSVerify bool

	// RequestTimeout specifies the timeout for HTTP requests (optional).
	// Defaults to 30 seconds if not specified.
	RequestTimeout int

	// MaxRetries specifies the maximum number of retry attempts for failed requests (optional).
	// Defaults to 3 if not specified.
	MaxRetries int

	// MinIOEndpoint is the MinIO endpoint for S3 operations (required).
	MinIOEndpoint string

	// MinIOAccessKey is the MinIO access key for S3 operations (required).
	MinIOAccessKey string

	// MinIOSecretKey is the MinIO secret key for S3 operations (required).
	MinIOSecretKey string

	// MinIORegion is the MinIO region for S3 operations (required).
	MinIORegion string
}

// ToConfiguration converts the ServiceAccount to a utils.Configuration.
// This is used internally when creating a client from a service account.
func (sa *ServiceAccount) ToConfiguration(opts ServiceAccountOptions) (utils.Configuration, error) {
	baseURL, realm, err := sa.ParseIssuer()
	if err != nil {
		return utils.Configuration{}, fmt.Errorf("failed to parse issuer: %w", err)
	}

	// Default ControlPlaneURL to BaseURL if not specified
	controlPlaneURL := opts.ControlPlaneURL
	if controlPlaneURL == "" {
		controlPlaneURL = opts.BaseURL
	}

	cfg := utils.Configuration{
		BaseURL:              opts.BaseURL,
		ControlPlaneURL:      controlPlaneURL,
		OrgID:                opts.OrgID,
		DataDockID:           opts.DataDockID,
		SkipTLSVerify:        opts.SkipTLSVerify,
		KeycloakBaseURL:      baseURL,
		KeycloakRealm:        realm,
		KeycloakClientID:     sa.ClientID,
		KeycloakClientSecret: sa.ClientSecret,
		MinIOEndpoint:        opts.MinIOEndpoint,
		MinIOAccessKey:       opts.MinIOAccessKey,
		MinIOSecretKey:       opts.MinIOSecretKey,
		MinIORegion:          opts.MinIORegion,
	}

	// Apply defaults for optional fields
	if opts.RequestTimeout > 0 {
		cfg.RequestTimeout = utils.SecondsToDuration(opts.RequestTimeout)
	} else {
		cfg.RequestTimeout = utils.DefaultRequestTimeout
	}

	if opts.MaxRetries > 0 {
		cfg.MaxRetries = opts.MaxRetries
	} else {
		cfg.MaxRetries = utils.DefaultMaxRetries
	}

	return cfg, nil
}
