# Application Examples

Proxeus plays the role of a bridge connecting two worlds, allowing you to digitize traditional processes and register information on a blockchain, as well as to make blockchain activity visible by generating human-readable documents.

Crypto Asset Reporting
======================

![](_media/old_proxeus/crypto_tax_reporting/1.svg)


### Abstract

This practical example shows how you can use Proxeus and its extension capabilities to build a dApp that generates a crypto asset report for any Ethereum wallet. It showcases Proxeus’ ability to access data sources within the Ethereum ecosystem as well as its effectiveness in generating standardized documents.


By creating a simple custom workflow node, enhanced a Proxeus workflow with the capability to access external data. This data is then used in a workflow that can be published and offered to other users for a fee using the XES token.


### The idea
Whereas a bank can provide you with a list of your accounts and assets, most cryptocurrency exchanges and wallets do not offer a comparable service yet. Systematically keeping track is a hassle. When it comes to the yearly tax declaration, token holders have to unlock their wallets, export or even write down all holdings and calculate the end-of-year value of each position. This requires knowhow and manual labor due to the absence of standardized processes and tools. Since the assets are managed on a public ledger and can be accessed fairly easily, it should be possible to automate some steps.


### The solution

All Ethereum-based tokens present on your Ethereum wallets can be queried using existing tools such as Etherscan. Combining this data with the end-of-year valuation of each token produces the list most tax authorities ask for: the name of the asset, the amount held at the end of the year and the value of the position held. It could even be inserted directly into the respective forms of the tax software if an interface is provided.

#### A crypto asset reporting workflow in Proxeus

The workflow should have elements that can do the following things:
‍
1) Ask the user to specify the wallet for which the report shall be created.
2) Retrieve the crypto assets of the selected wallet.
3) Look up the valuation of each asset.
4) Create a PDF report with this information, optionally using a template provided by the user’s tax authority

For steps 1) and 4) we can simply use out-of-the-box functionality of Proxeus. For points 2) and 3) we can implement custom nodes using Proxeus’ external nodes interface. We’ll explain the details a bit further down in this article.



#### About the external nodes interface

The external node interface allows anybody to develop custom-made workflow nodes that can be added without having to modify the Proxeus software. These custom nodes register themselves with the external nodes interface and are then automatically made available for use in workflows by Proxeus Core. Community members can exchange custom nodes simply by making their code or Docker images available. As the nodes can read and write workflow data and are unrestricted in their communication with external systems, they unlock enormous potential and an easy path for extending the core.

