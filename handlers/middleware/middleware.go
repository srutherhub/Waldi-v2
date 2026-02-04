package handlers

import "net/http"

func SetCacheHeader(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=3600")
		next(w, r)
	}
}

func SetStaticCacheHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		
		next.ServeHTTP(w, r)
	})
}
