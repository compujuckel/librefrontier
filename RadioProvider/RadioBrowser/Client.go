package RadioBrowser

import (
	"encoding/json"
	"github.com/pkg/errors"
	"gopkg.in/resty.v1"
	"librefrontier/RadioProvider"
	"log"
)

type Client struct {
}

func (r *Client) GetCountries() ([]RadioProvider.Country, error) {
	resp, err := resty.R().Get("http://www.radio-browser.info/webservice/json/countries")
	if err != nil {
		return nil, errors.Wrap(err, "get countries")
	}

	var countries []RadioProvider.Country

	err = json.Unmarshal(resp.Body(), &countries)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal countries")
	}

	log.Printf("Result: %v", countries)

	return countries, nil
}

func (r *Client) GetStationsByCountry(countryId string) ([]RadioProvider.Station, error) {
	resp, err := resty.R().Get("http://www.radio-browser.info/webservice/json/stations/bycountry/" + countryId)
	if err != nil {
		return nil, errors.Wrap(err, "get stations")
	}

	var stations []RadioProvider.Station

	err = json.Unmarshal(resp.Body(), &stations)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal stations")
	}

	log.Printf("Result: %v", stations)

	return stations, nil
}

func (r *Client) GetStationById(stationId string) (RadioProvider.Station, error) {
	resp, err := resty.R().Get("http://www.radio-browser.info/webservice/json/stations/byid/" + stationId)
	if err != nil {
		return RadioProvider.Station{}, errors.Wrap(err, "get station")
	}

	var stations []RadioProvider.Station

	err = json.Unmarshal(resp.Body(), &stations)
	if err != nil {
		return RadioProvider.Station{}, errors.Wrap(err, "unmarshal station")
	}

	log.Printf("Result: %v", stations)

	if len(stations) > 0 {
		return stations[0], nil
	}

	return RadioProvider.Station{}, errors.New("No station found")
}
