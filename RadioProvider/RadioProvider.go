package RadioProvider

type Country struct {
	Name         string `json:"name"`
	Id           string `json:"value"`
	StationCount string `json:"stationcount"`
}

type Station struct {
	Name      string `json:"name"`
	Id        string `json:"id"`
	StreamUrl string `json:"url"`
	Codec     string `json:"codec"`
	Bitrate   string `json:"bitrate"`
	Homepage  string `json:"homepage"`
	Country   string `json:"country"`
	Genre     string `json:"tags"`
}

type RadioProvider interface {
	GetCountries() ([]Country, error)
	GetStationsByCountry(countryId string) ([]Station, error)
	GetStationById(stationId string) (Station, error)
}
