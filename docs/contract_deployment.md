# Smart Contract Deployment

> This section describes the necessary steps to deploy the ProxeusFS using Remix & Metamask.

The latest ProxeusFS contract can be found in github: https://github.com/ProxeusApp/proxeus-contract

## Before Deployment
 * Ensure that the documentation in code and Readme.md is updated
 * Ensure that the tests are updated to provide 100% coverage. 
 * Ensure that the tests report no errors. See [README.md](https://github.com/ProxeusApp/proxeus-contract/blob/master/README.md) on how to run tests.

## Deployment procedure

### First time setup
 1. Install the Metamask extension in your browser: https://metamask.io/
 2. Create or Import an Ethereum Account in Metamask

### Deployment
 1. Switch to goerli/mainnet/polygon-mumbai/polygon-mainnet in Metamask depending on where you would like to deploy the contract
 2. Make sure your Ethereum Account has some Ether to fund the transactions (Goerli faucet: https://goerlifaucet.com/)
 3. Open Remix: https://remix.ethereum.org
 4. Upload the ProxeusFS.sol file
 5. Open the "Plugin Manager" (Sidebar) and enable Plugins "Solidity Compiler" & "Deploy & Run Transactions"
 6. Under "Solidity Compiler", select a compatible "Compiler" matching the pragma definition of the code. Example: `pragma solidity 0.5.1;` -> `0.5.10+commit.5a6ea5b1`
 7. Compile the contract
 8. Note that if you made changes to the function signatures you should copy & store the ABI (available after compilation) at this point for later inclusion in the core project to ensure compatibility between the smart contract and the golang code.
 9. Under "Deploy & Run Transactions", select "Injected Web3" as "Environment" to link Remix with Metamask
 10. Ensure that you have selected the right account in both Metamask and as "Account" in Remix
 11. Ensure that you have set the "Gas Limit" high enough depending on the size of the contract code.
 12. Select the correct Contract from the Dropdown and "Deploy"
 13. Confirm the Transaction in Metamask & wait for it to be mined.
 14. Once the deployment transaction is mined, the contract will be visible under "Deployed Contract"
 15. Copy & Store the newly deployed contract's address
