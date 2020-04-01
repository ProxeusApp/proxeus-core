# Custom Workflow Nodes


Custom workflow nodes are the primary method to extend Proxeus workflow to any use case.

They must be running as an HTTP API, expose the following endpoints


| Method | Endpoint         | Description                                                                                   |
|--------|------------------|-----------------------------------------------------------------------------------------------|
| POST   | /node/:id/next   | This will receive the data send by Proxeus, will process it and return new data |
| GET    | /node/:id/config | Read an eventual node configuration                                                                                              |
| POST   | /node/:id/config | Change an eventual node configuration                                                                                              |
| POST   | /node/:id/remove | Callback whenever the node is removed from Proxeus' side                                                                                              |
| POST   | /node/:id/close  | Callback whenever the node is closed from Proxeus' side                                                                                             |
| GET    | /health          | Proxeus will ping this healthcheck to know whether the node is running. A simple 200 should be returned"               |

and register to Proxeus' via HTTP.

A `POST` to `PROXEUS_URL/api/admin/external/register` sending a JSON in following format should be made to let Proxeus know there's a new node.

```json
{
    "id": "any id",
    "name": "Node name",
    "detail": "Description",
    "url": "http://YourNodeUrl.com",
    "secret": "secret key"
}
```

To simplify this we provide (Golang only) some functions.

### Register node

Following code will register the node and retry a few times on failure.
```gotemplate
externalnode.Register(proxeusUrl, serviceName, serviceUrl, jwtSecret, description)
```

For other util functions refer to `github.com/ProxeusApp/proxeus-core/externalnode` package


## Examples

You can find some implementation examples here:

* [Balance Retriever](https://github.com/ProxeusApp/node-balance-retriever) retrieves ETH + ERC20 tokens balances
* [Crypto Forex Rates](https://github.com/ProxeusApp/node-crypto-forex-rates) retrieves USD prices of given tokens
* [Proof of Existence](https://github.com/ProxeusApp/node-proof-of-existence) proves existence of a tweet
* [Mail sender](https://github.com/ProxeusApp/node-mail-sender) sends emails