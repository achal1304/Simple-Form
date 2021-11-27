package main

import "net/http"

func (app *application) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Origin")
		w.Header().Add("Vary", "Access-Control-Request-Method")

		w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, PUT, PATCH, DELETE,POST")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")

		next.ServeHTTP(w, r)
	})
}
