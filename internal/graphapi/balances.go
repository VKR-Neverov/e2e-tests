package graphapi

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/samber/lo"
	"net/http"
)

type BalancesFilter struct {
	AddressEq string `json:"addressEq,omitempty"`
	ChainEq   string `json:"chainEq,omitempty"`
	TokenEq   string `json:"tokenEq,omitempty"`
}

type Balances struct {
	Balance    float64 `json:"balance"`
	UsdBalance float64 `json:"usdBalance"`
}

type BalancesResponse struct {
	Errors  Errors               `json:"errors"`
	Data    map[string]*Balances `json:"data"`
	TraceID string
}

func (w *BalancesResponse) Error() Errors {
	return w.Errors
}

func (w *BalancesResponse) SetTraceID(traceID string) {
	w.TraceID = traceID
}

func (api *GraphAPI) Balances(ctx context.Context, filter BalancesFilter) (*Balances, error) {
	query := `
		query Balances($filter: BalanceFilter!) {
        	balance(filter: $filter) {
            	balance
            	usdBalance
        	}
    	}
	`

	reqBody, err := json.Marshal(GraphQLRequestBody{
		OperationName: lo.ToPtr("Balances"),
		Query:         &query,
		Variables: map[string]interface{}{
			"filter": filter,
		},
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, api.endpoint, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	var graphqlResponse BalancesResponse
	err = api.Do(req, &graphqlResponse)
	if err != nil {
		return nil, err
	}

	return graphqlResponse.Data["balance"], nil
}
