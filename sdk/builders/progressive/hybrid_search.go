package progressive

import (
	"context"
	"fmt"

	"github.com/nudibranches-tech/hyperfluid-sdk-go/sdk/builders"
	"github.com/nudibranches-tech/hyperfluid-sdk-go/sdk/builders/fluent"
	"github.com/nudibranches-tech/hyperfluid-sdk-go/sdk/utils"
)

// HybridSearchBuilder provides a progressive hybrid search interface starting from a DataDock.
type HybridSearchBuilder struct {
	client builders.ClientInterface

	dataDockID     string
	ftsQuery       string
	vectorQuery    string
	catalogName    string
	schemaName     string
	tableName      string
	columnsToIndex []string
	columnWeights  map[string]float32
	fusion         *fluent.FusionConfig
	ftsLimit       int
	vectorLimit    int
	limitVal       int
}

// Catalog sets the catalog name for the search.
func (b *HybridSearchBuilder) Catalog(name string) *HybridSearchBuilder {
	b.catalogName = name
	return b
}

// Schema sets the schema name for the search.
func (b *HybridSearchBuilder) Schema(name string) *HybridSearchBuilder {
	b.schemaName = name
	return b
}

// Table sets the table name for the search.
func (b *HybridSearchBuilder) Table(name string) *HybridSearchBuilder {
	b.tableName = name
	return b
}

// Columns sets the columns to index for the search.
func (b *HybridSearchBuilder) Columns(columns ...string) *HybridSearchBuilder {
	b.columnsToIndex = append(b.columnsToIndex, columns...)
	return b
}

// ColumnWeights sets per-column FTS weights for biased keyword matching.
func (b *HybridSearchBuilder) ColumnWeights(weights map[string]float32) *HybridSearchBuilder {
	b.columnWeights = weights
	return b
}

// Fusion sets the fusion configuration (strategy, rrf_k, alpha).
func (b *HybridSearchBuilder) Fusion(config fluent.FusionConfig) *HybridSearchBuilder {
	b.fusion = &config
	return b
}

// FTSLimit sets the number of FTS candidates before fusion.
func (b *HybridSearchBuilder) FTSLimit(n int) *HybridSearchBuilder {
	b.ftsLimit = n
	return b
}

// VectorLimit sets the number of vector candidates before fusion.
func (b *HybridSearchBuilder) VectorLimit(n int) *HybridSearchBuilder {
	b.vectorLimit = n
	return b
}

// Limit sets the maximum number of final results to return.
func (b *HybridSearchBuilder) Limit(n int) *HybridSearchBuilder {
	b.limitVal = n
	return b
}

// Execute executes the hybrid search query and returns the results.
func (b *HybridSearchBuilder) Execute(ctx context.Context) (*fluent.HybridSearchResults, error) {
	if b.ftsQuery == "" {
		return nil, fmt.Errorf("%w: FTS query is required", utils.ErrInvalidRequest)
	}
	if b.vectorQuery == "" {
		return nil, fmt.Errorf("%w: vector query is required", utils.ErrInvalidRequest)
	}
	if b.dataDockID == "" {
		return nil, fmt.Errorf("%w: data dock ID is required", utils.ErrInvalidRequest)
	}
	if b.catalogName == "" {
		return nil, fmt.Errorf("%w: catalog name is required", utils.ErrInvalidRequest)
	}
	if b.schemaName == "" {
		return nil, fmt.Errorf("%w: schema name is required", utils.ErrInvalidRequest)
	}
	if b.tableName == "" {
		return nil, fmt.Errorf("%w: table name is required", utils.ErrInvalidRequest)
	}
	if len(b.columnsToIndex) == 0 {
		return nil, fmt.Errorf("%w: at least one column must be specified", utils.ErrInvalidRequest)
	}

	requestBody := map[string]interface{}{
		"fts_query":        b.ftsQuery,
		"vector_query":     b.vectorQuery,
		"data_dock_id":     b.dataDockID,
		"catalog":          b.catalogName,
		"schema":           b.schemaName,
		"table":            b.tableName,
		"columns_to_index": b.columnsToIndex,
		"limit":            b.limitVal,
	}

	if b.columnWeights != nil {
		requestBody["column_weights"] = b.columnWeights
	}
	if b.fusion != nil {
		requestBody["fusion"] = b.fusion
	}
	if b.ftsLimit > 0 {
		requestBody["fts_limit"] = b.ftsLimit
	}
	if b.vectorLimit > 0 {
		requestBody["vector_limit"] = b.vectorLimit
	}

	endpoint := fmt.Sprintf("%s/api/hybrid-search", b.client.GetConfig().BaseURL)
	body := utils.JsonMarshal(requestBody)

	resp, err := b.client.Do(ctx, "POST", endpoint, body)
	if err != nil {
		return nil, err
	}

	if resp.Status != utils.StatusOK {
		return nil, fmt.Errorf("%w: %s", utils.ErrAPIError, resp.Error)
	}

	results := &fluent.HybridSearchResults{}
	if err := utils.UnmarshalData(resp.Data, results); err != nil {
		return nil, fmt.Errorf("failed to unmarshal hybrid search results: %w", err)
	}

	return results, nil
}
