package sdk

import (
	"bifrost-for-developers/sdk/utils"
	"context"
	"fmt"
	"net/url"
)

// Catalog represents a Hyperfluid data catalog.
type Catalog struct {
	Name   string
	client *Client
}

// Table retrieves a table from the catalog.
func (c *Catalog) Table(schemaName string, tableName string) *Table {
	return &Table{
		Name:        tableName,
		SchemaName:  schemaName,
		CatalogName: c.Name,
		client:      c.client,
	}
}

// Table represents a table in a Hyperfluid schema.
type Table struct {
	Name        string
	SchemaName  string
	CatalogName string
	client      *Client
}

// GetData retrieves data from the table.
func (t *Table) GetData(ctx context.Context, params url.Values) (*utils.Response, error) {
	endpoint := fmt.Sprintf(
		"%s/%s/openapi/%s/%s/%s",
		t.client.config.BaseURL,
		t.client.config.OrgID,
		t.CatalogName,
		t.SchemaName,
		t.Name,
	)

	if len(params) > 0 {
		endpoint += "?" + params.Encode()
	}

	return t.client.do(ctx, "GET", endpoint, nil)
}
