package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func deployToRealBackend() {
	fmt.Println("Real backend")
	backend, sk := getRealBackend()
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
	tx, err = contract.Transfer(transactor, recipient, val)
	reciept, err := bind.WaitMined(context.Background(), backend, tx)
	if reciept.Status == types.ReceiptStatusSuccessful {
		fmt.Printf("Transfer successful")
	}
	testEthclient(backend, sk)
}

func testEthclient(backend *ethclient.Client, sk *ecdsa.PrivateKey) {
	blockNr, _ := backend.BlockNumber(context.Background())
	fmt.Printf("BlockNr: %v\n", blockNr)
	balance, _ := backend.BalanceAt(context.Background(), crypto.PubkeyToAddress(sk.PublicKey), nil)
	fmt.Printf("Balance of faucet account: %v\n", balance)
}

func getRealBackend() (*ethclient.Client, *ecdsa.PrivateKey) {
	// eth.sendTransaction({from:personal.listAccounts[0], to:"0xb02A2EdA1b317FBd16760128836B0Ac59B560e9D", value: "100000000000000"})

	sk := crypto.ToECDSAUnsafe(common.FromHex(SK))
	if crypto.PubkeyToAddress(sk.PublicKey).Hex() != ADDR {
		panic(fmt.Sprintf("wrong address want %s got %s", crypto.PubkeyToAddress(sk.PublicKey).Hex(), ADDR))
	}
	cl, err := ethclient.Dial("http://127.0.0.1:8545")
	if err != nil {
		panic(err)
	}
	return cl, sk
}
