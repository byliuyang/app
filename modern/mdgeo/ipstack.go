package mdgeo

import (
	"fmt"
	"net/http"

	"github.com/short-d/app/fw"
)

var _ fw.GeoLocation = (*IPStack)(nil)

const baseURL = "http://api.ipstack.com"

type jsonLanguage struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type jsonLocation struct {
	Capital     string         `json:"capital"`
	Languages   []jsonLanguage `json:"languages"`
	CallingCode string         `json:"calling_code"`
	IsEU        bool           `json:"is_eu"`
}

type jsonTimeZone struct {
	ID               string `json:"id"`
	Code             string `json:"code"`
	IsDaylightSaving bool
}

type jsonCurrency struct {
	Code   string `json:"code"`
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
}

type jsonResponse struct {
	ContinentCode string       `json:"continent_code"`
	ContinentName string       `json:"continent_name"`
	CountryCode   string       `json:"country_code"`
	CountryName   string       `json:"country_name"`
	RegionCode    string       `json:"region_code"`
	RegionName    string       `json:"region_name"`
	City          string       `json:"city"`
	ZipCode       string       `json:"zip"`
	Longitude     float64      `json:"longitude"`
	Latitude      float64      `json:"latitude"`
	Location      jsonLocation `json:"location"`
	TimeZone      jsonTimeZone `json:"time_zone"`
	Currency      jsonCurrency `json:"currency"`
}

type IPStack struct {
	apiKey      string
	httpRequest fw.HTTPRequest
	logger      fw.Logger
}

func (I IPStack) GetLocation(ipAddress string) (fw.Location, error) {
	url := fmt.Sprintf("%s/%s?access_key=%s", baseURL, ipAddress, I.apiKey)
	res := jsonResponse{}
	err := I.httpRequest.JSON(http.MethodGet, url, map[string]string{}, "", &res)
	if err != nil {
		I.logger.Error(err)
		return fw.Location{}, err
	}

	var languages []fw.Language
	for _, jsonLanguage := range res.Location.Languages {
		language := fw.Language{
			Code: jsonLanguage.Code,
			Name: jsonLanguage.Name,
		}
		languages = append(languages, language)
	}

	return fw.Location{
		Continent: fw.Continent{
			Code: res.ContinentCode,
			Name: res.ContinentName,
		},
		Country: fw.Country{
			Code: res.CountryCode,
			Name: res.CountryName,
		},
		Region: fw.Region{
			Code: res.RegionCode,
			Name: res.RegionName,
		},
		City: res.City,
		Currency: fw.Currency{
			Code: res.Currency.Code,
			Name: res.Currency.Name,
		},
		Languages:       languages,
		IsEuropeanUnion: res.Location.IsEU,
	}, nil
}

func NewIPStack(apiKey string, httpRequest fw.HTTPRequest, logger fw.Logger) IPStack {
	return IPStack{
		apiKey:      apiKey,
		httpRequest: httpRequest,
		logger:      logger,
	}
}
