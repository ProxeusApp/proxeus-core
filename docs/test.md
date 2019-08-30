# Tests

## Update mocks
install the gomock package and the mockgen tool:
```
go install github.com/golang/mock/mockgen
```

Generate Mocks
(mocks have to be same package and path as original class, else there are issues with imports)
```
mockgen -package storm -source sys/db/storm/workflow_payments.go --destination sys/db/storm/workflow_payments_mock.go
mockgen -package storm -source sys/db/storm/user.go --destination sys/db/storm/user_mock.go
mockgen -package storm -source sys/db/storm/workflow.go --destination sys/db/storm/workflow_mock.go
mockgen -package blockchain -source main/handlers/blockchain/adapter.go --destination  main/handlers/blockchain/adapter_mock.go
```


## Run all the tests

```
make test
```


