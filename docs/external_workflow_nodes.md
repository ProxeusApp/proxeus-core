# External Workflow Nodes


External workflow nodes are the primary method to extend Proxeus workflow to any use case.
They are implemented as external HTTP servers that interact with the Proxeus Core using a
simple API.

![external_node](_media/external_node_half.png)

During workflow execution, when the state transitions to an external node, the following steps are executed:
* Proxeus Core sends the full workflow state the external node,
* the external node reads and updates the workflow state, and returns it to the core,
* Proxeus Core updates the workflow state. 

**Very Important note: Since the full workflow state is sent to the external node, Proxeus operator must 
use great care as this could leak data.  External node implementation must only be 
implemented using trusted node on a trusted network.  In practice, operator must understand the 
functionality of the external node and only used them in their own docker compose deployment.
Never use external nodes on the internet.**

## API

### External node API

External nodes must implement the external node API:

| Method | Endpoint         | Body | Description |
|--------|------------------|------|-------------|
| POST   | /node/:id/next   | JSON Object | This will receive the data send by Proxeus, will process it and return new data |
| GET    | /node/:id/config | HTML | Return the node configuration page |
| POST   | /node/:id/config | JSON Object | Change an eventual node configuration |
| POST   | /node/:id/remove | - | Callback whenever the node is removed from Proxeus' side |
| POST   | /node/:id/close  | - | Callback whenever the node is closed from Proxeus' side  |
| GET    | /health          | JSON Object | Proxeus will ping this healthcheck to know whether the node is running. A simple 200 should be returned" |

The JSON data exchanged during a Post to next or when reading/updating the config is a generic JSON Object.

### Proxeus Core API

To be available inside Proxeus, an external node must first register to the running Proxeus service using the `/api/admin/external/register`
endpoint. In addition, the core provides a configuration store interface to simplify external nodes implementation even if they need 
some configuration.  The configuration interface acts a memento: the configuration is not interpreted and returned verbatim.

| Method | Endpoint         | Body | Description |
|--------|------------------|------|--------------|
| POST   | /api/admin/external/register   | JSON | Register a new external node to be made available in Proxeus |
| POST   | /api/admin/external/config/:id   | JSON | Store the JSON data for the provided id |
| GET   | /api/admin/external/config/:id   | JSON | Return the JSON data for the provided id |

#### Registration

During registration, the external node must provide the following data:

```json
{
    "id": "any id",
    "name": "Node name",
    "detail": "Description",
    "url": "http://YourNodeUrl.com",
    "secret": "secret key"
}
```

| Property | Description |
|--------|-------------|
| id | a unique id for the external node |
| name | user friendly node name |
| detail | description of the functionality of the node |
| url | the URL of the external node |
| secret | a shared secret |

## Go Library

We provide a Golang library with helper functions here: http://github.com/ProxeusApp/node-go

## Examples

Please refer the Proxeus node [docker-compose.yml](https://github.com/ProxeusApp/proxeus-core/blob/master/docker-compose.yml) file 
for an example on how to configure an instance of the Proxeus Core with a number of external nodes.

The implementation of the external nodes used can be found here:

* [Balance Retriever](https://github.com/ProxeusApp/node-balance-retriever) retrieves ETH + ERC20 tokens balances
* [Crypto Forex Rates](https://github.com/ProxeusApp/node-crypto-forex-rates) retrieves USD prices of given tokens
* [Proof of Existence](https://github.com/ProxeusApp/node-proof-of-existence) proves existence of a tweet
* [Mail sender](https://github.com/ProxeusApp/node-mail-sender) sends emails
* [JSON sender](https://github.com/ProxeusApp/node-json-sender) Sends form data to a REST endpoint via POST request

