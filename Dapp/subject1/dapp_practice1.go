package main

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	const infuraURL = "https://sepolia.infura.io/v3/d3d4c98c0ebd4077b98b7436fa59a243"
	blockNumber := big.NewInt(5412345)

	client, err := ethclient.Dial(infuraURL)
	if err != nil {
		log.Fatalf("连接失败: %v", err)
	}
	defer client.Close()

	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Fatalf("区块查询失败: %v", err)
	}

	fmt.Println("区块查询成功:")
	fmt.Printf("区块高度: %d\n", block.Number().Uint64())
	fmt.Printf("区块哈希: %s\n", block.Hash().Hex())
	fmt.Printf("时间戳: %d\n", block.Time())
	fmt.Printf("交易数量: %d\n", len(block.Transactions()))
	fmt.Printf("矿工地址: %s\n", block.Coinbase().Hex())
	fmt.Printf("区块难度: %d\n", block.Difficulty().Uint64())
}
