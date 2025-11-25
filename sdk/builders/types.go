package builders

import (
	"context"

	"github.com/nudibranches-tech/bifrost-hyperfluid-sdk-dev/sdk/utils"
)

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

type Builder interface {
	validate() error
}

type Executor interface {
	Get(ctx context.Context) (*utils.Response, error)
}
