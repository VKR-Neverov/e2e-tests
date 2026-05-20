package service

import (
	"context"
	"fmt"

	"github.com/fidesy-pay/e2e-tests/internal/constants"
	"github.com/fidesy-pay/e2e-tests/internal/graphapi"
)

func PayInvoice(ctx context.Context, address string, chain string, tokenAmount float64) (string, error) {
	api := graphapi.New(constants.GraphQLEnpoint)

	loginResp, err := api.Login(ctx, graphapi.LoginInput{
		Username: "highbank",
		Password: "highbank",
	})
	if err != nil {
		return "", fmt.Errorf("api.Login: %w", err)
	}

	api.SetToken(loginResp.Token)

	gasLimit := int64(100000)
	hash := ""
	for gasLimit < 500000 {
		hashResp, err := api.Transfer(ctx, graphapi.TransferInput{
			ReceiverAddress: address,
			Chain:           chain,
			TokenAmount:     tokenAmount,
			GasLimit:        &gasLimit,
		})
		if err != nil {
			gasLimit += 50000
			continue
		}

		hash = hashResp.Hash
		break
	}

	return hash, nil
}
