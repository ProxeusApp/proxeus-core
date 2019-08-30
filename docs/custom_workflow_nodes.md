# Custom Workflow Nodes

Custom workflow nodes are the primary method to extend Proxeus workflow to any use cases.

They are written in Golang and can
* read the workflow data, 
* execute any code and
* updarte the workflow data.

Every node implement the following interface:

```
NodeIF interface {
		//Execute is being called when the node becomes the current target, doesn't matter whether it goes forward or backward.
		//If the node was executed before and returned proceed = false, the same instance is being used again.
		Execute(node *Node) (proceed bool, err error)
		//Remove is being called when this node is not part of the path anymore. In other words, when it doesn't exist in the State
		Remove(node *Node)
		//Close will be called on the instance before the next node is being executed or when the end is reached
		//Close is always the last called method either Execute()* -> Close() OR Execute()* -> Remove() -> Close().
		//Execute() can be called n times before Close() or Remove() as it is controlled by the node impl.
		//When it is called, the instance is released from the engine.
		Close()
	}
```



## Workflow Context

Each workflow node must have access to the `DocumentFlowInstance` context which gives full access to workflow data.

## Examples
You can find implementation example of this interface in the Proxeus repository under the [proxeus-core/main/app](https://github.com/ProxeusApp/proxeus-core/tree/master/main/app)
directory:

* [mail_sender.go](https://github.com/ProxeusApp/proxeus-core/tree/master/main/app/mail_sender.go) shows how to read data from the workflow,
* [price_retriever.go](https://github.com/ProxeusApp/proxeus-core/tree/master/main/app/price_retriever.go) shows how to update workflow data.






