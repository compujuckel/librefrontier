package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"go.uber.org/fx"
	"librefrontier/RadioProvider"
	"net/http"
	"strconv"
)

type ApiServer struct {
	db    *Database
	cfg   *Config
	xml   *XmlBuilder
	radio RadioProvider.RadioProvider
}

func NewApiController(lc fx.Lifecycle, config *Config, database *Database, xmlBuilder *XmlBuilder, radioProvider RadioProvider.RadioProvider) *ApiServer {
	a := ApiServer{}
	a.cfg = config
	a.db = database
	a.xml = xmlBuilder
	a.radio = radioProvider

	r := mux.NewRouter()

	r.HandleFunc("/setupapp/karcher/asp/BrowseXML/loginXML.asp", a.login).
		Queries("token", "{token}")
	r.HandleFunc("/setupapp/karcher/asp/BrowseXML/loginXML.asp", a.gofile).
		Queries("gofile", "").
		Queries("dlang", "{dlang}")
	r.HandleFunc("/setupapp/karcher/asp/BrowseXML/Search.asp", a.search).
		Queries("sSearchtype", "{sSearchtype}").
		Queries("Search", "{Search}").
		Queries("mac", "{mac}").
		Queries("dlang", "{dlang}")
	r.HandleFunc("/countries", a.getCountries).
		Queries("startItems", "{startItems}").
		Queries("endItems", "{endItems}")
	r.HandleFunc("/country/{country}", a.getStationsByCountry).
		Queries("startItems", "{startItems}").
		Queries("endItems", "{endItems}")
	r.HandleFunc("/stations/popular", a.getMostPopularStations).
		Queries("startItems", "{startItems}").
		Queries("endItems", "{endItems}")
	r.HandleFunc("/stations/liked", a.getMostLikedStations).
		Queries("startItems", "{startItems}").
		Queries("endItems", "{endItems}")
	r.HandleFunc("/stations/search", a.searchStations).
		Queries("search", "{search}").
		Queries("startItems", "{startItems}").
		Queries("endItems", "{endItems}")
	r.HandleFunc("/station/{station}/play", a.getStreamUrl)
	r.HandleFunc("/station/{station}", a.getStationDetail).
		Queries("mac", "{mac}")
	r.HandleFunc("/favorite/add/{station}", a.addFavorite).
		Queries("mac", "{mac}")
	r.HandleFunc("/favorite/remove/{station}", a.removeFavorite).
		Queries("mac", "{mac}")
	r.HandleFunc("/favorites", a.getFavorites).
		Queries("mac", "{mac}").
		Queries("startItems", "{startItems}").
		Queries("endItems", "{endItems}")

	server := http.Server{
		Addr:    ":80",
		Handler: r,
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

func (a *ApiServer) login(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Printf("login token = %s\n", vars["token"])

	fmt.Fprint(w, "<EncryptedToken>3a3f5ac48a1dab4e</EncryptedToken>")
}

func (a *ApiServer) gofile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Printf("gofile dlang = %s", vars["dlang"])

	items := []Item{
		{
			ItemType:     "Dir",
			Title:        "Favorites",
			UrlDir:       a.cfg.apiBaseUrl + "/favorites",
			UrlDirBackUp: a.cfg.apiBaseUrl + "/favorites",
		}, {
			ItemType:     "Dir",
			Title:        "By Country",
			UrlDir:       a.cfg.apiBaseUrl + "/countries",
			UrlDirBackUp: a.cfg.apiBaseUrl + "/countries",
		}, {
			ItemType:     "Dir",
			Title:        "Most popular",
			UrlDir:       a.cfg.apiBaseUrl + "/stations/popular",
			UrlDirBackUp: a.cfg.apiBaseUrl + "/stations/popular",
		}, {
			ItemType:     "Dir",
			Title:        "Most liked",
			UrlDir:       a.cfg.apiBaseUrl + "/stations/liked",
			UrlDirBackUp: a.cfg.apiBaseUrl + "/stations/liked",
		}, {
			ItemType:        "Search",
			SearchURL:       a.cfg.apiBaseUrl + "/stations/search?sSearchType=2",
			SearchURLBackUp: a.cfg.apiBaseUrl + "/stations/search?sSearchType=2",
			SearchCaption:   "Search stations",
			SearchTextbox:   "",
			SearchGo:        "Search",
			SearchCancel:    "%search-cancel%",
		}, {
			ItemType:     "Dir",
			Title:        "LibreFrontier PoC",
			UrlDir:       a.cfg.apiBaseUrl + "/",
			UrlDirBackUp: a.cfg.apiBaseUrl + "/",
		},
	}

	menu := ListOfItems{
		ItemCount: len(items),
		Items:     items,
	}

	a.xml.WriteToWire(w, menu)
}

func (a *ApiServer) search(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	a.db.CreateDevice(vars["mac"])

	log.Printf("search mac = %s Search = %s sSearchtype = %s\n", vars["mac"], vars["Search"], vars["sSearchtype"])

	station, err := a.radio.GetStationById(vars["Search"])
	if err != nil {
		w.WriteHeader(404)
		return
	}
	list := a.xml.CreateStationsList([]RadioProvider.Station{station}, 0, 0, true)

	a.xml.WriteToWire(w, list)
}

func (a *ApiServer) getCountries(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	iStart, err := strconv.Atoi(vars["startItems"])
	if err != nil {
		log.Error("Error converting str to int", err)
		w.WriteHeader(400)
		return
	}

	iEnd, err := strconv.Atoi(vars["endItems"])
	if err != nil {
		log.Error("Error converting str to int", err)
		w.WriteHeader(400)
		return
	}

	countries, err := a.radio.GetCountries()
	if err != nil {
		w.WriteHeader(500)
		return
	}
	list := a.xml.CreateCountryList(countries, iStart-1, iEnd)

	a.xml.WriteToWire(w, list)
}

func (a *ApiServer) getStationsByCountry(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	iStart, err := strconv.Atoi(vars["startItems"])
	if err != nil {
		log.Error("Error converting str to int", err)
		w.WriteHeader(400)
		return
	}

	iEnd, err := strconv.Atoi(vars["endItems"])
	if err != nil {
		log.Error("Error converting str to int", err)
		w.WriteHeader(400)
		return
	}

	stations, err := a.radio.GetStationsByCountry(vars["country"])
	if err != nil {
		w.WriteHeader(500)
		return
	}
	list := a.xml.CreateStationsList(stations, iStart-1, iEnd, false)

	a.xml.WriteToWire(w, list)
}

func (a *ApiServer) getMostPopularStations(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	iStart, err := strconv.Atoi(vars["startItems"])
	if err != nil {
		log.Error("Error converting str to int", err)
		w.WriteHeader(400)
		return
	}

	iEnd, err := strconv.Atoi(vars["endItems"])
	if err != nil {
		log.Error("Error converting str to int", err)
		w.WriteHeader(400)
		return
	}

	stations, err := a.radio.GetMostPopularStations(100)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	list := a.xml.CreateStationsList(stations, iStart-1, iEnd, false)

	a.xml.WriteToWire(w, list)
}

func (a *ApiServer) getMostLikedStations(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	iStart, err := strconv.Atoi(vars["startItems"])
	if err != nil {
		log.Error("Error converting str to int", err)
		w.WriteHeader(400)
		return
	}

	iEnd, err := strconv.Atoi(vars["endItems"])
	if err != nil {
		log.Error("Error converting str to int", err)
		w.WriteHeader(400)
		return
	}

	stations, err := a.radio.GetMostLikedStations(100)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	list := a.xml.CreateStationsList(stations, iStart-1, iEnd, false)

	a.xml.WriteToWire(w, list)
}

func (a *ApiServer) searchStations(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	iStart, err := strconv.Atoi(vars["startItems"])
	if err != nil {
		log.Error("Error converting str to int", err)
		w.WriteHeader(400)
		return
	}

	iEnd, err := strconv.Atoi(vars["endItems"])
	if err != nil {
		log.Error("Error converting str to int", err)
		w.WriteHeader(400)
		return
	}

	stations, err := a.radio.SearchStations(vars["search"])
	if err != nil {
		w.WriteHeader(500)
		return
	}
	list := a.xml.CreateStationsList(stations, iStart-1, iEnd, false)

	a.xml.WriteToWire(w, list)
}

func (a *ApiServer) getStationDetail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	station, err := a.radio.GetStationById(vars["station"])
	if err != nil {
		w.WriteHeader(404)
		return
	}

	id, err := strconv.ParseUint(station.Id, 10, 32)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	fav := a.db.IsFavorite(vars["mac"], id)
	list := a.xml.CreateStationDetail(station, fav)

	a.xml.WriteToWire(w, list)
}

