package RadioBrowser

import (
	"encoding/json"
	"gopkg.in/resty.v1"
	"librefrontier/RadioProvider"
	"log"
)

type RadioBrowserRadioProvider struct {
}

func (r *RadioBrowserRadioProvider) GetCountries() []RadioProvider.Country {
	resp, err := resty.R().Get("http://www.radio-browser.info/webservice/json/countries")
	if err != nil {
		log.Fatal("Error getting country list", err)
	}

	var countries []RadioProvider.Country

	err = json.Unmarshal(resp.Body(), &countries)
	if err != nil {
		log.Fatal("Error parsing json", err)
	}

	log.Printf("Result: %v", countries)

	return countries
}

func (r *RadioBrowserRadioProvider) GetStationsByCountry(countryId string) []RadioProvider.Station {
	resp, err := resty.R().Get("http://www.radio-browser.info/webservice/json/stations/bycountry/" + countryId)
	if err != nil {
		log.Fatal("Error getting station list", err)
	}

	var stations []RadioProvider.Station

	err = json.Unmarshal(resp.Body(), &stations)
	if err != nil {
		log.Fatal("Error parsing json", err)
	}

	log.Printf("Result: %v", stations)

	return stations
}

func (r *RadioBrowserRadioProvider) GetStationById(stationId string) RadioProvider.Station {
	resp, err := resty.R().Get("http://www.radio-browser.info/webservice/json/stations/byid/" + stationId)
	if err != nil {
		log.Fatal("Error getting station", err)
	}

	var stations []RadioProvider.Station

	err = json.Unmarshal(resp.Body(), &stations)
	if err != nil {
		log.Fatal("Error parsing json", err)
	}

	log.Printf("Result: %v", stations)

	return stations[0]
}
