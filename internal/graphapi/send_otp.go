package graphapi

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/samber/lo"
)

type SendOtpInput struct {
	Email string `json:"email"`
}

type SendOtpResponse struct {
	Errors  Errors                     `json:"errors"`
	Data    map[string]*SendOtpPayload "data"
	TraceID string
}

type SendOtpPayload struct {
	EmailCodeID string `json:"email_code_id"`
}

func (r *SendOtpResponse) Error() Errors {
	return r.Errors
}

func (r *SendOtpResponse) SetTraceID(traceID string) {
	r.TraceID = traceID
}

func (api *GraphAPI) SendEmailOtp(ctx context.Context, input SendOtpInput) (*SendOtpPayload, error) {
	mutation := `
		  mutation SendCode($input: SendEmailOtpInput!) {
        	sendEmailOtp(input: $input) {
                email_code_id
        	}

    	}
	`

	variables := map[string]interface{}{
		"input": input,
	}

	reqBody, err := json.Marshal(GraphQLRequestBody{
		OperationName: lo.ToPtr("SendCode"),
		Query:         &mutation,
		Variables:     variables,
	})
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, api.endpoint, bytes.NewBuffer(reqBody))

	graphqlResponse := new(SendOtpResponse)
	err = api.Do(req, graphqlResponse)
	if err != nil {
		return nil, err
	}

	return graphqlResponse.Data["sendEmailOtp"], nil
}
