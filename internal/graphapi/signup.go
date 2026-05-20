package graphapi

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/samber/lo"
)

type SignUpInput struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	EmailCodeID string `json:"email_code_id"`
	EmailCode   string `json:"email_code"`
	Email       string `json:"email"`
}

type SignUpResponse struct {
	Errors Errors `json:"errors"`
	Data   map[string]map[string]interface{} "data"
	TraceID string
}

func (r *SignUpResponse) Error() Errors {
	return r.Errors
}

func (r *SignUpResponse) SetTraceID(traceID string) {
	r.TraceID = traceID
}

func (api *GraphAPI) SignUp(ctx context.Context, input SignUpInput) (*SignUpResponse, error) {
	mutation := `
		mutation SignUp($input: SignUpInput!) {
			signUp(input: $input) {
				success
			}
		}
	`

	variables := map[string]interface{}{
		"input": input,
	}

	reqBody, err := json.Marshal(GraphQLRequestBody{
		OperationName: lo.ToPtr("SignUp"),
		Query:         &mutation,
		Variables:     variables,
	})
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, api.endpoint, bytes.NewBuffer(reqBody))

	graphqlResponse := new(SignUpResponse)
	err = api.Do(req, graphqlResponse)
	if err != nil {
		return graphqlResponse, err
	}

	return graphqlResponse, nil
}
