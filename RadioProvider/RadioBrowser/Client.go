package RadioBrowser

import (
	"encoding/json"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gopkg.in/resty.v1"
	"librefrontier/RadioProvider"
	"net/url"
	"strconv"
)

type Client struct {
}

func NewRadioBrowserClient() RadioProvider.RadioProvider {
	return &Client{}
}

var _ RadioProvider.RadioProvider = (*Client)(nil)

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

	log.Debugf("Result: %v", countries)

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

	log.Debugf("Result: %v", stations)

	return stations, nil
}

func (r *Client) GetMostPopularStations(count int) ([]RadioProvider.Station, error) {
	resp, err := resty.R().Get("http://www.radio-browser.info/webservice/json/stations/topclick/" + strconv.Itoa(count))
	if err != nil {
		return nil, errors.Wrap(err, "get stations")
	}

	var stations []RadioProvider.Station

	err = json.Unmarshal(resp.Body(), &stations)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal stations")
	}

	log.Debugf("Result: %v", stations)

	return stations, nil
}

func (r *Client) GetMostLikedStations(count int) ([]RadioProvider.Station, error) {
	resp, err := resty.R().Get("http://www.radio-browser.info/webservice/json/stations/topvote/" + strconv.Itoa(count))
	if err != nil {
		return nil, errors.Wrap(err, "get stations")
	}

	var stations []RadioProvider.Station

	err = json.Unmarshal(resp.Body(), &stations)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal stations")
	}

	log.Debugf("Result: %v", stations)

	return stations, nil
}

func (r *Client) SearchStations(search string) ([]RadioProvider.Station, error) {
	resp, err := resty.R().Get("http://www.radio-browser.info/webservice/json/stations/byname/" + url.PathEscape(search))
	if err != nil {
		return nil, errors.Wrap(err, "get stations")
	}

	var stations []RadioProvider.Station

	err = json.Unmarshal(resp.Body(), &stations)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal stations")
	}

	log.Debugf("Result: %v", stations)

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

	log.Debugf("Result: %v", stations)

	if len(stations) > 0 {
		return stations[0], nil
	}

	return RadioProvider.Station{}, errors.New("No station found")
}
