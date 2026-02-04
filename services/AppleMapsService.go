package services

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AppleMapsClient struct {
	Ams               IAppleMapsService
	geocodeUrl        string
	reverseGeoCodeUrl string
	searchUrl         string
}

func NewAppleMapsClient(ams IAppleMapsService) *AppleMapsClient {
	return &AppleMapsClient{Ams: ams, geocodeUrl: "https://maps-api.apple.com/v1/geocode", reverseGeoCodeUrl: "https://maps-api.apple.com/v1/reverseGeocode", searchUrl: "https://maps-api.apple.com/v1/search"}
}

func (mc *AppleMapsClient) ReverseGeocode(lat, lon float64) (string, error) {
	u, err := url.Parse(mc.reverseGeoCodeUrl)

	if err != nil {
		return "", fmt.Errorf("AppleMapsClient.ReverseGeocode: failed to parse Url")
	}

	latStr := strconv.FormatFloat(lat, 'f', 8, 64)
	lonStr := strconv.FormatFloat(lon, 'f', 8, 64)

	query := u.Query()
	query.Set("loc", latStr+","+lonStr)
	u.RawQuery = query.Encode()

	token := mc.Ams.GetAmToken()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return "", fmt.Errorf("AppleMapsClient.ReverseGeocode: Failed to get address")
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Do(req)

	if err != nil {
		return "", fmt.Errorf("AppleMapsClient.ReverseGeocode: Failed to initialize client")
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("AppleMapsClient.ReverseGeocode: Reverse geocode failed: %s", body)
	}

	var result GeocodeResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", err
	}

	var addressStr string

	for _, val := range result.Results[0].FormattedAddressLines {
		addressStr = addressStr + " " + strings.TrimSpace(val)
	}

	return strings.TrimSpace(addressStr), nil
}

func (mc *AppleMapsClient) Geocode(address string) (float64, float64, error) {
	u, err := url.Parse(mc.geocodeUrl)

	if err != nil {
		return 0, 0, err
	}

	query := u.Query()
	query.Set("q", address)
	u.RawQuery = query.Encode()

	token := mc.Ams.GetAmToken()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)

	if err != nil {
		return 0, 0, err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Do(req)

	if err != nil {
		return 0, 0, err
	}

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return 0, 0, fmt.Errorf("Geocode failed: %s", body)
	}

	var result GeocodeResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return 0, 0, err
	}

	return result.Results[0].Coordinate.Latitude, result.Results[0].Coordinate.Longitude, nil
}

func (mc *AppleMapsClient) Search(lat, lon float64) {
	u, err := url.Parse(mc.searchUrl)

	if err != nil {
		return
	}

	latStr := strconv.FormatFloat(lat, 'f', 8, 64)
	lonStr := strconv.FormatFloat(lon, 'f', 8, 64)

	lat1 := lat + 0.015
	lon1 := lon + 0.015
	lat2 := lat - 0.015
	lon2 := lon - 0.015

	lat1Str := strconv.FormatFloat(lat1, 'f', 8, 64)
	lon1Str := strconv.FormatFloat(lon1, 'f', 8, 64)
	lat2Str := strconv.FormatFloat(lat2, 'f', 8, 64)
	lon2Str := strconv.FormatFloat(lon2, 'f', 8, 64)

	query := u.Query()
	query.Set("q", "Cafe")
	query.Set("includePoiCategories", "Bakery,Beach,Cafe,FoodMarket,Hiking,Landmark,Library,Museum,NationalMonument,Park,Restaurant,Store")
	query.Set("userLocation", latStr+","+lonStr)
	query.Set("resultTypeFilter", "Poi")
	query.Set("searchRegionPriority", "required")
	query.Set("searchRegion", lat1Str+","+lon1Str+","+lat2Str+","+lon2Str)

	u.RawQuery = query.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)

	if err != nil {
		return
	}

	token := mc.Ams.GetAmToken()
	req.Header.Set("Authorization", "Bearer "+token)

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Do(req)

	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		fmt.Println(string(body))
		return
	}
	body, _ := io.ReadAll(resp.Body)

	fmt.Println(string(body))
}

type IAppleMapsService interface {
	GenerateMapsToken() string
	GenerateToken() string
	GetAmToken() string
	GetToken() string
	IsValidToken() bool
}

type AppleMapsService struct {
	Token    string
	AmToken  string
	AmExpiry int64
	Header   TokenHeader
	Payload  TokenPayload
	tokenUrl string
}

func NewAppleMapsService() *AppleMapsService {
	return &AppleMapsService{tokenUrl: "https://maps-api.apple.com/v1/token"}
}

func (am *AppleMapsService) GetAmToken() string {
	if am.IsValidToken() {
		return am.AmToken
	} else {
		return am.GenerateMapsToken()
	}
}

func (am *AppleMapsService) GetToken() string {
	if am.IsValidToken() {
		return am.Token
	} else {
		return am.GenerateToken()
	}
}

func (am *AppleMapsService) GenerateToken() string {
	MAPS_ID := os.Getenv("APPLE_MAPS_ID")
	TEAM_ID := os.Getenv("APPLE_TEAM_ID")
	PKEY := os.Getenv("APPLE_MAPS_PKEY")

	now := time.Now().Unix()
	exp := now + 60*60

	header := TokenHeader{Alg: "ES256", Kid: MAPS_ID, Typ: "JWT"}

	headerMap := map[string]any{
		"alg": header.Alg,
		"kid": header.Kid,
		"typ": header.Typ,
	}

	payload := TokenPayload{Iss: TEAM_ID, Iat: now, Exp: exp}

	t := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"iss": payload.Iss,
		"iat": payload.Iat,
		"exp": payload.Exp,
	})
	t.Header = headerMap

	signingKey, err := loadPrivateKeyFromString(PKEY)

	if err != nil {
		return ""
	}

	token, err := t.SignedString(signingKey)

	if err != nil {
		return ""
	}

	am.Token = token
	am.Header = header
	am.Payload = payload

	return token
}

func (am *AppleMapsService) IsValidToken() bool {
	now := time.Now().Unix()

	if am.Token == "" || am.AmToken == "" {
		return false
	}

	if now > am.AmExpiry {
		return false
	}

	if now > am.Payload.Exp {
		return false
	} else {
		return true
	}
}

func (am *AppleMapsService) GenerateMapsToken() string {
	tokenEndpoint := am.tokenUrl
	token := am.GetToken()

	req, err := http.NewRequest(http.MethodGet, tokenEndpoint, nil)

	if err != nil {
		return ""
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Do(req)

	if err != nil {
		return ""
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		fmt.Println("Apple API response:", bodyString)
		return ""
	}

	var result AppleTokenResponse

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return ""
	}

	am.AmToken = result.AccessToken
	am.AmExpiry = time.Now().Unix() + int64(result.ExpiresInSeconds)

	return result.AccessToken
}

func loadPrivateKeyFromString(pemString string) (*ecdsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(pemString))
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	ecdsaKey, ok := privateKey.(*ecdsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("not an ECDSA private key")
	}

	return ecdsaKey, nil
}

type TokenHeader struct {
	Alg string `json:"alg"`
	Kid string `json:"kid"`
	Typ string `json:"typ"`
}

type TokenPayload struct {
	Iss string `json:"iss"`
	Iat int64  `json:"iat"`
	Exp int64  `json:"exp"`
}

type AppleTokenResponse struct {
	AccessToken      string `json:"accessToken"`
	ExpiresInSeconds int    `json:"expiresInSeconds"`
}

type GeocodeResponse struct {
	Results []GeocodeResult `json:"results"`
}

type GeocodeResult struct {
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
