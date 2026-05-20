package graphapi

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/samber/lo"
	"net/http"
)

type SetEmailInput struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

type SetEmailResponse struct {
	Errors  Errors                             `json:"errors"`
	Data    map[string]map[string]*FlowPayload "data"
	TraceID string
}

func (r *SetEmailResponse) Error() Errors {
	return r.Errors
}

func (r *SetEmailResponse) SetTraceID(traceID string) {
	r.TraceID = traceID
}

func (api *GraphAPI) SetEmail(ctx context.Context, input SetEmailInput) (*FlowPayload, error) {
	mutation := `
		mutation SetEmail($input: SetEmailInput) {
        registrationMutations {
            setEmail(input: $input) {
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
		OperationName: lo.ToPtr("SetEmail"),
		Query:         &mutation,
		Variables:     variables,
	})
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, api.endpoint, bytes.NewBuffer(reqBody))

	graphqlResponse := new(SetEmailResponse)
	err = api.Do(req, graphqlResponse)
	if err != nil {
		return nil, err
	}

	return graphqlResponse.Data["registrationMutations"]["setEmail"], nil
}
