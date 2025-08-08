package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	client, err := ethclient.Dial("https://sepolia.infura.io/v3/d3d4c98c0ebd4077b98b7436fa59a243")
	if err != nil {
		log.Fatal(err)
	}

	privateKey, err := crypto.HexToECDSA("MY_PRIVATE_KEY")
	if err != nil {
		log.Fatal(err)
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("无法获取公钥")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	toAddress := common.HexToAddress("0x9823480714b0a06A5a09F23S825EaCfc32S70b3c")

	value := big.NewInt(1000000000000000)

	gasLimit := uint64(21000)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	tx := types.NewTransaction(
		nonce,
		toAddress,
		value,
		gasLimit,
		gasPrice,
		nil,
	)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("交易已发送! 交易哈希: %s\n", signedTx.Hash().Hex())
}
