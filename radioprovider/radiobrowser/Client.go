package radiobrowser

import (
	"encoding/json"
	"github.com/compujuckel/librefrontier/radioprovider"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gopkg.in/resty.v1"
	"net/url"
	"strconv"
)

type Client struct {
}

func NewRadioBrowserClient() radioprovider.RadioProvider {
	return &Client{}
}

var _ radioprovider.RadioProvider = (*Client)(nil)

func (r *Client) GetCountries() ([]radioprovider.Country, error) {
	resp, err := resty.R().Get("http://www.radio-browser.info/webservice/json/countries")
	if err != nil {
		return nil, errors.Wrap(err, "get countries")
	}

	var countries []radioprovider.Country

	err = json.Unmarshal(resp.Body(), &countries)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal countries")
	}

	log.Debugf("Result: %v", countries)

	return countries, nil
}

func (r *Client) GetStationsByCountry(countryId string) ([]radioprovider.Station, error) {
	resp, err := resty.R().Get("http://www.radio-browser.info/webservice/json/stations/bycountry/" + countryId)
	if err != nil {
		return nil, errors.Wrap(err, "get stations")
	}

	var stations []radioprovider.Station

	err = json.Unmarshal(resp.Body(), &stations)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal stations")
	}

	log.Debugf("Result: %v", stations)

	return stations, nil
}

func (r *Client) GetMostPopularStations(count int) ([]radioprovider.Station, error) {
	resp, err := resty.R().Get("http://www.radio-browser.info/webservice/json/stations/topclick/" + strconv.Itoa(count))
	if err != nil {
		return nil, errors.Wrap(err, "get stations")
	}

	var stations []radioprovider.Station

	err = json.Unmarshal(resp.Body(), &stations)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal stations")
	}

	log.Debugf("Result: %v", stations)

	return stations, nil
}

func (r *Client) GetMostLikedStations(count int) ([]radioprovider.Station, error) {
	resp, err := resty.R().Get("http://www.radio-browser.info/webservice/json/stations/topvote/" + strconv.Itoa(count))
	if err != nil {
		return nil, errors.Wrap(err, "get stations")
	}

	var stations []radioprovider.Station

	err = json.Unmarshal(resp.Body(), &stations)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal stations")
	}

	log.Debugf("Result: %v", stations)

	return stations, nil
}

func (r *Client) SearchStations(search string) ([]radioprovider.Station, error) {
	resp, err := resty.R().Get("http://www.radio-browser.info/webservice/json/stations/byname/" + url.PathEscape(search))
	if err != nil {
		return nil, errors.Wrap(err, "get stations")
	}

	var stations []radioprovider.Station

	err = json.Unmarshal(resp.Body(), &stations)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal stations")
	}

	log.Debugf("Result: %v", stations)

	return stations, nil
}

func (r *Client) GetStationById(stationId string) (radioprovider.Station, error) {
	resp, err := resty.R().Get("http://www.radio-browser.info/webservice/json/stations/byid/" + stationId)
	if err != nil {
		return radioprovider.Station{}, errors.Wrap(err, "get station")
	}

	var stations []radioprovider.Station

	err = json.Unmarshal(resp.Body(), &stations)
	if err != nil {
		return radioprovider.Station{}, errors.Wrap(err, "unmarshal station")
	}

	log.Debugf("Result: %v", stations)

	if len(stations) > 0 {
		return stations[0], nil
	}

	return radioprovider.Station{}, errors.New("No station found")
}
