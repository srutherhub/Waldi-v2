package handlers

import (
	"net/http"
	"os"
	pages "waldi-v2/components/pages"
	"waldi-v2/services"
)

func Homepage(w http.ResponseWriter, r *http.Request) {
	page := pages.Home()
	pages.Index(page).Render(r.Context(), w)
}

func Resultpage(a *services.AddressService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		mapKey := os.Getenv("APPLE_MAPKIT_KEY")

		id := r.PathValue("id")
		lat, lon, _ := a.DecodeCoords(id)

		address, _ := a.GetAddressFromCoords(lat, lon)

		props := pages.ResultProps{ApiKey: mapKey, Lat: lat, Lon: lon, Address: address}

		page := pages.Result(props)
		pages.Index(page).Render(r.Context(), w)
	}
}
