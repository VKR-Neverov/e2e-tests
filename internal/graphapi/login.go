package graphapi

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/samber/lo"
	"net/http"
)

type LoginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Errors  Errors                            `json:"errors"`
	Data    map[string]map[string]interface{} "data"
	Token   string
	TraceID string
}

func (r *LoginResponse) Error() Errors {
	return r.Errors
}

func (r *LoginResponse) SetTraceID(traceID string) {
	r.TraceID = traceID
}

func (api *GraphAPI) Login(ctx context.Context, input LoginInput) (*LoginResponse, error) {
	mutation := `
		mutation Login($input: LoginInput!) {
			login(input: $input) {
				token
			}
		}
	`

	variables := map[string]interface{}{
		"input": input,
	}

	reqBody, err := json.Marshal(GraphQLRequestBody{
		OperationName: lo.ToPtr("Login"),
		Query:         &mutation,
		Variables:     variables,
	})
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, api.endpoint, bytes.NewBuffer(reqBody))

	graphqlResponse := new(LoginResponse)
	err = api.Do(req, graphqlResponse)
	if err != nil {
		return graphqlResponse, err
	}

	// Check for data in the GraphQL response
	if data, ok := graphqlResponse.Data["login"]; ok {
		// Check for the "token" field
		if token, ok := data["token"].(string); ok {
			graphqlResponse.Token = token
		}
	}

	return graphqlResponse, nil
}
