# Configuration

When starting the server for the first time, you will need to configure it.

Open http://localhost:1323/init to configure your server.  Use the
configuration below as an example.

## Document Service URL
```
https://dev-ds.proxeus.com/
```

## Infura API Key
Generate a  [infura API Key](https://infura.io)

## Blockchain contract address (goerli)
```
0x66FF4FBF80D4a3C85a54974446309a2858221689
```
(alternatively deploy own smart contract from ProxeusFS.sol)

## Email from
```
no-reply@proxeus.com
```
## Sparkpost API Key

Set up a free account on [SparkPost](https://www.sparkpost.com)

## Initial Email

Use your email address for the root user and choose a secure password.
```
youremail@address.com
```

## Encryption Secret Key

This is a salt to hash your user's passwords in the database. You can use any value, with preference to hard generated strings. Make sure that it is exactly 32 characters long. Do not change the key on a running instance. This can only be set using an environment variable.

# Full Configuration

You can get the full list of configuration parameters using the `-h` parameter:

```
./artifacts/proxeus -h
Usage of ./artifacts/proxeus:
  -AirdropAmountEther string
    	Amount of Ether to airdrop to newly registered users. (PROXEUS_AIRDROP_AMOUNT_ETHER) (default "0")
  -AirdropAmountXES string
    	Amount of XES to airdrop to newly registered users. (PROXEUS_AIRDROP_AMOUNT_XES) (default "0")
  -AirdropEnabled string
    	Enables/Disables the XES & Ether airdrop feature on test network. (PROXEUS_AIRDROP_ENABLED) (default "false")
  -AirdropWalletfile string
    	Path to File containing Private Key of the Wallet to fund Airdrops of XES and Ether. (PROXEUS_AIRDROP_WALLETFILE)
  -AirdropWalletkey string
    	Path to File containing the Key for the Airdrop Private Key. (PROXEUS_AIRDROP_WALLETKEY)
  -AllowHttp string
    	Allow the use of HTTP =NOT FOR PRODUCTION=. (PROXEUS_ALLOW_HTTP) (default "false")
  -AutoTLS
    	Automatically generate Let's Encrypt certificate (Server must be reachable on port 443 from public internet). (PROXEUS_AUTO_TLS)
  -BlockchainContractAddress string
    	Ethereum contract address which will be used to register files and verify them. (PROXEUS_BLOCKCHAIN_CONTRACT_ADDRESS)
  -BlockchainNet string
    	Ethereum blockchain net like mainnet/goerli/polygon-mumbai/polygon-mainnet. (PROXEUS_BLOCKCHAIN_NET) (default "goerli")
  -CacheExpiry string
    	Common cache expiry which will be used for email tokens or similar. (PROXEUS_CACHE_EXPIRY) (default "24h")
  -DataDir string
    	Database directory path. All data will be stored here. (PROXEUS_DATA_DIR) (default "/data/hosted")
  -DatabaseEngine string
    	Selects database engine, supported values: storm, mongo. (PROXEUS_DATABASE_ENGINE) (default "storm")
  -DatabaseURI string
    	Sets database connection string, not required for embedded databases. (PROXEUS_DATABASE_URI)
  -DefaultRole string
    	Default role that is going to be used for new registrations. Value is case insensitive. (PROXEUS_DEFAULT_ROLE) (default "creator")
  -DefaultWorkflowIds string
    	Workflow IDs to set to clone and add to a new user (PROXEUS_DEFAULT_WORKFLOW_IDS)
  -DocumentServiceUrl string
    	Document Service URL which will be used to render documents. (PROXEUS_DOCUMENT_SERVICE_URL) (default "http://document-service:2115/")
  -EmailFrom string
    	Email that is being used to send out emails. (PROXEUS_EMAIL_FROM)
  -EthClientURL string
    	Ethereum client URL (PROXEUS_ETH_CLIENT_URL) (default "https://goerli.infura.io/v3/")
  -EthWebSocketURL string
    	Ethereum websocket URL (PROXEUS_ETH_WEB_SOCKET_URL) (default "wss://goerli.infura.io/ws/v3/")
  -InfuraApiKey string
    	API Key to access Infura node. (PROXEUS_INFURA_API_KEY)
  -LogPath string
    	Location of the log file of this service. (PROXEUS_LOG_PATH) (default "./log")
  -PlatformDomain string
    	Platform Domain used to for links to Platform (PROXEUS_PLATFORM_DOMAIN)
  -ServiceAddress string
    	address and port of this service (PROXEUS_SERVICE_ADDRESS) (default ":1323")
  -SessionExpiry string
    	Session expiry like 1h as one hour, 1m as one minute or 1s as one second. (PROXEUS_SESSION_EXPIRY) (default "1h")
  -SettingsFile string
    	Path to the settings file (PROXEUS_SETTINGS_FILE) (default "~/.proxeus/settings/main.json")
  -SparkpostApiKey string
    	Sparkpost API key which will be used to send out emails. (PROXEUS_SPARKPOST_API_KEY)
  -TestMode string
    	Run the server in test mode =NOT FOR PRODUCTION=. (PROXEUS_TEST_MODE) (default "false")
  -XESContractAddress string
    	 (PROXEUS_XESCONTRACT_ADDRESS) (default "0x15FeA089CC48B4f4596242c138156e3B53579B37")
```
