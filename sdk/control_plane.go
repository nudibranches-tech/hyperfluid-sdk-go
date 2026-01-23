package sdk

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"sync"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"

	"github.com/nudibranches-tech/hyperfluid-sdk-go/sdk/controlplaneapiclient"
)

// ControlPlaneClient wraps the generated OpenAPI client with automatic OAuth2 token management.
type ControlPlaneClient struct {
	*controlplaneapiclient.ClientWithResponses
	httpClient *http.Client
	tokenURL   string
}

// controlPlaneClientCache stores lazily-initialized control plane clients per SDK client.
var (
	controlPlaneClients = make(map[*Client]*ControlPlaneClient)
	controlPlaneMu      sync.RWMutex
)

// ControlPlane returns a Control Plane API client with automatic OAuth2 authentication.
// The client is lazily initialized and cached for subsequent calls.
//
// The client uses the OAuth2 Client Credentials flow to automatically obtain and refresh
// tokens using the Keycloak credentials from the SDK configuration.
//
// Example:
//
//	client, _ := sdk.NewClientFromServiceAccountFile("/path/to/sa.json", opts)
//	cp := client.ControlPlane()
//
//	// List all data docks
//	resp, err := cp.ListDataDocksWithResponse(ctx)
//	if err != nil {
//	    log.Fatalf("Failed to list data docks: %v", err)
//	}
//	for _, dock := range *resp.JSON200 {
//	    fmt.Printf("DataDock: %s\n", dock.Name)
//	}
func (c *Client) ControlPlane() (*ControlPlaneClient, error) {
	// Check cache first
	controlPlaneMu.RLock()
	if cp, ok := controlPlaneClients[c]; ok {
		controlPlaneMu.RUnlock()
		return cp, nil
	}
	controlPlaneMu.RUnlock()

	// Initialize new client
	controlPlaneMu.Lock()
	defer controlPlaneMu.Unlock()

	// Double-check after acquiring write lock
	if cp, ok := controlPlaneClients[c]; ok {
		return cp, nil
	}

	cp, err := newControlPlaneClient(c)
	if err != nil {
		return nil, err
	}

	controlPlaneClients[c] = cp
	return cp, nil
}

// newControlPlaneClient creates a new ControlPlaneClient with OAuth2 authentication.
func newControlPlaneClient(c *Client) (*ControlPlaneClient, error) {
	if c.config.ControlPlaneURL == "" {
		return nil, fmt.Errorf("ControlPlaneURL is not configured")
	}

	if c.config.KeycloakClientID == "" || c.config.KeycloakClientSecret == "" {
		return nil, fmt.Errorf("keycloak client credentials are not configured")
	}

	if c.config.KeycloakBaseURL == "" || c.config.KeycloakRealm == "" {
		return nil, fmt.Errorf("keycloak base URL or realm is not configured")
	}

	tokenURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token",
		c.config.KeycloakBaseURL, c.config.KeycloakRealm)

	// Configure OAuth2 Client Credentials
	oauthConfig := &clientcredentials.Config{
		ClientID:     c.config.KeycloakClientID,
		ClientSecret: c.config.KeycloakClientSecret,
		TokenURL:     tokenURL,
		Scopes:       []string{}, // Add scopes if needed
	}

	// Create a base HTTP client with TLS configuration
	baseTransport := http.DefaultTransport.(*http.Transport).Clone()
	if c.config.SkipTLSVerify {
		baseTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	// Create context with custom HTTP client for OAuth2 token requests
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, &http.Client{
		Transport: baseTransport,
		Timeout:   c.config.RequestTimeout,
	})

	// Get OAuth2 HTTP client with automatic token refresh
	httpClient := oauthConfig.Client(ctx)
	httpClient.Timeout = c.config.RequestTimeout

	// Create the generated OpenAPI client
	apiClient, err := controlplaneapiclient.NewClientWithResponses(
		c.config.ControlPlaneURL,
		controlplaneapiclient.WithHTTPClient(httpClient),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create control plane client: %w", err)
	}

	return &ControlPlaneClient{
		ClientWithResponses: apiClient,
		httpClient:          httpClient,
		tokenURL:            tokenURL,
	}, nil
}

// Close releases resources associated with the ControlPlaneClient.
// This removes the client from the cache.
func (cp *ControlPlaneClient) Close() {
	controlPlaneMu.Lock()
	defer controlPlaneMu.Unlock()

	// Remove from cache by finding the key
	for client, cached := range controlPlaneClients {
		if cached == cp {
			delete(controlPlaneClients, client)
			break
		}
	}
}
