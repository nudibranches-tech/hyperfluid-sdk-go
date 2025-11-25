package progressive

import (
	"context"
	"fmt"
	"net/url"

	"github.com/nudibranches-tech/bifrost-hyperfluid-sdk-dev/sdk/builders"
	"github.com/nudibranches-tech/bifrost-hyperfluid-sdk-dev/sdk/utils"
)

// TableQueryBuilder combines table navigation with query building.
// This is the final level where you can build queries AND execute them.
// Inherits all query building methods from the original QueryBuilder.
type TableQueryBuilder struct {
	client builders.ClientInterface
	orgID  string

	// Table location
	catalogName string
	schemaName  string
	tableName   string

	// Query parameters (same as QueryBuilder)
	selectCols []string
	filters    []builders.Filter
	orderBy    []builders.OrderClause
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
	t.filters = append(t.filters, builders.Filter{
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
	t.orderBy = append(t.orderBy, builders.OrderClause{
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
		t.client.GetConfig().BaseURL,
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

	return t.client.Do(ctx, "GET", endpoint, nil)
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
