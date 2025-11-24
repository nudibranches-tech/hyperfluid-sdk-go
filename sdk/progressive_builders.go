package sdk

import (
	"context"
	"fmt"
	"net/url"

	"github.com/nudibranches-tech/bifrost-hyperfluid-sdk-dev/sdk/utils"
)

// Progressive Builders - Each level has its own type with specific methods
// This forces the correct order: Org → Harbor → DataDock → Catalog → Schema → Table

// ============================================================================
// Level 1: Organization Builder
// ============================================================================

// OrgBuilder represents an organization context.
// Available methods:
//   - Harbor(id) - Navigate to a specific harbor
//   - ListHarbors(ctx) - List all harbors in this org
//   - CreateHarbor(ctx, name) - Create a new harbor
//   - ListDataDocks(ctx) - List all datadocks across all harbors
type OrgBuilder struct {
	client *Client
	orgID  string
}

// Harbor navigates to a specific harbor in this organization.
func (o *OrgBuilder) Harbor(harborID string) *HarborBuilder {
	return &HarborBuilder{
		client:   o.client,
		orgID:    o.orgID,
		harborID: harborID,
	}
}

// ListHarbors retrieves all harbors in this organization.
func (o *OrgBuilder) ListHarbors(ctx context.Context) (*utils.Response, error) {
	endpoint := fmt.Sprintf("%s/%s/harbors",
		o.client.config.BaseURL,
		url.PathEscape(o.orgID),
	)
	return o.client.do(ctx, "GET", endpoint, nil)
}

// CreateHarbor creates a new harbor in this organization.
func (o *OrgBuilder) CreateHarbor(ctx context.Context, name string) (*utils.Response, error) {
	endpoint := fmt.Sprintf("%s/%s/harbors",
		o.client.config.BaseURL,
		url.PathEscape(o.orgID),
	)
	body := utils.JsonMarshal(map[string]interface{}{
		"name": name,
	})
	return o.client.do(ctx, "POST", endpoint, body)
}

// ListDataDocks retrieves all datadocks across all harbors in this organization.
func (o *OrgBuilder) ListDataDocks(ctx context.Context) (*utils.Response, error) {
	endpoint := fmt.Sprintf("%s/%s/data-docks",
		o.client.config.BaseURL,
		url.PathEscape(o.orgID),
	)
	return o.client.do(ctx, "GET", endpoint, nil)
}

// RefreshAllDataDocks triggers a catalog refresh on all datadocks in this organization.
func (o *OrgBuilder) RefreshAllDataDocks(ctx context.Context) (*utils.Response, error) {
	endpoint := fmt.Sprintf("%s/%s/data-docks/refresh",
		o.client.config.BaseURL,
		url.PathEscape(o.orgID),
	)
	return o.client.do(ctx, "POST", endpoint, nil)
}

// ============================================================================
// Level 2: Harbor Builder
// ============================================================================

// HarborBuilder represents a harbor context.
// Available methods:
//   - DataDock(id) - Navigate to a specific datadock
//   - ListDataDocks(ctx) - List all datadocks in this harbor
//   - CreateDataDock(ctx, config) - Create a new datadock
//   - Delete(ctx) - Delete this harbor
type HarborBuilder struct {
	client   *Client
	orgID    string
	harborID string
}

// DataDock navigates to a specific datadock in this harbor.
func (h *HarborBuilder) DataDock(dataDockID string) *DataDockBuilder {
	return &DataDockBuilder{
		client:     h.client,
		orgID:      h.orgID,
		harborID:   h.harborID,
		dataDockID: dataDockID,
	}
}

// ListDataDocks retrieves all datadocks in this harbor.
func (h *HarborBuilder) ListDataDocks(ctx context.Context) (*utils.Response, error) {
	endpoint := fmt.Sprintf("%s/harbors/%s/data-docks",
		h.client.config.BaseURL,
		url.PathEscape(h.harborID),
	)
	return h.client.do(ctx, "GET", endpoint, nil)
}

// CreateDataDock creates a new datadock in this harbor.
func (h *HarborBuilder) CreateDataDock(ctx context.Context, config map[string]interface{}) (*utils.Response, error) {
	// Ensure harbor_id is set
	config["harbor_id"] = h.harborID

	endpoint := fmt.Sprintf("%s/data-docks", h.client.config.BaseURL)
	body := utils.JsonMarshal(config)
	return h.client.do(ctx, "POST", endpoint, body)
}

// Delete removes this harbor.
func (h *HarborBuilder) Delete(ctx context.Context) (*utils.Response, error) {
	endpoint := fmt.Sprintf("%s/harbors/%s",
		h.client.config.BaseURL,
		url.PathEscape(h.harborID),
	)
	return h.client.do(ctx, "DELETE", endpoint, nil)
}

// ============================================================================
// Level 3: DataDock Builder
// ============================================================================

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
	client     *Client
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
		d.client.config.BaseURL,
		url.PathEscape(d.dataDockID),
	)
	return d.client.do(ctx, "GET", endpoint, nil)
}

