package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/fidesy-pay/e2e-tests/internal/graphapi"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

const endpoint = "http://facade.pay.fidesy.tech/query"

func main() {
	ctx := context.Background()

	wg := sync.WaitGroup{}

	//LoadTesting(ctx, "Login", Login)
	//LoadTesting(ctx, "Me", Me)

	wg.Add(1)
	go func() {
		defer wg.Done()
		LoadTestingWithRPS(ctx, "Profile", Profile)
	}()


	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	LoadTestingWithRPS(ctx, "SignUp", SignUp)
	// }()

	wg.Wait()
}

func LoadTestingWithRPS(ctx context.Context, name string, callback func(ctx context.Context) func() error) {
	requests := 10000
	rps := 100

	callbackFunc := callback(ctx)

	start := time.Now()
	for i := 0; i < requests; i++ {
		if i%100 == 0 {
			log.Printf("%d/%d", i, requests)
		}

		go func() {
			err := callbackFunc()
			if err != nil {
				fmt.Printf("callback: %v", err)
			}
		}()

		time.Sleep(time.Duration(1000/rps) * time.Millisecond)
	}

	log.Printf("%s Requests %d, workers %d, time %s, RPS %.0f", name, requests, rps, time.Since(start), float64(requests)/time.Since(start).Seconds())
}

func LoadTesting(ctx context.Context, name string, callback func(ctx context.Context) func() error) {
	requests := 100
	workers := 20

	workerPool := make(chan struct{}, workers)
	wg := sync.WaitGroup{}

	callbackFunc := callback(ctx)

	start := time.Now()
	for i := 0; i < requests; i++ {
		if i%100 == 0 {
			log.Printf("%d/%d", i, requests)
		}

		workerPool <- struct{}{}
		wg.Add(1)
		go func() {
			defer func() {
				wg.Done()
				<-workerPool
			}()

			err := callbackFunc()
			if err != nil {
				fmt.Printf("callback: %v", err)
			}
		}()
	}

	wg.Wait()
	log.Printf("%s Requests %d, workers %d, time %s, RPS %.0f", name, requests, workers, time.Since(start), float64(requests)/time.Since(start).Seconds())
}

func Login(ctx context.Context) func() error {
	return func() error {
		api := graphapi.New(endpoint)

		loginResp, err := api.Login(ctx, graphapi.LoginInput{
			Username: "weekend",
			Password: "weekend",
		})
		if err != nil {
			return fmt.Errorf("api.Login: %v", err)
		}

		if loginResp.Token == "" {
			return errors.New("token is empty")
		}

		return nil
	}
}

func Me(ctx context.Context) func() error {
	api := graphapi.New(endpoint)

	username := fmt.Sprintf("load_testing_%s", uuid.NewString())

	signUpResp, err := api.SignUp(ctx, graphapi.SignUpInput{
		Username: username,
		Password: "SecretPassword",
	})
	if err != nil {
		log.Fatalf("api.SignUp: %v", err)
	}

	api.SetToken(signUpResp.Token)

	return func() error {
		me, err := api.Me(ctx)
		if err != nil {
			return fmt.Errorf("api.Me: %w", err)
		}

		if me == nil || me.APIKey == "" {
			return errors.New("me data is empty")
		}

		return nil
	}
}

func Profile(ctx context.Context) func() error {
	api := graphapi.New(endpoint)

	//username := fmt.Sprintf("load_testing_%s", uuid.NewString())
	//
	//signUpResp, err := api.SignUp(ctx, graphapi.SignUpInput{
	//	Username: username,
	//	Password: "SecretPassword",
	//})
	//if err != nil {
	//	log.Fatalf("api.SignUp: %v", err)
	//}

	resp, err := api.Login(ctx, graphapi.LoginInput{
		Username: "weekend",
		Password: "weekend",
	})
	if err != nil {
		panic(err)
	}

	api.SetToken(resp.Token)

	return func() error {
		errGroup := errgroup.Group{}

		errGroup.Go(func() error {
			me, err := api.Me(ctx)
			if err != nil {
				return fmt.Errorf("api.Me: %w", err)
			}

			if me == nil || me.APIKey == "" {
				return errors.New("me data is empty")
			}

			return nil
		})

		errGroup.Go(func() error {
			_, err = api.ListInvoices(ctx)
			if err != nil {
				return fmt.Errorf("api.ListInvoices: %w", err)
			}

			return nil
		})

		errGroup.Go(func() error {
			_, err = api.MainBalance(ctx)
			if err != nil {
				return fmt.Errorf("api.MainBalance: %w", err)
			}

			return nil
		})

		wallets, err := api.ListWallets(ctx)
		if err != nil {
			return fmt.Errorf("api.ListWallets: %w", err)
		}

		if len(wallets) == 0 {
			return errors.New("no wallets found")
		}

		errGroup.Go(func() error {
			_, err = api.Balances(ctx, graphapi.BalancesFilter{
				AddressEq: wallets[0].Address,
				ChainEq:   "arbitrum",
				TokenEq:   "ethereum",
			})
			if err != nil {
				return fmt.Errorf("api.Balances: %w", err)
			}

			return nil
		})

		errGroup.Go(func() error {
			_, err = api.Balances(ctx, graphapi.BalancesFilter{
				AddressEq: wallets[0].Address,
				ChainEq:   "polygon",
				TokenEq:   "matic-network",
			})
			if err != nil {
				return fmt.Errorf("api.Balances: %w", err)
			}

			return nil
		})

		errGroup.Go(func() error {
			_, err = api.ListTransactions(ctx)
			if err != nil {
				return fmt.Errorf("api.ListTransactions: %w", err)
			}

			return nil
		})

		return errGroup.Wait()
	}
}

func SignUp(ctx context.Context) func() error {
	return func() error {
		api := graphapi.New(endpoint)

		username := fmt.Sprintf("load_testing_%s", uuid.NewString())
		email := fmt.Sprintf("load_testing_%s", uuid.NewString())

		flow, err := api.CreateRegistration(ctx, graphapi.CreateRegistrationInput{})
		if err != nil {
			return fmt.Errorf("api.CreateRegistration: %w", err)
		}

		flow, err = api.SetEmail(ctx, graphapi.SetEmailInput{
			ID:    flow.ID,
			Email: email,
		})
		if err != nil {
			return fmt.Errorf("api.SetEmail: %w", err)
		}

		flow, err = api.ConfirmEmail(ctx, graphapi.ConfirmEmailInput{
			ID:   flow.ID,
			Code: "111111",
		})
		if err != nil {
			return fmt.Errorf("api.ConfirmEmail: %w", err)
		}

		flow, err = api.SetCredentials(ctx, graphapi.SetCredentialsInput{
			ID:       flow.ID,
			Username: username,
			Password: "Password",
		})
		if err != nil {
			return fmt.Errorf("api.SetCredentials: %w", err)
		}

		if flow.State != "COMPLETED" {
			return errors.New("flow is not completed")
		}

		return nil
	}
}
