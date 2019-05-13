package main

import (
	"context"
	"github.com/compujuckel/librefrontier/common"
	"github.com/compujuckel/librefrontier/common/radioprovider"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"go.uber.org/fx"
	"net/http"
	"strconv"
)

type ApiServer struct {
	db    *common.Database
	cfg   *common.Config
	xml   *common.XmlBuilder
	gin   *gin.Engine
	radio radioprovider.RadioProvider
}

type DeviceInfo struct {
	Mac      string `form:"mac" binding:"required"`
	Language string `form:"dlang"`
	Fver     string `form:"fver"`
	Vendor   string `form:"ven"`
}

type PaginatedRequest struct {
	Device *DeviceInfo
	Start  int `form:"startItems" binding:"required"`
	End    int `form:"endItems" binding:"required"`
}

type SearchRequestSingle struct {
	Device     *DeviceInfo
	SearchType int    `form:"sSearchtype" binding:"required"`
	Search     string `form:"Search" binding:"required"`
}

type SearchRequest struct {
	Device     *DeviceInfo
	Start      int    `form:"startItems"`
	End        int    `form:"endItems"`
	SearchType int    `form:"sSearchtype" binding:"required"`
	Search     string `form:"search" binding:"required"`
}

func NewApiController(lc fx.Lifecycle, config *common.Config, database *common.Database, xmlBuilder *common.XmlBuilder, radioProvider radioprovider.RadioProvider) *ApiServer {
	a := ApiServer{}
	a.cfg = config
	a.db = database
	a.xml = xmlBuilder
	a.radio = radioProvider
	a.gin = gin.Default()

	a.gin.GET("/setupapp/karcher/asp/BrowseXML/loginXML.asp", a.fsLoginXML)
	a.gin.GET("/setupapp/karcher/asp/BrowseXML/Search.asp", a.fsSearch)
	a.gin.GET("/countries", a.getCountries)
	a.gin.GET("/country/:country", a.getStationsByCountry)
	a.gin.GET("/stations/popular", a.getMostPopularStations)
	a.gin.GET("/stations/liked", a.getMostLikedStations)
	a.gin.GET("/stations/search", a.searchStations)
	a.gin.GET("/station/:station/play", a.getStreamUrl)
	a.gin.GET("/station/:station", a.getStationDetail)
	a.gin.GET("/favorite/add/:station", a.addFavorite)
	a.gin.GET("/favorite/remove/:station", a.removeFavorite)
	a.gin.GET("/favorites", a.getFavorites)
	a.gin.GET("/empty", a.getEmpty)

	server := http.Server{
		Addr:    ":80",
		Handler: a.gin,
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

	return &a
}

func (a *ApiServer) fsLoginXML(c *gin.Context) {
	log.Printf("fs_loginXML")

	if c.Query("token") == "0" {
		// TODO investigate how this is used
		c.String(http.StatusOK, "<EncryptedToken>3a3f5ac48a1dab4e</EncryptedToken>")
		return
	}

	items := []common.Item{
		{
			ItemType:     "Dir",
			Title:        "Favorites",
			UrlDir:       a.cfg.GetApiBaseUrl() + "/favorites",
			UrlDirBackUp: a.cfg.GetApiBaseUrl() + "/favorites",
		}, {
			ItemType:     "Dir",
			Title:        "By Country",
			UrlDir:       a.cfg.GetApiBaseUrl() + "/countries",
			UrlDirBackUp: a.cfg.GetApiBaseUrl() + "/countries",
		}, {
			ItemType:     "Dir",
			Title:        "Most popular",
			UrlDir:       a.cfg.GetApiBaseUrl() + "/stations/popular",
			UrlDirBackUp: a.cfg.GetApiBaseUrl() + "/stations/popular",
		}, {
			ItemType:     "Dir",
			Title:        "Most liked",
			UrlDir:       a.cfg.GetApiBaseUrl() + "/stations/liked",
			UrlDirBackUp: a.cfg.GetApiBaseUrl() + "/stations/liked",
		}, {
			ItemType:        "Search",
			SearchURL:       a.cfg.GetApiBaseUrl() + "/stations/search?sSearchtype=2",
			SearchURLBackUp: a.cfg.GetApiBaseUrl() + "/stations/search?sSearchtype=2",
			SearchCaption:   "Search stations",
			SearchTextbox:   "",
			SearchGo:        "Search",
			SearchCancel:    "%search-cancel%",
		}, {
			ItemType:     "Dir",
			Title:        "LibreFrontier PoC",
			UrlDir:       a.cfg.GetApiBaseUrl() + "/empty",
			UrlDirBackUp: a.cfg.GetApiBaseUrl() + "/empty",
		},
	}

	menu := common.ListOfItems{
		ItemCount: len(items),
		Items:     items,
	}

	// sadly we cannot use c.XML here because it does not write the XML header
	a.xml.WriteToWire(c.Writer, menu)
}

func (a *ApiServer) fsSearch(c *gin.Context) {
	var r SearchRequestSingle

	if c.Bind(&r) != nil {
		return
	}

	a.db.CreateDevice(r.Device.Mac)

	log.Printf("search mac = %s Search = %s sSearchtype = %s\n", r.Device.Mac, r.Search, r.SearchType)

	station, err := a.radio.GetStationById(r.Search)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	list := a.xml.CreateStationsList([]radioprovider.Station{station}, 0, 0, true)

	a.xml.WriteToWire(c.Writer, list)
}

func (a *ApiServer) getEmpty(c *gin.Context) {
	a.xml.WriteToWire(c.Writer, common.ListOfItems{})
}

func (a *ApiServer) getCountries(c *gin.Context) {
	var p PaginatedRequest
	if c.Bind(&p) != nil {
		return
	}

	countries, err := a.radio.GetCountries()
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	list := a.xml.CreateCountryList(countries, p.Start-1, p.End)

	a.xml.WriteToWire(c.Writer, list)
}

func (a *ApiServer) getStationsByCountry(c *gin.Context) {
	var p PaginatedRequest
	if c.Bind(&p) != nil {
		return
	}

	stations, err := a.radio.GetStationsByCountry(c.Param("country"))
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	list := a.xml.CreateStationsList(stations, p.Start-1, p.End, false)

	a.xml.WriteToWire(c.Writer, list)
}

func (a *ApiServer) getMostPopularStations(c *gin.Context) {
	var p PaginatedRequest
	if c.Bind(&p) != nil {
		return
	}

	stations, err := a.radio.GetMostPopularStations(100)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	list := a.xml.CreateStationsList(stations, p.Start-1, p.End, false)

	a.xml.WriteToWire(c.Writer, list)
}

func (a *ApiServer) getMostLikedStations(c *gin.Context) {
	var p PaginatedRequest
	if c.Bind(&p) != nil {
		return
	}

	stations, err := a.radio.GetMostLikedStations(100)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	list := a.xml.CreateStationsList(stations, p.Start-1, p.End, false)

	a.xml.WriteToWire(c.Writer, list)
}

func (a *ApiServer) searchStations(c *gin.Context) {
	var s SearchRequest
	if c.Bind(&s) != nil {
		return
	}

	stations, err := a.radio.SearchStations(s.Search)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	list := a.xml.CreateStationsList(stations, s.Start-1, s.End, false)

	a.xml.WriteToWire(c.Writer, list)
}

func (a *ApiServer) getStationDetail(c *gin.Context) {
	var d DeviceInfo
	if c.Bind(&d) != nil {
		return
	}

	station, err := a.radio.GetStationById(c.Param("station"))
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	id, err := strconv.ParseUint(station.Id, 10, 32)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	fav := a.db.IsFavorite(d.Mac, id)
	list := a.xml.CreateStationDetail(station, fav)

	a.xml.WriteToWire(c.Writer, list)
}

func (a *ApiServer) getStreamUrl(c *gin.Context) {
	station, err := a.radio.GetStationById(c.Param("station"))
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.String(http.StatusOK, station.StreamUrl)
}

func (a *ApiServer) addFavorite(c *gin.Context) {
	var d DeviceInfo
	if c.Bind(&d) != nil {
		return
	}

	station, err := a.radio.GetStationById(c.Param("station"))
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	id, err := strconv.ParseInt(station.Id, 10, 64)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	a.db.AddFavorite(d.Mac, id, station.Name)
	log.Infof("Added favorite %s for mac %s", station.Name, d.Mac)

	list := a.xml.CreateStationDetail(station, false)

	a.xml.WriteToWire(c.Writer, list)
}

func (a *ApiServer) removeFavorite(c *gin.Context) {
	var d DeviceInfo
	if c.Bind(&d) != nil {
		return
	}

	station, err := a.radio.GetStationById(c.Param("station"))
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	id, err := strconv.ParseUint(station.Id, 10, 64)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	a.db.RemoveFavorite(d.Mac, id)
	log.Infof("Removed favorite %s for mac %s", station.Name, d.Mac)

	list := a.xml.CreateStationDetail(station, false)

	a.xml.WriteToWire(c.Writer, list)
}

func (a *ApiServer) getFavorites(c *gin.Context) {
	var p PaginatedRequest
	if c.Bind(&p) != nil {
		return
	}

	stations := a.db.GetFavoriteStations(p.Device.Mac)

	list := a.xml.CreateStationsList(stations, p.Start-1, p.End, false)

	a.xml.WriteToWire(c.Writer, list)
}
