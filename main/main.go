package main

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"path"

	"github.com/ethereum/go-ethereum/accounts/abi"

	"git.proxeus.com/core/central/lib/wallet"

	"git.proxeus.com/core/central/main/handlers/blockchain"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"strings"

	"git.proxeus.com/core/central/lib/default_server"
	cfg "git.proxeus.com/core/central/main/config"
	"git.proxeus.com/core/central/main/handlers"
	"git.proxeus.com/core/central/main/handlers/api"
	"git.proxeus.com/core/central/main/handlers/assets"
	"git.proxeus.com/core/central/main/www"
	"git.proxeus.com/core/central/sys"
	"git.proxeus.com/core/central/sys/i18n"
	"git.proxeus.com/core/central/sys/validate"
)

// ServerVersion is added to http headers and can be set during making a build
var ServerVersion = "build-unknown"

func xVersionHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("X-Version", ServerVersion)
		return next(c)
	}
}

var embedded *www.Embedded

func main() {
	e := echo.New()
	e.HTTPErrorHandler = www.DefaultHTTPErrorHandler
	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return h(&www.Context{Context: c})
		}
	})
	e.Pre(xVersionHeader)
	c := middleware.DefaultSecureConfig
	c.XFrameOptions = ""
	e.Pre(middleware.SecureWithConfig(c))

	e.GET("/static/*", StaticHandler)

	www.SetupSession(e)
	system, err := sys.NewWithSettings(cfg.Config.Settings)
	if err != nil {
		panic(err)
	}
	embedded = &www.Embedded{Asset: assets.Asset}
	sys.ReadAllFile = func(path string) ([]byte, error) {
		return embedded.Asset(path)
	}
	www.SetSystem(system)

	go func() { //parse i18n from the UI assets to provide them under the translation section
		i18nUIParser := i18n.NewUIParser()
		dir := "static/assets/js"
		list, _ := assets.AssetDir(dir)
		for _, p := range list {
			bts, _ := assets.Asset(dir + "/" + p)
			i18nUIParser.Parse(bts)
		}
		trans := i18nUIParser.Translations()
		langs, _ := system.DB.I18n.GetAllLangs()
		//include the lang codes as keys as well to translate the lang label
		for _, l := range langs {
			trans[l.Code] = l.Code
		}

		//include the validation messages
		for _, msg := range validate.AllMessages() {
			trans[msg] = msg
		}

		lang, _ := system.DB.I18n.GetFallback()
		allTrans, _ := system.DB.I18n.GetAll(lang)
		for k, v := range trans {
			if _, exists := allTrans[k]; !exists {
				fmt.Println("inserting initial translation:", k, "::::", v)
				_ = system.DB.I18n.Put(lang, k, v)
			}
		}
	}()

	secure := www.NewSecurity()

	// Routes
	e.Pre(middleware.Secure())

	api.ServerVersion = ServerVersion

	handlers.MainHostedAPI(e, secure, system)

	XESABI, err := abi.JSON(strings.NewReader(wallet.XesMainTokenABI))
	if err != nil {
		panic(err)
	}

	xesAdapter := blockchain.NewAdapter(cfg.Config.XESContractAddress, XESABI)

	bcListenerPayment := blockchain.NewPaymentListener(xesAdapter, cfg.Config.EthWebSocketURL,
		cfg.Config.EthClientURL, system.DB.WorkflowPaymentsDB)
	ctxPay := context.Background()
	go bcListenerPayment.Listen(ctxPay)

	ProxeusFSABI, err := abi.JSON(strings.NewReader(wallet.ProxeusFSABI))
	if err != nil {
		panic(err)
	}

	bcListenerSignature := blockchain.NewSignatureListener(cfg.Config.EthWebSocketURL,
		cfg.Config.EthClientURL, system.GetSettings().BlockchainContractAddress, system.DB.SignatureRequestsDB, system.DB.User, system.EmailSender, ProxeusFSABI)
	ctxSig := context.Background()
	go bcListenerSignature.Listen(ctxSig)

	default_server.StartServer(e, cfg.Config.ServiceAddress, false)
	system.Shutdown()
}

// StaticHandler servers static files from bindata assets
func StaticHandler(c echo.Context) error {
	url := c.Request().URL.String()
	i := strings.Index(url, "?")
	if i != -1 {
		url = url[:i]
	}
	ct := ""
	header := c.Response().Header()
	ext := path.Ext(url)
	if ext == ".js" || ext == ".css" {
		header.Set("Cache-Control", "public,max-age=31536000")
	} else {
		header.Set("Cache-Control", "public,max-age=72000")
	}
	b, err := embedded.FindAssetWithCT(url, &ct)
	if err == nil {
		return c.Blob(http.StatusOK, ct, b)
	}
	return echo.ErrNotFound
}
