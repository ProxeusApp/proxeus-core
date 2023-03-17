module github.com/ProxeusApp/proxeus-core

go 1.16

require (
	github.com/DataDog/zstd v1.4.8 // indirect
	github.com/SparkPost/gosparkpost v0.2.0
	github.com/asdine/storm v0.0.0-20190418133842-e0f77eada154
	github.com/c2h5oh/datasize v0.0.0-20200825124411-48ed595a09d2
	github.com/cespare/cp v1.1.1 // indirect
	github.com/davecgh/go-spew v1.1.1
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/disintegration/imaging v1.6.2
	github.com/dop251/goja v0.0.0-20220408131256-ffe77e20c6f1
	github.com/ethereum/go-ethereum v1.10.22
	github.com/fatih/structs v1.1.0 // indirect
	github.com/golang/mock v1.6.0
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/google/renameio v1.0.1
	github.com/gorilla/sessions v1.2.1
	github.com/h2non/filetype v1.1.3
	github.com/imkira/go-interpol v1.1.0 // indirect
	github.com/labstack/gommon v0.3.1
	github.com/mattn/go-runewidth v0.0.10 // indirect
	github.com/onsi/gomega v1.19.0
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/pkg/errors v0.9.1
	github.com/robertkrimen/otto v0.0.0-20211024170158-b87d35c0b86f
	github.com/satori/go.uuid v1.2.0
	github.com/stretchr/testify v1.7.2
	github.com/tidwall/pretty v1.1.0 // indirect
	github.com/vmihailenco/msgpack v4.0.4+incompatible // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb // indirect
	github.com/xeipuuv/gojsonschema v1.2.0 // indirect
	go.etcd.io/bbolt v1.3.6 // indirect
	go.mongodb.org/mongo-driver v1.11.1
	golang.org/x/crypto v0.0.0-20220622213112-05595931fe9d
	golang.org/x/net v0.4.0
	gopkg.in/gavv/httpexpect.v2 v2.3.1
)

require (
	github.com/Sereal/Sereal v0.0.0-20200820125258-a016b7cda3f3 // indirect
	github.com/dlclark/regexp2 v1.4.1-0.20220329233251-d0559a0de6e3 // indirect
	github.com/fasthttp/websocket v1.4.3-rc.10 // indirect
	github.com/gballet/go-libpcsclite v0.0.0-20191108122812-4678299bea08 // indirect
	github.com/go-bindata/go-bindata/v3 v3.1.3 // indirect
	github.com/go-stack/stack v1.8.1 // indirect
	github.com/gorilla/context v1.1.1 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/kisielk/errcheck v1.6.2 // indirect
	github.com/klauspost/compress v1.15.1 // indirect
	github.com/labstack/echo v3.3.10+incompatible
	github.com/labstack/echo-contrib v0.12.0
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/prometheus/tsdb v0.10.0 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/rjeczalik/notify v0.9.2 // indirect
	github.com/sergi/go-diff v1.2.0 // indirect
	github.com/shirou/gopsutil v3.21.11+incompatible // indirect
	github.com/status-im/keycard-go v0.0.0-20200402102358-957c09536969 // indirect
	github.com/tklauser/go-sysconf v0.3.10 // indirect
	github.com/tyler-smith/go-bip39 v1.1.0 // indirect
	github.com/valyala/fasthttp v1.35.0 // indirect
	github.com/wadey/gocovmerge v0.0.0-20160331181800-b5bfa59ec0ad // indirect
	github.com/xdg-go/scram v1.1.1 // indirect
	github.com/youmark/pkcs8 v0.0.0-20201027041543-1326539a0a0a // indirect
	github.com/yuin/goldmark v1.5.3 // indirect
	github.com/yusufpapurcu/wmi v1.2.2 // indirect
	golang.org/x/image v0.0.0-20220413100746-70e8d0d3baa9 // indirect
	golang.org/x/lint v0.0.0-20210508222113-6edffad5e616 // indirect
	golang.org/x/tools v0.4.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
)

replace (
	github.com/ProxeusApp/proxeus-core => ./
	github.com/labstack/echo-contrib => github.com/labstack/echo-contrib v0.0.0-20180222075343-7d9d9632a4aa // fix https://github.com/ProxeusApp/proxeus-core/issues/216
	gopkg.in/gavv/httpexpect.v2 => github.com/gavv/httpexpect/v2 v2.2.0 // fix https://github.com/gavv/httpexpect/issues/60
	gopkg.in/urfave/cli.v1 => github.com/urfave/cli v1.22.5 // fix https://github.com/ProxeusApp/proxeus-core/issues/213
)
