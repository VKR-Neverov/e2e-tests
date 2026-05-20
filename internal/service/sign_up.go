package service

import (
	"context"
	"fmt"

	"github.com/fidesy-pay/e2e-tests/internal/graphapi"
	"github.com/fidesy-pay/e2e-tests/internal/model"
	"github.com/google/uuid"
)

func SignUp(ctx context.Context, api *graphapi.GraphAPI) (*model.User, error) {
	username := fmt.Sprintf("test_%s", uuid.NewString())
	email := fmt.Sprintf("test_%s", uuid.NewString())
	password := uuid.NewString()

	sendOtpResult, err := api.SendEmailOtp(ctx, graphapi.SendOtpInput{Email: email})
	if err != nil {
		return nil, fmt.Errorf("api.SendEmailOtp: %w", err)
	}

	_, err = api.SignUp(ctx, graphapi.SignUpInput{
		Username:    username,
		Password:    password,
		EmailCodeID: sendOtpResult.EmailCodeID,
		EmailCode:   "111111",
		Email:       email,
	})
	if err != nil {
		return nil, fmt.Errorf("api.SignUp: %w", err)
	}

	resp, err := api.Login(ctx, graphapi.LoginInput{
		Username: username,
		Password: password,
	})
	if err != nil {
		return nil, fmt.Errorf("api.Login: %w", err)
	}

	api.SetToken(resp.Token)

	client, err := api.Me(ctx)
	if err != nil {
		return nil, fmt.Errorf("api.Me: %w", err)
	}

	return &model.User{
		Username: username,
		Password: password,
		Email:    email,
		APIKey:   client.APIKey,
	}, nil
}
