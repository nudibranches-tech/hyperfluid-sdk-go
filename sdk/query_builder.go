package sdk

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/nudibranches-tech/bifrost-hyperfluid-sdk-dev/sdk/utils"
)

// QueryBuilder provides a fluent interface for building and executing queries.
type QueryBuilder struct {
	client *Client
	errors []error

	// Hierarchy
	orgID       string
	catalogName string
	schemaName  string
	tableName   string

	// Query parameters
	selectCols []string
	filters    []Filter
	orderBy    []OrderClause
	limitVal   int
	offsetVal  int
	rawParams  url.Values
}

// Filter represents a WHERE clause condition.
type Filter struct {
	Column   string
	Operator string // =, >, <, >=, <=, !=, LIKE, IN
	Value    interface{}
}

// OrderClause represents an ORDER BY clause.
type OrderClause struct {
	Column    string
	Direction string // ASC or DESC
}

// newQueryBuilder creates a new QueryBuilder instance.
func newQueryBuilder(client *Client) *QueryBuilder {
	return &QueryBuilder{
		client:    client,
		errors:    []error{},
		orgID:     client.config.OrgID, // Use default from config
		rawParams: url.Values{},
	}
}

// Org sets the organization ID for the query.
// If not called, uses the OrgID from client configuration.
func (qb *QueryBuilder) Org(orgID string) *QueryBuilder {
	if orgID == "" {
		qb.errors = append(qb.errors, fmt.Errorf("organization ID cannot be empty"))
	}
	qb.orgID = orgID
	return qb
}

// Catalog sets the catalog name for the query.
func (qb *QueryBuilder) Catalog(name string) *QueryBuilder {
	if name == "" {
		qb.errors = append(qb.errors, fmt.Errorf("catalog name cannot be empty"))
	}
	qb.catalogName = name
	return qb
}

// Schema sets the schema name for the query.
func (qb *QueryBuilder) Schema(name string) *QueryBuilder {
	if name == "" {
		qb.errors = append(qb.errors, fmt.Errorf("schema name cannot be empty"))
	}
	qb.schemaName = name
	return qb
}

// Table sets the table name for the query.
func (qb *QueryBuilder) Table(name string) *QueryBuilder {
	if name == "" {
		qb.errors = append(qb.errors, fmt.Errorf("table name cannot be empty"))
	}
	qb.tableName = name
	return qb
}

// Select specifies which columns to retrieve.
// Can be called multiple times to add more columns.
func (qb *QueryBuilder) Select(columns ...string) *QueryBuilder {
	qb.selectCols = append(qb.selectCols, columns...)
	return qb
}

// Where adds a filter condition to the query.
// Supported operators: =, >, <, >=, <=, !=, LIKE, IN
func (qb *QueryBuilder) Where(column, operator string, value interface{}) *QueryBuilder {
	validOperators := map[string]bool{
		"=": true, ">": true, "<": true, ">=": true, "<=": true,
		"!=": true, "LIKE": true, "IN": true,
	}

	if !validOperators[operator] {
		qb.errors = append(qb.errors, fmt.Errorf("invalid operator '%s'", operator))
	}

	qb.filters = append(qb.filters, Filter{
		Column:   column,
		Operator: operator,
		Value:    value,
	})
	return qb
}

// OrderBy adds an ORDER BY clause to the query.
// Direction should be "ASC" or "DESC" (defaults to "ASC" if empty).
func (qb *QueryBuilder) OrderBy(column, direction string) *QueryBuilder {
	if direction == "" {
		direction = "ASC"
	}

	direction = strings.ToUpper(direction)
	if direction != "ASC" && direction != "DESC" {
		qb.errors = append(qb.errors, fmt.Errorf("invalid order direction '%s', must be ASC or DESC", direction))
		return qb
	}

	qb.orderBy = append(qb.orderBy, OrderClause{
		Column:    column,
		Direction: direction,
	})
	return qb
}

// Limit sets the maximum number of rows to return.
func (qb *QueryBuilder) Limit(n int) *QueryBuilder {
	if n < 0 {
		qb.errors = append(qb.errors, fmt.Errorf("limit cannot be negative"))
		return qb
	}
	qb.limitVal = n
	return qb
}

// Offset sets the number of rows to skip.
func (qb *QueryBuilder) Offset(n int) *QueryBuilder {
	if n < 0 {
		qb.errors = append(qb.errors, fmt.Errorf("offset cannot be negative"))
		return qb
	}
	qb.offsetVal = n
	return qb
}

// RawParams allows adding custom query parameters.
// This is an escape hatch for advanced use cases.
func (qb *QueryBuilder) RawParams(params url.Values) *QueryBuilder {
	for key, values := range params {
		for _, value := range values {
			qb.rawParams.Add(key, value)
		}
	}
	return qb
}

