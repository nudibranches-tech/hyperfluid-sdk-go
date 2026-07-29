package progressive

import (
	"context"
	"fmt"
	"net/url"

	"github.com/nudibranches-tech/hyperfluid-sdk-go/sdk/builders"
	"github.com/nudibranches-tech/hyperfluid-sdk-go/sdk/utils"
)

type dataContainerOverview struct {
	ID   string `json:"id"`
	Kind string `json:"kind"`
}

// IcebergContainerBuilder represents a named Iceberg data container for a
// datadock.
//
// Available methods:
//   - DeleteSchema(ctx, name) - Delete a schema by name
type IcebergContainerBuilder struct {
	client        builders.ClientInterface
	orgID         string
	harborID      string
	dataDockID    string
	containerName string
}

// resolveContainerID finds this datadock's Iceberg container id by name.
func (f *IcebergContainerBuilder) resolveContainerID(ctx context.Context) (string, error) {
	query := url.Values{}
	query.Set("name", f.containerName)
	endpoint := fmt.Sprintf("%s/data-docks/%s/data-containers?%s",
		f.client.GetConfig().BaseURL,
		url.PathEscape(f.dataDockID),
		query.Encode(),
	)
	resp, err := f.client.Do(ctx, "GET", endpoint, nil)
	if err != nil {
		return "", err
	}

	var containers []dataContainerOverview
	if err := utils.UnmarshalData(resp.Data, &containers); err != nil {
		return "", fmt.Errorf("failed to parse data containers: %w", err)
	}

	for _, c := range containers {
		if c.Kind == "Iceberg" {
			return c.ID, nil
		}
	}

	return "", fmt.Errorf("%w: no Iceberg data container named %q found for data dock %s",
		utils.ErrNotFound, f.containerName, f.dataDockID)
}

// DeleteSchema deletes an Iceberg schema by name.
func (f *IcebergContainerBuilder) DeleteSchema(ctx context.Context, schemaName string) (*utils.Response, error) {
	dataContainerID, err := f.resolveContainerID(ctx)
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("%s/data-containers/iceberg/%s/schemas/%s",
		f.client.GetConfig().BaseURL,
		url.PathEscape(dataContainerID),
		url.PathEscape(schemaName),
	)
	return f.client.Do(ctx, "DELETE", endpoint, nil)
}
