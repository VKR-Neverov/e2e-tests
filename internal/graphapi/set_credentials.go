package graphapi

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/samber/lo"
	"net/http"
)

type SetCredentialsInput struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type SetCredentialsResponse struct {
	Errors  Errors                             `json:"errors"`
	Data    map[string]map[string]*FlowPayload "data"
	TraceID string
}

func (r *SetCredentialsResponse) Error() Errors {
	return r.Errors
}

func (r *SetCredentialsResponse) SetTraceID(traceID string) {
	r.TraceID = traceID
}

func (api *GraphAPI) SetCredentials(ctx context.Context, input SetCredentialsInput) (*FlowPayload, error) {
	mutation := `
		 mutation SetCredentials($input: SetCredentialsInput!) {
        registrationMutations {
            setCredentials(input: $input) {
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
		OperationName: lo.ToPtr("SetCredentials"),
		Query:         &mutation,
		Variables:     variables,
	})
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, api.endpoint, bytes.NewBuffer(reqBody))

	graphqlResponse := new(SetCredentialsResponse)
	err = api.Do(req, graphqlResponse)
	if err != nil {
		return nil, err
	}

	return graphqlResponse.Data["registrationMutations"]["setCredentials"], nil
}
