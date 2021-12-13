module github.com/ProxeusApp/proxeus-core

go 1.16

require (
	github.com/DataDog/zstd v1.4.8 // indirect
	github.com/SparkPost/gosparkpost v0.2.0
	github.com/asdine/storm v0.0.0-20190418133842-e0f77eada154
	github.com/btcsuite/btcd v0.21.0-beta // indirect
	github.com/c2h5oh/datasize v0.0.0-20200825124411-48ed595a09d2
	github.com/cespare/cp v1.1.1 // indirect
	github.com/davecgh/go-spew v1.1.1
	github.com/deckarep/golang-set v1.7.1 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/disintegration/imaging v1.6.2
	github.com/dop251/goja v0.0.0-20210322220816-6fc852574a34
	github.com/ethereum/go-ethereum v1.10.9
	github.com/fatih/structs v1.1.0 // indirect
	github.com/golang/mock v1.5.0
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/google/renameio v1.0.0
	github.com/gorilla/sessions v1.2.1
	github.com/h2non/filetype v1.1.1
	github.com/imkira/go-interpol v1.1.0 // indirect
	github.com/jackpal/go-nat-pmp v1.0.2 // indirect
	github.com/labstack/gommon v0.3.0
	github.com/mattn/go-runewidth v0.0.10 // indirect
	github.com/onsi/gomega v1.11.0
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/pkg/errors v0.9.1
	github.com/robertkrimen/otto v0.0.0-20200922221731-ef014fd054ac
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/satori/go.uuid v1.2.0
	github.com/stretchr/testify v1.7.0
	github.com/tidwall/pretty v1.1.0 // indirect
	github.com/vmihailenco/msgpack v4.0.4+incompatible // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb // indirect
	github.com/xeipuuv/gojsonschema v1.2.0 // indirect
	go.etcd.io/bbolt v1.3.5 // indirect
	go.mongodb.org/mongo-driver v1.7.3
	golang.org/x/crypto v0.0.0-20210513164829-c07d793c2f9a
	golang.org/x/net v0.0.0-20211209124913-491a49abca63
	golang.org/x/sys v0.0.0-20211209171907-798191bca915 // indirect
	gopkg.in/gavv/httpexpect.v2 v2.2.0
	gopkg.in/sourcemap.v1 v1.0.5 // indirect
)

require (
	github.com/Sereal/Sereal v0.0.0-20200820125258-a016b7cda3f3 // indirect
	github.com/StackExchange/wmi v0.0.0-20210224194228-fe8f1750fd46 // indirect
	github.com/fasthttp/websocket v1.4.3-rc.10 // indirect
	github.com/gballet/go-libpcsclite v0.0.0-20191108122812-4678299bea08 // indirect
	github.com/go-ole/go-ole v1.2.5 // indirect
	github.com/golang/protobuf v1.5.1 // indirect
	github.com/gorilla/context v1.1.1 // indirect
	github.com/labstack/echo v3.3.10+incompatible
	github.com/labstack/echo-contrib v0.9.0
	github.com/onsi/ginkgo v1.15.2 // indirect
	github.com/prometheus/tsdb v0.10.0 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/rjeczalik/notify v0.9.2 // indirect
	github.com/sergi/go-diff v1.1.0 // indirect
	github.com/status-im/keycard-go v0.0.0-20200402102358-957c09536969 // indirect
	github.com/tyler-smith/go-bip39 v1.1.0 // indirect
	github.com/youmark/pkcs8 v0.0.0-20201027041543-1326539a0a0a // indirect
	golang.org/x/image v0.0.0-20210220032944-ac19c3e999fb // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	gopkg.in/urfave/cli.v1 v1.22.5 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

replace (
	github.com/ProxeusApp/proxeus-core => ./
	github.com/labstack/echo-contrib => github.com/labstack/echo-contrib v0.0.0-20180222075343-7d9d9632a4aa // fix https://github.com/ProxeusApp/proxeus-core/issues/216
	gopkg.in/gavv/httpexpect.v2 => github.com/gavv/httpexpect/v2 v2.2.0 // fix https://github.com/gavv/httpexpect/issues/60
	gopkg.in/urfave/cli.v1 => github.com/urfave/cli v1.22.5 // fix https://github.com/ProxeusApp/proxeus-core/issues/213
)
