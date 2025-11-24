package sdk

import (
	"github.com/nudibranches-tech/bifrost-hyperfluid-sdk-dev/sdk/utils"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	// "time" // time import is no longer needed
)

// authMutex protects token updates to prevent race conditions during refresh.
var authMutex sync.Mutex

func (c *Client) hasKeycloakPasswordGrantCredentials() bool {
	return c.config.KeycloakUsername != "" && c.config.KeycloakPassword != ""
}

func (c *Client) hasKeycloakClientCredentials() bool {
	return c.config.KeycloakClientID != "" && c.config.KeycloakClientSecret != ""
}

func (c *Client) isKeycloakAuthMethodConfigured() bool {
	return c.hasKeycloakPasswordGrantCredentials() || c.hasKeycloakClientCredentials()
}

// refreshToken attempts to refresh the access token using available Keycloak credentials.
func (c *Client) refreshToken(ctx context.Context) (string, error) {
	authMutex.Lock()
	defer authMutex.Unlock()

	// Note: This is a simplified implementation.
	// In production, you should:
	// 1. Parse JWT to check expiry
	// 2. Only refresh if token is actually expired or about to expire
	// 3. Store token expiry timestamp separately
	//
	// For now, we always refresh when this is called (typically on 401 errors)

	if c.hasKeycloakClientCredentials() {
		newToken, err := c.refreshAccessTokenClientCredentials(ctx)
		if err == nil {
			c.config.Token = newToken
			return newToken, nil
		}
		// Log error but try password grant as fallback if configured
		fmt.Printf("Client Credentials Grant failed: %v, attempting password grant...\n", err)
	}

	if c.hasKeycloakPasswordGrantCredentials() {
		newToken, err := c.refreshAccessTokenPasswordGrant(ctx)
		if err == nil {
			c.config.Token = newToken
			return newToken, nil
		}
		return "", fmt.Errorf("%w: password grant failed: %w", utils.ErrAuthenticationFailed, err)
	}

	return "", utils.ErrInvalidConfiguration
}

// refreshAccessTokenClientCredentials performs the Client Credentials Grant flow.
func (c *Client) refreshAccessTokenClientCredentials(ctx context.Context) (string, error) {
	form := url.Values{
		"grant_type":    {"client_credentials"},
		"client_id":     {c.config.KeycloakClientID},
		"client_secret": {c.config.KeycloakClientSecret},
	}
	return c.exchangeKeycloakToken(ctx, form)
}

// refreshAccessTokenPasswordGrant performs the Resource Owner Password Credentials Grant flow.
func (c *Client) refreshAccessTokenPasswordGrant(ctx context.Context) (string, error) {
	form := url.Values{
		"grant_type": {"password"},
		"client_id":  {c.config.KeycloakClientID},
		"username":   {c.config.KeycloakUsername},
		"password":   {c.config.KeycloakPassword},
	}
	return c.exchangeKeycloakToken(ctx, form)
}

// exchangeKeycloakToken sends the request to Keycloak's token endpoint.
func (c *Client) exchangeKeycloakToken(ctx context.Context, form url.Values) (string, error) {
	if c.config.KeycloakBaseURL == "" || c.config.KeycloakRealm == "" {
		return "", fmt.Errorf("%w: Keycloak base URL or realm not configured", utils.ErrInvalidConfiguration)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", c.config.KeycloakBaseURL, c.config.KeycloakRealm),
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		return "", fmt.Errorf("%w: cannot create Keycloak request: %w", utils.ErrInvalidRequest, err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Use a dedicated HTTP client for Keycloak to avoid potential deadlocks
	// if the main client's transport relies on token refresh itself.
	keycloakClient := &http.Client{
		Timeout: c.config.RequestTimeout, // Use the same timeout as main requests
	}
	if c.config.SkipTLSVerify {
		keycloakClient.Transport = &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	}

	resp, err := keycloakClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("%w: cannot reach Keycloak: %w", utils.ErrAuthenticationFailed, err)
	}

	// Read body and close immediately
	body, _ := io.ReadAll(resp.Body) // io.ReadAll already handles errors internally to return empty slice
	_ = resp.Body.Close()            // Always close after reading (error ignored - we already have the body)

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%w: Keycloak token exchange failed (%d): %s", utils.ErrAuthenticationFailed, resp.StatusCode, body)
	}

	var parsed map[string]any
	if err := json.Unmarshal(body, &parsed); err != nil {
		return "", fmt.Errorf("%w: invalid Keycloak response: %w", utils.ErrAuthenticationFailed, err)
	}
	token, ok := parsed["access_token"].(string)
	if !ok || token == "" {
		return "", fmt.Errorf("%w: missing access_token in Keycloak response", utils.ErrAuthenticationFailed)
	}

	return token, nil
}
