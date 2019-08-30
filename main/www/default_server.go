package www

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"path"
	"syscall"

	"golang.org/x/crypto/acme/autocert"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/natefinch/lumberjack"
)

func Setup(logFileLocation string) *echo.Echo {
	e := echo.New()

	// logging setup
	{
		e.Debug = true
		var lw io.Writer
		lw = &lumberjack.Logger{
			Filename: logFileLocation,
			MaxSize:  100, // MB
			MaxAge:   120, // days
		}
		// test it
		_, err := lw.Write([]byte("log init\n"))
		if err != nil {
			log.Printf("File logging disabled due to: <%s>\n", err)
			// fallback to std
			lw = os.Stdout
		} else {
			log.Printf("Logging to: %s\n", logFileLocation)
		}
		e.Logger.SetOutput(lw)
		log.SetOutput(lw)
		e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{Output: lw}))
	}

	// very important
	e.Use(middleware.Recover())

	e.HTTPErrorHandler = DefaultHTTPErrorHandler
	e.HideBanner = true

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
