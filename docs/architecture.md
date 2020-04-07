# Architecture

## High level overview

Proxeus has a simple layered architecture.  A client, either the Proxeus web application or other client 
executable access Proxeus Core through an HTTP interface to create and execute workflows.
Proxeus Core itself will use external resources such as the Ethereum blockchain, a database and external workflow 
node implementation to serve the client.

Between client requests and the external resources, we have a classical layered architecture with its API, 
Service and System layers.  The System layer is itself responsible to encapsulate the blockchain, database, 
external node interfaces and for a SPI (service programing interface).  Developers can easily create new implementations of those system service layers.

![architecture_overview](_media/architecture_overview_half.png)

### Blockchain Interface

Currently, Proxeus is available with Ethereum integration.  Ethereum can be used for user authentication, workflow payment in XES and document signature.

### Database Interfaces

Proxeus is available with two database integrations: 
* BoltDB
* Mongodb

Other database can be integrated by implementing the [DB](https://github.com/ProxeusApp/proxeus-core/blob/master/storage/database/db/interface.go) interface.

### External Nodes Interfaces

External nodes are the main extension method to add new services to Proxeus workflow.  They are implemented as HTTP servers implementing a simple generic 
API.  
