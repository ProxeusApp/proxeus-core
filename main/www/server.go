package www

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path"
	"syscall"

	"golang.org/x/crypto/acme/autocert"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func Setup(serverVersion string) *echo.Echo {
	e := echo.New()
	e.HTTPErrorHandler = DefaultHTTPErrorHandler

	// Pre routing middleware
	e.Pre(xVersionHeader(serverVersion))
	c := middleware.DefaultSecureConfig
	c.XFrameOptions = ""
	e.Pre(middleware.SecureWithConfig(c))
	e.Pre(middleware.Secure())

	// Post routing middleware
	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return h(&Context{Context: c})
		}
	})
	//Simple Request Logging
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "[echo] ${time_rfc3339} client=${remote_ip}, method=${method}, uri=${uri}, status=${status}\n",
	}))

	//Request Logging with User Info and Body on Error
	e.Use(middleware.BodyDump(func(e echo.Context, reqBody, resBody []byte) {
		c := e.(*Context)
		s := c.Session(false)
		if s == nil {
			return
		}
		if s.S.ID != "" {
			id := s.UserID()
			user, err := c.System().DB.User.Get(s, id)
			if err != nil {
				return
			}
			userName := user.Name
			userAddr := user.EthereumAddr
			log.Println("[echo] Method: "+e.Request().Method, "Status:", e.Response().Status, "User: "+userAddr, "("+userName+")", "URI: "+e.Request().RequestURI)
			if len(reqBody) > 0 && c.Response().Status != 200 && c.Response().Status != 404 {
				fmt.Printf("[echo][errorrequest] %s\n", reqBody)
			}
		}

	}))

	e.Use(SessionMiddleware())
	e.Use(SessionTokenAuth)

	return e
}

func StartServer(e *echo.Echo, addr string, autoTLS bool) {
	if autoTLS {
		dirCache := path.Join(os.TempDir(), ".cache")
		e.AutoTLSManager.Cache = autocert.DirCache(dirCache)
	}
	quit := make(chan os.Signal)

	// Start server
	go func() {
		if autoTLS {
			fmt.Println("starting https at", addr)
			if err := e.StartAutoTLS(addr); err != nil {
				log.Println(err)
			}
		} else {
			fmt.Println("starting plain http at", addr)
			if err := e.Start(addr); err != nil {
				log.Println(err)
			}
		}
	}()

	signal.Notify(quit, os.Interrupt)
	signal.Notify(quit, syscall.SIGTERM)
	<-quit
	if err := e.Shutdown(context.Background()); err != nil {
		log.Fatal(err)
	}
}

func xVersionHeader(version string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("X-Version", version)
			return next(c)
		}
	}
}
