package sdk

import (
	"github.com/nudibranches-tech/bifrost-hyperfluid-sdk-dev/sdk/utils"
	"net/http"
)

// Client is the main entry point for the SDK.
type Client struct {
	config     utils.Configuration
	httpClient *http.Client
}

// NewClient creates a new Bifrost client.
func NewClient(config utils.Configuration) *Client {
	// we create a copy of the configuration to avoid side effects
	cfg := config
	return &Client{
		config: cfg,
		httpClient: utils.CreateHTTPClientWithSettings(
			cfg.SkipTLSVerify,
			cfg.RequestTimeout,
		),
	}
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
func (c *Client) Query() *QueryBuilder {
	return newQueryBuilder(c)
}

// Catalog starts a new fluent query with the catalog name.
// This is a shortcut for client.Query().Catalog(name).
func (c *Client) Catalog(name string) *QueryBuilder {
	return newQueryBuilder(c).Catalog(name)
}

// Org starts a new fluent API navigation with a specific organization.
// This uses the progressive builder pattern for type-safe navigation:
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
func (c *Client) Org(orgID string) *OrgBuilder {
	return &OrgBuilder{
		client: c,
		orgID:  orgID,
	}
}

// OrgFromConfig creates an OrgBuilder using the OrgID from the client configuration.
// This is a convenience method when you always use the same organization.
func (c *Client) OrgFromConfig() *OrgBuilder {
	return &OrgBuilder{
		client: c,
		orgID:  c.config.OrgID,
	}
}
