package fluent

import (
	"context"
	"fmt"

	"github.com/nudibranches-tech/hyperfluid-sdk-go/sdk/utils"
)

// FusionStrategy defines how FTS and vector results are combined.
type FusionStrategy string

const (
	FusionStrategyRRF    FusionStrategy = "rrf"
	FusionStrategyLinear FusionStrategy = "linear"
)

// FusionConfig controls the fusion step in hybrid search.
type FusionConfig struct {
	Strategy  FusionStrategy `json:"strategy"`
	RRFK      float32        `json:"rrf_k"`
	Alpha     float32        `json:"alpha"`
	FusionKey string         `json:"fusion_key,omitempty"`
}

// HybridSearchResult represents a single hybrid search result with a fused score.
type HybridSearchResult struct {
	Record  DocumentRecord `json:"record"`
	Score   float64        `json:"score"`
	Sources []string       `json:"sources"`
}

// HybridSearchResults holds the response from a hybrid search query.
type HybridSearchResults struct {
	Results     []HybridSearchResult `json:"results"`
	Total       int                  `json:"total"`
	TimeTakenMs int                  `json:"took_ms"`
}

// HybridSearchBuilder provides a fluent interface for building and executing hybrid search queries.
type HybridSearchBuilder struct {
	client interface {
		Do(ctx context.Context, method, endpoint string, body []byte) (*utils.Response, error)
		GetConfig() utils.Configuration
	}
	errors []error

	searchQuery    string
	dataDockID     string
	catalogName    string
	schemaName     string
	tableName      string
	columnsToIndex []string
	columnWeights  map[string]float32
	fusion         *FusionConfig
	ftsLimit       int
	vectorLimit    int
	limitVal       int
}

// NewHybridSearchBuilder creates a new HybridSearchBuilder instance.
func NewHybridSearchBuilder(client interface {
	Do(ctx context.Context, method, endpoint string, body []byte) (*utils.Response, error)
	GetConfig() utils.Configuration
}) *HybridSearchBuilder {
	return &HybridSearchBuilder{
		client:         client,
		errors:         []error{},
		dataDockID:     client.GetConfig().DataDockID,
		columnsToIndex: []string{},
		limitVal:       20,
	}
}

// Query sets the search query string.
func (b *HybridSearchBuilder) Query(query string) *HybridSearchBuilder {
	if query == "" {
		b.errors = append(b.errors, fmt.Errorf("search query cannot be empty"))
	}
	b.searchQuery = query
	return b
}

// DataDock sets the data dock ID for the search.
func (b *HybridSearchBuilder) DataDock(dataDockID string) *HybridSearchBuilder {
	if dataDockID == "" {
		b.errors = append(b.errors, fmt.Errorf("data dock ID cannot be empty"))
	}
	b.dataDockID = dataDockID
	return b
}

// Catalog sets the catalog name for the search.
func (b *HybridSearchBuilder) Catalog(name string) *HybridSearchBuilder {
	if name == "" {
		b.errors = append(b.errors, fmt.Errorf("catalog name cannot be empty"))
	}
	b.catalogName = name
	return b
}

// Schema sets the schema name for the search.
func (b *HybridSearchBuilder) Schema(name string) *HybridSearchBuilder {
	if name == "" {
		b.errors = append(b.errors, fmt.Errorf("schema name cannot be empty"))
	}
	b.schemaName = name
	return b
}

// Table sets the table name for the search.
func (b *HybridSearchBuilder) Table(name string) *HybridSearchBuilder {
	if name == "" {
		b.errors = append(b.errors, fmt.Errorf("table name cannot be empty"))
	}
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
func (b *HybridSearchBuilder) Fusion(config FusionConfig) *HybridSearchBuilder {
	b.fusion = &config
	return b
}

// FTSLimit sets the number of FTS candidates before fusion (default: 100, max: 1000).
func (b *HybridSearchBuilder) FTSLimit(n int) *HybridSearchBuilder {
	if n <= 0 || n > 1000 {
		b.errors = append(b.errors, fmt.Errorf("fts_limit must be between 1 and 1000"))
		return b
	}
	b.ftsLimit = n
	return b
}

// VectorLimit sets the number of vector candidates before fusion (default: 100, max: 1000).
func (b *HybridSearchBuilder) VectorLimit(n int) *HybridSearchBuilder {
	if n <= 0 || n > 1000 {
		b.errors = append(b.errors, fmt.Errorf("vector_limit must be between 1 and 1000"))
		return b
	}
	b.vectorLimit = n
	return b
}

// Limit sets the maximum number of final results to return.
func (b *HybridSearchBuilder) Limit(n int) *HybridSearchBuilder {
	if n <= 0 || n > 100 {
		b.errors = append(b.errors, fmt.Errorf("limit must be between 1 and 100"))
		return b
	}
	b.limitVal = n
	return b
}

func (b *HybridSearchBuilder) validate() error {
	if len(b.errors) > 0 {
		return fmt.Errorf("hybrid search builder validation failed: %s", b.errors[0].Error())
	}
	if b.searchQuery == "" {
		return fmt.Errorf("%w: search query is required", utils.ErrInvalidRequest)
	}
	if b.dataDockID == "" {
		return fmt.Errorf("%w: data dock ID is required", utils.ErrInvalidRequest)
	}
	if b.catalogName == "" {
		return fmt.Errorf("%w: catalog name is required", utils.ErrInvalidRequest)
	}
	if b.schemaName == "" {
		return fmt.Errorf("%w: schema name is required", utils.ErrInvalidRequest)
	}
	if b.tableName == "" {
		return fmt.Errorf("%w: table name is required", utils.ErrInvalidRequest)
	}
	if len(b.columnsToIndex) == 0 {
		return fmt.Errorf("%w: at least one column must be specified", utils.ErrInvalidRequest)
	}
	return nil
}

// Execute executes the hybrid search query and returns the results.
func (b *HybridSearchBuilder) Execute(ctx context.Context) (*HybridSearchResults, error) {
	if err := b.validate(); err != nil {
		return nil, err
	}

	requestBody := map[string]interface{}{
		"query":            b.searchQuery,
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

	results := &HybridSearchResults{}
	if err := utils.UnmarshalData(resp.Data, results); err != nil {
		return nil, fmt.Errorf("failed to unmarshal hybrid search results: %w", err)
	}

	return results, nil
}
