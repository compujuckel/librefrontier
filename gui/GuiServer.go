package main

import (
	"context"
	"github.com/compujuckel/librefrontier/common"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"go.uber.org/fx"
	"net/http"
)

type GuiServer struct {
	db  *common.Database
	cfg *common.Config
	gin *gin.Engine
}

func (g *GuiServer) GetFavorites(c *gin.Context) {
	favorites := g.db.GetFavoriteStations("...")

	c.HTML(http.StatusOK, "favorites.html.tpl", gin.H{
		"favorites": favorites,
	})
}

func NewGuiController(lc fx.Lifecycle, config *common.Config, database *common.Database) *GuiServer {
	g := GuiServer{}
	g.cfg = config
	g.db = database
	g.gin = gin.Default()

	g.gin.LoadHTMLGlob("templates/*")
	g.gin.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html.tpl", gin.H{})
	})
	g.gin.GET("/favorites", g.GetFavorites)

	server := http.Server{
		Addr:    ":8080",
		Handler: g.gin,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Info("Starting HTTP server.")
			// In production, we'd want to separate the Listen and Serve phases for
			// better error-handling.
			go server.ListenAndServe()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("Stopping HTTP server.")
			return server.Shutdown(ctx)
		},
	})

	return &g
}
