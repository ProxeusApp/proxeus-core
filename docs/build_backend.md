# Back End Build

## Building

Please read the [Build All](build_all.md) section before continuing.

`make init`

Should you get 

After that, you're free to run 

`make server`

to build the application.

## Running
An executable will be created under `./artifacts/`. 
To run it, simply run `./artifacts/proxeus`

That should start Proxeus' server and server it on port `:1323/`
The first time, access http://localhost:1323/init and configure your settings.


The document-service dependency is needed as well for some operations. Please refer to the documentation on its repository. 

## Tests
Unit & integration tests can be run with the command

`make test`

## Code formatting & Validation

`goimports` is used to fix/sort Go import lines and format code. 
A validation of the code can be requested with `make validate`

Follow [Building Front End](frontend.md) for the web client.
