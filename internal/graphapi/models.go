package graphapi

type GraphQLRequestBody struct {
	OperationName *string     `json:"operationName"`
	Query         *string     `json:"query"`
	Variables     interface{} `json:"variables"`
}

type Errors []struct {
	Message string   `json:"message"`
	Path    []string `json:"path"`
}

type Response interface {
	Error() Errors
	SetTraceID(traceID string)
}
