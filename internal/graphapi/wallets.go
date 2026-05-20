package graphapi

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/samber/lo"
	"net/http"
)

type Wallet struct {
	Address string  `json:"address"`
	Balance float64 `json:"balance"`
	Chain   string  `json:"chain"`
}

type WalletsResponse struct {
	Errors  Errors                          `json:"errors"`
	Data    map[string]map[string][]*Wallet `json:"data"`
	TraceID string
}

func (w *WalletsResponse) Error() Errors {
	return w.Errors
}

func (w *WalletsResponse) SetTraceID(traceID string) {
	w.TraceID = traceID
}

func (api *GraphAPI) ListWallets(ctx context.Context) ([]*Wallet, error) {
	query := `
		query Wallets {
			wallets {
				items {
					address
					chain
				}
			}
		}
	`

	reqBody, err := json.Marshal(GraphQLRequestBody{
		OperationName: lo.ToPtr("Wallets"),
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

	var graphqlResponse WalletsResponse
	err = api.Do(req, &graphqlResponse)
	if err != nil {
		return nil, err
	}

	return graphqlResponse.Data["wallets"]["items"], nil
}
