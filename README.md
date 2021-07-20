# ETHcc web3.go sample project

## Download solc 
`wget https://github.com/ethereum/solidity/releases/download/v0.8.6/solc-static-linux`
`chmod +x solc-static-linux`

## Generate the contract bindings
`go generate`

## Build the project
`go build`

## Real backend
For the real backend you need to run a geth node in dev mode.

### Dev mode
Start the geth node in dev mode:
`geth --dev --http console`

### Fund the deployer for the real backend
Execute the following in the console:
`eth.sendTransaction({from:personal.listAccounts[0], to:"0xb02A2EdA1b317FBd16760128836B0Ac59B560e9D", value: "100000000000000"})`

