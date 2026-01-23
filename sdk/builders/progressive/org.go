package progressive

import (
	"context"
	"fmt"
	"net/url"

	"github.com/nudibranches-tech/hyperfluid-sdk-go/sdk/builders"
	"github.com/nudibranches-tech/hyperfluid-sdk-go/sdk/utils"
)

// OrgBuilder represents an organization context.
// Available methods:
//   - Harbor(id) - Navigate to a specific harbor
//   - ListHarbors(ctx) - List all harbors in this org
//   - CreateHarbor(ctx, name) - Create a new harbor
//   - ListDataDocks(ctx) - List all datadocks across all harbors
type OrgBuilder struct {
	Client builders.ClientInterface
	OrgID  string
}

// Harbor navigates to a specific harbor in this organization.
func (o *OrgBuilder) Harbor(harborID string) *HarborBuilder {
	return &HarborBuilder{
		client:   o.Client,
		orgID:    o.OrgID,
		harborID: harborID,
	}
}

// ListHarbors retrieves all harbors in this organization.
func (o *OrgBuilder) ListHarbors(ctx context.Context) (*utils.Response, error) {
	endpoint := fmt.Sprintf("%s/%s/harbors",
		o.Client.GetConfig().BaseURL,
		url.PathEscape(o.OrgID),
	)
	return o.Client.Do(ctx, "GET", endpoint, nil)
}

// CreateHarbor creates a new harbor in this organization.
func (o *OrgBuilder) CreateHarbor(ctx context.Context, name string) (*utils.Response, error) {
	endpoint := fmt.Sprintf("%s/%s/harbors",
		o.Client.GetConfig().BaseURL,
		url.PathEscape(o.OrgID),
	)
	body := utils.JsonMarshal(map[string]interface{}{
		"name": name,
	})
	return o.Client.Do(ctx, "POST", endpoint, body)
}

// ListDataDocks retrieves all datadocks across all harbors in this organization.
func (o *OrgBuilder) ListDataDocks(ctx context.Context) (*utils.Response, error) {
	endpoint := fmt.Sprintf("%s/%s/data-docks",
		o.Client.GetConfig().BaseURL,
		url.PathEscape(o.OrgID),
	)
	return o.Client.Do(ctx, "GET", endpoint, nil)
}

// RefreshAllDataDocks triggers a catalog refresh on all datadocks in this organization.
func (o *OrgBuilder) RefreshAllDataDocks(ctx context.Context) (*utils.Response, error) {
	endpoint := fmt.Sprintf("%s/%s/data-docks/refresh",
		o.Client.GetConfig().BaseURL,
		url.PathEscape(o.OrgID),
	)
	return o.Client.Do(ctx, "POST", endpoint, nil)
}
