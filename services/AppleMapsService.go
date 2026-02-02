package services

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AppleMapsService struct {
	Token    string
	AmToken  string
	AmExpiry int64
	Header   TokenHeader
	Payload  TokenPayload
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

func NewAppleMapsService() *AppleMapsService {
	return &AppleMapsService{}
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
		fmt.Println(err.Error())
	}

	token, err := t.SignedString(signingKey)

	if err != nil {
		fmt.Println(err.Error())
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
	tokenEndpoint := os.Getenv("APPLE_TOKEN_URL")
	token := am.GetToken()

	req, err := http.NewRequest(http.MethodGet, tokenEndpoint, nil)

	if err != nil {
		fmt.Println(err.Error())
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Do(req)
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
