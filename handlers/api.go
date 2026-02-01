package handlers

import (
	"net/http"
	"waldi-v2/components"
)

func AddressForm(w http.ResponseWriter, r *http.Request) {
	address := r.FormValue("address")
	components.AddressForm(address).Render(r.Context(), w)
}

func BrowserLocation(w http.ResponseWriter, r *http.Request) {
	lat := r.FormValue("lat")
	lon := r.FormValue("lon")
	components.AddressForm(lat+", "+lon).Render(r.Context(), w)
}
