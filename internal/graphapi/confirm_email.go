package graphapi

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/samber/lo"
	"net/http"
)

type ConfirmEmailInput struct {
	ID   string `json:"id"`
	Code string `json:"code"`
}

type ConfirmEmailResponse struct {
	Errors  Errors                             `json:"errors"`
	Data    map[string]map[string]*FlowPayload "data"
	TraceID string
}

func (r *ConfirmEmailResponse) Error() Errors {
	return r.Errors
}

func (r *ConfirmEmailResponse) SetTraceID(traceID string) {
	r.TraceID = traceID
}

func (api *GraphAPI) ConfirmEmail(ctx context.Context, input ConfirmEmailInput) (*FlowPayload, error) {
	mutation := `
		mutation ConfirmEmail($input: ConfirmEmailInput!) {
        registrationMutations {
            confirmEmail(input: $input) {
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
		OperationName: lo.ToPtr("ConfirmEmail"),
		Query:         &mutation,
		Variables:     variables,
	})
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, api.endpoint, bytes.NewBuffer(reqBody))

	graphqlResponse := new(ConfirmEmailResponse)
	err = api.Do(req, graphqlResponse)
	if err != nil {
		return nil, err
	}

	return graphqlResponse.Data["registrationMutations"]["confirmEmail"], nil
}
