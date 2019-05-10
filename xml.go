package main

import (
	"encoding/xml"
	"librefrontier/RadioProvider"
	"log"
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

func CreateCountryList(countries []RadioProvider.Country, start int, end int) ListOfItems {
	result := ListOfItems{
		ItemCount: len(countries),
	}

	end++

	if end > result.ItemCount {
		end = result.ItemCount
	}

	log.Printf("countries %d - %d\n", start, end)

	var items []Item

	items = append(items, Item{
		ItemType:          "Previous",
		UrlPrevious:       "http://192.168.178.156/setupapp/karcher/asp/BrowseXML/loginXML.asp?gofile=",
		UrlPreviousBackUp: "http://192.168.178.156/setupapp/karcher/asp/BrowseXML/loginXML.asp?gofile=",
	})
	for i := start; i < end; i++ {
		items = append(items, Item{
			ItemType:     "Dir",
			Title:        countries[i].Name,
			UrlDir:       "http://192.168.178.156/country/" + url.PathEscape(countries[i].Id),
			UrlDirBackUp: "http://192.168.178.156/country/" + url.PathEscape(countries[i].Id),
		})
	}

	result.Items = items

	return result
}

func CreateStationsList(stations []RadioProvider.Station, start int, end int) ListOfItems {
	result := ListOfItems{
		ItemCount: len(stations),
	}

	end++

	if end > result.ItemCount {
		end = result.ItemCount
	}

	log.Printf("stations %d - %d\n", start, end)

	var items []Item

	items = append(items, Item{
		ItemType:          "Previous",
		UrlPrevious:       "http://192.168.178.156/setupapp/karcher/asp/BrowseXML/loginXML.asp?gofile=",
		UrlPreviousBackUp: "http://192.168.178.156/setupapp/karcher/asp/BrowseXML/loginXML.asp?gofile=",
	})
	for i := start; i < end; i++ {
		items = append(items, CreateStationItem(stations[i]))
	}

	result.Items = items

	return result
}

func CreateStationItem(station RadioProvider.Station) Item {
	return Item{
		ItemType:    "Station",
		StationName: station.Name, // + " (" + station.Codec + " " + station.Bitrate + ")",
		StationId:   station.Id,
		//StationLocation:  station.Country,
		//StationDesc:      station.Homepage,
		//StationBandWidth: station.Bitrate,
		//StationMime:      station.Codec,
		//StationFormat:    station.Genre,
		StationUrl: "http://192.168.178.156/station/" + station.Id + "/play",
		//StationUrl: station.StreamUrl,
	}
}

func WriteToWire(w http.ResponseWriter, items ListOfItems) {
	result, err := xml.Marshal(items)
	if err != nil {
		log.Fatal("Error in xml.Marshal", err)
	}

	hdr := []byte(xml.Header)

	contentLength := len(hdr) + len(result)
	w.Header().Set("Content-Length", strconv.Itoa(contentLength))

	_, err = w.Write(hdr)
	if err != nil {
		log.Fatal("Error writing to wire", err)
	}

	_, err = w.Write(result)
	if err != nil {
		log.Fatal("Error writing to wire", err)
	}
}
