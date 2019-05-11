package main

import (
	log "github.com/sirupsen/logrus"
	"go.uber.org/fx"
	"librefrontier/RadioProvider/RadioBrowser"
	"os"
)

func Startup(a *ApiServer) {

}

func main() {
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.TextFormatter{
		ForceColors: true,
	})
	log.Info("Main Startup")

	app := fx.New(
		fx.Provide(
			NewEnvConfig,
			NewXmlBuilder,
			RadioBrowser.NewRadioBrowserClient,
			NewDatabase,
			NewApiController,
		),
		fx.Invoke(Startup),
	)

	app.Run()
}
