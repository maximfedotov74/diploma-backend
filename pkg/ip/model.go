package ip

type ipLocationData struct {
	PostalCode     string `json:"postal_code"`
	Country        string `json:"country"`
	RegionWithType string `json:"region_with_type"`
	CityWithType   string `json:"city_with_type"`
	City           string `json:"city"`
	GeoLat         string `json:"geo_lat"`
	GeoLon         string `json:"geo_lon"`
}

type ipLocation struct {
	Value             string         `json:"value"`
	UnrestrictedValue string         `json:"unrestricted_value"`
	Data              ipLocationData `json:"data"`
}

type IpLocationResponse struct {
	Location ipLocation `json:"location"`
}
