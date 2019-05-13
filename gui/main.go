package main

import (
	"github.com/compujuckel/librefrontier/common"
	log "github.com/sirupsen/logrus"
	"go.uber.org/fx"
	"os"
)

func Startup(g *GuiServer) {

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
			common.NewDatabase,
			NewGuiController,
		),
		fx.Invoke(Startup),
	)

	app.Run()
}
