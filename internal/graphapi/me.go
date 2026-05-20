package graphapi

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/samber/lo"
	"net/http"
)

type MeResponse struct {
	Errors   Errors                            `json:"errors"`
	Data     map[string]map[string]interface{} "data"
	APIKey   string
	ClientID string
	TraceID  string
}

func (r *MeResponse) Error() Errors {
	return r.Errors
}

func (r *MeResponse) SetTraceID(traceID string) {
	r.TraceID = traceID
}

func (api *GraphAPI) Me(ctx context.Context) (*MeResponse, error) {
	query := `
		query Me {
			me {
				id
				api_key
			}
		}
	`

	reqBody, err := json.Marshal(GraphQLRequestBody{
		OperationName: lo.ToPtr("Me"),
		Query:         &query,
		Variables:     nil,
	})
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, api.endpoint, bytes.NewBuffer(reqBody))

	graphqlResponse := new(MeResponse)
	err = api.Do(req, graphqlResponse)
	if err != nil {
		return graphqlResponse, err
	}

	// Check for data in the GraphQL response
	if data, ok := graphqlResponse.Data["me"]; ok {
		// Check for the "token" field
		if apiKey, ok := data["api_key"].(string); ok {
			graphqlResponse.APIKey = apiKey
		}

		if clientID, ok := data["id"].(string); ok {
			graphqlResponse.ClientID = clientID
		}
	}

	return graphqlResponse, nil
}
