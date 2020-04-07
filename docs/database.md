# Database

Proxeus comes with two database interface implementations:
* BoltDB
* Mongodb

Other database can be integrated by implementing the [DB](https://github.com/ProxeusApp/proxeus-core/blob/master/storage/database/db/interface.go) interface.

## BoltDB

BoltDB is the default database integration.  It uses the [Storm](https://github.com/asdine/storm) toolkit to provide indexes, improved methods 
to store and fetch data, and an advanced query system.

BoltDB directly stores its data on the filesystem.

In addition of being the default database integration, BoltDB exports are also used as exchange format during export and import even when using other database 
integration like Mongodb.

Please refer to the [storm](https://github.com/ProxeusApp/proxeus-core/blob/master/storage/database/db/storm.go) integration.

## Mongodb

For larger deployment and to provide addition scalability and resiliency, Proxeus comes with a Mongodb integration.  Due to Proxeus use of transactions,
Proxeus requires a Mongodb replica set. 

Please refer to the [mongodb](https://github.com/ProxeusApp/proxeus-core/blob/master/storage/database/db/mongo.go) integration.
