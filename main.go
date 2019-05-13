package main

import (
	"github.com/compujuckel/librefrontier/radioprovider/radiobrowser"
	log "github.com/sirupsen/logrus"
	"go.uber.org/fx"
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
			radiobrowser.NewRadioBrowserClient,
			NewDatabase,
			NewApiController,
		),
		fx.Invoke(Startup),
	)

	app.Run()
}
