package services

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type AddressService struct {
	Am *AppleMapsService
}

func NewAddressService(am *AppleMapsService) *AddressService {
	return &AddressService{Am: am}
}

type ReverseGeocodeResponse struct {
	Results []ReverseGeocodeResult `json:"results"`
}

type ReverseGeocodeResult struct {
	Coordinate            Coordinate        `json:"coordinate"`
	DisplayMapRegion      MapRegion         `json:"displayMapRegion"`
	Name                  string            `json:"name"`
	FormattedAddressLines []string          `json:"formattedAddressLines"`
	StructuredAddress     StructuredAddress `json:"structuredAddress"`
	Country               string            `json:"country"`
	CountryCode           string            `json:"countryCode"`
}

type Coordinate struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type MapRegion struct {
	SouthLatitude float64 `json:"southLatitude"`
	WestLongitude float64 `json:"westLongitude"`
	NorthLatitude float64 `json:"northLatitude"`
	EastLongitude float64 `json:"eastLongitude"`
}

type StructuredAddress struct {
	AdministrativeArea     string   `json:"administrativeArea"`
	AdministrativeAreaCode string   `json:"administrativeAreaCode"`
	Locality               string   `json:"locality"`
	PostCode               string   `json:"postCode"`
	Thoroughfare           string   `json:"thoroughfare"`
	FullThoroughfare       string   `json:"fullThoroughfare"`
	AreasOfInterest        []string `json:"areasOfInterest"`
}

func (a *AddressService) EncodeCoords(lat, lon float64) string {
	coords := fmt.Sprintf("%.8f,%.8f", lat, lon)

	encoded := base64.URLEncoding.EncodeToString([]byte(coords))

	return encoded
}

func (a *AddressService) DecodeCoords(encoded string) (float64, float64, error) {
	var lat float64
	var lon float64

	decoded, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		return 0, 0, err
	}

	fmt.Sscanf(string(decoded), "%f,%f", &lat, &lon)
	return lat, lon, nil
}

func (a *AddressService) GetAddressFromCoords(lat, lon float64) (string, error) {
	reverseGeoCodeUrl := os.Getenv("APPLE_RGEOCODE_URL")
	u, err := url.Parse(reverseGeoCodeUrl)

	if err != nil {
		return "", fmt.Errorf(err.Error())
	}

	latStr := strconv.FormatFloat(lat, 'f', 8, 64)
	lonStr := strconv.FormatFloat(lon, 'f', 8, 64)

	query := u.Query()
	query.Set("loc", latStr+","+lonStr)
	u.RawQuery = query.Encode()

	token := a.Am.GetAmToken()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return "", fmt.Errorf("Failed to get address")
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Do(req)
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Reverse geocode failed: %s", body)
	}

	var result ReverseGeocodeResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", err
	}

	var addressStr string

	for _,val := range result.Results[0].FormattedAddressLines {
		addressStr=addressStr+" "+strings.TrimSpace(val)
	}
	
	return strings.TrimSpace(addressStr), nil
}

func (a *AddressService) GetCoordsFromAddress(address string) (lat float64, lon float64, err error) {

	return 999, 999, nil
}
