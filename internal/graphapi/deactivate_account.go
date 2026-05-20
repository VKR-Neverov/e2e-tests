package graphapi

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/samber/lo"
	"net/http"
)

type DeactivateClientResponse struct {
	Errors  Errors                                       `json:"errors"`
	Data    map[string]map[string]map[string]interface{} "data"
	TraceID string
	Invoice *Invoice
}

func (r *DeactivateClientResponse) Error() Errors {
	return r.Errors
}

func (r *DeactivateClientResponse) SetTraceID(traceID string) {
	r.TraceID = traceID
}

func (api *GraphAPI) DeactivateClient(ctx context.Context) error {
	mutation := `
		mutation Deactivate {
  			clientMutations {
				deactivate {
					id
				}
			}
		}
	`

	variables := map[string]interface{}{}

	reqBody, err := json.Marshal(GraphQLRequestBody{
		OperationName: lo.ToPtr("Deactivate"),
		Query:         &mutation,
		Variables:     variables,
	})
	if err != nil {
		return err
	}

	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, api.endpoint, bytes.NewBuffer(reqBody))

	graphqlResponse := new(DeactivateClientResponse)
	err = api.Do(req, graphqlResponse)
	if err != nil {
		return err
	}

	return nil
}