The Proxeus Golang library for external nodes is documented in its [Github Repostory](https://github.com/ProxeusApp/node-go.git)

## Check out our demo Workflow

The complete crypto asset report workflow is available on our demo [platform](https://morrison.proxeus.org/). It will be copied to your account automatically when you sign up.

You can start by just running the workflow: User View > Documents > “New document”.

In the admin panel you can analyze how we configured the workflow and as you have your own copy, you are free to make any changes to it.

## How to rebuilt it yourself

If this is your first contact with Proxeus, we recommend you familiarize yourself with the platform first before trying to tackle this example. The “Demo” section on the website provides you with a quick overview.

#### 1) Create a user form

It should inform the user what needs to be done (enter a wallet address) and collect the required input; in this case the wallet address. In our example we used a title element and a simple text input. In the design view on the left it looks like this:

![](_media/old_proxeus/crypto_tax_reporting/2.png)

A good form validates user input to prevent user mistakes. With regular expression we can check if the input is a valid Ethereum public address: ^0x[a-fA-F0-9]{40}$

The components configuration on the right side of the form builder will then look as follows:

![](_media/old_proxeus/crypto_tax_reporting/3.png)

The name of the input component, ethAddress, is important. It is the variable name to access that input later in the workflow.

#### 2) Retrieve balances

Now we can retrieve the ETH and ERC-20 balances of the input Ethereum address using a custom workflow node. Our example implementation can be found on GitHub: [node-balance-retriever](https://github.com/ProxeusApp/node-balance-retriever.git).

If you wish to build it yourself, we recommend studying the documentation of the [external nodes interface](https://github.com/ProxeusApp/node-go.git). For Golang there is an official implementation of Ethereum: [go-ethereum](https://github.com/ethereum/go-ethereum.git). You could run your own Ethereum node for quick access to the data - or use a service like Infura.

#### 3) Determine asset valuation

After finding out what crypto assets the user holds, the next step is to determine the value of each position. Price aggregators such as CoinMarketCap, Livecoinwatch or CryptoCompare are usually a reliable source of information as they do not rely on a single exchange only. For your example we’ve used the [API of CryptoCompare](https://min-api.cryptocompare.com/).

Our example implementation can be found on GitHub: [node-crypto-forex-rates](https://github.com/ProxeusApp/node-crypto-forex-rates.git)


#### 4) Design a report template

Now you can start designing the template for your report. Create a layout to your liking and add the placeholders to be replaced by the Proxeus document service. The amount of Ether in the wallet, for example, will be available through the variable name “ETH” while its price is “USD_ETH”. You would add it to the template as

```{{input.ETH}} and {{input.USD_ETH}}```

In our example you’ll also have XES, MKR, BAT, USDC, REP and OMG at your disposal.

To display the total value of a position, you’ll have to multiply the two values amount and price:


```{{input.ETH * input.USD_ETH}}.```

The sum of all tokens could be displayed as follows:
The ERC tokens amount to a total of USD 

```{{((input.XES * input.USD_XES)+(input.MKR * input.USD_MKR)+(input.BAT * input.USD_BAT)+(input.USDC * input.USD_USDC)+(input.REP * input.USD_REP)+(input.OMG * input.USD_OMG))}}.```

You can download our example template here.

#### 5) Complete the workflow design

Now all the elements of the workflow are ready and just have to be added to a workflow and put in the right order. The template comes first, followed by the form, which asks the user for the wallet address. Then the first custom node retrieves the balances, for which the second module then fetches the valuation. Lastly the data is inserted into the template and rendered as a PDF.


External nodes can have a configuration UI that can be reached via the workflow editor and a double-click on the node’s symbol. The “Crypto to Fiat Forex Rates” node, one of our examples, has the following UI:


The UI is defined in the node’s code as HTML:


```
const configHTML = `
<!DOCTYPE html>
<html>
<body>
<form action="/node/{{.Id}}/config?auth={{.AuthToken}}" method="post">
Convert to Fiat currency: <input type="text" size="2" name="FiatCurrency" value="{{.FiatCurrency}}">
<input type="submit" value="Submit">
</form>
</body>
</html>
`
```

#### 6) Publish and monetize your workflow (optional)

If you have the knowledge how to compile a good report on crypto assets, this may be valuable to other members of the crypto community. Why not offer it to them for a fee in XES? In the workflow editor, you can set the price. Save, hit “Publish” and start advertising! All fees will automatically be paid to the wallet that you have connected in your user profile on the platform.


## Limitations

1) To ensure that this example can be reproduced with reasonable effort, we have opted for two simplifications.
We retrieve the assets held at the time of execution. This means that the report has to be created at the end of the year or before making trades in the new year. Retrieving the amount held at Dec 31 at midnight would require detailed historical data - not just the block headers - which is possible only through a full archive node of Ethereum. Most users do not have access to such a node as it is usually a premium service by providers like Infura.

    If you have access to a full archive node, you can retrieve the balance as follows using the [Ethereum JavaScript API:](https://github.com/web3/web3.js.git)

```eth.getBalance("<ADDRESS_HERE>", <BLOCK_NUMBER>);```
‍
‍
For ERC-20, you need the token’s contract address (as retrieved in our example) and the contract’s ABI. The ABI is instantiated as follows:

```
var tokenContract = eth.contract([{
"type":"function",
"name":"balanceOf",
"constant":true,
"payable":false,
"Inputs":[{"name":"","type":"address"}],
"outputs":[{"name":"","type":"uint256","value":"0"}] 
}]);‍ 
```


‍
The balance can then be retrieved like this:

```
> var erc20ContractAddress = "<ADDRESS_OF_TOKEN'S_CONTRACT>";
> var account = "YOUR_ADDRESS";
> tokenContract.at(erc20ContractAddress).balanceOf(account);
```

2) We fetch the price of each asset at the time of execution. Ideally the report would perform a lookup on the official list provided by the tax authority. This differs from country to country and can usually be looked up on the website of the tax authority.
We hope you’ve enjoyed this example. If you have any questions, the Proxeus community is always happy to help. When looking for help, please make sure to provide a good description of your problem and what you’re trying to achieve.

