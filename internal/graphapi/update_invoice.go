package graphapi

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/samber/lo"
	"net/http"
)

type UpdateInvoiceInput struct {
	ID    string `json:"id"`
	Chain string `json:"chain"`
	Token string `json:"token"`
}

type UpdateInvoiceResponse struct {
	Errors  Errors                                    `json:"errors"`
	Data    map[string]map[string]map[string]*Invoice "data"
	TraceID string
	Invoice *Invoice
}

func (r *UpdateInvoiceResponse) Error() Errors {
	return r.Errors
}

func (r *UpdateInvoiceResponse) SetTraceID(traceID string) {
	r.TraceID = traceID
}

func (api *GraphAPI) UpdateInvoice(ctx context.Context, input UpdateInvoiceInput) (*UpdateInvoiceResponse, error) {
	mutation := `
		mutation UpdateInvoice($input: UpdateInvoiceInput!) {
  			invoiceMutations {
    			updateInvoice(input: $input) {
      				invoice {
        				address
						chain
						created_at
						id
						status
						token
						token_amount
						usd_amount
      				}
    			}
  			}
		}
	`

	variables := map[string]interface{}{
		"input": input,
	}

	reqBody, err := json.Marshal(GraphQLRequestBody{
		OperationName: lo.ToPtr("UpdateInvoice"),
		Query:         &mutation,
		Variables:     variables,
	})
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, api.endpoint, bytes.NewBuffer(reqBody))

	graphqlResponse := new(UpdateInvoiceResponse)
	err = api.Do(req, graphqlResponse)
	if err != nil {
		return nil, err
	}

	graphqlResponse.Invoice = graphqlResponse.Data["invoiceMutations"]["updateInvoice"]["invoice"]

	return graphqlResponse, nil
}
