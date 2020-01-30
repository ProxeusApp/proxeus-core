package pricegetter

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	tpl "html/template"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/ProxeusApp/proxeus-core/main/priceservice"
	"github.com/ProxeusApp/proxeus-core/sys/model"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

const jwtSecret = "my secret"
const proxeusUrl = "http://127.0.0.1:1323"
const serviceUrl = "127.0.0.1:8011"
const authKey = "auth"

var dummyConfigStore = map[string]*configData{}

type configData struct {
	FromCurrency string
	ToCurrency   string
}

func health(c echo.Context) error {
	return c.String(http.StatusOK, "I'm ok")
}

func config(c echo.Context) error {
	id, err := nodeID(c)
	if err != nil {
		return err
	}
	conf := getConfig(id)
	var buf bytes.Buffer
	err = configPage.Execute(&buf, map[string]string{
		"Id":           id,
		"AuthToken":    c.QueryParam(authKey),
		"FromCurrency": conf.FromCurrency,
		"ToCurrency":   conf.ToCurrency,
	})
	if err != nil {
		return err
	}
	return c.Stream(http.StatusOK, "text/html", &buf)
}

func setConfig(c echo.Context) error {
	id, err := nodeID(c)
	if err != nil {
		return err
	}
	conf := &configData{
		FromCurrency: strings.TrimSpace(c.FormValue("FromCurrency")),
		ToCurrency:   strings.TrimSpace(c.FormValue("ToCurrency")),
	}
	if conf.FromCurrency == "" || conf.ToCurrency == "" {
		return errors.New("empty currency")
	}
	dummyConfigStore[id] = conf
	return config(c)
}

func getConfig(id string) *configData {
	conf, ok := dummyConfigStore[id]
	if !ok {
		conf = &configData{
			FromCurrency: "CHF",
			ToCurrency:   "XES",
		}
	}
	return conf
}

func nop(_ echo.Context) error {
	return nil
}

func next(c echo.Context) error {
	id, err := nodeID(c)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}
	conf := getConfig(id)

	ratio, err := convert(conf.FromCurrency, conf.ToCurrency)
	if err != nil {
		return err
	}

	re, err := regexp.Compile("[0-9]+ " + conf.FromCurrency)
	if err != nil {
		return err
	}

	var replaceErr error
	replaced := re.ReplaceAllFunc(body, func(r []byte) []byte {
		if replaceErr != nil {
			return r
		}
		d := strings.Split(string(r), " ")
		val := d[0]
		valInt, err := strconv.Atoi(val)
		if err != nil {
			replaceErr = err
			return r
		}
		return []byte(fmt.Sprintf("%.3f %s", ratio*float64(valInt), conf.ToCurrency))
	})
	if replaceErr != nil {
		return replaceErr
	}
	return c.String(http.StatusOK, string(replaced))
}

func convert(from, to string) (float64, error) {
	s := priceservice.NewCryptoComparePriceService("API_KEY",
		"https://min-api.cryptocompare.com")
	return s.GetPriceInFor(to, from)
}

func nodeID(c echo.Context) (string, error) {
	id := c.Param("id")
	if id == "" {
		return "", errors.New("empty id")
	}
	t := c.Get("user").(*jwt.Token)
	if id != t.Claims.(jwt.MapClaims)["jti"].(string) {
		return "", errors.New("id mismatch")
	}
	return id, nil
}

var configPage *tpl.Template

func Run() {
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Recover())
	e.GET("/health", health)
	{
		g := e.Group("/node/:id")
		conf := middleware.DefaultJWTConfig
		conf.SigningKey = []byte(jwtSecret)
		conf.TokenLookup = "query:" + authKey
		g.Use(middleware.JWTWithConfig(conf))

		g.GET("/config", config)
		g.POST("/config", setConfig)
		g.POST("/next", next)
		g.POST("/remove", nop)
		g.POST("/close", nop)
	}
	parseTemplates()
	register()
	e.Start(serviceUrl)
}

func register() {
	client := http.Client{Timeout: 5 * time.Second}
	for {
		n := model.ExternalNode{
			ID:     "priceGetter",
			Name:   "priceGetter",
			Detail: "Converts crypto-currencies",
			Url:    "http://" + serviceUrl,
			Secret: jwtSecret,
		}
		buf, err := json.Marshal(n)
		if err != nil {
			panic(err.Error())
		}
		r, err := client.Post(proxeusUrl+"/api/admin/external/register",
			"application/json", bytes.NewBuffer(buf))
		if err == nil && r.StatusCode == http.StatusOK {
			log.Print("[nodeservice] ", n.Name, " registered")
			return
		}
		log.Print("[nodeservice] error registering ", n.Name, " err ", err)
		time.Sleep(5 * time.Second)
	}
}

func parseTemplates() {
	var err error
	configPage, err = tpl.New("").Parse(configHTML)
	if err != nil {
		panic(err.Error())
	}
}

const configHTML = `
<!DOCTYPE html>
<html>
<body>
<form action="/node/{{.Id}}/config?auth={{.AuthToken}}" method="post">
Convert fom currency: <input type="text" size="2" name="FromCurrency" value="{{.FromCurrency}}">
to currency: <input type="text" size="2" name="ToCurrency" value="{{.ToCurrency}}"><br/>
<input type="submit" value="Submit">
</form>
</body>
</html>
`
