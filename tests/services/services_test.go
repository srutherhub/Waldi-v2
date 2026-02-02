package tests

import (
	"testing"
	"waldi-v2/services"
)

func TestAddressService(t *testing.T) {
	a := &services.AddressService{}

	encoded := a.EncodeCoords(43.0, -89.0)

	if encoded == "" {
		t.Fatal("AddressService.EncodeCoords failed to encode 43,-89: ")
	}

	_, _, err := a.DecodeCoords(encoded)

	if err != nil {
		t.Fatal("AddressService.DecodeCoords failed to decode encoded id: " + err.Error())
	}

}
