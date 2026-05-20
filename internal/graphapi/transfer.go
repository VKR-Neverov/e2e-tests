package graphapi

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/samber/lo"
)

type TransferInput struct {
	ReceiverAddress string  `json:"receiver_address"`
	TokenAmount     float64 `json:"token_amount"`
	Chain           string  `json:"chain"`
	GasLimit        *int64  `json:"gas_limit"`
}

type TransferResponse struct {
	Errors  Errors                  `json:"errors"`
	Data    map[string]*HashPayload "data"
	TraceID string
}

type HashPayload struct {
	Hash string `json:"hash"`
}

func (r *TransferResponse) Error() Errors {
	return r.Errors
}

func (r *TransferResponse) SetTraceID(traceID string) {
	r.TraceID = traceID
}

func (api *GraphAPI) Transfer(ctx context.Context, input TransferInput) (*HashPayload, error) {
	mutation := `
		mutation Transfer($input: TransferInput!) {
			transfer(input: $input) {
		    	hash
		    }
		}
	`

	variables := map[string]interface{}{
		"input": input,
	}

	reqBody, err := json.Marshal(GraphQLRequestBody{
		OperationName: lo.ToPtr("Transfer"),
		Query:         &mutation,
		Variables:     variables,
	})
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, api.endpoint, bytes.NewBuffer(reqBody))

	graphqlResponse := new(TransferResponse)
	err = api.Do(req, graphqlResponse)
	if err != nil {
		return nil, err
	}

	return graphqlResponse.Data["transfer"], nil
}
