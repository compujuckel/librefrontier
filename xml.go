package main

import (
	"encoding/xml"
	log "github.com/sirupsen/logrus"
	"librefrontier/RadioProvider"
	"net/http"
	"net/url"
	"strconv"
)

type ListOfItems struct {
	XMLName   xml.Name `xml:"ListOfItems"`
	ItemCount int      `xml:"ItemCount"`
	Items     []Item   `xml:"Item"`
}

type Item struct {
	XMLName  xml.Name `xml:"Item"`
	ItemType string   `xml:"ItemType"`

	Title        string `xml:"Title,omitempty"`
	UrlDir       string `xml:"UrlDir,omitempty"`
	UrlDirBackUp string `xml:"UrlDirBackUp,omitempty"`

	UrlPrevious       string `xml:"UrlPrevious,omitempty"`
	UrlPreviousBackUp string `xml:"UrlPreviousBackUp,omitempty"`

	StationId        string `xml:"StationId,omitempty"`
	StationName      string `xml:"StationName,omitempty"`
	StationUrl       string `xml:"StationUrl,omitempty"`
	StationDesc      string `xml:"StationDesc,omitempty"`
	StationFormat    string `xml:"StationFormat,omitempty"`
	StationLocation  string `xml:"StationLocation,omitempty"`
	Logo             string `xml:"Logo,omitempty"`
	StationBandWidth string `xml:"StationBandWidth,omitempty"`
	StationMime      string `xml:"StationMime,omitempty"`
	Relia            string `xml:"Relia,omitempty"`
	Bookmark         string `xml:"Bookmark,omitempty"`

	SearchURL       string `xml:"SearchURL,omitempty"`
	SearchURLBackUp string `xml:"SearchURLBackUp,omitempty"`
	SearchCaption   string `xml:"SearchCaption,omitempty"`
	SearchTextbox   string `xml:"SearchTextbox,omitempty"`
	SearchGo        string `xml:"SearchGo,omitempty"`
	SearchCancel    string `xml:"SearchCancel,omitempty"`
}

type XmlBuilder struct {
	cfg *Config
}

func NewXmlBuilder(config *Config) *XmlBuilder {
	x := XmlBuilder{
		cfg: config,
	}

	return &x
}

func (x *XmlBuilder) CreateCountryList(countries []RadioProvider.Country, start int, end int) ListOfItems {
	result := ListOfItems{
		ItemCount: len(countries),
	}

	end++

	if end > result.ItemCount {
		end = result.ItemCount
	}

	log.Debugf("countries %d - %d\n", start, end)

	var items []Item

	items = append(items, Item{
		ItemType:          "Previous",
		UrlPrevious:       x.cfg.apiBaseUrl + "/setupapp/karcher/asp/BrowseXML/loginXML.asp?gofile=",
		UrlPreviousBackUp: x.cfg.apiBaseUrl + "/setupapp/karcher/asp/BrowseXML/loginXML.asp?gofile=",
	})
	for i := start; i < end; i++ {
		items = append(items, Item{
			ItemType:     "Dir",
			Title:        countries[i].Name,
			UrlDir:       x.cfg.apiBaseUrl + "/country/" + url.PathEscape(countries[i].Id),
			UrlDirBackUp: x.cfg.apiBaseUrl + "/country/" + url.PathEscape(countries[i].Id),
		})
	}

	result.Items = items

	return result
}

func (x *XmlBuilder) CreateStationsList(stations []RadioProvider.Station, start int, end int) ListOfItems {
	result := ListOfItems{
		ItemCount: len(stations),
	}

	end++

	if end > result.ItemCount {
		end = result.ItemCount
	}

	log.Debugf("stations %d - %d\n", start, end)

	var items []Item

	items = append(items, Item{
		ItemType:          "Previous",
		UrlPrevious:       x.cfg.apiBaseUrl + "/setupapp/karcher/asp/BrowseXML/loginXML.asp?gofile=",
		UrlPreviousBackUp: x.cfg.apiBaseUrl + "/setupapp/karcher/asp/BrowseXML/loginXML.asp?gofile=",
	})
	for i := start; i < end; i++ {
		items = append(items, x.CreateStationItem(stations[i]))
	}

	result.Items = items

	return result
}

func (x *XmlBuilder) CreateStationItem(station RadioProvider.Station) Item {
	return Item{
		ItemType:    "Station",
		StationName: station.Name,
		StationId:   station.Id,
		//StationLocation:  station.Country,
		//StationDesc:      station.Homepage,
		//StationBandWidth: station.Bitrate,
		//StationMime:      station.Codec,
		//StationFormat:    station.Genre,
		StationUrl: x.cfg.apiBaseUrl + "/station/" + station.Id + "/play",
		//StationUrl: station.StreamUrl,
	}
}

func (x *XmlBuilder) WriteToWire(w http.ResponseWriter, items ListOfItems) {
	result, err := xml.Marshal(items)
	if err != nil {
		log.Error("Error in xml.Marshal", err)
	}

	hdr := []byte(xml.Header)

	contentLength := len(hdr) + len(result)
	w.Header().Set("Content-Length", strconv.Itoa(contentLength))

	_, err = w.Write(hdr)
	if err != nil {
		log.Error("Error writing to wire", err)
	}

	_, err = w.Write(result)
	if err != nil {
		log.Error("Error writing to wire", err)
	}
}
