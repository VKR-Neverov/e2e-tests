package graphapi

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/samber/lo"
	"net/http"
)

type CreateRegistrationInput struct{}

type CreateRegistrationResponse struct {
	Errors  Errors                             `json:"errors"`
	Data    map[string]map[string]*FlowPayload "data"
	TraceID string
}

type FlowPayload struct {
	ID    string
	State string
	Type  string
}

func (r *CreateRegistrationResponse) Error() Errors {
	return r.Errors
}

func (r *CreateRegistrationResponse) SetTraceID(traceID string) {
	r.TraceID = traceID
}

func (api *GraphAPI) CreateRegistration(ctx context.Context, input CreateRegistrationInput) (*FlowPayload, error) {
	mutation := `
		mutation CreateRegistration {
        registrationMutations {
            createRegistration {
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
		OperationName: lo.ToPtr("CreateRegistration"),
		Query:         &mutation,
		Variables:     variables,
	})
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, api.endpoint, bytes.NewBuffer(reqBody))

	graphqlResponse := new(CreateRegistrationResponse)
	err = api.Do(req, graphqlResponse)
	if err != nil {
		return nil, err
	}

	return graphqlResponse.Data["registrationMutations"]["createRegistration"], nil
}
