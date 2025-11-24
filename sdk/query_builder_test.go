package sdk

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/nudibranches-tech/bifrost-hyperfluid-sdk-dev/sdk/utils"
)

func TestQueryBuilder_BasicChaining(t *testing.T) {
	client := newTestClient(utils.Configuration{
		Token: "test-token",
		OrgID: "default-org",
	}, func(req *http.Request) (*http.Response, error) {
		// Verify the URL was constructed correctly
		expectedPath := "/default-org/openapi/test-catalog/test-schema/test-table"
		if !strings.Contains(req.URL.Path, expectedPath) {
			t.Errorf("Expected path to contain '%s', got '%s'", expectedPath, req.URL.Path)
		}

		// Verify query parameters
		query := req.URL.Query()
		if query.Get("_limit") != "10" {
			t.Errorf("Expected _limit=10, got %s", query.Get("_limit"))
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"data": "success"}`)),
		}, nil
	})

	resp, err := client.
		Catalog("test-catalog").
		Schema("test-schema").
		Table("test-table").
		Limit(10).
		Get(context.Background())

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if resp.Status != utils.StatusOK {
		t.Errorf("Expected status OK, got %s", resp.Status)
	}
}

func TestQueryBuilder_WithSelect(t *testing.T) {
	client := newTestClient(utils.Configuration{
		Token: "test-token",
		OrgID: "test-org",
	}, func(req *http.Request) (*http.Response, error) {
		query := req.URL.Query()
		selectParam := query.Get("select")
		if selectParam != "id,name,email" {
			t.Errorf("Expected select=id,name,email, got %s", selectParam)
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`[]`)),
		}, nil
	})

	_, err := client.
		Catalog("cat").
		Schema("schema").
		Table("users").
		Select("id", "name", "email").
		Get(context.Background())

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestQueryBuilder_WithMultipleSelects(t *testing.T) {
	client := newTestClient(utils.Configuration{
		Token: "test-token",
		OrgID: "test-org",
	}, func(req *http.Request) (*http.Response, error) {
		query := req.URL.Query()
		selectParam := query.Get("select")
		if selectParam != "id,name,email,phone" {
			t.Errorf("Expected select=id,name,email,phone, got %s", selectParam)
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`[]`)),
		}, nil
	})

	_, err := client.
		Catalog("cat").
		Schema("schema").
		Table("users").
		Select("id", "name").
		Select("email", "phone").
		Get(context.Background())

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestQueryBuilder_WithFilters(t *testing.T) {
	client := newTestClient(utils.Configuration{
		Token: "test-token",
		OrgID: "test-org",
	}, func(req *http.Request) (*http.Response, error) {
		// Check for filter parameters
		if !strings.Contains(req.URL.RawQuery, "age") {
			t.Error("Expected age filter in query")
		}
		if !strings.Contains(req.URL.RawQuery, "status") {
			t.Error("Expected status filter in query")
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`[]`)),
		}, nil
	})

	_, err := client.
		Catalog("cat").
		Schema("schema").
		Table("users").
		Where("age", ">", 18).
		Where("status", "=", "active").
		Get(context.Background())

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestQueryBuilder_WithOrderBy(t *testing.T) {
	client := newTestClient(utils.Configuration{
		Token: "test-token",
		OrgID: "test-org",
	}, func(req *http.Request) (*http.Response, error) {
		query := req.URL.Query()
		orderParam := query.Get("order")
		if orderParam != "created_at.desc,name.asc" {
			t.Errorf("Expected order=created_at.desc,name.asc, got %s", orderParam)
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`[]`)),
		}, nil
	})

	_, err := client.
		Catalog("cat").
		Schema("schema").
		Table("users").
		OrderBy("created_at", "DESC").
		OrderBy("name", "ASC").
		Get(context.Background())

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestQueryBuilder_WithPagination(t *testing.T) {
	client := newTestClient(utils.Configuration{
		Token: "test-token",
		OrgID: "test-org",
	}, func(req *http.Request) (*http.Response, error) {
		query := req.URL.Query()
		if query.Get("_limit") != "25" {
			t.Errorf("Expected _limit=25, got %s", query.Get("_limit"))
		}
		if query.Get("_offset") != "50" {
			t.Errorf("Expected _offset=50, got %s", query.Get("_offset"))
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`[]`)),
		}, nil
	})

	_, err := client.
		Catalog("cat").
		Schema("schema").
		Table("users").
		Limit(25).
		Offset(50).
		Get(context.Background())

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

// Note: Custom org test removed - now use Progressive API with client.Org(customID)
// The Org() method now returns OrgBuilder, not QueryBuilder

func TestQueryBuilder_ValidationErrors(t *testing.T) {
	client := newTestClient(utils.Configuration{
		Token: "test-token",
		OrgID: "test-org",
	}, nil)

	tests := []struct {
		name        string
		buildQuery  func() *QueryBuilder
		expectError bool
		errorMsg    string
	}{
		{
			name: "missing catalog",
			buildQuery: func() *QueryBuilder {
				return client.Query().Schema("schema").Table("table")
			},
			expectError: true,
			errorMsg:    "catalog name is required",
		},
		{
			name: "missing schema",
			buildQuery: func() *QueryBuilder {
				return client.Query().Catalog("cat").Table("table")
			},
			expectError: true,
			errorMsg:    "schema name is required",
		},
		{
			name: "missing table",
			buildQuery: func() *QueryBuilder {
				return client.Query().Catalog("cat").Schema("schema")
			},
			expectError: true,
			errorMsg:    "table name is required",
		},
		{
			name: "empty catalog name",
			buildQuery: func() *QueryBuilder {
				return client.Query().Catalog("").Schema("schema").Table("table")
			},
			expectError: true,
			errorMsg:    "catalog name cannot be empty",
		},
		{
			name: "negative limit",
			buildQuery: func() *QueryBuilder {
				return client.Query().Catalog("cat").Schema("schema").Table("table").Limit(-1)
			},
			expectError: true,
			errorMsg:    "limit cannot be negative",
		},
		{
			name: "negative offset",
			buildQuery: func() *QueryBuilder {
				return client.Query().Catalog("cat").Schema("schema").Table("table").Offset(-10)
			},
			expectError: true,
			errorMsg:    "offset cannot be negative",
		},
		{
			name: "invalid operator",
			buildQuery: func() *QueryBuilder {
				return client.Query().Catalog("cat").Schema("schema").Table("table").Where("col", "??", "val")
			},
			expectError: true,
			errorMsg:    "invalid operator",
		},
		{
			name: "invalid order direction",
			buildQuery: func() *QueryBuilder {
				return client.Query().Catalog("cat").Schema("schema").Table("table").OrderBy("col", "INVALID")
			},
			expectError: true,
			errorMsg:    "must be ASC or DESC",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qb := tt.buildQuery()
			_, err := qb.Get(context.Background())

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got nil")
				} else if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error to contain '%s', got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
			}
		})
	}
}

func TestQueryBuilder_RawParams(t *testing.T) {
	client := newTestClient(utils.Configuration{
		Token: "test-token",
		OrgID: "test-org",
	}, func(req *http.Request) (*http.Response, error) {
		query := req.URL.Query()
		if query.Get("custom_param") != "custom_value" {
			t.Errorf("Expected custom_param=custom_value, got %s", query.Get("custom_param"))
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`[]`)),
		}, nil
	})

	rawParams := make(map[string][]string)
	rawParams["custom_param"] = []string{"custom_value"}

	_, err := client.
		Catalog("cat").
		Schema("schema").
		Table("table").
		RawParams(rawParams).
		Get(context.Background())

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestQueryBuilder_ComplexQuery(t *testing.T) {
	client := newTestClient(utils.Configuration{
		Token: "test-token",
		OrgID: "test-org",
	}, func(req *http.Request) (*http.Response, error) {
		query := req.URL.Query()

		// Verify all parameters are present
		if query.Get("select") != "id,name,email" {
			t.Errorf("Unexpected select parameter: %s", query.Get("select"))
		}
		if query.Get("_limit") != "100" {
			t.Errorf("Unexpected limit: %s", query.Get("_limit"))
		}
		if query.Get("_offset") != "200" {
			t.Errorf("Unexpected offset: %s", query.Get("_offset"))
		}
		if query.Get("order") != "created_at.desc" {
			t.Errorf("Unexpected order: %s", query.Get("order"))
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`[{"id":1,"name":"Test","email":"test@example.com"}]`)),
		}, nil
	})

	resp, err := client.
		Catalog("sales").
		Schema("public").
		Table("customers").
		Select("id", "name", "email").
		Where("status", "=", "active").
		Where("age", ">", 18).
		OrderBy("created_at", "DESC").
		Limit(100).
		Offset(200).
		Get(context.Background())

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if resp.Status != utils.StatusOK {
		t.Errorf("Expected status OK, got %s", resp.Status)
	}

	// Verify response data
	data, ok := resp.Data.([]interface{})
	if !ok {
		t.Fatalf("Expected data to be array, got %T", resp.Data)
	}
	if len(data) != 1 {
		t.Errorf("Expected 1 row, got %d", len(data))
	}
}

func TestQueryBuilder_URLEscaping(t *testing.T) {
	client := newTestClient(utils.Configuration{
		Token: "test-token",
		OrgID: "test-org",
	}, func(req *http.Request) (*http.Response, error) {
		// Verify special characters are properly escaped
		path := req.URL.Path
		if strings.Contains(path, "../") {
			t.Error("Path should not contain unescaped ../")
		}
		// Path should be properly encoded
		if !strings.Contains(path, "test%2Fcatalog") && !strings.Contains(path, "test/catalog") {
			t.Errorf("Expected escaped path, got %s", path)
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`[]`)),
		}, nil
	})

	_, err := client.
		Catalog("test/catalog").
		Schema("test schema").
		Table("test-table").
		Get(context.Background())

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestQueryBuilder_OrderByDefaultDirection(t *testing.T) {
	client := newTestClient(utils.Configuration{
		Token: "test-token",
		OrgID: "test-org",
	}, func(req *http.Request) (*http.Response, error) {
		query := req.URL.Query()
		orderParam := query.Get("order")
		// Empty direction should default to ASC
		if orderParam != "name.asc" {
			t.Errorf("Expected order=name.asc (default), got %s", orderParam)
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`[]`)),
		}, nil
	})

	_, err := client.
		Catalog("cat").
		Schema("schema").
		Table("users").
		OrderBy("name", ""). // Empty direction should default to ASC
		Get(context.Background())

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}
