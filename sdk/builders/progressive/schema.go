package progressive

import (
	"context"
	"fmt"
	"net/url"

	"github.com/nudibranches-tech/hyperfluid-sdk-go/sdk/builders"
)

// SchemaBuilder represents a schema context.
// Available methods:
//   - Table(name) - Navigate to a specific table (returns TableQueryBuilder for querying)
//   - ListTables(ctx) - List all tables in this schema
type SchemaBuilder struct {
	client      builders.ClientInterface
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
		filters:    []builders.Filter{},
		orderBy:    []builders.OrderClause{},
		rawParams:  url.Values{},
	}
}

// ListTables retrieves all tables in this schema.
func (s *SchemaBuilder) ListTables(ctx context.Context) ([]string, error) {
	// Get full catalog metadata
	endpoint := fmt.Sprintf("%s/data-docks/%s/catalog",
		s.client.GetConfig().BaseURL,
		url.PathEscape(s.dataDockID),
	)

	resp, err := s.client.Do(ctx, "GET", endpoint, nil)
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
