package main

import (
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"librefrontier/RadioProvider"
	"librefrontier/RadioProvider/RadioBrowser"
	"net/http"
	"os"
	"strconv"
)

var baseUrl = os.Getenv("LF_BASE_URL")

func login(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Printf("login token = %s\n", vars["token"])

	fmt.Fprint(w, "<EncryptedToken>3a3f5ac48a1dab4e</EncryptedToken>")
}

func gofile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Printf("gofile dlang = %s", vars["dlang"])

	items := []Item{
		{
			ItemType:     "Dir",
			Title:        "By Country",
			UrlDir:       baseUrl + "/countries",
			UrlDirBackUp: baseUrl + "/countries",
		}, {
			ItemType:     "Dir",
			Title:        "Most popular",
			UrlDir:       baseUrl + "/stations/popular",
			UrlDirBackUp: baseUrl + "/stations/popular",
		}, {
			ItemType:     "Dir",
			Title:        "Most liked",
			UrlDir:       baseUrl + "/stations/liked",
			UrlDirBackUp: baseUrl + "/stations/liked",
		}, {
			ItemType:        "Search",
			SearchURL:       baseUrl + "/stations/search?sSearchType=2",
			SearchURLBackUp: baseUrl + "/stations/search?sSearchType=2",
			SearchCaption:   "Search stations",
			SearchTextbox:   "",
			SearchGo:        "Search",
			SearchCancel:    "%search-cancel%",
		}, {
			ItemType:     "Dir",
			Title:        "LibreFrontier PoC",
			UrlDir:       baseUrl + "/",
			UrlDirBackUp: baseUrl + "/",
		},
	}

	menu := ListOfItems{
		ItemCount: len(items),
		Items:     items,
	}

	WriteToWire(w, menu)
}

func search(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	log.Printf("search mac = %s Search = %s sSearchtype = %s\n", vars["mac"], vars["Search"], vars["sSearchtype"])

	rp := RadioBrowser.Client{}

	station, err := rp.GetStationById(vars["Search"])
	if err != nil {
		w.WriteHeader(404)
		return
	}
	list := CreateStationsList([]RadioProvider.Station{station}, 0, 0)

	WriteToWire(w, list)
}

func getCountries(w http.ResponseWriter, r *http.Request) {
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

	rp := RadioBrowser.Client{}

	countries, err := rp.GetCountries()
	if err != nil {
		w.WriteHeader(500)
		return
	}
	list := CreateCountryList(countries, iStart-1, iEnd)

	WriteToWire(w, list)
}

func getStationsByCountry(w http.ResponseWriter, r *http.Request) {
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

	rp := RadioBrowser.Client{}

	stations, err := rp.GetStationsByCountry(vars["country"])
	if err != nil {
		w.WriteHeader(500)
		return
	}
	list := CreateStationsList(stations, iStart-1, iEnd)

	WriteToWire(w, list)
}

func getMostPopularStations(w http.ResponseWriter, r *http.Request) {
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

	rp := RadioBrowser.Client{}

	stations, err := rp.GetMostPopularStations(100)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	list := CreateStationsList(stations, iStart-1, iEnd)

	WriteToWire(w, list)
}

func getMostLikedStations(w http.ResponseWriter, r *http.Request) {
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

	rp := RadioBrowser.Client{}

	stations, err := rp.GetMostLikedStations(100)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	list := CreateStationsList(stations, iStart-1, iEnd)

	WriteToWire(w, list)
}

func searchStations(w http.ResponseWriter, r *http.Request) {
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

	rp := RadioBrowser.Client{}

	stations, err := rp.SearchStations(vars["search"])
	if err != nil {
		w.WriteHeader(500)
		return
	}
	list := CreateStationsList(stations, iStart-1, iEnd)

	WriteToWire(w, list)
}

func getStreamUrl(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	rp := RadioBrowser.Client{}

	station, err := rp.GetStationById(vars["station"])
	if err != nil {
		w.WriteHeader(404)
		return
	}

	w.Write([]byte(station.StreamUrl))
}

func main() {
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.TextFormatter{
		ForceColors: true,
	})
	log.Info("Main Startup")

	r := mux.NewRouter()

	r.HandleFunc("/setupapp/karcher/asp/BrowseXML/loginXML.asp", login).
		Queries("token", "{token}")
	r.HandleFunc("/setupapp/karcher/asp/BrowseXML/loginXML.asp", gofile).
		Queries("gofile", "").
		Queries("dlang", "{dlang}")
	r.HandleFunc("/setupapp/karcher/asp/BrowseXML/Search.asp", search).
		Queries("sSearchtype", "{sSearchtype}").
		Queries("Search", "{Search}").
		Queries("mac", "{mac}").
		Queries("dlang", "{dlang}")
	r.HandleFunc("/countries", getCountries).
		Queries("startItems", "{startItems}").
		Queries("endItems", "{endItems}")
	r.HandleFunc("/country/{country}", getStationsByCountry).
		Queries("startItems", "{startItems}").
		Queries("endItems", "{endItems}")
	r.HandleFunc("/stations/popular", getMostPopularStations).
		Queries("startItems", "{startItems}").
		Queries("endItems", "{endItems}")
	r.HandleFunc("/stations/liked", getMostLikedStations).
		Queries("startItems", "{startItems}").
		Queries("endItems", "{endItems}")
	r.HandleFunc("/stations/search", searchStations).
		Queries("search", "{search}").
		Queries("startItems", "{startItems}").
		Queries("endItems", "{endItems}")
	r.HandleFunc("/station/{station}/play", getStreamUrl)

	// ?sSearchtype=3&Search=75692&mac=b640a0c203b5ee50dac407aff8713da4&dlang=eng&fver=6&ven=teufel2

	err := http.ListenAndServe("0.0.0.0:80", r)
	if err != nil {
		log.Fatalf("Could not start web server: %s", err)
	}
}
