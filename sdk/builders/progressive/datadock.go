package progressive

import (
	"context"
	"fmt"
	"net/url"

	"github.com/nudibranches-tech/hyperfluid-sdk-go/sdk/builders"
	"github.com/nudibranches-tech/hyperfluid-sdk-go/sdk/utils"
)

// DataDockBuilder represents a datadock context.
// Available methods:
//   - Catalog(name) - Navigate to a specific catalog
//   - GetCatalog(ctx) - Get the full catalog metadata
//   - RefreshCatalog(ctx) - Trigger catalog introspection
//   - WakeUp(ctx) - Bring datadock online
//   - Sleep(ctx) - Put datadock to sleep
//   - Get(ctx) - Get datadock details
//   - Update(ctx, config) - Update datadock configuration
//   - Delete(ctx) - Delete this datadock
type DataDockBuilder struct {
	client     builders.ClientInterface
	orgID      string
	harborID   string
	dataDockID string
}

// Catalog navigates to a specific catalog in this datadock.
func (d *DataDockBuilder) Catalog(catalogName string) *CatalogBuilder {
	return &CatalogBuilder{
		client:      d.client,
		orgID:       d.orgID,
		dataDockID:  d.dataDockID,
		catalogName: catalogName,
	}
}

// GetCatalog retrieves the full catalog metadata (schemas, tables, columns).
func (d *DataDockBuilder) GetCatalog(ctx context.Context) (*utils.Response, error) {
	endpoint := fmt.Sprintf("%s/data-docks/%s/catalog",
		d.client.GetConfig().BaseURL,
		url.PathEscape(d.dataDockID),
	)
	return d.client.Do(ctx, "GET", endpoint, nil)
}

// RefreshCatalog triggers catalog introspection and updates metadata.
func (d *DataDockBuilder) RefreshCatalog(ctx context.Context) (*utils.Response, error) {
	endpoint := fmt.Sprintf("%s/data-docks/%s/catalog/refresh",
		d.client.GetConfig().BaseURL,
		url.PathEscape(d.dataDockID),
	)
	return d.client.Do(ctx, "POST", endpoint, nil)
}

// WakeUp brings the datadock online (for TrinoInternal/MinioInternal).
func (d *DataDockBuilder) WakeUp(ctx context.Context) (*utils.Response, error) {
	endpoint := fmt.Sprintf("%s/data-docks/%s/wake-up",
		d.client.GetConfig().BaseURL,
		url.PathEscape(d.dataDockID),
	)
	return d.client.Do(ctx, "POST", endpoint, nil)
}

// Sleep puts the datadock to sleep (cost optimization).
func (d *DataDockBuilder) Sleep(ctx context.Context) (*utils.Response, error) {
	endpoint := fmt.Sprintf("%s/data-docks/%s/sleep",
		d.client.GetConfig().BaseURL,
		url.PathEscape(d.dataDockID),
	)
	return d.client.Do(ctx, "POST", endpoint, nil)
}

// Get retrieves datadock details.
func (d *DataDockBuilder) Get(ctx context.Context) (*utils.Response, error) {
	endpoint := fmt.Sprintf("%s/data-docks/%s",
		d.client.GetConfig().BaseURL,
		url.PathEscape(d.dataDockID),
	)
	return d.client.Do(ctx, "GET", endpoint, nil)
}

// Update modifies datadock configuration.
func (d *DataDockBuilder) Update(ctx context.Context, config map[string]interface{}) (*utils.Response, error) {
	endpoint := fmt.Sprintf("%s/data-docks/%s",
		d.client.GetConfig().BaseURL,
		url.PathEscape(d.dataDockID),
	)
	body := utils.JsonMarshal(config)
	return d.client.Do(ctx, "PATCH", endpoint, body)
}

// Delete removes this datadock.
func (d *DataDockBuilder) Delete(ctx context.Context) (*utils.Response, error) {
	endpoint := fmt.Sprintf("%s/data-docks/%s",
		d.client.GetConfig().BaseURL,
		url.PathEscape(d.dataDockID),
	)
	return d.client.Do(ctx, "DELETE", endpoint, nil)
}

// Search starts a search builder for this datadock.
// Returns a SearchBuilder that can be used to build and execute full-text search queries.
func (d *DataDockBuilder) Search(query string) *SearchBuilder {
	return &SearchBuilder{
		client:         d.client,
		dataDockID:     d.dataDockID,
		searchQuery:    query,
		columnsToIndex: []string{},
		limitVal:       20, // Default limit
	}
}
