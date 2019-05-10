package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"librefrontier/RadioProvider"
	"librefrontier/RadioProvider/RadioBrowser"
	"log"
	"net/http"
	"strconv"
)

func login(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Printf("login token = %s\n", vars["token"])

	fmt.Fprint(w, "<EncryptedToken>3a3f5ac48a1dab4e</EncryptedToken>")
}

func gofile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Printf("gofile dlang = %s", vars["dlang"])

	menu := ListOfItems{
		ItemCount: 4,
		Items: []Item{
			{
				ItemType:     "Dir",
				Title:        "By Country",
				UrlDir:       "http://192.168.178.156/countries",
				UrlDirBackUp: "http://192.168.178.156/countries",
			}, {
				ItemType:     "Dir",
				Title:        "Most popular",
				UrlDir:       "http://192.168.178.156/stations/popular",
				UrlDirBackUp: "http://192.168.178.156/stations/popular",
			}, {
				ItemType:     "Dir",
				Title:        "Most liked",
				UrlDir:       "http://192.168.178.156/stations/liked",
				UrlDirBackUp: "http://192.168.178.156/stations/liked",
			}, {
				ItemType:     "Dir",
				Title:        "LibreFrontier PoC",
				UrlDir:       "http://192.168.178.156/",
				UrlDirBackUp: "http://192.168.178.156/",
			},
		},
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
		log.Fatal("Error converting str to int", err)
	}

	iEnd, err := strconv.Atoi(vars["endItems"])
	if err != nil {
		log.Fatal("Error converting str to int", err)
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
		log.Fatal("Error converting str to int", err)
	}

	iEnd, err := strconv.Atoi(vars["endItems"])
	if err != nil {
		log.Fatal("Error converting str to int", err)
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
		log.Fatal("Error converting str to int", err)
	}

	iEnd, err := strconv.Atoi(vars["endItems"])
	if err != nil {
		log.Fatal("Error converting str to int", err)
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
		log.Fatal("Error converting str to int", err)
	}

	iEnd, err := strconv.Atoi(vars["endItems"])
	if err != nil {
		log.Fatal("Error converting str to int", err)
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
	r.HandleFunc("/station/{station}/play", getStreamUrl)

	// ?sSearchtype=3&Search=75692&mac=b640a0c203b5ee50dac407aff8713da4&dlang=eng&fver=6&ven=teufel2

	err := http.ListenAndServe("0.0.0.0:80", r)
	if err != nil {
		log.Fatalf("Could not start web server: %s", err)
	}
}
