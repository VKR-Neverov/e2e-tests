package graphapi

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/samber/lo"
	"net/http"
)

type MainBalance struct {
	UsdBalance float64 `json:"usdBalance"`
}

type MainBalanceResponse struct {
	Errors  Errors                  `json:"errors"`
	Data    map[string]*MainBalance `json:"data"`
	TraceID string
}

func (w *MainBalanceResponse) Error() Errors {
	return w.Errors
}

func (w *MainBalanceResponse) SetTraceID(traceID string) {
	w.TraceID = traceID
}

func (api *GraphAPI) MainBalance(ctx context.Context) (*MainBalance, error) {
	query := `
		query MainBalance {
        	mainBalance {
            	usdBalance
        	}
    	}
	`

	reqBody, err := json.Marshal(GraphQLRequestBody{
		OperationName: lo.ToPtr("MainBalance"),
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

	var graphqlResponse MainBalanceResponse
	err = api.Do(req, &graphqlResponse)
	if err != nil {
		return nil, err
	}

	return graphqlResponse.Data["mainBalance"], nil
}
