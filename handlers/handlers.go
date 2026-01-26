package handlers

import (
	"net/http"
	"waldi-v2/components/pages"
)

func Homepage(w http.ResponseWriter, r *http.Request) {
	page := components.Home()
	components.Index(page).Render(r.Context(), w)
}
