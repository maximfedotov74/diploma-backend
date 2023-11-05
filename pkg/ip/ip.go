package ip

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type IpService struct {
	DadataApiKey string
	httpClient   http.Client
}

func NewIpService(key string) *IpService {
	return &IpService{DadataApiKey: key, httpClient: http.Client{}}
}

func (is *IpService) GetGeolocation(ip string) (*IpLocationResponse, error) {
	url := fmt.Sprintf("https://suggestions.dadata.ru/suggestions/api/4_1/rs/iplocate/address?ip=%s", ip)
	method := "GET"
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", is.DadataApiKey))
	response, err := is.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	ipResponse := IpLocationResponse{}

	err = json.Unmarshal(body, &ipResponse)
	if err != nil {
		return nil, err
	}

	return &ipResponse, nil
}
