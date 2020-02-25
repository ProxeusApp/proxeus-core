package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path"

	"github.com/ProxeusApp/proxeus-core/main/handlers/workflow"

	"github.com/ProxeusApp/proxeus-core/service"

	"github.com/ProxeusApp/proxeus-core/main/handlers/api"
	"github.com/ProxeusApp/proxeus-core/main/handlers/payment"

	"github.com/labstack/echo"

	"strings"

	cfg "github.com/ProxeusApp/proxeus-core/main/config"
	"github.com/ProxeusApp/proxeus-core/main/handlers"
	"github.com/ProxeusApp/proxeus-core/main/handlers/assets"
	"github.com/ProxeusApp/proxeus-core/main/www"
	"github.com/ProxeusApp/proxeus-core/storage"
	"github.com/ProxeusApp/proxeus-core/sys"
	"github.com/ProxeusApp/proxeus-core/sys/i18n"
	"github.com/ProxeusApp/proxeus-core/sys/validate"
)

// ServerVersion is added to http headers and can be set during making a build
var ServerVersion = "build-unknown"

var embedded *www.Embedded

func main() {
	cfg.Init()
	system, err := sys.NewWithSettings(cfg.Config.SettingsFile, &cfg.Config.Settings)
	if err != nil {
		panic(err)
	}

	if system.TestMode {
		fmt.Println()
		fmt.Println("#######################################################")
		fmt.Println("# STARTING PROXEUS IN TEST MODE - NOT FOR PRODUCTION #")
		fmt.Println("#######################################################")
		fmt.Println()
	}

	if system.AllowHttp {
		fmt.Println()
		fmt.Println("#######################################################")
		fmt.Println("# ALLOWING HTTP - NOT FOR PRODUCTION                  #")
		fmt.Println("#######################################################")
		fmt.Println()
	}

	fmt.Println()
	fmt.Println("#######################################################")
	fmt.Printf("Configuration: %#v\n", cfg.Config)
	fmt.Printf("system settings: %#v\n", system.GetSettings())
	fmt.Println("#######################################################")
	fmt.Println()

	//Important: Pass system to services (and not e.g. system.DB.WorkflowPayments because system.DB variable is replaced on calling api/handlers.PostInit()
	service.Init(system)
	userService := service.NewUserService()
	paymentService := service.NewPaymentService(userService)
	workflowService := service.NewWorkflowService(userService)
	nodeService := service.NewNodeService(workflowService)
	fileService := service.NewFileService()
	documentService := service.NewDocumentService(userService, fileService)
	templateDocumentService := service.NewTemplateDocumentService()
	userDocumentService := service.NewUserDocumentService(userService, fileService, templateDocumentService)

	payment.Init(paymentService, userService)
	api.Init(paymentService, userService, workflowService, documentService, userDocumentService, fileService, templateDocumentService)
	workflow.Init(workflowService, userService, nodeService)

	www.SetSystem(system)

	embedded = &www.Embedded{Asset: assets.Asset}
	sys.ReadAllFile = func(path string) ([]byte, error) {
		return embedded.Asset(path)
	}

	err = initI18n(system.DB.I18n)
	if err != nil {
		fmt.Printf("Error while initialising i18n: %s\n", err)
	}

	e := www.Setup(ServerVersion)

	// Static route
	e.GET("/static/*", StaticHandler)

	// Initial config middleware
	configured, err := system.Configured()
	if err != nil && !os.IsNotExist(err) {
		panic(err)
	}
	if !configured {
		e.Use(www.NewInitialHandler(configured).Handler)
	}

	// Main routes
	handlers.MainHostedAPI(e, www.NewSecurity(), ServerVersion)

	www.StartServer(e, cfg.Config.ServiceAddress, false)
	system.Shutdown()
}

func initI18n(db storage.I18nIF) error {
	i18nUIParser := i18n.NewUIParser()
	dir := "static/assets/js"
	list, _ := assets.AssetDir(dir)

	for _, p := range list {
		bts, _ := assets.Asset(dir + "/" + p)
		i18nUIParser.Parse(bts)
	}
	trans := i18nUIParser.Translations()
	langs, _ := db.GetAllLangs()
	//include the lang codes as keys as well to translate the lang label
	for _, l := range langs {
		trans[l.Code] = l.Code
	}

	//include the validation messages
	for _, msg := range validate.AllMessages() {
		trans[msg] = msg
	}

	lang, _ := db.GetFallback()
	allTrans, _ := db.GetAll(lang)
	for k, v := range trans {
		if _, exists := allTrans[k]; !exists {
			_ = db.Put(lang, k, v)
		}
	}
	err := db.PutLang("en", true)
	if err != nil {
		return err
	}

	return nil
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
	b, err := embedded.FindAssetWithContentType(url, &ct)
	if err == nil {
		return c.Blob(http.StatusOK, ct, b)
	}
	return echo.ErrNotFound
}
