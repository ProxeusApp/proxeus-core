# Proxeus platform

## Payment

### Run payment tests
```
make test-payment
```

### Requirements
* The user pays for each start of a workflow (Exception: user is owner of the workflow or the workflow price is 0).
* The payment is made before the workflow is initiated (e.g. before the first form is displayed).
* A successful payment (confirmed transaction) starts the workflow. A failed transaction shows an error message on the payment page.
* Each workflow can have a different price. It is set in the workflow builder by the workflow owner and can be changed anytime.
* The default price is 0 XES. If the price is 0, the payment is skipped and the workflow starts right away. The price relevant for the workflow execution is the one that was valid at the time of transaction.
* The price is displayed where you choose the workflow to be started and on the payment page.

### Technical concept
* A Workflow-Payments is an ERC-20 Transfer from the buying user's eth-address to the selling user's eth-address
* When a payment is submitted to the blockchain the platform-backend and the platform-frontend listens to the `Transfer`-Event, that is emitted by the blockchain as soon as the payment is confirmed.
* When the backend receives the event it persists the successful payment.
* When the frontend receives the event it sends a "AddWorkflowPayment"-Request to the backend and claims the payment by sending the transaction hash of the payment and workflowId of the workflow the user wants to start. In case the request fails the frontend retries the request up to 10 times in an interval of 2 seconds.
* The backend checks for the validity of the transaction parameters. This check includes checking whether from-address, to-address and xes-amount match and that the payment has not been claimed before.
* The payment is persisted as long as the workflow has not been finished by the user.
* When the user finishes the workflow the payment is removed. When the user starts the same or any other workflow a new payment is due.

### Known Issues
Due to the distributed and shared nature of the blockchain, in case multiple Proxeus Platform instances would be running, all payments would be shared between these Proxeus Platform instances. Therefore in the rare case where the same buyer and seller (same eth-address in metamask) use various Platform instances it would be possible to pay for a workflow on one instance and use the workflow without payment on another instance. We mitigate the risk of exploiting this behaviour by checking the from-address, to-address and xes-amount in the backend of the Platform. Thanks to this measure the described issue could potentially only arise, in a scenario where the same buyer and same seller of a workflow would be registered on multiple Proxeus Platform Instances and buyer and seller would both have to be registered on 2 or more of the same Proxeus Platform Instances. In addition to that the price of the workflow would have to exactly match the price on the other instances.
