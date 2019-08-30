## Smart contract Golang Binding
When changes are made to the smart contracts function signatures it is necessary to regenerate the Golang bindings so that the core projects stays compatible.
For more info check the [Smart Contract Readme](https://github.com/ProxeusApp/proxeus-contract).

## To generate a contract.go from an ABI
 - Go to https://geth.ethereum.org/downloads/ and download the most recent Geth & Tools
 - Copy "abigen" to $GOPATH/bin
 - Run `$GOPATH/bin/abigen --abi [ABI FILE] --pkg wallet --type [name of contract] --out [output file]`
Example:
```
$GOPATH/bin/abigen --abi ~/Projects/Proxeus/wallet-blockchain/DocumentRegistry.abi --pkg wallet --type DocumentRegistry --out ~/Projects/Proxeus/wallet-blockchain/documentRegistry.go
```
