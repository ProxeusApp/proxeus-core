package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path"

	"github.com/labstack/echo"

	"strings"

	cfg "git.proxeus.com/core/central/main/config"
	"git.proxeus.com/core/central/main/handlers"
	"git.proxeus.com/core/central/main/handlers/assets"
	"git.proxeus.com/core/central/main/www"
	"git.proxeus.com/core/central/sys"
	"git.proxeus.com/core/central/sys/i18n"
	"git.proxeus.com/core/central/sys/validate"
)

// ServerVersion is added to http headers and can be set during making a build
var ServerVersion = "build-unknown"

var embedded *www.Embedded

func main() {
	system, err := sys.NewWithSettings(cfg.Config.Settings)
	if err != nil {
		panic(err)
	}

	if system.TestMode {
		fmt.Println("#######################################################")
		fmt.Println("# STARTING PROXEUS IN TEST MODE - NOT FOR PRODUCTION #")
		fmt.Println("#######################################################")
	}

	www.SetSystem(system)

	embedded = &www.Embedded{Asset: assets.Asset}
	sys.ReadAllFile = func(path string) ([]byte, error) {
		return embedded.Asset(path)
	}

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
		err = system.DB.I18n.PutLang("en", true)
		if err != nil {
			fmt.Println("Error activating fallback lang: ", err)
		}
	}()

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
