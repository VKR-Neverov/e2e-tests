package restapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	apiKeyHeader = "api-key"
)

type RestAPI struct {
	endpoint string
	client   *http.Client
}

func New(endpoint string) *RestAPI {
	return &RestAPI{
		endpoint: endpoint,
		client:   &http.Client{},
	}
}

type CreateInvoiceResponse struct {
	ID string `json:"invoice_id"`
}

func (api *RestAPI) CreateInvoice(ctx context.Context, apiKey string, usdAmount float64) (*CreateInvoiceResponse, error) {
	reqBody, _ := json.Marshal(map[string]interface{}{
		"usd_amount": usdAmount,
	})

	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/invoice", api.endpoint), bytes.NewBuffer(reqBody))
	req.Header.Set(apiKeyHeader, apiKey)

	resp, err := api.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client.Do: %w", err)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	createInvoiceResponse := new(CreateInvoiceResponse)
	if err = json.Unmarshal(body, &createInvoiceResponse); err != nil {
		return nil, err
	}

	return createInvoiceResponse, nil
}
