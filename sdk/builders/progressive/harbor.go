package progressive

import (
	"context"
	"fmt"
	"net/url"

	"github.com/nudibranches-tech/hyperfluid-sdk-go/sdk/builders"
	"github.com/nudibranches-tech/hyperfluid-sdk-go/sdk/utils"
)

// HarborBuilder represents a harbor context.
// Available methods:
//   - DataDock(id) - Navigate to a specific datadock
//   - ListDataDocks(ctx) - List all datadocks in this harbor
//   - CreateDataDock(ctx, config) - Create a new datadock
//   - Delete(ctx) - Delete this harbor
type HarborBuilder struct {
	client   builders.ClientInterface
	orgID    string
	harborID string
}

// DataDock navigates to a specific datadock in this harbor.
func (h *HarborBuilder) DataDock(dataDockID string) *DataDockBuilder {
	return &DataDockBuilder{
		client:     h.client,
		orgID:      h.orgID,
		harborID:   h.harborID,
		dataDockID: dataDockID,
	}
}

// ListDataDocks retrieves all datadocks in this harbor.
func (h *HarborBuilder) ListDataDocks(ctx context.Context) (*utils.Response, error) {
	endpoint := fmt.Sprintf("%s/harbors/%s/data-docks",
		h.client.GetConfig().BaseURL,
		url.PathEscape(h.harborID),
	)
	return h.client.Do(ctx, "GET", endpoint, nil)
}

// CreateDataDock creates a new datadock in this harbor.
func (h *HarborBuilder) CreateDataDock(ctx context.Context, config map[string]interface{}) (*utils.Response, error) {
	// Ensure harbor_id is set
	config["harbor_id"] = h.harborID

	endpoint := fmt.Sprintf("%s/data-docks", h.client.GetConfig().BaseURL)
	body := utils.JsonMarshal(config)
	return h.client.Do(ctx, "POST", endpoint, body)
}

// Delete removes this harbor.
func (h *HarborBuilder) Delete(ctx context.Context) (*utils.Response, error) {
	endpoint := fmt.Sprintf("%s/harbors/%s",
		h.client.GetConfig().BaseURL,
		url.PathEscape(h.harborID),
	)
	return h.client.Do(ctx, "DELETE", endpoint, nil)
}
