package core

import (
	"bifrost-for-developers/sdk/core/utils"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func GetToken(context context.Context) string {
	var currentConfiguration = GetConfiguration()
	if currentConfiguration == nil {
		fmt.Printf("Error getting configuration: %v", utils.ErrorUninitialized)
		return ""
	}

	if currentConfiguration != nil && currentConfiguration.Token != "" {
		print("Using existing token from configuration")
		return currentConfiguration.Token
	}

	print("Fetching new token from Keycloak")
	newToken, err := FetchKeycloakToken(context)
	if err != nil {
		fmt.Printf("Error fetching Keycloak token: %v", err)
		return utils.GetEnvironmentVariable("HYPERFLUID_TOKEN", "")
	}

	return newToken
}

func FetchKeycloakToken(context context.Context) (string, error) {
	var currentConfiguration = GetConfiguration()
	if currentConfiguration == nil {
		return "", fmt.Errorf(utils.ErrorUninitialized.Error())
	}

	form := url.Values{
		"grant_type": {"password"},
		"client_id":  {currentConfiguration.KeycloakClientID},
		"username":   {currentConfiguration.KeycloakUsername},
		"password":   {currentConfiguration.KeycloakPassword},
	}

	req, err := http.NewRequestWithContext(
		context,
		"POST",
		fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", currentConfiguration.KeycloakBaseURL, currentConfiguration.KeycloakRealm),
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		return "", fmt.Errorf("cannot create Keycloak request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	if currentConfiguration.SkipTLSVerify {
		client.Transport = &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("cannot reach Keycloak: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("keycloak auth failed (%d): %s", resp.StatusCode, body)
	}

	var parsed map[string]any
	if err := json.Unmarshal(body, &parsed); err != nil {
		return "", fmt.Errorf("invalid Keycloak response: %w", err)
	}
	token, ok := parsed["access_token"].(string)
	if !ok || token == "" {
		return "", fmt.Errorf("missing access_token in Keycloak response")
	}

	UpdateToken(token)
	return token, nil
}

func UpdateToken(newToken string) {
	var currentConfiguration = GetConfiguration()
	if currentConfiguration == nil {
		fmt.Printf("Error getting configuration: %v", utils.ErrorUninitialized)
		return
	}
	globalConfiguration.Token = newToken
	fmt.Printf("Updated token: %s", newToken)
}
