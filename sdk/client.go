package sdk

import (
	"context"
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

// NewClient creates a new Bifrost client.
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
