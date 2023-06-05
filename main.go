package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	// Prompt user for Infura project ID
	var infuraProjectID string
	fmt.Print("Enter your Infura project ID: ")
	fmt.Scanln(&infuraProjectID)

	// Prompt user for Uniswap V3 pool address
	var poolAddress string
	fmt.Print("Enter the Uniswap V3 pool address: ")
	fmt.Scanln(&poolAddress)

	// Configure Ethereum client
	ethClient, err := ethclient.Dial("https://mainnet.infura.io/v3/" + infuraProjectID)
	if err != nil {
		log.Fatal(err)
	}

	// Create a new instance of the Uniswap V3 pool contract
	poolContract, err := NewUniswapV3Pool(common.HexToAddress(poolAddress), ethClient)
	if err != nil {
		log.Fatal(err)
	}

	// Start monitoring the pool
	fmt.Println("Monitoring Uniswap V3 pool...")
	monitorPool(poolContract)
}

func monitorPool(poolContract *UniswapV3Pool) {
	// Create a context and cancel function for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Subscribe to relevant events
	filterQuery := ethereum.FilterQuery{
		Addresses: []common.Address{poolContract.contractAddress},
	}
	logs := make(chan types.Log)
	sub, err := poolContract.ethClient.SubscribeFilterLogs(ctx, filterQuery, logs)
	if err != nil {
		log.Fatal(err)
	}

	// Start listening for events
	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case eventLog := <-logs:
			// Handle event
			handleEvent(eventLog)
		}
	}
}

func handleEvent(eventLog types.Log) {
	// Parse the event data
	event := &SwapEvent{} // Replace with your event struct
	abiObj, err := abi.JSON(strings.NewReader(poolABI))
	if err != nil {
		log.Fatal(err)
	}

	err = abiObj.Unpack(event, "YourEvent", eventLog.Data)
	if err != nil {
		log.Fatal(err)
	}

	// Process the event data
	fmt.Println("Received event:")
	fmt.Printf("FromToken: %s\n", event.FromToken.Hex())
	fmt.Printf("ToToken: %s\n", event.ToToken.Hex())
}

// Replace with your actual Uniswap V3 pool contract ABI
const poolABI = `
	[{"constant":true,"inputs":[{"name":"tick","type":"int24"}],"name":"secondsInside","outputs":[{"name":"","type":"uint32"}],"payable":false,"stateMutability":"view","type":"function"}]
`

// Replace with your actual Uniswap V3 pool contract struct definition
type UniswapV3Pool struct {
	contractAddress common.Address
	ethClient       *ethclient.Client
}

// NewUniswapV3Pool creates a new instance of the Uniswap V3 pool contract
func NewUniswapV3Pool(address common.Address, ethClient *ethclient.Client) (*UniswapV3Pool, error) {
	return &UniswapV3Pool{
		contractAddress: address,
		ethClient:       ethClient,
	}, nil
}

// Random event/ Add your own custom event
type SwapEvent struct {
	FromToken common.Address
	ToToken   common.Address
}