// validate checks that all required fields are set.
func (qb *QueryBuilder) validate() error {
	// Check for accumulated errors during building
	if len(qb.errors) > 0 {
		var errMsgs []string
		for _, err := range qb.errors {
			errMsgs = append(errMsgs, err.Error())
		}
		return fmt.Errorf("query builder validation failed: %s", strings.Join(errMsgs, "; "))
	}

	// Check required fields
	if qb.orgID == "" {
		return fmt.Errorf("%w: organization ID is required", utils.ErrInvalidRequest)
	}
	if qb.catalogName == "" {
		return fmt.Errorf("%w: catalog name is required", utils.ErrInvalidRequest)
	}
	if qb.schemaName == "" {
		return fmt.Errorf("%w: schema name is required", utils.ErrInvalidRequest)
	}
	if qb.tableName == "" {
		return fmt.Errorf("%w: table name is required", utils.ErrInvalidRequest)
	}

	return nil
}

// buildEndpoint constructs the API endpoint URL.
func (qb *QueryBuilder) buildEndpoint() string {
	// Use url.PathEscape for each segment to prevent injection
	return fmt.Sprintf(
		"%s/%s/openapi/%s/%s/%s",
		strings.TrimRight(qb.client.config.BaseURL, "/"),
		url.PathEscape(qb.orgID),
		url.PathEscape(qb.catalogName),
		url.PathEscape(qb.schemaName),
		url.PathEscape(qb.tableName),
	)
}

// buildParams constructs the query parameters.
func (qb *QueryBuilder) buildParams() url.Values {
	params := url.Values{}

	// Copy raw params first (they can be overridden)
	for key, values := range qb.rawParams {
		for _, value := range values {
			params.Add(key, value)
		}
	}

	// Add SELECT columns
	if len(qb.selectCols) > 0 {
		params.Set("select", strings.Join(qb.selectCols, ","))
	}

	// Add WHERE filters
	// Note: This assumes the API supports filter parameters
	// Adjust based on actual API capabilities
	for _, filter := range qb.filters {
		paramName := fmt.Sprintf("%s[%s]", filter.Column, filter.Operator)
		params.Add(paramName, fmt.Sprintf("%v", filter.Value))
	}

	// Add ORDER BY
	if len(qb.orderBy) > 0 {
		var orderParts []string
		for _, order := range qb.orderBy {
			if order.Direction == "DESC" {
				orderParts = append(orderParts, fmt.Sprintf("%s.desc", order.Column))
			} else {
				orderParts = append(orderParts, fmt.Sprintf("%s.asc", order.Column))
			}
		}
		params.Set("order", strings.Join(orderParts, ","))
	}

	// Add LIMIT
	if qb.limitVal > 0 {
		params.Set("_limit", strconv.Itoa(qb.limitVal))
	}

	// Add OFFSET
	if qb.offsetVal > 0 {
		params.Set("_offset", strconv.Itoa(qb.offsetVal))
	}

	return params
}

// Get executes the query and returns the results.
// This is the terminal operation that actually makes the API request.
func (qb *QueryBuilder) Get(ctx context.Context) (*utils.Response, error) {
	// Validate the query
	if err := qb.validate(); err != nil {
		return nil, err
	}

	// Build endpoint and parameters
	endpoint := qb.buildEndpoint()
	params := qb.buildParams()

	// Add parameters to endpoint
	if len(params) > 0 {
		endpoint += "?" + params.Encode()
	}

	// Execute the request
	return qb.client.do(ctx, "GET", endpoint, nil)
}

// Count returns the count of rows matching the query.
// Similar to Get() but requests only the count.
func (qb *QueryBuilder) Count(ctx context.Context) (int, error) {
	// Validate the query
	if err := qb.validate(); err != nil {
		return 0, err
	}

	// Build endpoint and parameters
	endpoint := qb.buildEndpoint()
	params := qb.buildParams()

	// Add count parameter (API-specific)
	params.Set("count", "exact")
	params.Set("_limit", "0")

	endpoint += "?" + params.Encode()

	// Execute the request
	resp, err := qb.client.do(ctx, "GET", endpoint, nil)
	if err != nil {
		return 0, err
	}

	// Extract count from response (adjust based on actual API response format)
	if countVal, ok := resp.Data.(map[string]interface{})["count"]; ok {
		if count, ok := countVal.(float64); ok {
			return int(count), nil
		}
	}

	return 0, fmt.Errorf("unable to extract count from response")
}

// Post executes a POST request to insert data.
func (qb *QueryBuilder) Post(ctx context.Context, data interface{}) (*utils.Response, error) {
	if err := qb.validate(); err != nil {
		return nil, err
	}

	endpoint := qb.buildEndpoint()
	body := utils.JsonMarshal(data)

	return qb.client.do(ctx, "POST", endpoint, body)
}

// Put executes a PUT request to update data.
func (qb *QueryBuilder) Put(ctx context.Context, data interface{}) (*utils.Response, error) {
	if err := qb.validate(); err != nil {
		return nil, err
	}

	endpoint := qb.buildEndpoint()
	params := qb.buildParams()

	if len(params) > 0 {
		endpoint += "?" + params.Encode()
	}

	body := utils.JsonMarshal(data)
	return qb.client.do(ctx, "PUT", endpoint, body)
}

// Delete executes a DELETE request.
func (qb *QueryBuilder) Delete(ctx context.Context) (*utils.Response, error) {
	if err := qb.validate(); err != nil {
		return nil, err
	}

	endpoint := qb.buildEndpoint()
	params := qb.buildParams()

	if len(params) > 0 {
		endpoint += "?" + params.Encode()
	}

	return qb.client.do(ctx, "DELETE", endpoint, nil)
}
