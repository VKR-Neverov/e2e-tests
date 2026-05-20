package graphapi

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/samber/lo"
)

type Transaction struct {
	ID string `json:"id,omitempty"`
}

type TransactionsResponse struct {
	Errors  Errors                               `json:"errors"`
	Data    map[string]map[string][]*Transaction `json:"data"`
	TraceID string
}

func (i *TransactionsResponse) Error() Errors {
	return i.Errors
}

func (i *TransactionsResponse) SetTraceID(traceID string) {
	i.TraceID = traceID
}

func (api *GraphAPI) ListTransactions(ctx context.Context) ([]*Transaction, error) {
	query := `
		query Transactions {
			transactions {
				items {
					id
				}
			}
		}
	`

	reqBody, err := json.Marshal(GraphQLRequestBody{
		OperationName: lo.ToPtr("Transactions"),
		Query:         &query,
		Variables:     nil,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, api.endpoint, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	var graphqlResponse TransactionsResponse
	err = api.Do(req, &graphqlResponse)
	if err != nil {
		return nil, err
	}

	return graphqlResponse.Data["transactions"]["items"], nil
}
