package services

import (
	"encoding/base64"
	"fmt"
)

type IAddressService interface {
	EncodeCoords(float64, float64) string
	DecodeCoords(string) (float64, float64, error)
	GetAddressFromCoords(float64, float64) (string, error)
	GetCoordsFromAddress(string) (float64, float64, error)
	GetNearbyLocations(float64, float64) ([]PointOfInterest, error)
}

type IAddressClient interface {
	Geocode(string) (float64, float64, error)
	ReverseGeocode(float64, float64) (string, error)
	Search(float64, float64) ([]PointOfInterest, error)
}

type AddressService struct {
	Mc IAddressClient
}

type PointOfInterest struct {
	Id        string
	IdType    string
	Name      string
	Category  PlaceCategory
	Address   string
	Longitude float64
	Latitude  float64
}

type PlaceCategory int

const (
	PlaceCategoryCafe PlaceCategory = iota
	PlaceCategoryRestaurant
	PlaceCategoryPark
	PlaceCategoryGrocery
)

func NewAddressService(mc IAddressClient) *AddressService {
	return &AddressService{Mc: mc}
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

	address, err := a.Mc.ReverseGeocode(lat, lon)

	if err != nil {
		return "", err
	}

	return address, nil
}

func (a *AddressService) GetCoordsFromAddress(address string) (lat float64, lon float64, err error) {

	lat, lon, err = a.Mc.Geocode(address)

	if err != nil {
		return 0, 0, err
	}

	return lat, lon, nil
}

func (a *AddressService) GetNearbyLocations(lat, lon float64) ([]PointOfInterest, error) {
	res, err := a.Mc.Search(lat, lon)
	if err != nil {
		return []PointOfInterest{}, nil
	}

	return res, nil
}
