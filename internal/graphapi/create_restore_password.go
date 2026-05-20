package graphapi

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/samber/lo"
	"net/http"
)

type CreateRestorePasswordInput struct {
	Email string `json:"email"`
}

type CreateRestorePasswordResponse struct {
	Errors  Errors                             `json:"errors"`
	Data    map[string]map[string]*FlowPayload "data"
	TraceID string
}

func (r *CreateRestorePasswordResponse) Error() Errors {
	return r.Errors
}

func (r *CreateRestorePasswordResponse) SetTraceID(traceID string) {
	r.TraceID = traceID
}

func (api *GraphAPI) CreateRestorePassword(ctx context.Context, input CreateRestorePasswordInput) (*FlowPayload, error) {
	mutation := `
		mutation CreateRestorePassword($input: CreateRestorePasswordInput!) {
        restorePasswordMutations {
            createRestorePassword(input: $input) {
                id
                state
                type
            }
        }
    }
	`

	variables := map[string]interface{}{
		"input": input,
	}

	reqBody, err := json.Marshal(GraphQLRequestBody{
		OperationName: lo.ToPtr("CreateRestorePassword"),
		Query:         &mutation,
		Variables:     variables,
	})
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, api.endpoint, bytes.NewBuffer(reqBody))

	graphqlResponse := new(CreateRestorePasswordResponse)
	err = api.Do(req, graphqlResponse)
	if err != nil {
		return nil, err
	}

	return graphqlResponse.Data["restorePasswordMutations"]["createRestorePassword"], nil
}
