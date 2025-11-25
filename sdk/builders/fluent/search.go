package fluent

import (
	"context"
	"fmt"

	"github.com/nudibranches-tech/bifrost-hyperfluid-sdk-dev/sdk/utils"
)

// SearchBuilder provides a fluent interface for building and executing full-text search queries.
type SearchBuilder struct {
	client interface {
		Do(ctx context.Context, method, endpoint string, body []byte) (*utils.Response, error)
		GetConfig() utils.Configuration
	}
	errors []error

	// Search parameters
	searchQuery    string
	dataDockID     string
	catalogName    string
	schemaName     string
	tableName      string
	columnsToIndex []string
	limitVal       int
}

// NewSearchBuilder creates a new SearchBuilder instance.
func NewSearchBuilder(client interface {
	Do(ctx context.Context, method, endpoint string, body []byte) (*utils.Response, error)
	GetConfig() utils.Configuration
}) *SearchBuilder {
	return &SearchBuilder{
		client:         client,
		errors:         []error{},
		dataDockID:     client.GetConfig().DataDockID, // Auto-set from config
		columnsToIndex: []string{},
		limitVal:       20, // Default limit
	}
}

// Query sets the search query string.
func (sb *SearchBuilder) Query(query string) *SearchBuilder {
	if query == "" {
		sb.errors = append(sb.errors, fmt.Errorf("search query cannot be empty"))
	}
	sb.searchQuery = query
	return sb
}

// DataDock sets the data dock ID for the search.
func (sb *SearchBuilder) DataDock(dataDockID string) *SearchBuilder {
	if dataDockID == "" {
		sb.errors = append(sb.errors, fmt.Errorf("data dock ID cannot be empty"))
	}
	sb.dataDockID = dataDockID
	return sb
}

// Catalog sets the catalog name for the search.
func (sb *SearchBuilder) Catalog(name string) *SearchBuilder {
	if name == "" {
		sb.errors = append(sb.errors, fmt.Errorf("catalog name cannot be empty"))
	}
	sb.catalogName = name
	return sb
}

// Schema sets the schema name for the search.
func (sb *SearchBuilder) Schema(name string) *SearchBuilder {
	if name == "" {
		sb.errors = append(sb.errors, fmt.Errorf("schema name cannot be empty"))
	}
	sb.schemaName = name
	return sb
}

// Table sets the table name for the search.
func (sb *SearchBuilder) Table(name string) *SearchBuilder {
	if name == "" {
		sb.errors = append(sb.errors, fmt.Errorf("table name cannot be empty"))
	}
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
	if n <= 0 {
		sb.errors = append(sb.errors, fmt.Errorf("limit must be greater than 0"))
		return sb
	}
	sb.limitVal = n
	return sb
}

// validate checks that all required fields are set.
func (sb *SearchBuilder) validate() error {
	// Check for accumulated errors during building
	if len(sb.errors) > 0 {
		var errMsgs []string
		for _, err := range sb.errors {
			errMsgs = append(errMsgs, err.Error())
		}
		return fmt.Errorf("search builder validation failed: %s", errMsgs[0])
	}

	// Check required fields
	if sb.searchQuery == "" {
		return fmt.Errorf("%w: search query is required", utils.ErrInvalidRequest)
	}
	if sb.dataDockID == "" {
		return fmt.Errorf("%w: data dock ID is required", utils.ErrInvalidRequest)
	}
	if sb.catalogName == "" {
		return fmt.Errorf("%w: catalog name is required", utils.ErrInvalidRequest)
	}
	if sb.schemaName == "" {
		return fmt.Errorf("%w: schema name is required", utils.ErrInvalidRequest)
	}
	if sb.tableName == "" {
		return fmt.Errorf("%w: table name is required", utils.ErrInvalidRequest)
	}
	if len(sb.columnsToIndex) == 0 {
		return fmt.Errorf("%w: at least one column must be specified", utils.ErrInvalidRequest)
	}

	return nil
}

// Execute executes the search query and returns the results.
func (sb *SearchBuilder) Execute(ctx context.Context) (*utils.Response, error) {
	// Validate the search
	if err := sb.validate(); err != nil {
		return nil, err
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
