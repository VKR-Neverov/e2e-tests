package graphapi

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/samber/lo"
	"net/http"
)

type SetPasswordInput struct {
	ID       string `json:"id"`
	Password string `json:"password"`
}

type SetPasswordResponse struct {
	Errors  Errors                             `json:"errors"`
	Data    map[string]map[string]*FlowPayload "data"
	TraceID string
}

func (r *SetPasswordResponse) Error() Errors {
	return r.Errors
}

func (r *SetPasswordResponse) SetTraceID(traceID string) {
	r.TraceID = traceID
}

func (api *GraphAPI) SetPassword(ctx context.Context, input SetPasswordInput) (*FlowPayload, error) {
	mutation := `
		 mutation SetPassword($input: SetPasswordInput!) {
        restorePasswordMutations {
            setPassword(input: $input) {
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
		OperationName: lo.ToPtr("SetPassword"),
		Query:         &mutation,
		Variables:     variables,
	})
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, api.endpoint, bytes.NewBuffer(reqBody))

	graphqlResponse := new(SetPasswordResponse)
	err = api.Do(req, graphqlResponse)
	if err != nil {
		return nil, err
	}

	return graphqlResponse.Data["restorePasswordMutations"]["setPassword"], nil
}
