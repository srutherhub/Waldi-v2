package tests

import (
	"testing"
	"waldi-v2/services"
)

func TestAddressService(t *testing.T) {
	am := services.NewAppleMapsService()
	amc :=services.NewAppleMapsClient(am)
	a := services.NewAddressService(amc)

	encoded := a.EncodeCoords(43.0, -89.0)

	if encoded == "" {
		t.Fatal("AddressService.EncodeCoords failed to encode 43,-89: ")
	}

	_, _, err := a.DecodeCoords(encoded)

	if err != nil {
		t.Fatal("AddressService.DecodeCoords failed to decode encoded id: " + err.Error())
	}

}


func TestAppleMapService(t *testing.T){
	am := services.NewAppleMapsService()

	token:=am.GenerateToken()

	if token=="" {
		t.Fatal("AppleMapsService.GenerateToken failed to generate token")
	}
}