package handlers

import (
	"net/http"
	pages "waldi-v2/components/pages"
)

func Homepage(w http.ResponseWriter, r *http.Request) {
	page := pages.Home()
	pages.Index(page).Render(r.Context(), w)
}
