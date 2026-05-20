package evm

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"strings"
)

type EVM struct {
	client *ethclient.Client
}

func New(ctx context.Context, wsRpcURL string) (*EVM, error) {
	client, err := ethclient.DialContext(ctx, wsRpcURL)
	if err != nil {
		return nil, fmt.Errorf("ethclient.Dial: %v", err)
	}

	return &EVM{
		client: client,
	}, nil
}

func (e *EVM) Transfer(
	ctx context.Context,
	privateKeyStr string,
	to string,
	amount int64,
) (string, error) {
	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(privateKeyStr, "0x"))
	if err != nil {
		return "", fmt.Errorf("crypto.HexToECDSA: %w", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", errors.New("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := e.client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		return "", fmt.Errorf("client.PendingNonceAt: %w", err)
	}

	gasLimit := uint64(21000)
	gasPrice, err := e.client.SuggestGasPrice(ctx)
	if err != nil {
		return "", fmt.Errorf("client.SuggestGasPrice: %w", err)
	}

	toAddress := common.HexToAddress(to)
	_amount := big.NewInt(amount)

	transactionFee := new(big.Int).Mul(gasPrice, new(big.Int).SetUint64(gasLimit))
	_amount = _amount.Sub(_amount, transactionFee)

	var data []byte

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &toAddress,
		Value:    _amount,
		Gas:      gasLimit,
		GasPrice: gasPrice,
		Data:     data,
	})

	chainID, err := e.client.NetworkID(ctx)
	if err != nil {
		return "", fmt.Errorf("client.NetworkID: %w", err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return "", fmt.Errorf("types.SignTx: %w", err)
	}

	err = e.client.SendTransaction(ctx, signedTx)
	if err != nil {
		return "", fmt.Errorf("client.SendTransaction: %w", err)
	}

	return signedTx.Hash().String(), nil
}
