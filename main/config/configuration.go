package config

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"flag"

	"github.com/ProxeusApp/proxeus-core/sys/model"
)

//This configuration can be used in two ways:
//1. Using the default meta in a struct
//2. Using the specified arguments in flag
type Configuration struct {
	SettingsFile    string `json:"settingsFile" default:"~/.proxeus/settings/main.json" usage:"Path to the settings file"`
	EthClientURL    string `json:"ethClientURL" default:"https://ropsten.infura.io/v3/" usage:"Ethereum client URL"`
	EthWebSocketURL string `json:"ethWebSocketURL" default:"wss://ropsten.infura.io/ws/v3/" usage:"Ethereum websocket URL"`

	ServiceAddress string `json:"serviceAddress" default:":1323" usage:"address and port of this service"`

	AutoTLS bool `json:"autoTLS" default:"false" usage:"Automatically generate Let's Encrypt certificate (Server must be reachable on port 443 from public internet)."`

	XESContractAddress string `json:"XESContractAddress" default:"0x84E0b37e8f5B4B86d5d299b0B0e33686405A3919"`

	AirdropWalletfile string `json:"airdropWalletfile" usage:"Path to File containing Private Key of the Wallet to fund Airdrops of XES and Ether."`
	AirdropWalletkey  string `json:"airdropWalletkey" usage:"Path to File containing the Key for the Airdrop Private Key."`

	model.Settings // extend cmd line args with settings
}

var Config Configuration

func init() {
	flagStruct(Config)
	pCfg := &Config
	if !strings.HasSuffix(os.Args[0], ".test") {
		flag.Parse()
	}
	flag.VisitAll(func(f *flag.Flag) {
		initStruct(f, reflect.ValueOf(pCfg).Elem())
	})

	fmt.Println("###########################################")
	fmt.Printf("CONFIG %#v\n", Config)
	fmt.Println("###########################################")
}

func initStruct(f *flag.Flag, v reflect.Value) {
	field := v.FieldByName(f.Name)
	strFlagVal := f.Value.String()
	if strFlagVal == f.DefValue {
		//if val same as default try from env var
		strVal := os.Getenv(f.Name)
		if strVal != "" {
			strFlagVal = strVal
		}
	}
	if field.Kind() == reflect.String {
		field.SetString(strFlagVal)
	} else if field.Kind() == reflect.Bool {
		bl, _ := strconv.ParseBool(strFlagVal)
		field.SetBool(bl)
	} else if field.Kind() != reflect.Invalid && field.Type() == reflect.TypeOf(model.CREATOR) {
		field.Set(reflect.ValueOf(model.StringToRole(strFlagVal)))
	} else if field.Kind() == reflect.Struct {
		initStruct(f, field)
	}
}

func flagStruct(strct interface{}) {
	roleType := reflect.TypeOf(model.CREATOR)
	v := reflect.ValueOf(strct)
	t := reflect.TypeOf(strct)
	if t.Kind() != reflect.Struct {
		return
	}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		fname := f.Name
		if f.Type.Kind() == reflect.Bool {
			bv, _ := strconv.ParseBool(f.Tag.Get("default"))
			flag.Bool(fname, bv, f.Tag.Get("usage"))
		} else if f.Type.Kind() == reflect.String || f.Type == roleType {
			flag.String(fname, f.Tag.Get("default"), f.Tag.Get("usage"))
		} else if f.Type.Kind() == reflect.Struct {
			flagStruct(v.Field(i).Interface())
		}
	}
}
