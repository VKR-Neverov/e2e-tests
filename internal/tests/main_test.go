package basic_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/fidesy-pay/e2e-tests/internal/constants"
	"github.com/fidesy-pay/e2e-tests/internal/graphapi"
	"github.com/fidesy-pay/e2e-tests/internal/model"
	"github.com/fidesy-pay/e2e-tests/internal/restapi"
	"github.com/fidesy-pay/e2e-tests/internal/service"
	"github.com/stretchr/testify/suite"
)

const (
	endpoint        = "http://pay.fidesy.tech:15000/query"
	restAPIEndpoint = "http://pay.fidesy.tech:15000/api"
)

type testSuite struct {
	suite.Suite

	api     *graphapi.GraphAPI
	restAPI *restapi.RestAPI
	user    *model.User
}

func TestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(testSuite))
}

func (s *testSuite) SetupSuite() {
	s.api = graphapi.New(endpoint)
	s.restAPI = restapi.New(restAPIEndpoint)

	user, err := service.SignUp(context.Background(), s.api)
	s.Require().NoError(err, fmt.Errorf("service.SignUp: %w", err))

	resp, err := s.api.Login(context.Background(), graphapi.LoginInput{
		Username: user.Username,
		Password: user.Password,
	})

	s.Require().NoError(err, fmt.Errorf("api.Login: %w", err))
	s.Require().NotEmpty(resp.Token)

	s.api.SetToken(resp.Token)

	client, err := s.api.Me(context.Background())
	s.Require().NoError(err, fmt.Errorf("api.Me: %w", err))

	user.APIKey = client.APIKey

	s.user = user
	s.api.SetToken(resp.Token)
}

func (s *testSuite) Test_Wallets() {
	r := s.Require()
	ctx := context.Background()

	api := graphapi.New(constants.GraphQLEnpoint)

	_, err := service.SignUp(ctx, api)
	r.NoError(err, fmt.Errorf("service.SignUp: %w", err))

	wallets, err := api.ListWallets(context.Background())

	r.NoError(err, fmt.Errorf("api.ListWallets: %w", err))
	r.Equal(3, len(wallets))
}

func (s *testSuite) Test_Invoices() {
	r := s.Require()
	ctx := context.Background()

	api := graphapi.New(constants.GraphQLEnpoint)

	user, err := service.SignUp(ctx, api)
	r.NoError(err, fmt.Errorf("service.SignUp: %w", err))

	invoices, err := api.ListInvoices(ctx)
	r.NoError(err, fmt.Errorf("api.ListInvoices: %w", err))
	r.Equal(0, len(invoices))

	invoiceResp, err := s.restAPI.CreateInvoice(ctx, user.APIKey, 10)
	r.NoError(err, fmt.Errorf("restAPI.CreateInvoice: %w", err))
	r.NotEmpty(invoiceResp.ID)

	invoices, err = api.ListInvoices(ctx)
	r.NoError(err, fmt.Errorf("api.ListInvoices: %w", err))
	r.Equal(1, len(invoices))
	r.Equal("NEW", invoices[0].Status)

	updateInvoiceResp, err := api.UpdateInvoice(ctx, graphapi.UpdateInvoiceInput{
		ID:    invoiceResp.ID,
		Chain: "polygon",
		Token: "matic-network",
	})
	r.NoError(err, fmt.Errorf("api.UpdateInvoice: %w", err))
	r.NotNil(updateInvoiceResp)
	r.Equal("polygon", updateInvoiceResp.Invoice.Chain)
	r.Equal("matic-network", updateInvoiceResp.Invoice.Token)

	invoices, err = api.ListInvoices(ctx)
	r.NoError(err, fmt.Errorf("api.ListInvoices: %w", err))
	r.Equal(1, len(invoices))
	r.Equal("PENDING", invoices[0].Status)
}

func (s *testSuite) Test_UpdateInvoice() {
	r := s.Require()
	ctx := context.Background()

	api := graphapi.New(constants.GraphQLEnpoint)

	user, err := service.SignUp(ctx, api)
	r.NoError(err, fmt.Errorf("service.SignUp: %w", err))

	testCases := []struct {
		name      string
		chain     string
		token     string
		usdAmount float64
	}{
		{
			name:      "ArbitrumEthereum",
			chain:     "arbitrum",
			token:     "ethereum",
			usdAmount: 10,
		},
		{
			name:      "PolygonMatic",
			chain:     "polygon",
			token:     "matic-network",
			usdAmount: 55.12,
		},
		{
			name:      "OptimismEthereum",
			chain:     "optimism",
			token:     "ethereum",
			usdAmount: 0.99,
		},
	}

	for _, tt := range testCases {
		tt := tt
		s.Run(tt.name, func() {

			invoiceResp, err := s.restAPI.CreateInvoice(ctx, user.APIKey, tt.usdAmount)
			r.NoError(err, fmt.Errorf("restAPI.CreateInvoice: %w", err))
			r.NotEmpty(invoiceResp.ID)

			updateInvoiceResp, err := api.UpdateInvoice(ctx, graphapi.UpdateInvoiceInput{
				ID:    invoiceResp.ID,
				Chain: tt.chain,
				Token: tt.token,
			})
			r.NoError(err, fmt.Errorf("api.UpdateInvoice: %w", err))

			invoice := updateInvoiceResp.Invoice

			r.NotEmpty(invoice.Address)
			r.Equal(tt.chain, invoice.Chain)
			r.Equal(tt.token, invoice.Token)
			r.NotZero(invoice.CreatedAt)
			r.Equal("PENDING", invoice.Status)
			r.NotZero(invoice.TokenAmount)
			r.Equal(tt.usdAmount, invoice.UsdAmount)
		})
	}
}

