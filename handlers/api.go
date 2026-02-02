package handlers

import (
	"net/http"
	"strconv"
	"waldi-v2/components"
	"waldi-v2/services"
)

func AddressForm(a *services.AddressService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		address := r.FormValue("address")

		lat, lon, err := a.GetCoordsFromAddress(address)

		if err != nil {
			w.Header().Set("HX-Redirect", "/")
			return
		}

		urlEncoded := a.EncodeCoords(lat, lon)

		w.Header().Set("HX-Redirect", "/result/"+urlEncoded)
		w.WriteHeader(http.StatusOK)
		return
	}
}

func BrowserLocation(a *services.AddressService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lat := r.FormValue("lat")
		lon := r.FormValue("lon")

		latfloat, _ := strconv.ParseFloat(lat, 64)
		lonfloat, _ := strconv.ParseFloat(lon, 64)

		address, _ := a.GetAddressFromCoords(latfloat, lonfloat)

		components.AddressForm(address).Render(r.Context(), w)
	}
}
