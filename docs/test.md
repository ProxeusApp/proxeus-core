# Tests


The following command runs all the tests without test coverage:

```
make init server test test-integration test-api test-ui
```

The following commands are used to get the test coverage:

```
make init server test test-integration test-api coverage=true
make coverage
cat artifacts/coverage.txt
```

You can access the HTML version of the coverage with by opening the file 
`artifacts/coverage.html` in your browser.
