package main

import (
	"github.com/compujuckel/librefrontier/common"
	"github.com/compujuckel/librefrontier/common/radioprovider/radiobrowser"
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
			common.NewEnvConfig,
			common.NewXmlBuilder,
			radiobrowser.NewRadioBrowserClient,
			common.NewDatabase,
			NewApiController,
		),
		fx.Invoke(Startup),
	)

	app.Run()
}
