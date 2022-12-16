package model

import "reflect"

type Settings struct {
	DocumentServiceUrl        string `json:"documentServiceUrl" validate:"required=true" default:"http://document-service:2115/" usage:"Document Service URL which will be used to render documents."`
	PlatformDomain            string `json:"platformDomain" validate:"required=true" default:"" usage:"Platform Domain used to for links to Platform"`
	DataDir                   string `json:"dataDir" default:"/data/hosted" usage:"Database directory path. All data will be stored here."`
	DefaultRole               string `json:"defaultRole" default:"creator" usage:"Default role that is going to be used for new registrations. Value is case insensitive."`
	SessionExpiry             string `json:"sessionExpiry" validate:"required=true" default:"1h" usage:"Session expiry like 1h as one hour, 1m as one minute or 1s as one second."`
	CacheExpiry               string `json:"cacheExpiry" validate:"required=true" default:"24h" usage:"Common cache expiry which will be used for email tokens or similar."`
	BlockchainNet             string `json:"blockchainNet" validate:"required=true" default:"goerli" usage:"Ethereum blockchain net like mainnet or goerli."`
	InfuraApiKey              string `json:"infuraApiKey" validate:"required=true" usage:"API Key to access Infura node."`
	BlockchainContractAddress string `json:"blockchainContractAddress" validate:"required=true" default:"" usage:"Ethereum contract address which will be used to register files and verify them."`
	SparkpostApiKey           string `json:"sparkpostApiKey" validate:"required=true" usage:"Sparkpost API key which will be used to send out emails."`
	EmailFrom                 string `json:"emailFrom" validate:"required=true,email=true" usage:"Email that is being used to send out emails."`
	LogPath                   string `json:"logPath" default:"./log" usage:"Location of the log file of this service."`
	DefaultWorkflowIds        string `json:"defaultWorkflowIds" usage:"Workflow IDs to set to clone and add to a new user"`
	AirdropEnabled            string `json:"airdropEnabled" validate:"required=true" default:"false" usage:"Enables/Disables the XES & Ether airdrop feature."`
	AirdropAmountXES          string `json:"airdropAmountXES" default:"0" usage:"Amount of XES to airdrop to newly registered users."`
	AirdropAmountEther        string `json:"airdropAmountEther" default:"0" usage:"Amount of Ether to airdrop to newly registered users."`
	DatabaseEngine            string `json:"databaseEngine" default:"storm" usage:"Selects database engine, supported values: storm, mongo."`
	DatabaseURI               string `json:"DatabaseURI" default:"" usage:"Sets database connection string, not required for embedded databases."`
	TestMode                  string `json:"testMode" default:"false" usage:"Run the server in test mode =NOT FOR PRODUCTION=."`
	AllowHttp                 string `json:"allowHttp" default:"false" usage:"Allow the use of HTTP =NOT FOR PRODUCTION=."`
}

func NewDefaultSettings() *Settings {
	stngs := &Settings{}
	t := reflect.TypeOf(Settings{})
	v := reflect.ValueOf(stngs)
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.Type.Kind() == reflect.String {
			v.Elem().Field(i).SetString(f.Tag.Get("default"))
		}
	}
	return stngs
}
