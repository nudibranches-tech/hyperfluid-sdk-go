package progressive

import (
	"context"
	"fmt"
	"net/url"

	"github.com/nudibranches-tech/hyperfluid-sdk-go/sdk/builders"
)

// CatalogBuilder represents a catalog context.
// Available methods:
//   - Schema(name) - Navigate to a specific schema
//   - ListSchemas(ctx) - List all schemas in this catalog
type CatalogBuilder struct {
	client      builders.ClientInterface
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
		c.client.GetConfig().BaseURL,
		url.PathEscape(c.dataDockID),
	)

	resp, err := c.client.Do(ctx, "GET", endpoint, nil)
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
