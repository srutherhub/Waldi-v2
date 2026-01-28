package handlers

import (
	"net/http"
	"waldi-v2/components"
)

func SubmitLocation(w http.ResponseWriter, r *http.Request) {
	address := r.FormValue("address")

	components.LocationForm(address).Render(r.Context(), w)
}
