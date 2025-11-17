package bifrost

import (
	"bifrost-for-developers/sdk/core"
	"bifrost-for-developers/sdk/core/utils"
)

type Configuration = utils.Configuration

type BifrostRequest = utils.BifrostRequest
type Result = utils.Result
type Response = utils.Response
type RequestType = utils.RequestType
type GraphQLPayload = utils.GraphQLPayload
type OpenAPIPayload = utils.OpenAPIPayload
type PostgresPayload = utils.PostgresPayload

const (
	RequestGraphQL  = utils.RequestGraphQL
	RequestOpenAPI  = utils.RequestOpenAPI
	RequestPostgres = utils.RequestPostgres
)

func Init() {
	core.InitFromEnv()
}

func GetGlobalConfiguration() *Configuration {
	return core.GetConfiguration()
}

func Request(bifrostRequest BifrostRequest) <-chan Result {
	resultChan := make(chan Result, 1)
	go func() {
		response, err := core.HandleRequest(bifrostRequest)
		resultChan <- Result{Response: response, Error: err}
		close(resultChan)
	}()
	return resultChan
}
