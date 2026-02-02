package handlers

import (
	"net/http"
	"os"
	"strconv"
	pages "waldi-v2/components/pages"
	"waldi-v2/services"
)

func Homepage(w http.ResponseWriter, r *http.Request) {
	page := pages.Home()
	pages.Index(page).Render(r.Context(), w)
}

func Resultpage(a *services.AddressService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		mapKey := os.Getenv("MAP_KEY")

		id := r.PathValue("id")
		lat, lon, _ := a.DecodeCoords(id)

		latStr := strconv.FormatFloat(lat, 'f', 8, 64)
		lonStr := strconv.FormatFloat(lon, 'f', 8, 64)
		address, _ := a.GetAddressFromCoords(lat, lon)

		props := pages.ResultProps{ApiKey: mapKey, Lat: latStr, Lon: lonStr, Address: address}

		page := pages.Result(props)
		pages.Index(page).Render(r.Context(), w)
	}
}
