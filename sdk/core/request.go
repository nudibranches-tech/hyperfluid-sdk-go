package core

import (
	"bifrost-for-developers/sdk/core/utils"
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"time"

	_ "github.com/lib/pq"
)

func HandleRequest(bifrostRequest utils.BifrostRequest) (*utils.Response, error) {

	if globalConfiguration == nil {
		return nil, utils.ErrorUninitialized
	}

	if globalConfiguration.OrgID == "" {
		return utils.ResponseError("missing organization ID (HYPERFLUID_ORG_ID)")
	}

	switch bifrostRequest.Type {
	case utils.RequestGraphQL:
		if bifrostRequest.GraphQLPayload == nil || bifrostRequest.GraphQLPayload.Query == "" {
			return nil, utils.ErrorInvalidRequest
		}
		return executeGraphQLRequest(bifrostRequest.GraphQLPayload)

	case utils.RequestOpenAPI:
		if bifrostRequest.OpenAPIPayload == nil || bifrostRequest.OpenAPIPayload.Catalog == "" {
			return nil, utils.ErrorInvalidRequest
		}
		return executeOpenAPIRequest(bifrostRequest.OpenAPIPayload)

	case utils.RequestPostgres:
		if bifrostRequest.PostgresPayload == nil || bifrostRequest.PostgresPayload.SQL == "" {
			return nil, utils.ErrorInvalidRequest
		}
		return executePostgresRequest(bifrostRequest.PostgresPayload)

	default:
		return nil, utils.ErrorUnsupportedRequestType
	}
}

func executeGraphQLRequest(graphQLPayload *utils.GraphQLPayload) (*utils.Response, error) {
	requestBody := map[string]any{"query": graphQLPayload.Query}
	if graphQLPayload.Variables != nil {
		requestBody["variables"] = graphQLPayload.Variables
	}

	endpointURL := fmt.Sprintf("%s/%s/graphql", globalConfiguration.BaseURL, globalConfiguration.OrgID)
	fmt.Printf("GraphQL Request: %s", endpointURL)
	return executeHTTPRequestWithRetryAndTokenRefresh(http.MethodPost, endpointURL, utils.JsonMarshal(requestBody))
}

func executeOpenAPIRequest(openAPIPayload *utils.OpenAPIPayload) (*utils.Response, error) {
	endpointURL := fmt.Sprintf(
		"%s/%s/openapi/%s/%s/%s",
		globalConfiguration.BaseURL,
		globalConfiguration.OrgID,
		openAPIPayload.Catalog,
		openAPIPayload.Schema,
		openAPIPayload.Table,
	)

	if len(openAPIPayload.Params) > 0 {
		queryValues := url.Values{}
		for key, value := range openAPIPayload.Params {
			queryValues.Set(key, value)

		}
		fmt.Printf("OpenAPI Request: %s", queryValues.Encode())
		endpointURL += "?" + queryValues.Encode()
	}
	httpMethod := openAPIPayload.Method
	if httpMethod == "" {
		httpMethod = http.MethodGet
	}

	return executeHTTPRequestWithRetryAndTokenRefresh(httpMethod, endpointURL, nil)
}

