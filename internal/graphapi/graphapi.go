package graphapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type GraphAPI struct {
	client   *http.Client
	token    string
	endpoint string
}

func New(endpoint string) *GraphAPI {
	return &GraphAPI{
		client:   &http.Client{},
		endpoint: endpoint,
	}
}

func (api *GraphAPI) SetToken(token string) {
	api.token = token
}

func (api *GraphAPI) Do(req *http.Request, dst Response) error {
	req.Header.Set("Authorization", "Bearer "+api.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := api.client.Do(req)
	defer func() {
		if resp != nil {
			dst.SetTraceID(resp.Header.Get("x-trace-id"))
		}
	}()
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(body, &dst); err != nil {
		return err
	}

	if err = checkErrors(dst.Error()); err != nil {
		return err
	}

	return nil
}

func checkErrors(errors Errors) error {
	if len(errors) > 0 {
		// Handle errors
		errorMessages := make([]string, len(errors))
		for i, err := range errors {
			errorMessages[i] = err.Message
		}
		return fmt.Errorf("GraphQL errors: %v", errorMessages)
	}

	return nil
}