// RefreshCatalog triggers catalog introspection and updates metadata.
func (d *DataDockBuilder) RefreshCatalog(ctx context.Context) (*utils.Response, error) {
	endpoint := fmt.Sprintf("%s/data-docks/%s/catalog/refresh",
		d.client.config.BaseURL,
		url.PathEscape(d.dataDockID),
	)
	return d.client.do(ctx, "POST", endpoint, nil)
}

// WakeUp brings the datadock online (for TrinoInternal/MinioInternal).
func (d *DataDockBuilder) WakeUp(ctx context.Context) (*utils.Response, error) {
	endpoint := fmt.Sprintf("%s/data-docks/%s/wake-up",
		d.client.config.BaseURL,
		url.PathEscape(d.dataDockID),
	)
	return d.client.do(ctx, "POST", endpoint, nil)
}

// Sleep puts the datadock to sleep (cost optimization).
func (d *DataDockBuilder) Sleep(ctx context.Context) (*utils.Response, error) {
	endpoint := fmt.Sprintf("%s/data-docks/%s/sleep",
		d.client.config.BaseURL,
		url.PathEscape(d.dataDockID),
	)
	return d.client.do(ctx, "POST", endpoint, nil)
}

// Get retrieves datadock details.
func (d *DataDockBuilder) Get(ctx context.Context) (*utils.Response, error) {
	endpoint := fmt.Sprintf("%s/data-docks/%s",
		d.client.config.BaseURL,
		url.PathEscape(d.dataDockID),
	)
	return d.client.do(ctx, "GET", endpoint, nil)
}

// Update modifies datadock configuration.
func (d *DataDockBuilder) Update(ctx context.Context, config map[string]interface{}) (*utils.Response, error) {
	endpoint := fmt.Sprintf("%s/data-docks/%s",
		d.client.config.BaseURL,
		url.PathEscape(d.dataDockID),
	)
	body := utils.JsonMarshal(config)
	return d.client.do(ctx, "PATCH", endpoint, body)
}

// Delete removes this datadock.
func (d *DataDockBuilder) Delete(ctx context.Context) (*utils.Response, error) {
	endpoint := fmt.Sprintf("%s/data-docks/%s",
		d.client.config.BaseURL,
		url.PathEscape(d.dataDockID),
	)
	return d.client.do(ctx, "DELETE", endpoint, nil)
}

// ============================================================================
// Level 4: Catalog Builder
// ============================================================================

// CatalogBuilder represents a catalog context.
// Available methods:
//   - Schema(name) - Navigate to a specific schema
//   - ListSchemas(ctx) - List all schemas in this catalog
type CatalogBuilder struct {
	client      *Client
	orgID       string
	dataDockID  string
	catalogName string
}

// Schema navigates to a specific schema in this catalog.
func (c *CatalogBuilder) Schema(schemaName string) *SchemaBuilder {
	return &SchemaBuilder{
		client:      c.client,
		orgID:       c.orgID,
		dataDockID:  c.dataDockID,
		catalogName: c.catalogName,
		schemaName:  schemaName,
	}
}