func (s *testSuite) Test_DeactivateClient() {
	s.T().Parallel()

	r := s.Require()
	ctx := context.Background()

	api := graphapi.New(endpoint)

	user, err := service.SignUp(ctx, api)
	r.NoError(err, fmt.Errorf("service.SignUp: %w", err))

	loginResp, err := api.Login(ctx, graphapi.LoginInput{
		Username: user.Username,
		Password: user.Password,
	})
	r.NoError(err, fmt.Errorf("api.Login: %w", err))
	r.NotEmpty(loginResp.Token)

	api.SetToken(loginResp.Token)

	err = api.DeactivateClient(ctx)
	r.NoError(err, fmt.Errorf("api.DeactivateClient: %w", err))

	time.Sleep(1000 * time.Millisecond)

	loginResp, err = api.Login(ctx, graphapi.LoginInput{
		Username: user.Username,
		Password: user.Password,
	})
	r.Error(err, "Client is not deactivated in AuthService after 1 second.")

	r.Equal("GraphQL errors: [authClient.Login: rpc error: code = NotFound desc = storage.GetUser: entity not found]", err.Error())
}

// func (s *testSuite) Test_RestorePassword() {
// 	s.T().Parallel()

// 	r := s.Require()
// 	ctx := context.Background()

// 	api := graphapi.New(endpoint)

// 	user, err := service.SignUp(ctx, api)
// 	r.NoError(err, fmt.Errorf("service.SignUp: %w", err))

// 	loginResp, err := api.Login(ctx, graphapi.LoginInput{
// 		Username: user.Username,
// 		Password: user.Password,
// 	})
// 	r.NoError(err, fmt.Errorf("api.Login: %w", err))
// 	r.NotEmpty(loginResp.Token)

// 	newPassword := uuid.NewString()

// 	flow, err := api.CreateRestorePassword(ctx, graphapi.CreateRestorePasswordInput{
// 		Email: user.Email,
// 	})
// 	r.NoError(err, fmt.Errorf("api.CreateRestorePassword: %w", err))
// 	r.NotEmpty(flow.ID)
// 	r.Equal("WAITING_EMAIL_CONFIRMATION", flow.State)

// 	flow, err = api.ConfirmEmail(ctx, graphapi.ConfirmEmailInput{
// 		ID:   flow.ID,
// 		Code: "111111",
// 	})
// 	r.NoError(err, fmt.Errorf("api.ConfirmEmail: %w", err))
// 	r.Equal("WAITING_PASSWORD", flow.State)

// 	flow, err = api.SetPassword(ctx, graphapi.SetPasswordInput{
// 		ID:       flow.ID,
// 		Password: newPassword,
// 	})
// 	r.NoError(err, fmt.Errorf("api.SetPassword: %w", err))
// 	r.Equal("COMPLETED", flow.State)

// 	loginResp, err = api.Login(ctx, graphapi.LoginInput{
// 		Username: user.Username,
// 		Password: newPassword,
// 	})
// 	r.NoError(err, fmt.Errorf("api.Login: %w", err))
// 	r.NotEmpty(loginResp.Token)
// }

// func (s *testSuite) Test_Transfer() {
// 	s.T().Parallel()

// 	r := s.Require()
// 	ctx := context.Background()

// 	api := graphapi.New(constants.GraphQLEnpoint)

// 	user, err := service.SignUp(ctx, api)
// 	r.NoError(err, fmt.Errorf("service.SignUp: %w", err))

// 	invoiceResp, err := s.restAPI.CreateInvoice(ctx, user.APIKey, 4.00)
// 	r.NoError(err, fmt.Errorf("restAPI.CreateInvoice: %w", err))

// 	updatedInvoice, err := api.UpdateInvoice(ctx, graphapi.UpdateInvoiceInput{
// 		ID:    invoiceResp.ID,
// 		Chain: "arbitrum",
// 		Token: "ethereum",
// 	})
// 	r.NoError(err, fmt.Errorf("api.UpdateInvoice: %w", err))

// 	invoice := updatedInvoice.Invoice

// 	_, err = service.PayInvoice(
// 		ctx,
// 		invoice.Address,
// 		invoice.Chain,
// 		1.2*invoice.TokenAmount,
// 	)
// 	r.NoError(err, fmt.Errorf("service.PayInvoice: %w", err))

// 	defer func() {
// 		counter := 0
// 		var gasLimit int64

// 		for counter < 10 {
// 			_, err = api.Transfer(ctx, graphapi.TransferInput{
// 				ReceiverAddress: constants.FundsWallet,
// 				TokenAmount:     invoice.TokenAmount,
// 				Chain:           "arbitrum",
// 				GasLimit:        &gasLimit,
// 			})
// 			if err == nil {
// 				return
// 			}

// 			counter++
// 			gasLimit += 50000
// 		}

// 		log.Println("does not return funds")
// 	}()

// 	retries := 0
// 	for {
// 		invoices, err := api.ListInvoices(ctx)
// 		r.NoError(err, fmt.Errorf("api.ListInvoices: %w", err))

// 		filteredInvoices := lo.Filter(invoices, func(inv *graphapi.Invoice, _ int) bool {
// 			return inv.ID == invoice.ID
// 		})

// 		inv := filteredInvoices[0]
// 		if inv.Status == "SUCCESS" {
// 			break
// 		}

// 		time.Sleep(time.Second)
// 		retries++

// 		if retries == 30 {
// 			r.Fail("Invoice not successful after 30 seconds")
// 		}
// 	}
// }