func executeHTTPRequestWithRetryAndTokenRefresh(httpMethod string, endpointURL string, requestBodyBytes []byte) (*utils.Response, error) {
	var lastError error
	var lastResponse *utils.Response

	for attemptNumber := 0; attemptNumber <= globalConfiguration.MaxRetries; attemptNumber++ {
		if attemptNumber > 0 {
			delayDuration := time.Duration(math.Pow(2, float64(attemptNumber-1))*100) * time.Millisecond
			time.Sleep(delayDuration)
		}

		executionContext, cancelExecution := context.WithTimeout(context.Background(), globalConfiguration.RequestTimeout)

		httpRequest, creationError := http.NewRequestWithContext(
			executionContext,
			httpMethod,
			endpointURL,
			bytes.NewBuffer(requestBodyBytes),
		)
		if creationError != nil {
			cancelExecution()
			return utils.ResponseError("Cannot create HTTP request: " + creationError.Error())
		}

		if globalConfiguration.Token == "" {
			cancelExecution()
			return utils.ResponseError("Missing token (HYPERFLUID_TOKEN or Keycloak required)")
		}

		httpRequest.Header.Set("Authorization", "Bearer "+globalConfiguration.Token)
		if requestBodyBytes != nil {
			httpRequest.Header.Set("Content-Type", "application/json")
		}

		httpClient := utils.CreateHTTPClientWithSettings(globalConfiguration.SkipTLSVerify, globalConfiguration.RequestTimeout)
		httpResponse, httpError := httpClient.Do(httpRequest)
		cancelExecution()

		if httpError != nil {
			lastError = httpError
			continue
		}

		responseBodyBytes, _ := io.ReadAll(httpResponse.Body)
		httpResponse.Body.Close()

		if httpResponse.StatusCode >= 300 {
			lastResponse = &utils.Response{
				Status:   utils.StatusError,
				Error:    string(responseBodyBytes),
				HTTPCode: httpResponse.StatusCode,
			}

			if httpResponse.StatusCode == http.StatusUnauthorized &&
				globalConfiguration.KeycloakUsername != "" &&
				globalConfiguration.KeycloakPassword != "" {
				if _, refreshError := FetchKeycloakToken(context.Background()); refreshError == nil {
					continue
				}
			}

			if httpResponse.StatusCode >= 400 && httpResponse.StatusCode < 500 {
				return lastResponse, nil
			}

			continue
		}

		var parsedBody any
		if parseError := json.Unmarshal(responseBodyBytes, &parsedBody); parseError != nil {
			lastError = parseError
			continue
		}

		return &utils.Response{
			Status:   utils.StatusOK,
			Data:     parsedBody,
			HTTPCode: httpResponse.StatusCode,
		}, nil
	}

	if lastResponse != nil && lastResponse.Error != "" {
		return lastResponse, fmt.Errorf("max retries exceeded: %s", lastResponse.Error)
	}

	return nil, fmt.Errorf("max retries exceeded: %v", lastError)
}

func executePostgresRequest(postgresPayload *utils.PostgresPayload) (*utils.Response, error) {
	if globalConfiguration.PostgresUser == "" || globalConfiguration.PostgresHost == "" || globalConfiguration.PostgresPort == 0 || globalConfiguration.PostgresDatabase == "" {
		return utils.ResponseError("missing PostgreSQL information (HYPERFLUID_POSTGRES_USER, HYPERFLUID_POSTGRES_HOST, HYPERFLUID_POSTGRES_PORT, HYPERFLUID_POSTGRES_DATABASE)")
	}
	if globalConfiguration.Token == "" {
		return utils.ResponseError("missing authentication token required for PostgreSQL")
	}

	databaseName := globalConfiguration.PostgresDatabase
	if databaseName == "" {
		databaseName = globalConfiguration.OrgID
	}
	postgresConnectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", globalConfiguration.PostgresHost, globalConfiguration.PostgresPort, globalConfiguration.PostgresUser, globalConfiguration.Token, databaseName)

	databaseConnection, connectionError := sql.Open("postgres", postgresConnectionString)
	if connectionError != nil {
		return utils.ResponseError("Cannot open PostgreSQL connection: " + connectionError.Error())
	}
	defer databaseConnection.Close()

	executionContext, cancelExecution := context.WithTimeout(context.Background(), globalConfiguration.RequestTimeout)
	defer cancelExecution()

	if pingError := databaseConnection.PingContext(executionContext); pingError != nil {
		return utils.ResponseError("PostgreSQL ping failed: " + pingError.Error())
	}

	rows, queryError := databaseConnection.QueryContext(executionContext, postgresPayload.SQL)
	if queryError != nil {
		return utils.ResponseError("PostgreSQL query failed: " + queryError.Error())
	}
	defer rows.Close()

	columnNames, _ := rows.Columns()
	var resultRows []map[string]any

	for rows.Next() {
		columnValues := make([]any, len(columnNames))
		valuePointers := make([]any, len(columnValues))
		for index := range columnValues {
			valuePointers[index] = &columnValues[index]
		}

		if scanError := rows.Scan(valuePointers...); scanError != nil {
			return utils.ResponseError("Failed to scan PostgreSQL row: " + scanError.Error())
		}

		rowMap := make(map[string]any)
		for index, columnName := range columnNames {
			if byteValue, isBytes := columnValues[index].([]byte); isBytes {
				rowMap[columnName] = string(byteValue)
			} else {
				rowMap[columnName] = columnValues[index]
			}
		}

		resultRows = append(resultRows, rowMap)
	}

	return utils.ResponseSuccess(resultRows), nil
}
