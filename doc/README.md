This folder is only for the license loader of the project.

Go to [docs/README](../docs/README.md) for user and developer instructions.

# How to generate the ACKNOWLEDGEMENTS file

The acknowledgement file is constructed automatically using the `create_acknowledgement.sh` bash script.

## Dependencies

Before running this script, you will need to install the following dependencies:

* hub command line tool from https://hub.github.com
* yq YAML command line tool from https://github.com/kislyuk/yq
* jq JSON command line tool from https://stedolan.github.io/jq
* curl command line tool from https://curl.haxx.se

## Github authentication

Run the following command to authenticate with Github:
```
hub api user
```

## Generate the file

Run the following command to generate the acknowledgement file:
```
./create_acknowledgment.sh > ACKNOWLEDGEMENTS
```
