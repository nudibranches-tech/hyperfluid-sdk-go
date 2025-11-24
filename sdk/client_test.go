package sdk

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/nudibranches-tech/bifrost-hyperfluid-sdk-dev/sdk/utils"
)

func TestNewClient(t *testing.T) {
	config := utils.Configuration{
		BaseURL: "http://localhost",
		OrgID:   "test-org",
		Token:   "test-token",
	}
	client := NewClient(config)

	if client == nil {
		t.Fatal("NewClient should not return nil")
	}
	if client.config.BaseURL != "http://localhost" {
		t.Errorf("Expected BaseURL to be 'http://localhost', got '%s'", client.config.BaseURL)
	}
}

func TestCatalogMethod(t *testing.T) {
	client := NewClient(utils.Configuration{OrgID: "test-org"})
	qb := client.Catalog("test-catalog")

	if qb == nil {
		t.Fatal("Catalog should not return nil")
	}
	if qb.catalogName != "test-catalog" {
		t.Errorf("Expected catalog name to be 'test-catalog', got '%s'", qb.catalogName)
	}
	if qb.client != client {
		t.Error("QueryBuilder client should be the same as the parent client")
	}
}

// mockRoundTripper is used to mock HTTP responses in tests.
type mockRoundTripper struct {
	roundTripFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.roundTripFunc(req)
}

func newTestClient(config utils.Configuration, handler func(req *http.Request) (*http.Response, error)) *Client {
	return &Client{
		config: config,
		httpClient: &http.Client{
			Transport: &mockRoundTripper{roundTripFunc: handler},
		},
	}
}

func TestFluentAPI_Success(t *testing.T) {
	client := newTestClient(utils.Configuration{Token: "test-token", OrgID: "test-org"}, func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"data": "test"}`)),
		}, nil
	})

	resp, err := client.Catalog("c").Schema("s").Table("t").Get(context.Background())

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if resp.Status != utils.StatusOK {
		t.Errorf("Expected status OK, got %s", resp.Status)
	}
	if resp.Data.(map[string]interface{})["data"] != "test" {
		t.Errorf("Unexpected response data: %v", resp.Data)
	}
}

func TestFluentAPI_NotFound(t *testing.T) {
	client := newTestClient(utils.Configuration{Token: "test-token", OrgID: "test-org"}, func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusNotFound,
			Body:       io.NopCloser(strings.NewReader("")),
		}, nil
	})

	_, err := client.Catalog("c").Schema("s").Table("t").Get(context.Background())

	if !errors.Is(err, utils.ErrNotFound) {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}

func TestFluentAPI_PermissionDenied(t *testing.T) {
	client := newTestClient(utils.Configuration{Token: "test-token", OrgID: "test-org"}, func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusForbidden,
			Body:       io.NopCloser(strings.NewReader("")),
		}, nil
	})

	_, err := client.Catalog("c").Schema("s").Table("t").Get(context.Background())

	if !errors.Is(err, utils.ErrPermissionDenied) {
		t.Errorf("Expected ErrPermissionDenied, got %v", err)
	}
}

func TestFluentAPI_ServerError_Retry(t *testing.T) {
	reqCount := 0
	client := newTestClient(utils.Configuration{Token: "test-token", OrgID: "test-org", MaxRetries: 1}, func(req *http.Request) (*http.Response, error) {
		reqCount++
		if reqCount == 1 {
			return &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       io.NopCloser(strings.NewReader("")),
			}, nil
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"data": "success"}`)),
		}, nil
	})

	resp, err := client.Catalog("c").Schema("s").Table("t").Get(context.Background())

	if err != nil {
		t.Fatalf("Expected no error on retry, got %v", err)
	}
	if resp.Status != utils.StatusOK {
		t.Errorf("Expected status OK on retry, got %s", resp.Status)
	}
	if reqCount != 2 {
		t.Errorf("Expected 2 requests, got %d", reqCount)
	}
}
