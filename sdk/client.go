package sdk

import (
	"context"
	"fmt"
	"net/http"

	"github.com/nudibranches-tech/bifrost-hyperfluid-sdk-dev/sdk/builders/fluent"
	"github.com/nudibranches-tech/bifrost-hyperfluid-sdk-dev/sdk/builders/progressive"
	"github.com/nudibranches-tech/bifrost-hyperfluid-sdk-dev/sdk/utils"
)

// Client is the main entry point for the SDK.
type Client struct {
	config     utils.Configuration
	httpClient *http.Client
}

// NewClient creates a new Bifrost client with the provided configuration.
func NewClient(config utils.Configuration) *Client {
	// Create a copy of the configuration to avoid side effects
	cfg := config
	return &Client{
		config: cfg,
		httpClient: utils.CreateHTTPClientWithSettings(
			cfg.SkipTLSVerify,
			cfg.RequestTimeout,
		),
	}
}

// NewClientFromServiceAccount creates a new Bifrost client using a ServiceAccount.
// This is the recommended way to create a client for service-to-service authentication.
//
// Example:
//
//	// Load service account from file (e.g., Kubernetes mounted secret)
//	sa, err := sdk.LoadServiceAccount("/var/run/secrets/hyperfluid/service_account.json")
//	if err != nil {
//	    log.Fatalf("Failed to load service account: %v", err)
//	}
//
//	// Create client
//	client, err := sdk.NewClientFromServiceAccount(sa, sdk.ServiceAccountOptions{
//	    BaseURL: "https://api.hyperfluid.cloud",
//	    OrgID:   "my-org-id",
//	})
//	if err != nil {
//	    log.Fatalf("Failed to create client: %v", err)
//	}
func NewClientFromServiceAccount(sa *ServiceAccount, opts ServiceAccountOptions) (*Client, error) {
	if sa == nil {
		return nil, fmt.Errorf("service account is nil")
	}

	if opts.BaseURL == "" {
		return nil, fmt.Errorf("BaseURL is required in ServiceAccountOptions")
	}

	cfg, err := sa.ToConfiguration(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create configuration from service account: %w", err)
	}

	return NewClient(cfg), nil
}

// NewClientFromServiceAccountFile creates a new Bifrost client by loading a ServiceAccount
// from a JSON file. This is a convenience function that combines LoadServiceAccount and
// NewClientFromServiceAccount.
//
// This is ideal for Kubernetes deployments where secrets are mounted as files:
//
//	client, err := sdk.NewClientFromServiceAccountFile(
//	    "/var/run/secrets/hyperfluid/service_account.json",
//	    sdk.ServiceAccountOptions{
//	        BaseURL: "https://api.hyperfluid.cloud",
//	    },
//	)
func NewClientFromServiceAccountFile(path string, opts ServiceAccountOptions) (*Client, error) {
	sa, err := LoadServiceAccount(path)
	if err != nil {
		return nil, err
	}
	return NewClientFromServiceAccount(sa, opts)
}

// NewClientFromServiceAccountJSON creates a new Bifrost client by parsing a ServiceAccount
// from a JSON string. This is useful when the service account is provided via environment
// variables.
//
// Example:
//
//	saJSON := os.Getenv("HYPERFLUID_SERVICE_ACCOUNT")
//	client, err := sdk.NewClientFromServiceAccountJSON(saJSON, sdk.ServiceAccountOptions{
//	    BaseURL: os.Getenv("HYPERFLUID_API_URL"),
//	})
func NewClientFromServiceAccountJSON(jsonStr string, opts ServiceAccountOptions) (*Client, error) {
	sa, err := LoadServiceAccountFromJSON(jsonStr)
	if err != nil {
		return nil, err
	}
	return NewClientFromServiceAccount(sa, opts)
}

// Do executes an HTTP request (implements the interface needed by builders)
func (c *Client) Do(ctx context.Context, method, endpoint string, body []byte) (*utils.Response, error) {
	return c.do(ctx, method, endpoint, body)
}

// GetConfig returns the client configuration (implements the interface needed by builders)
func (c *Client) GetConfig() utils.Configuration {
	return c.config
}

// Query creates a new QueryBuilder for fluent query construction.
// Example:
//
//	resp, err := client.Query().
//	    Catalog("sales").
//	    Schema("public").
//	    Table("orders").
//	    Limit(10).
//	    Get(ctx)
func (c *Client) Query() *fluent.QueryBuilder {
	return fluent.NewQueryBuilder(c)
}

func (c *Client) S3() (*fluent.S3Builder, error) {
	return fluent.NewS3Builder(c)
}

// Catalog starts a new fluent query with the catalog name.
// This is a shortcut for client.Query().DataDock(defaultID).Catalog(name).
// Uses DataDockID from config if available.
func (c *Client) Catalog(name string) *fluent.QueryBuilder {
	qb := fluent.NewQueryBuilder(c)
	// Auto-set DataDockID from config if available
	if c.config.DataDockID != "" {
		qb = qb.DataDock(c.config.DataDockID)
	}
	return qb.Catalog(name)
}

// Org uses the progressive builder pattern for type-safe navigation:
//
//	client.Org(id).Harbor(id).DataDock(id).Catalog(name).Schema(name).Table(name)
//
// Each level provides contextual methods:
//   - Org: ListHarbors(), CreateHarbor(), ListDataDocks()
//   - Harbor: ListDataDocks(), CreateDataDock(), Delete()
//   - DataDock: GetCatalog(), RefreshCatalog(), WakeUp(), Sleep()
//   - Catalog: Schema(), ListSchemas()
//   - Schema: Table(), ListTables()
//   - Table: Select(), Where(), Limit(), Get()
func (c *Client) Org(orgID string) *progressive.OrgBuilder {
	return &progressive.OrgBuilder{
		Client: c,
		OrgID:  orgID,
	}
}

// DataDock starts a new fluent query with the data dock ID.
// This allows starting queries with: client.DataDock(id).Catalog(...).Schema(...).Table(...)
// This is for FLUENT API (data queries).
// Example:
//
//	resp, err := client.DataDock("datadock-id").
//	    Catalog("sales").
//	    Schema("public").
//	    Table("orders").
//	    Limit(10).
//	    Get(ctx)
func (c *Client) DataDock(dataDockID string) *fluent.QueryBuilder {
	return fluent.NewQueryBuilder(c).DataDock(dataDockID)
}

// OrgFromConfig creates an OrgBuilder using the OrgID from the client configuration.
// This is a convenience method when you always use the same organization.
func (c *Client) OrgFromConfig() *progressive.OrgBuilder {
	return &progressive.OrgBuilder{
		Client: c,
		OrgID:  c.config.OrgID,
	}
}

// Search creates a new SearchBuilder for full-text search queries.
// Example:
//
//	resp, err := client.Search().
//	    Query("machine learning").
//	    DataDock("data-dock-id").
//	    Catalog("catalog").
//	    Schema("public").
//	    Table("documents").
//	    Columns("title", "content", "summary").
//	    Limit(10).
//	    Execute(ctx)
func (c *Client) Search() *fluent.SearchBuilder {
	return fluent.NewSearchBuilder(c)
}
