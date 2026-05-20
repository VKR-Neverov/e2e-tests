package graphapi

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/samber/lo"
	"net/http"
	"time"
)

type Invoice struct {
	ID          string    `json:"id,omitempty"`
	Address     string    `json:"address,omitempty"`
	Chain       string    `json:"chain,omitempty"`
	Token       string    `json:"token,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	Status      string    `json:"status,omitempty"`
	TokenAmount float64   `json:"token_amount,omitempty"`
	UsdAmount   float64   `json:"usd_amount,omitempty"`
}

type InvoicesResponse struct {
	Errors  Errors                           `json:"errors"`
	Data    map[string]map[string][]*Invoice `json:"data"`
	TraceID string
}

func (i *InvoicesResponse) Error() Errors {
	return i.Errors
}

func (i *InvoicesResponse) SetTraceID(traceID string) {
	i.TraceID = traceID
}

func (api *GraphAPI) ListInvoices(ctx context.Context) ([]*Invoice, error) {
	query := `
		query Invoices($filter: InvoicesFilter) {
			invoices(filter: $filter) {
				items {
					id
					address
					chain
					token
					created_at
					status	
					token_amount
					usd_amount
				}
			}
		}
	`

	reqBody, err := json.Marshal(GraphQLRequestBody{
		OperationName: lo.ToPtr("Invoices"),
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

	var graphqlResponse InvoicesResponse
	err = api.Do(req, &graphqlResponse)
	if err != nil {
		return nil, err
	}

	return graphqlResponse.Data["invoices"]["items"], nil
}
