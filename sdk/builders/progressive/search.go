package progressive

import (
	"context"
	"fmt"

	"github.com/nudibranches-tech/hyperfluid-sdk-go/sdk/builders"
	"github.com/nudibranches-tech/hyperfluid-sdk-go/sdk/utils"
)

// SearchBuilder provides a progressive search interface starting from a DataDock.
type SearchBuilder struct {
	client builders.ClientInterface

	// Pre-set from DataDock
	dataDockID  string
	searchQuery string

	// To be set progressively
	catalogName    string
	schemaName     string
	tableName      string
	columnsToIndex []string
	limitVal       int
}

// Catalog sets the catalog name for the search.
func (sb *SearchBuilder) Catalog(name string) *SearchBuilder {
	sb.catalogName = name
	return sb
}

// Schema sets the schema name for the search.
func (sb *SearchBuilder) Schema(name string) *SearchBuilder {
	sb.schemaName = name
	return sb
}

// Table sets the table name for the search.
func (sb *SearchBuilder) Table(name string) *SearchBuilder {
	sb.tableName = name
	return sb
}

// Columns sets the columns to index for the search.
// Can be called multiple times to add more columns.
func (sb *SearchBuilder) Columns(columns ...string) *SearchBuilder {
	sb.columnsToIndex = append(sb.columnsToIndex, columns...)
	return sb
}

// Limit sets the maximum number of results to return.
func (sb *SearchBuilder) Limit(n int) *SearchBuilder {
	sb.limitVal = n
	return sb
}

// Execute executes the search query and returns the results.
func (sb *SearchBuilder) Execute(ctx context.Context) (*utils.Response, error) {
	// Validate required fields
	if sb.searchQuery == "" {
		return nil, fmt.Errorf("%w: search query is required", utils.ErrInvalidRequest)
	}
	if sb.dataDockID == "" {
		return nil, fmt.Errorf("%w: data dock ID is required", utils.ErrInvalidRequest)
	}
	if sb.catalogName == "" {
		return nil, fmt.Errorf("%w: catalog name is required", utils.ErrInvalidRequest)
	}
	if sb.schemaName == "" {
		return nil, fmt.Errorf("%w: schema name is required", utils.ErrInvalidRequest)
	}
	if sb.tableName == "" {
		return nil, fmt.Errorf("%w: table name is required", utils.ErrInvalidRequest)
	}
	if len(sb.columnsToIndex) == 0 {
		return nil, fmt.Errorf("%w: at least one column must be specified", utils.ErrInvalidRequest)
	}

	// Build the request body
	requestBody := map[string]interface{}{
		"query":            sb.searchQuery,
		"data_dock_id":     sb.dataDockID,
		"catalog":          sb.catalogName,
		"schema":           sb.schemaName,
		"table":            sb.tableName,
		"limit":            sb.limitVal,
		"columns_to_index": sb.columnsToIndex,
	}

	// Build endpoint
	endpoint := fmt.Sprintf("%s/api/search", sb.client.GetConfig().BaseURL)

	// Marshal request body
	body := utils.JsonMarshal(requestBody)

	// Execute the request
	return sb.client.Do(ctx, "POST", endpoint, body)
}
