package builders

import (
	"context"

	"github.com/nudibranches-tech/bifrost-hyperfluid-sdk-dev/sdk/utils"
)

// ClientInterface defines the methods that builders need from the SDK client.
// This avoids circular imports between sdk and builders packages.
type ClientInterface interface {
	Do(ctx context.Context, method, endpoint string, body []byte) (*utils.Response, error)
	GetConfig() utils.Configuration
}