func (a *ApiServer) getStreamUrl(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	station, err := a.radio.GetStationById(vars["station"])
	if err != nil {
		w.WriteHeader(404)
		return
	}

	w.Write([]byte(station.StreamUrl))
}

func (a *ApiServer) addFavorite(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	station, err := a.radio.GetStationById(vars["station"])
	if err != nil {
		w.WriteHeader(404)
		return
	}

	id, err := strconv.ParseInt(station.Id, 10, 64)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	a.db.AddFavorite(vars["mac"], id, station.Name)
	log.Infof("Added favorite %s for mac %s", station.Name, vars["mac"])

	list := a.xml.CreateStationDetail(station, false)

	a.xml.WriteToWire(w, list)
}

func (a *ApiServer) removeFavorite(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	station, err := a.radio.GetStationById(vars["station"])
	if err != nil {
		w.WriteHeader(404)
		return
	}

	id, err := strconv.ParseUint(station.Id, 10, 64)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	a.db.RemoveFavorite(vars["mac"], id)
	log.Infof("Removed favorite %s for mac %s", station.Name, vars["mac"])

	list := a.xml.CreateStationDetail(station, false)

	a.xml.WriteToWire(w, list)
}

func (a *ApiServer) getFavorites(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	iStart, err := strconv.Atoi(vars["startItems"])
	if err != nil {
		log.Error("Error converting str to int", err)
		w.WriteHeader(400)
		return
	}

	iEnd, err := strconv.Atoi(vars["endItems"])
	if err != nil {
		log.Error("Error converting str to int", err)
		w.WriteHeader(400)
		return
	}

	stations := a.db.GetFavoriteStations(vars["mac"])

	list := a.xml.CreateStationsList(stations, iStart-1, iEnd, false)

	a.xml.WriteToWire(w, list)
}
