package services

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"
)

type GeoResult struct {
	Country string `json:"country"`
	City    string `json:"city"`
}

type GeoService struct {
	client *http.Client
}

func NewGeoService() *GeoService {
	return &GeoService{
		client: &http.Client{Timeout: 2 * time.Second},
	}
}

// Lookup returns country and city for an IP address using ip-api.com (free tier).
func (g *GeoService) Lookup(ip string) (*GeoResult, error) {
	// strip port if present
	if host, _, err := net.SplitHostPort(ip); err == nil {
		ip = host
	}

	// skip private/local IPs
	if ip == "" || ip == "127.0.0.1" || ip == "::1" || strings.HasPrefix(ip, "192.168.") || strings.HasPrefix(ip, "10.") {
		return &GeoResult{Country: "local", City: "local"}, nil
	}

	resp, err := g.client.Get(fmt.Sprintf("http://ip-api.com/json/%s?fields=country,city,countryCode", ip))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Country     string `json:"country"`
		City        string `json:"city"`
		CountryCode string `json:"countryCode"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &GeoResult{
		Country: result.CountryCode,
		City:    result.City,
	}, nil
}
