package config

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"unicode"

	"flag"

	"github.com/ProxeusApp/proxeus-core/sys/model"
)

const ProxeusEnvPrefix = "PROXEUS_"

//This configuration can be used in two ways:
//1. Using the default meta in a struct
//2. Using the specified arguments in flag
type Configuration struct {
	SettingsFile    string `json:"settingsFile" default:"~/.proxeus/settings/main.json" usage:"Path to the settings file"`
	EthClientURL    string `json:"ethClientURL" default:"https://goerli.infura.io/v3/" usage:"Ethereum client URL"`
	EthWebSocketURL string `json:"ethWebSocketURL" default:"wss://goerli.infura.io/ws/v3/" usage:"Ethereum websocket URL"`

	ServiceAddress string `json:"serviceAddress" default:":1323" usage:"address and port of this service"`

	AutoTLS bool `json:"autoTLS" default:"false" usage:"Automatically generate Let's Encrypt certificate (Server must be reachable on port 443 from public internet)."`

	XESContractAddress string `json:"XESContractAddress" default:"0x15FeA089CC48B4f4596242c138156e3B53579B37"`

	AirdropWalletfile string `json:"airdropWalletfile" usage:"Path to File containing Private Key of the Wallet to fund Airdrops of XES and Ether."`
	AirdropWalletkey  string `json:"airdropWalletkey" usage:"Path to File containing the Key for the Airdrop Private Key."`

	model.Settings // extend cmd line args with settings
}

var Config *Configuration

func Init() {
	if Config != nil {
		return
	}
	Config = New()
}

func New() *Configuration {
	var c Configuration

	flagFromStruct(flag.CommandLine, os.Environ(), &c)
	flag.Parse()

	return &c
}

func flagFromStruct(fs *flag.FlagSet, env []string, strct interface{}) {
	v := reflect.ValueOf(strct)
	if v.Kind() != reflect.Ptr {
		return
	}

	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return
	}
	t := v.Type()

	em := envMap(env)
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		fname := f.Name
		if f.Type.Kind() == reflect.Bool {
			d, u := defaultAndUsage(f, fname, em)
			bv, _ := strconv.ParseBool(d)
			fs.BoolVar(v.Field(i).Addr().Interface().(*bool), fname, bv, u)
		} else if f.Type.Kind() == reflect.String {
			d, u := defaultAndUsage(f, fname, em)
			fs.StringVar(v.Field(i).Addr().Interface().(*string), fname, d, u)
		} else if f.Type.Kind() == reflect.Struct {
			flagFromStruct(fs, env, v.Field(i).Addr().Interface())
		}
	}
}

func envMap(env []string) map[string]string {
	m := map[string]string{}
	for _, e := range env {
		kv := strings.Split(e, "=")
		m[kv[0]] = kv[1]
	}

	return m
}

func defaultAndUsage(f reflect.StructField, fname string, envMap map[string]string) (string, string) {
	envName := fieldToEnv(fname)
	if envValue := envMap[envName]; envValue != "" {
		return envValue, fmt.Sprintf("%s (%s=%s)", f.Tag.Get("usage"), envName, envValue)
	} else {
		return f.Tag.Get("default"), fmt.Sprintf("%s (%s)", f.Tag.Get("usage"), envName)
	}
}

func fieldToEnv(f string) string {
	var e strings.Builder

	e.WriteString(ProxeusEnvPrefix)
	previous := '_'
	for _, c := range f {
		if c == '-' {
			c = '_'
		}
		if previous == '_' && c == '_' {
			continue
		}
		if previous != '_' && !unicode.IsUpper(previous) && unicode.IsUpper(c) {
			e.WriteByte('_')
		}

		e.WriteRune(unicode.ToUpper(c))

		previous = c
	}

	return e.String()
}