// ListSchemas retrieves all schemas in this catalog.
// This parses the catalog metadata to extract schemas.
func (c *CatalogBuilder) ListSchemas(ctx context.Context) ([]string, error) {
	// Get full catalog metadata
	endpoint := fmt.Sprintf("%s/data-docks/%s/catalog",
		c.client.config.BaseURL,
		url.PathEscape(c.dataDockID),
	)

	resp, err := c.client.do(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	// Extract schemas for this catalog
	var schemas []string
	if catalogs, ok := resp.Data.(map[string]interface{})["catalogs"].([]interface{}); ok {
		for _, cat := range catalogs {
			if catMap, ok := cat.(map[string]interface{}); ok {
				if catMap["catalog_name"] == c.catalogName {
					if schemaList, ok := catMap["schemas"].([]interface{}); ok {
						for _, s := range schemaList {
							if sMap, ok := s.(map[string]interface{}); ok {
								if name, ok := sMap["schema_name"].(string); ok {
									schemas = append(schemas, name)
								}
							}
						}
					}
				}
			}
		}
	}

	return schemas, nil
}

// ============================================================================
// Level 5: Schema Builder
// ============================================================================

// SchemaBuilder represents a schema context.
// Available methods:
//   - Table(name) - Navigate to a specific table (returns TableQueryBuilder for querying)
//   - ListTables(ctx) - List all tables in this schema
type SchemaBuilder struct {
	client      *Client
	orgID       string
	dataDockID  string
	catalogName string
	schemaName  string
}

// Table navigates to a specific table in this schema.
// Returns a TableQueryBuilder which supports both queries and operations.
func (s *SchemaBuilder) Table(tableName string) *TableQueryBuilder {
	return &TableQueryBuilder{
		client:      s.client,
		orgID:       s.orgID,
		catalogName: s.catalogName,
		schemaName:  s.schemaName,
		tableName:   tableName,
		// Query builder fields
		selectCols: []string{},
		filters:    []Filter{},
		orderBy:    []OrderClause{},
		rawParams:  url.Values{},
	}
}

// ListTables retrieves all tables in this schema.
func (s *SchemaBuilder) ListTables(ctx context.Context) ([]string, error) {
	// Get full catalog metadata
	endpoint := fmt.Sprintf("%s/data-docks/%s/catalog",
		s.client.config.BaseURL,
		url.PathEscape(s.dataDockID),
	)

	resp, err := s.client.do(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	// Extract tables for this schema
	var tables []string
	if catalogs, ok := resp.Data.(map[string]interface{})["catalogs"].([]interface{}); ok {
		for _, cat := range catalogs {
			if catMap, ok := cat.(map[string]interface{}); ok {
				if catMap["catalog_name"] == s.catalogName {
					if schemaList, ok := catMap["schemas"].([]interface{}); ok {
						for _, sch := range schemaList {
							if schMap, ok := sch.(map[string]interface{}); ok {
								if schMap["schema_name"] == s.schemaName {
									if tableList, ok := schMap["tables"].([]interface{}); ok {
										for _, t := range tableList {
											if tMap, ok := t.(map[string]interface{}); ok {
												if name, ok := tMap["table_name"].(string); ok {
													tables = append(tables, name)
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}

	return tables, nil
}

// ============================================================================
// Level 6: Table Query Builder
// ============================================================================

// TableQueryBuilder combines table navigation with query building.
// This is the final level where you can build queries AND execute them.
// Inherits all query building methods from the original QueryBuilder.
type TableQueryBuilder struct {
	client *Client
	orgID  string

	// Table location
	catalogName string
	schemaName  string
	tableName   string

	// Query parameters (same as QueryBuilder)
	selectCols []string
	filters    []Filter
	orderBy    []OrderClause
	limitVal   int
	offsetVal  int
	rawParams  url.Values
}

// Query building methods - same as original QueryBuilder
// These return *TableQueryBuilder for chaining

func (t *TableQueryBuilder) Select(columns ...string) *TableQueryBuilder {
	t.selectCols = append(t.selectCols, columns...)
	return t
}

func (t *TableQueryBuilder) Where(column, operator string, value interface{}) *TableQueryBuilder {
	t.filters = append(t.filters, Filter{
		Column:   column,
		Operator: operator,
		Value:    value,
	})
	return t
}

func (t *TableQueryBuilder) OrderBy(column, direction string) *TableQueryBuilder {
	if direction == "" {
		direction = "ASC"
	}
	t.orderBy = append(t.orderBy, OrderClause{
		Column:    column,
		Direction: direction,
	})
	return t
}

func (t *TableQueryBuilder) Limit(n int) *TableQueryBuilder {
	t.limitVal = n
	return t
}

func (t *TableQueryBuilder) Offset(n int) *TableQueryBuilder {
	t.offsetVal = n
	return t
}

func (t *TableQueryBuilder) RawParams(params url.Values) *TableQueryBuilder {
	for key, values := range params {
		for _, value := range values {
			t.rawParams.Add(key, value)
		}
	}
	return t
}

// Execution method - builds the query and executes it

func (t *TableQueryBuilder) Get(ctx context.Context) (*utils.Response, error) {
	// Build endpoint using Bifrost OpenAPI format
	endpoint := fmt.Sprintf(
		"%s/%s/openapi/%s/%s/%s",
		t.client.config.BaseURL,
		url.PathEscape(t.orgID),
		url.PathEscape(t.catalogName),
		url.PathEscape(t.schemaName),
		url.PathEscape(t.tableName),
	)

	// Build query parameters using the same logic as QueryBuilder
	params := t.buildParams()

	if len(params) > 0 {
		endpoint += "?" + params.Encode()
	}

	return t.client.do(ctx, "GET", endpoint, nil)
}

// buildParams constructs query parameters (same as QueryBuilder)
func (t *TableQueryBuilder) buildParams() url.Values {
	params := url.Values{}

	// Copy raw params first
	for key, values := range t.rawParams {
		for _, value := range values {
			params.Add(key, value)
		}
	}

	// Add SELECT columns
	if len(t.selectCols) > 0 {
		params.Set("select", fmt.Sprintf("%s", t.selectCols))
	}

	// Add WHERE filters
	for _, filter := range t.filters {
		paramName := fmt.Sprintf("%s[%s]", filter.Column, filter.Operator)
		params.Add(paramName, fmt.Sprintf("%v", filter.Value))
	}

	// Add ORDER BY
	if len(t.orderBy) > 0 {
		var orderParts []string
		for _, order := range t.orderBy {
			if order.Direction == "DESC" {
				orderParts = append(orderParts, fmt.Sprintf("%s.desc", order.Column))
			} else {
				orderParts = append(orderParts, fmt.Sprintf("%s.asc", order.Column))
			}
		}
		params.Set("order", fmt.Sprintf("%s", orderParts))
	}

	// Add LIMIT
	if t.limitVal > 0 {
		params.Set("_limit", fmt.Sprintf("%d", t.limitVal))
	}

	// Add OFFSET
	if t.offsetVal > 0 {
		params.Set("_offset", fmt.Sprintf("%d", t.offsetVal))
	}

	return params
}
