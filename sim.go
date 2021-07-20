package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
)

func deployToSimBackend() {
	backend, sk := getSimBackend()
	deployer := crypto.PubkeyToAddress(sk.PublicKey)
	transactor, err := bind.NewKeyedTransactorWithChainID(sk, big.NewInt(1337))
	if err != nil {
		panic(err)
	}
	supply := big.NewInt(1000)
	_, tx, contract, err := DeploySmallContract(transactor, backend, supply)
	if err != nil {
		panic(err)
	}
	backend.Commit()
	addr, err := bind.WaitDeployed(context.Background(), backend, tx)
	if err != nil {
		panic(err)
	}
	fmt.Println("Contract deployed")
	_ = addr
	// Interact with the smart contract
	// Retrieve the symbol (a free data retrieval call)
	sym, err := contract.Symbol(nil)
	fmt.Printf("Symbol: %v\n", sym)
	// Get the balance of the deployer
	bal, err := contract.Balance(nil, deployer)
	fmt.Printf("Balance deployer: %v\n", bal)
	// Transfer some tokens to 1234
	val := big.NewInt(100)
	recipient := common.BigToAddress(big.NewInt(1234))
	_, err = contract.Transfer(transactor, recipient, val)
	backend.Commit()
	fmt.Println("Funds transferred")
	// Get the balance after the transfer
	bal, err = contract.Balance(nil, deployer)
	fmt.Printf("Balance deployer: %v\n", bal)
	bal, err = contract.Balance(nil, recipient)
	fmt.Printf("Balance recipient: %v\n", bal)

	// Test the event system
	testEventFilter(contract, deployer)
	subscribeFilterEvent(backend, transactor, contract, deployer)
	testEventFilter(contract, deployer)
}

func testEventFilter(contract *SmallContract, deployer common.Address) {
	iterator, err := contract.FilterEvent(nil, []common.Address{deployer}, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("Querying old events")
	for iterator.Next() {
		event := iterator.Event
		if event != nil {
			fmt.Printf("Found event: transfer %v tokens from %v to %v\n", event.Tokens, event.From, event.To)
		}
	}
}

func subscribeFilterEvent(backend *backends.SimulatedBackend, transactor *bind.TransactOpts, contract *SmallContract, deployer common.Address) {
	// Set up the event subscription
	eventChan := make(chan *SmallContractEvent)
	sub, err := contract.WatchEvent(nil, eventChan, []common.Address{deployer}, nil)
	if err != nil {
		panic(err)
	}
	defer sub.Unsubscribe()
	fmt.Println("Set up event subscription")
	// Send transaction
	val := big.NewInt(200)
	recipient := common.BigToAddress(big.NewInt(2047))
	_, err = contract.Transfer(transactor, recipient, val)
	backend.Commit()
	// Check the event channel for the event
	select {
	case event := <-eventChan:
		if event != nil {
			fmt.Printf("Found event: transfer %v tokens from %v to %v\n", event.Tokens, event.From, event.To)
		}
	default:
		fmt.Println("No event found")
	}
}

func getSimBackend() (*backends.SimulatedBackend, *ecdsa.PrivateKey) {
	sk, err := crypto.GenerateKey()
	if err != nil {
		panic(err)
	}
	faucetAddr := crypto.PubkeyToAddress(sk.PublicKey)
	addr := map[common.Address]core.GenesisAccount{
		faucetAddr: {Balance: new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 256), big.NewInt(9))},
	}
	alloc := core.GenesisAlloc(addr)
	return backends.NewSimulatedBackend(alloc, 80000000), sk
}
